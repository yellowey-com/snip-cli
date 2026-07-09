package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellowey-com/snip-cli/pkg/cli"
	"github.com/yellowey-com/snip-cli/pkg/storage"
	"github.com/yellowey-com/snip-cli/pkg/ui"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	dirPath := homeDir + "/.config/snip/snippets"
	_ = os.MkdirAll(dirPath, 0o755)

	filterQuery, filename, shouldRunUI := cli.Execute(os.Args[1:], dirPath)
	if !shouldRunUI {
		return
	}

	var items []list.Item

	if filename != "" {
		content, err := storage.ReadSnippet(dirPath, filename)
		if err != nil {
			fmt.Printf("Error: File '%s' not found.\n", filename)
			os.Exit(1)
		}
		for _, snip := range storage.ParseSnippetFile(content) {
			items = append(items, ui.NewItem(snip, filename))
		}
	} else {
		files, _ := storage.ListSnippets(dirPath)
		for _, file := range files {
			content, _ := storage.ReadSnippet(dirPath, file)
			for _, snip := range storage.ParseSnippetFile(content) {
				items = append(items, ui.NewItem(snip, file))
			}
		}
	}

	if len(items) == 0 {
		fmt.Println("No snippets found.")
		os.Exit(0)
	}

	p := tea.NewProgram(ui.NewModel(items, dirPath, filterQuery), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if finalModel, ok := m.(ui.Model); ok && finalModel.Selected != "" {
		finalCommand, err := cli.ResolvePlaceholders(finalModel.Selected)
		if err != nil {
			fmt.Println("\n✕ Cancelled")
			os.Exit(0)
		}

		if finalModel.Execute {
			cmd := exec.Command("/bin/sh", "-c", finalCommand)
			cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
			cmd.Run()
		} else {
			clipboard.WriteAll(finalCommand)
			fmt.Printf("✓ Copied: %s\n", finalCommand)
		}
	}
}
