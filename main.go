package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shceph/nasm-fmt/formatter"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("No file provided. Aborting...")
	}

	// tokens := formatter.Tokenize(file)
	// formatter.PrintTokens(tokens)

	output, err := formatter.Format(os.Args[1], formatter.DefaultFormatOpts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
