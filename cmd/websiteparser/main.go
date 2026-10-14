package main

import (
	"fmt"
	"net/http"
	"github.com/Catalyn45/laxpars"
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

	parser := laxpars.NewParser(resp.Body)

	root := parser.Parse()

	fmt.Println(root.StringTree())
}