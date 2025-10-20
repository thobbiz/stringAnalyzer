package main

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const ErrInvalidType = "json: cannot unmarshal number into Go struct field stringRequest.value of type string"

type stringRequest struct {
	Text string `json:"value"`
}

type getStringRequest struct {
	Text string `uri:"string_value" binding:"required"`
}

type getStringWithParamRequest struct {
	IsPalindrome      *bool  `form:"is_palindrome" json:"is_palindrome"`
	MinLength         *int   `form:"min_length" json:"min_length"`
	MaxLength         *int   `form:"max_length" json:"max_length"`
	WordCount         *int   `form:"word_count" json:"word_count"`
	ContainsCharacter string `form:"contains_character" json:"contains_character"`
}

type getStringWithParamResult struct {
	Data    []ObjectString            `json:"data"`
	Count   int                       `json:"count"`
	Filters getStringWithParamRequest `json:"filters_applied"`
}

type naturalLanguageRequest struct {
	Original      string                    `json:"original"`
	ParsedFilters getStringWithParamRequest `json:"parsed_filters"`
}

type naturalLanguageResponse struct {
	Data    []ObjectString         `json:"data"`
	Count   int                    `json:"count"`
	Filters naturalLanguageRequest `json:"interpreted_query"`
}

type deleteStringRequest struct {
	Text string `uri:"string_value" binding:"required"`
}

func createString(ctx *gin.Context) {
	var req stringRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if err.Error() == ErrInvalidType {
			ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	re := regexp.MustCompile(`\s+`)
	text := re.ReplaceAllString(req.Text, "")
	if text == "" {
		err := errors.New("invalid request body")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := ObjectString{
		Id:    hash(req.Text),
		Value: req.Text,
		Properties: Properties{
			StringLength:     length(req.Text),
			IsPalindrome:     is_palindrome(req.Text),
			UniqueCharacater: uniqueCharacters(req.Text),
			WordCount:        wordCount(req.Text),
			Hash:             hash(req.Text),
			CharFreqMap:      characterFreq(req.Text),
		},
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if _, exists := ObjectDB[req.Text]; exists {
		err := errors.New("String already exists")
		ctx.JSON(http.StatusConflict, errorResponse(err))
		return
	}
	ObjectDB[req.Text] = arg
	ctx.JSON(http.StatusCreated, arg)
}

func getString(ctx *gin.Context) {
	var req getStringRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		if err.Error() == ErrInvalidType {
			ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if value, exists := ObjectDB[req.Text]; !exists {
		err := errors.New("String doesn't exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else {
		ctx.JSON(http.StatusOK, value)
		return
	}
}

func getStringWithParams(ctx *gin.Context) {
	var req getStringWithParamRequest
	var result getStringWithParamResult
	if err := ctx.ShouldBindQuery(&req); err != nil {
		if err.Error() == ErrInvalidType {
			ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	for _, value := range ObjectDB {
		if valid(value, req) {
			result.Data = append(result.Data, value)
		}
	}

	result.Count = len(result.Data)
	result.Filters = req

	ctx.JSON(http.StatusOK, result)
}

func naturalLanguageString(ctx *gin.Context) {
	var req naturalLanguageRequest
	var result naturalLanguageResponse

	if err := ctx.ShouldBindUri(&req); err != nil {
		if err.Error() == ErrInvalidType {
			ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	filters := parseNaturalLanguageText(req.Original)
	req.ParsedFilters = filters

	for _, value := range ObjectDB {
		if valid(value, req.ParsedFilters) {
			result.Data = append(result.Data, value)
		}
	}

	result.Count = len(result.Data)
	result.Filters = req

	ctx.JSON(http.StatusOK, result)
}

func deleteString(ctx *gin.Context) {
	var req deleteStringRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		if err.Error() == ErrInvalidType {
			ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, exists := ObjectDB[req.Text]; !exists {
		err := errors.New("String doesn't exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else {
		delete(ObjectDB, req.Text)
		ctx.JSON(http.StatusNoContent, gin.H{})
		return
	}
}

func parseNaturalLanguageText(text string) getStringWithParamRequest {
	request := naturalLanguageRequest{}
	lowerText := strings.ToLower(text)

	if strings.Contains(lowerText, "palindrome") || strings.Contains(lowerText, "palindromic") {
		isPalindrome := true
		request.ParsedFilters.IsPalindrome = &isPalindrome
	}

	wordCountPatterns := []struct {
		pattern *regexp.Regexp
		value   int
	}{
		{regexp.MustCompile(`\bsingle\s+word\b`), 1},
		{regexp.MustCompile(`\bone\s+word\b`), 1},
		{regexp.MustCompile(`\btwo\s+words?\b`), 2},
		{regexp.MustCompile(`\bthree\s+words?\b`), 3},
		{regexp.MustCompile(`\bfour\s+words?\b`), 4},
		{regexp.MustCompile(`\bfive\s+words?\b`), 5},
	}

	for _, p := range wordCountPatterns {
		if p.pattern.MatchString(lowerText) {
			request.ParsedFilters.WordCount = &p.value
			break
		}
	}

	numWordPattern := regexp.MustCompile(`\b(\d+)\s+words?\b`)
	if matches := numWordPattern.FindStringSubmatch(lowerText); len(matches) > 1 {
		if count, err := strconv.Atoi(matches[1]); err == nil {
			request.ParsedFilters.WordCount = &count
		}
	}

	// Check for length constraints
	minLengthPattern := regexp.MustCompile(`(?:longer than|more than|at least)\s+(\d+)\s+(?:characters|chars|letters)`)
	if matches := minLengthPattern.FindStringSubmatch(lowerText); len(matches) > 1 {
		if length, err := strconv.Atoi(matches[1]); err == nil {
			request.ParsedFilters.MinLength = &length
		}
	}

	maxLengthPattern := regexp.MustCompile(`(?:shorter than|less than|at most)\s+(\d+)\s+(?:characters|chars|letters)`)
	if matches := maxLengthPattern.FindStringSubmatch(lowerText); len(matches) > 1 {
		if length, err := strconv.Atoi(matches[1]); err == nil {
			request.ParsedFilters.MaxLength = &length
		}
	}

	containsPattern := regexp.MustCompile(`contains?\s+["']?(\w+)["']?`)
	if matches := containsPattern.FindStringSubmatch(lowerText); len(matches) > 1 {
		request.ParsedFilters.ContainsCharacter = matches[1]
	}
	return request.ParsedFilters
}

func containsCharacter(character string, text string) bool {
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, "")

	for _, ch := range text {
		chStr := string(ch)
		if chStr == character {
			return true
		}
	}

	return false
}

func valid(value ObjectString, req getStringWithParamRequest) bool {
	p := value.Properties

	if req.IsPalindrome != nil && p.IsPalindrome != *req.IsPalindrome {
		return false
	}
	if req.MinLength != nil && p.StringLength < *req.MinLength {
		return false
	}
	if req.MaxLength != nil && p.StringLength > *req.MaxLength {
		return false
	}
	if req.WordCount != nil && p.WordCount != *req.WordCount {
		return false
	}
	if req.ContainsCharacter != "" && !containsCharacter(req.ContainsCharacter, value.Value) {
		return false
	}

	return true
}
