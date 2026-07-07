package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellowey-com/snip-cli/pkg/storage"
	"github.com/yellowey-com/snip-cli/pkg/ui"
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
		if len(snippets) == 0 {
			fmt.Printf("No snippets found in file '%s'.\n", filename)
			os.Exit(0)
		}

		p := tea.NewProgram(ui.NewModel(snippets), tea.WithAltScreen())
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Error running UI: %v\n", err)
			os.Exit(1)
		}

		if finalModel, ok := m.(ui.Model); ok && finalModel.Selected != "" {
			fmt.Println(finalModel.Selected)
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
