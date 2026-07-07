package storage

import (
	"bufio"
	"os"
	"strings"
)

type Snippet struct {
	Description string
	Command     string
}

func ListSnippets(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var snippetFiles []string
	for _, file := range files {
		if !file.IsDir() {
			snippetFiles = append(snippetFiles, file.Name())
		}
	}

	return snippetFiles, nil
}

func ReadSnippet(dirPath, filename string) (string, error) {
	filePath := dirPath + "/" + filename
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ParseSnippetFile(content string) []Snippet {
	var list []Snippet
	var currentDesc string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if after, found := strings.CutPrefix(line, "##"); found {
			currentDesc = after
			continue
		}

		if line != "" && !strings.HasPrefix(line, "#") && currentDesc != "" {
			list = append(list, Snippet{
				Description: currentDesc,
				Command:     line,
			})
			currentDesc = ""
		}
	}
	return list
}
