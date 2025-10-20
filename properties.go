package main

import (
	"regexp"
	"strings"
)

func length(text string) int {
	return len(text)
}

func is_palindrome(text string) bool {
	low := 0
	high := length(text) - 1
	for low < high {
		if text[low] != text[high] {
			return false
		}
		low++
		high--
	}
	return true
}

func uniqueCharacters(text string) int {
	re := regexp.MustCompile(`s+`)
	text = re.ReplaceAllString(text, "")

	countMap := make(map[rune]bool)
	for _, ch := range text {
		countMap[ch] = true
	}
	return len(countMap)
}

func wordCount(text string) int {
	text = strings.TrimSpace(text)

	if len(text) == 0 {
		return 0
	}

	wordCount := 0
	inWord := false

	for _, ch := range text {
		if ch != ' ' && ch != '\t' && ch != '\n' {
			if !inWord {
				wordCount++
				inWord = true
			}
		} else {
			inWord = false
		}
	}

	return wordCount
}
