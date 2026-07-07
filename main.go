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

		snippets := storage.ParseSnippetFile(content)

		fmt.Printf("Parsed %d snippets from %s:\n", len(snippets), filename)
		fmt.Println("==================================================")
		for i, snip := range snippets {
			fmt.Printf("[%d] Desc:    %s\n", i+1, snip.Description)
			fmt.Printf("    Command: %s\n\n", snip.Command)
		}
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
