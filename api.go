package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type stringRequest struct {
	Text string `json:"value"`
}

func createString(ctx *gin.Context) {
	var req stringRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	ObjectDB["req.Text"] = arg
	ctx.JSON(http.StatusCreated, arg)
}
