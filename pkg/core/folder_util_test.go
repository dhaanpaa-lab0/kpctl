package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHomeDir(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Tilde only",
			input:    "~",
			expected: tempHome,
		},
		{
			name:     "Tilde with slash and path",
			input:    "~/sub/dir",
			expected: filepath.Join(tempHome, "sub", "dir"),
		},
		{
			name:     "Tilde without slash",
			input:    "~sub/dir",
			expected: filepath.Join(tempHome, "sub", "dir"),
		},
		{
			name:     "Standard absolute path",
			input:    "/tmp/some/path",
			expected: "/tmp/some/path",
		},
		{
			name:     "Standard relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandHomeDir(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ExpandHomeDir(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateDirectoryTree(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	targetDir := "~/a/b/c/d"
	err := CreateDirectoryTree(targetDir)
	if err != nil {
		t.Fatalf("CreateDirectoryTree(%q) returned error: %v", targetDir, err)
	}

	expectedPath := filepath.Join(tempHome, "a", "b", "c", "d")
	info, err := os.Stat(expectedPath)
	if err != nil {
		t.Fatalf("os.Stat(%q) failed: %v", expectedPath, err)
	}
	if !info.IsDir() {
		t.Errorf("expected %q to be a directory", expectedPath)
	}
}

func TestCreateFolderTree(t *testing.T) {
	tempDir := t.TempDir()
	targetDir := filepath.Join(tempDir, "one", "two", "three")

	err := CreateFolderTree(targetDir)
	if err != nil {
		t.Fatalf("CreateFolderTree(%q) returned error: %v", targetDir, err)
	}

	info, err := os.Stat(targetDir)
	if err != nil {
		t.Fatalf("os.Stat(%q) failed: %v", targetDir, err)
	}
	if !info.IsDir() {
		t.Errorf("expected %q to be a directory", targetDir)
	}
}
