package storage

import (
	"reflect"
	"testing"
)

func TestParseSnippetFile(t *testing.T) {
	content := `
# My Snippets

## List files
ls -la

## Remove file
rm -f file.txt
`
	expected := []Snippet{
		{Description: "List files", Command: "ls -la"},
		{Description: "Remove file", Command: "rm -f file.txt"},
	}

	result := ParseSnippetFile(content)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestParseSnippetFile_Empty(t *testing.T) {
	content := "# Just a header\n\n# Comment"
	result := ParseSnippetFile(content)

	if len(result) != 0 {
		t.Errorf("expected 0 snippets, got %d", len(result))
	}
}
