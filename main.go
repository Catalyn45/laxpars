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
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status:", resp.Status)
		return
	}

	lexer := NewLexer(resp.Body)
	parser := NewParser(lexer)

	root := parser.Parse()

	fmt.Println(root.StringTree())
}