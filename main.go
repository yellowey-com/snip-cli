package main

import (
	"fmt"
	"os"

	"github.com/yellowey-com/snip-cli/pkg/storage"
)

func main() {
	dirPath := "snippets"

	if len(os.Args) > 1 {
		filename := os.Args[1]
		content, err := storage.ReadSnippet(dirPath, filename)
		if err != nil {
			fmt.Printf("Error: Could not read file '%s'. Make sure it exists.\n", filename)
			os.Exit(1)
		}

		fmt.Println(content)
		return
	}

	files, err := storage.ListSnippets(dirPath)
	if err != nil {
		fmt.Printf("Error: Could not read directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Available snippet categories:")
	fmt.Println("----------------------------")

	for _, file := range files {
		fmt.Printf("- %s\n", file)
	}
}
