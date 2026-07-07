package storage

import (
	"os"
)

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
