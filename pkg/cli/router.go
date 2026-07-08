package cli

import (
	"fmt"
)

func Execute(args []string, dirPath string) {
	if len(args) < 1 {
		fmt.Println("Usage: snip <command> [args]")
		return
	}

	switch args[0] {
	case "find":
		if len(args) < 2 {
			fmt.Println("Usage: snip find <query>")
			return
		}
		FindCommand(dirPath, args[1])
	case "run":
		if len(args) < 2 {
			fmt.Println("Usage: snip run <description>")
			return
		}
		RunCommand(dirPath, args[1])
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
	}
}
