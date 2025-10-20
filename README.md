## String Analyzer

StringAnalyzer is a backend service built with Golang that collects strings and returns fact about each string to the user.

## Features:
- rate limiting using `golang.org/x/time/rate`
- error handling

## Stack and Tools:
- Golang
- Gin Web Framework
- net/http

## Get Started
### Prerequisites
- G0 1.21+
- Internet Connection
- Modules:
  ```bash
  go get github.com/gin-gonic/gin

## Usage
- Start the server:
  ```bash
  go run main.go
- Send a string:
  ```bash
  curl http://localhost:7070/strings

- Get a string:
  ```bash
  curl http://localhost:7070/strings/:string_value
  
- Delete a string:
  ```bash
  curl http://localhost:7070/strings/:string_value
  
- Get a string by nlp:
  ```bash
  curl http://localhost:7070/strings/filter-by-natural-language
  
- Get a set of strings with specific properties:
  ```bash
  curl http://localhost:7070/strings?query

