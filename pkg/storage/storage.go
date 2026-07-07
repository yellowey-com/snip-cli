package storage

import (
	"bufio"
	"fmt"
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

func AppendSnippet(dirPath, filename, desc, command string) error {
	filePath := dirPath + "/" + filename

	snippetBlock := fmt.Sprintf("\n## %s\n\n%s\n", strings.TrimSpace(desc), strings.TrimSpace(command))

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(snippetBlock)
	return err
}

func RemoveSnippet(dirPath, filename, desc string) error {
	content, err := ReadSnippet(dirPath, filename)
	if err != nil {
		return err
	}

	snippets := ParseSnippetFile(content)
	var updated []Snippet

	for _, snip := range snippets {
		if snip.Description != desc {
			updated = append(updated, snip)
		}
	}

	return saveAllSnippets(dirPath, filename, updated)
}

func EditSnippet(dirPath, filename, desc, newCommand string) error {
	content, err := ReadSnippet(dirPath, filename)
	if err != nil {
		return err
	}

	snippets := ParseSnippetFile(content)
	changed := false

	for i, snip := range snippets {
		if snip.Description == desc {
			snippets[i].Command = newCommand
			changed = true
			break
		}
	}

	if !changed {
		return fmt.Errorf("snippet with description '%s' not found", desc)
	}

	return saveAllSnippets(dirPath, filename, snippets)
}

func saveAllSnippets(dirPath, filename string, snippets []Snippet) error {
	filePath := dirPath + "/" + filename
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString("# " + strings.TrimSuffix(filename, ".md") + " Snippets\n")

	for _, snip := range snippets {
		block := fmt.Sprintf("\n## %s\n\n%s\n", snip.Description, snip.Command)
		_, _ = writer.WriteString(block)
	}

	return writer.Flush()
}
