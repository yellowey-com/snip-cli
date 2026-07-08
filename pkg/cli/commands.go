package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yellowey-com/snip-cli/pkg/storage"
)

const dirPath = "./snippets"

func FindCommand(query string) {
	snippets, err := storage.LoadSnippets(dirPath)
	if err != nil {
		fmt.Printf("Error loading snippets: %v\n", err)
		return
	}

	found := false
	for _, s := range snippets {
		if strings.Contains(strings.ToLower(s.Description), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(s.Command), strings.ToLower(query)) {
			fmt.Printf("[%s] %s -> %s\n", s.Category, s.Description, s.Command)
			found = true
		}
	}
	if !found {
		fmt.Println("No snippets found.")
	}
}

func RunCommand(descQuery string) {
	snippets, _ := storage.LoadSnippets(dirPath)

	for _, s := range snippets {
		if strings.Contains(strings.ToLower(s.Description), strings.ToLower(descQuery)) {
			fmt.Printf("Running: %s\n", s.Command)
			cmd := exec.Command("sh", "-c", s.Command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			return
		}
	}
	fmt.Println("Snippet not found.")
}
