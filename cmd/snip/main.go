package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellowey-com/snip-cli/pkg/cli"
	"github.com/yellowey-com/snip-cli/pkg/storage"
	"github.com/yellowey-com/snip-cli/pkg/ui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "reload" {
		reloadSelf()
		return
	}
	RunUI()
}

func reloadSelf() {
	build := exec.Command("go", "install", "./cmd/snip")
	build.Stdout, build.Stderr = os.Stdout, os.Stderr
	if err := build.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}

	gopath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		fmt.Printf("Error resolving GOPATH: %v\n", err)
		os.Exit(1)
	}
	binPath := filepath.Join(strings.TrimSpace(string(gopath)), "bin", "snip")

	args := append([]string{binPath}, os.Args[2:]...)
	if err := syscall.Exec(binPath, args, os.Environ()); err != nil {
		fmt.Printf("Exec failed: %v\n", err)
		os.Exit(1)
	}
}

func RunUI() {
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
