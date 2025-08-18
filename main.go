package main

import (
	"log"
	"os"

	"github.com/shceph/nasm-fmt/formatter"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("No file provided. Aborting...")
	}

	_, err := formatter.Format(os.Args[1], formatter.DefaultFormatOpts)

	if err != nil {
		log.Fatal(err)
	}
}
