package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellowey-com/snip-cli/pkg/storage"
	"github.com/yellowey-com/snip-cli/pkg/ui"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	dirPath := homeDir + "/.config/snip/snippets"
	_ = os.MkdirAll(dirPath, 0o755)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snip add <file.md> <description> <command>")
				os.Exit(1)
			}
			err := storage.AppendSnippet(dirPath, os.Args[2], os.Args[3], os.Args[4])
			if err != nil {
				fmt.Printf("Error adding snippet: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ Snippet added to %s\n", os.Args[2])
			return

		case "remove":
			if len(os.Args) < 4 {
				fmt.Println("Usage: snip remove <file.md> <description>")
				os.Exit(1)
			}
			err := storage.RemoveSnippet(dirPath, os.Args[2], os.Args[3])
			if err != nil {
				fmt.Printf("Error removing snippet: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ Snippet removed from %s\n", os.Args[2])
			return

		case "edit":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snip edit <file.md> <description> <new_command>")
				os.Exit(1)
			}
			err := storage.EditSnippet(dirPath, os.Args[2], os.Args[3], os.Args[4])
			if err != nil {
				fmt.Printf("Error editing snippet: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ Snippet updated in %s\n", os.Args[2])
			return
		}
	}

	var items []list.Item

	if len(os.Args) > 1 {
		filename := os.Args[1]
		content, err := storage.ReadSnippet(dirPath, filename)
		if err != nil {
			fmt.Printf("Error: Could not read file '%s'. Make sure it exists.\n", filename)
			os.Exit(1)
		}

		snippets := storage.ParseSnippetFile(content)
		for _, snip := range snippets {
			items = append(items, ui.NewItem(snip, filename))
		}
	} else {
		files, err := storage.ListSnippets(dirPath)
		if err != nil {
			fmt.Printf("Error: Could not read directory: %v\n", err)
			os.Exit(1)
		}

		for _, file := range files {
			content, err := storage.ReadSnippet(dirPath, file)
			if err != nil {
				continue
			}
			snippets := storage.ParseSnippetFile(content)
			for _, snip := range snippets {
				items = append(items, ui.NewItem(snip, file))
			}
		}
	}

	if len(items) == 0 {
		fmt.Println("No snippets found. Use 'snip add <file.md> <desc> <cmd>' to add some.")
		os.Exit(0)
	}

	p := tea.NewProgram(ui.NewModel(items), tea.WithAltScreen())
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
}
