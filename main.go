package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atotto/clipboard"
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
			if finalModel.Execute {
				fmt.Printf("Executing: %s\n\n", finalModel.Selected)
				cmd := exec.Command("/bin/sh", "-c", finalModel.Selected)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin

				if err := cmd.Run(); err != nil {
					fmt.Printf("\nExecution failed: %v\n", err)
					os.Exit(1)
				}
			} else {
				err := clipboard.WriteAll(finalModel.Selected)
				if err != nil {
					fmt.Printf("Error copying to clipboard: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("✓ Copied to clipboard: %s\n", finalModel.Selected)
			}
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
