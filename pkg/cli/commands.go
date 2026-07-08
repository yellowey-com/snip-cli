package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yellowey-com/snip-cli/pkg/storage"
)

type foundSnippet struct {
	storage.Snippet
	Category string
}

func loadAllSnippets(dirPath string) ([]foundSnippet, error) {
	files, err := storage.ListSnippets(dirPath)
	if err != nil {
		return nil, err
	}

	var found []foundSnippet
	for _, file := range files {
		content, err := storage.ReadSnippet(dirPath, file)
		if err != nil {
			continue
		}
		for _, s := range storage.ParseSnippetFile(content) {
			found = append(found, foundSnippet{Snippet: s, Category: file})
		}
	}
	return found, nil
}

func FindCommand(dirPath, query string) {
	snippets, err := loadAllSnippets(dirPath)
	if err != nil {
		fmt.Printf("Error loading snippets: %v\n", err)
		return
	}

	query = strings.ToLower(query)
	found := false
	for _, s := range snippets {
		if strings.Contains(strings.ToLower(s.Description), query) ||
			strings.Contains(strings.ToLower(s.Command), query) {
			fmt.Printf("[%s] %s -> %s\n", s.Category, s.Description, s.Command)
			found = true
		}
	}
	if !found {
		fmt.Println("No snippets found.")
	}
}

func RunCommand(dirPath, descQuery string) {
	snippets, err := loadAllSnippets(dirPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	descQuery = strings.ToLower(descQuery)
	for _, s := range snippets {
		if strings.Contains(strings.ToLower(s.Description), descQuery) {
			fmt.Printf("Running: %s\n", s.Command)

			cmd := exec.Command("sh", "-c", s.Command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			if err := cmd.Run(); err != nil {
				fmt.Printf("Execution error: %v\n", err)
			}
			return
		}
	}
	fmt.Println("Snippet not found.")
}
