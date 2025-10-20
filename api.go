package main

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

const ErrInvalidType = "json: cannot unmarshal number into Go struct field stringRequest.value of type string"

type stringRequest struct {
	Text string `json:"value"`
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

type getStringRequest struct {
	Text string `uri:"string_value" binding:"required"`
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

type getStringWithParamRequest struct {
	IsPalindrome      *bool  `form:"is_palindrome"`
	MinLength         *int   `form:"min_length"`
	MaxLength         *int   `form:"max_length"`
	WordCount         *int   `form:"word_count"`
	ContainsCharacter string `form:"contains_character"`
}

type getStringWithParamResult struct {
	Data []ObjectString `form:"data"`
}

func getStringWithParams(ctx *gin.Context) {
	var req getStringWithParamRequest
	var result getStringWithParamResult
	if err := ctx.ShouldBindUri(&req); err != nil {
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

	ctx.JSON(http.StatusOK, result)
}

func valid(value ObjectString, req getStringWithParamRequest) bool {
	if (value.Properties.IsPalindrome == *req.IsPalindrome) &&
		(value.Properties.StringLength >= *req.MinLength) &&
		(value.Properties.StringLength <= *req.MaxLength) &&
		(value.Properties.WordCount == *req.WordCount) && containsCharacter(req.ContainsCharacter, value.Value) {
		return true
	} else {
		return false
	}
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
