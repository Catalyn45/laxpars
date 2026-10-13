package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get(os.Args[1])
	if err != nil {
		fmt.Println("Request error:", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return
	}

	lexer := NewLexer(resp.Body)
	parser := NewParser(lexer)

	parser.Parse()
}