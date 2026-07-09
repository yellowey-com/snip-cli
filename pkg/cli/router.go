package cli

import (
	"fmt"
	"os"

	"github.com/yellowey-com/snip-cli/pkg/storage"
)

func Execute(args []string, dirPath string) (filterQuery string, filename string, shouldRunUI bool) {
	if len(args) < 1 {
		return "", "", true
	}

	switch args[0] {
	case "--list":
		files, _ := storage.ListSnippets(dirPath)
		for _, file := range files {
			content, _ := storage.ReadSnippet(dirPath, file)
			for _, snip := range storage.ParseSnippetFile(content) {
				fmt.Println(snip.Description)
			}
		}
		return "", "", false

	case "completion":
		fmt.Println("find run add remove edit")
		return "", "", false

	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: snip run <description>")
			os.Exit(1)
		}
		RunCommand(dirPath, args[1])
		return "", "", false

	case "find":
		if len(args) < 2 {
			fmt.Println("Usage: snip find <query>")
			os.Exit(1)
		}
		return args[1], "", true

	case "add":
		if len(args) < 4 {
			fmt.Println("Usage: snip add <file.md> <description> <command>")
			os.Exit(1)
		}
		storage.AppendSnippet(dirPath, args[1], args[2], args[3])
		fmt.Printf("✓ Snippet added\n")
		return "", "", false

	case "remove":
		if len(args) < 3 {
			fmt.Println("Usage: snip remove <file.md> <description>")
			os.Exit(1)
		}
		storage.RemoveSnippet(dirPath, args[1], args[2])
		fmt.Printf("✓ Snippet removed\n")
		return "", "", false

	case "edit":
		if len(args) < 4 {
			fmt.Println("Usage: snip edit <file.md> <description> <new_command>")
			os.Exit(1)
		}
		storage.EditSnippet(dirPath, args[1], args[2], args[3])
		fmt.Printf("✓ Snippet updated\n")
		return "", "", false

	default:
		return "", args[0], true
	}
}
