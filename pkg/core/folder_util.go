package core

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandHomeDir replaces a leading '~' in the given path with the user's home directory.
func ExpandHomeDir(path string) (string, error) {
	if rest, ok := strings.CutPrefix(path, "~"); ok {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if rest == "" {
			return homeDir, nil
		}
		if strings.HasPrefix(rest, "/") || strings.HasPrefix(rest, "\\") {
			return filepath.Join(homeDir, rest[1:]), nil
		}
		return filepath.Join(homeDir, rest), nil
	}
	return path, nil
}

// CreateDirectoryTree creates a directory tree at the specified path.
// If the path starts with '~', the tilde is replaced with the user's home folder.
func CreateDirectoryTree(path string) error {
	expandedPath, err := ExpandHomeDir(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(expandedPath, 0755)
}

// CreateFolderTree is an alias for CreateDirectoryTree.
func CreateFolderTree(path string) error {
	return CreateDirectoryTree(path)
}
