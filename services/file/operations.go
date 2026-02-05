package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

var root, _ = os.UserHomeDir()

// Renaming file function
func RenameFile(sourcePath string, newName string) error {
	// New file name can't be empty
	if newName == "" || strings.Contains(newName, "/") {
		return os.ErrInvalid
	}

	// File must live inside root
	cleanPath := filepath.Clean(sourcePath)

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

// Delete file function
func DeleteFile(path string) error {
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

	// Is the destination exist?
	if _, err := os.Stat(absPath); err != nil {
		return os.ErrNotExist
	}

	// Perform the delete operation
	return os.Remove(absPath)
}

// Delete folder function
func DeleteFolder(path string) error {
	// File must live inside root
	cleanPath := filepath.Clean(path)

	absPath := filepath.Join(root, cleanPath)

	relPath, err := filepath.Rel(root, absPath)
	if err != nil {
		return os.ErrPermission
	}

	if relPath == "." || strings.HasPrefix(relPath, "..") {
		return os.ErrPermission
	}

	// Is the destination exist?
	if _, err := os.Stat(absPath); err != nil {
		return os.ErrNotExist
	}

	// Perform the delete operation
	return os.RemoveAll(absPath)
}

// Copy file function
func CopyFile(sourcePath string, destinationDir string) error {
	// Clean the paths
	cleanSourcePath := filepath.Clean(sourcePath)
	cleanDestinationDir := filepath.Clean(destinationDir)

	// Construct absolute paths
	destinationDirLocation := filepath.Join(root, cleanDestinationDir)
	sourceLocation := filepath.Join(root, cleanSourcePath)

	// Ensure source and destination are within root
	fileRelPath, err := filepath.Rel(root, sourceLocation)
	if err != nil || strings.HasPrefix(fileRelPath, "..") {
		return os.ErrPermission
	}

	desRelPath, err := filepath.Rel(root, destinationDirLocation)
	if err != nil || strings.HasPrefix(desRelPath, "..") {
		return os.ErrPermission
	}

	destinationFile := filepath.Join(destinationDirLocation, filepath.Base(sourcePath))

	// Check source file existence and is it not a directory?
	if f, err := os.Stat(sourceLocation); err != nil {
		return os.ErrNotExist
	} else if f.IsDir() {
		return os.ErrInvalid
	}

	// Check the destination file and is it a directory?
	if _, err := os.Stat(destinationFile); err == nil {
		return os.ErrExist
	}

	info, err := os.Stat(destinationDirLocation)
	if err != nil {
		return os.ErrNotExist
	}
	if !info.IsDir() {
		return os.ErrInvalid
	}

	// Perform the copy operation
	srcFile, err := os.Open(sourceLocation)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
