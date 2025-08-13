package main

import (
	"fmt"
	"os"

	"github.com/shceph/nasm-fmt/formatter"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No file provided. Aborting...")
		return
	}

	file, err := os.Open(os.Args[1])

	if err != nil {
		panic(err)
	}

	tokens := formatter.Tokenize(file)
	formatter.PrintTokens(tokens)
}
