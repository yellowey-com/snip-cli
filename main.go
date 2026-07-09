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

	var filterQuery string

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--list":
			files, _ := storage.ListSnippets(dirPath)
			for _, file := range files {
				content, _ := storage.ReadSnippet(dirPath, file)
				for _, snip := range storage.ParseSnippetFile(content) {
					fmt.Println(snip.Description)
				}
			}
			return
		case "completion":
			fmt.Println("find run add remove edit")
			return
		case "run":
			if len(os.Args) < 3 {
				fmt.Println("Usage: snip run <description>")
				os.Exit(1)
			}
			cli.RunCommand(dirPath, os.Args[2])
			return
		case "find":
			if len(os.Args) < 3 {
				fmt.Println("Usage: snip find <query>")
				os.Exit(1)
			}
			filterQuery = os.Args[2]
		case "add":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snip add <file.md> <description> <command>")
				os.Exit(1)
			}
			storage.AppendSnippet(dirPath, os.Args[2], os.Args[3], os.Args[4])
			fmt.Printf("✓ Snippet added\n")
			return
		case "remove":
			if len(os.Args) < 4 {
				fmt.Println("Usage: snip remove <file.md> <description>")
				os.Exit(1)
			}
			storage.RemoveSnippet(dirPath, os.Args[2], os.Args[3])
			fmt.Printf("✓ Snippet removed\n")
			return
		case "edit":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snip edit <file.md> <description> <new_command>")
				os.Exit(1)
			}
			storage.EditSnippet(dirPath, os.Args[2], os.Args[3], os.Args[4])
			fmt.Printf("✓ Snippet updated\n")
			return
		}
	}

	var items []list.Item

	if filterQuery == "" && len(os.Args) > 1 {
		filename := os.Args[1]
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
		if finalModel.Execute {
			cmd := exec.Command("/bin/sh", "-c", finalModel.Selected)
			cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
			cmd.Run()
		} else {
			clipboard.WriteAll(finalModel.Selected)
			fmt.Printf("✓ Copied: %s\n", finalModel.Selected)
		}
	}
}
