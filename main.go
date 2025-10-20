package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/strings", createString)
	router.Run(":7070")
}

var ObjectDB map[string]ObjectString

type ObjectString struct {
	Id         string     `json:"id"`
	Value      string     `json:"string"`
	Properties Properties `json:"properties"`
	CreatedAt  string     `json:"created_at"`
}

type Properties struct {
	StringLength     int            `json:"length"`
	IsPalindrome     bool           `json:"is_palindrome"`
	UniqueCharacater int            `json:"unique_characters"`
	WordCount        int            `json:"word_count"`
	Hash             string         `json:"sha256_hash"`
	CharFreqMap      map[string]int `json:"character_frequency_map"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
