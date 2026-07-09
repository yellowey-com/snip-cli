package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var placeholderRegex = regexp.MustCompile(`<([^>]+)>`)

func ResolvePlaceholders(command string) (string, error) {
	matches := placeholderRegex.FindAllStringSubmatch(command, -1)
	if len(matches) == 0 {
		return command, nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	updatedCommand := command

	for _, match := range matches {
		fullMatch := match[0]
		placeholderName := match[1]

		fmt.Printf("Enter value for [%s]: ", placeholderName)
		if !scanner.Scan() {
			return "", fmt.Errorf("cancelled")
		}

		input := scanner.Text()
		updatedCommand = strings.Replace(updatedCommand, fullMatch, input, 1)
	}

	return updatedCommand, nil
}
