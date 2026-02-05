package file

import (
	"os"
	"strings"
	"path/filepath"
)

var root, _ = os.UserHomeDir()

func RenameFile(path string, newName string) error {
	// New file name can't be empty
	if newName == "" || strings.Contains(newName, "/") {
		return os.ErrInvalid
	}

	// File must live inside root
	cleanPath := filepath.Clean(path)

	absPath := filepath.Join(root, cleanPath)

	relPath, err := filepath.Rel(root, absPath)
	if err != nil {
		return os.ErrPermission
	}

	if strings.HasPrefix(relPath, "..") {
		return os.ErrPermission
	}

	parentDir := filepath.Dir(absPath)

	newPath := filepath.Join(parentDir, newName)

	// Is the destination already exist?
	if _, err := os.Stat(newPath); err == nil {
		return os.ErrExist
	}

	// Perform the rename operation
	return os.Rename(absPath, newPath)
}