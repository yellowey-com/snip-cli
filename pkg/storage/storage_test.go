package storage

import (
	"os"
	"path/filepath"
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

func TestAppendSnippet(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "test.md"

	err := AppendSnippet(tmpDir, filename, "New Command", "echo 'hello'")
	if err != nil {
		t.Fatalf("failed to append snippet: %v", err)
	}

	content, err := ReadSnippet(tmpDir, filename)
	if err != nil {
		t.Fatalf("failed to read snippet: %v", err)
	}

	snippets := ParseSnippetFile(content)
	if len(snippets) != 1 {
		t.Fatalf("expected 1 snippet, got %d", len(snippets))
	}

	if snippets[0].Description != "New Command" || snippets[0].Command != "echo 'hello'" {
		t.Errorf("unexpected snippet content: %+v", snippets[0])
	}
}

func TestRemoveSnippet(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "test.md"

	initialContent := "# Test\n\n## First\n\ncmd1\n\n## Second\n\ncmd2\n"
	err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(initialContent), 0o644)
	if err != nil {
		t.Fatalf("failed to setup test file: %v", err)
	}

	err = RemoveSnippet(tmpDir, filename, "First")
	if err != nil {
		t.Fatalf("failed to remove snippet: %v", err)
	}

	content, err := ReadSnippet(tmpDir, filename)
	if err != nil {
		t.Fatalf("failed to read snippet file: %v", err)
	}

	snippets := ParseSnippetFile(content)
	if len(snippets) != 1 {
		t.Fatalf("expected 1 snippet, got %d", len(snippets))
	}

	if snippets[0].Description != "Second" {
		t.Errorf("expected only 'Second' snippet to remain, got %s", snippets[0].Description)
	}
}

func TestEditSnippet(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "test.md"

	initialContent := "# Test\n\n## Command To Edit\n\nold_command\n"
	err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(initialContent), 0o644)
	if err != nil {
		t.Fatalf("failed to setup test file: %v", err)
	}

	err = EditSnippet(tmpDir, filename, "Command To Edit", "new_command")
	if err != nil {
		t.Fatalf("failed to edit snippet: %v", err)
	}

	content, err := ReadSnippet(tmpDir, filename)
	if err != nil {
		t.Fatalf("failed to read snippet file: %v", err)
	}

	snippets := ParseSnippetFile(content)
	if len(snippets) != 1 {
		t.Fatalf("expected 1 snippet, got %d", len(snippets))
	}

	if snippets[0].Command != "new_command" {
		t.Errorf("expected command to be 'new_command', got '%s'", snippets[0].Command)
	}
}
