package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetAllowedRoots returns all allowed root directories for file operations.
// This includes the user's home directory and all mounted drives.
func GetAllowedRoots() []string {
	var roots []string

	// Add home directory
	if home, err := GetHomeDirectory(); err == nil {
		roots = append(roots, home)
	}

	// Add mounted drives (same filtering logic as DrivesDataParser)
	mounts, err := GetMountedDrives()
	if err == nil {
		for _, mount := range mounts {
			// Filter virtual filesystems
			if mount.FSType == "tmpfs" || mount.FSType == "devtmpfs" || mount.FSType == "proc" ||
				mount.FSType == "sysfs" || mount.FSType == "cgroup" || mount.FSType == "overlay" ||
				mount.FSType == "rootfs" || mount.FSType == "cgroup2" || mount.FSType == "debugfs" ||
				mount.FSType == "tracefs" || mount.FSType == "configfs" || mount.FSType == "binfmt_misc" ||
				mount.FSType == "fusectl" || mount.FSType == "hugetlbfs" || mount.FSType == "mqueue" ||
				mount.FSType == "pstore" || mount.FSType == "securityfs" || mount.FSType == "efivarfs" ||
				mount.FSType == "bpf" {
				continue
			}

			// Exclude system mount points
			if strings.HasPrefix(mount.MountPoint, "/proc") || strings.HasPrefix(mount.MountPoint, "/sys") ||
				strings.HasPrefix(mount.MountPoint, "/dev") || strings.HasPrefix(mount.MountPoint, "/run") ||
				strings.HasPrefix(mount.MountPoint, "/usr") || strings.HasPrefix(mount.MountPoint, "/var/run") ||
				strings.HasPrefix(mount.MountPoint, "/var/lib/docker") || strings.HasPrefix(mount.MountPoint, "/init") {
				continue
			}

			// Include root, /mnt/*, /media/*
			if mount.MountPoint == "/" || strings.HasPrefix(mount.MountPoint, "/mnt/") ||
				strings.HasPrefix(mount.MountPoint, "/media") {
				roots = append(roots, mount.MountPoint)
			}
		}
	}

	return roots
}

// resolveAndValidatePath converts a path to absolute and validates it's within allowed roots.
// For paths starting with /, it validates against all allowed roots.
// For relative paths, it resolves against the home directory.
// Returns the absolute path and the root it belongs to, or an error.
func resolveAndValidatePath(path string) (absPath string, root string, err error) {
	cleanPath := filepath.Clean(path)

	// Determine if path is absolute or relative
	if filepath.IsAbs(cleanPath) {
		absPath = cleanPath
	} else {
		// Relative path - resolve against home directory
		home, err := GetHomeDirectory()
		if err != nil {
			return "", "", os.ErrPermission
		}
		absPath = filepath.Join(home, cleanPath)
	}

	// Validate path is within an allowed root
	allowedRoots := GetAllowedRoots()
	for _, allowedRoot := range allowedRoots {
		relPath, err := filepath.Rel(allowedRoot, absPath)
		if err == nil && !strings.HasPrefix(relPath, "..") && relPath != ".." {
			return absPath, allowedRoot, nil
		}
	}

	return "", "", os.ErrPermission
}

// Renaming file function
func RenameFile(sourcePath string, newName string) error {
	// New file name can't be empty or contain path separators
	if newName == "" || strings.Contains(newName, "/") {
		return os.ErrInvalid
	}

	// Validate and resolve path
	absPath, _, err := resolveAndValidatePath(sourcePath)
	if err != nil {
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
	// Validate and resolve path
	absPath, _, err := resolveAndValidatePath(path)
	if err != nil {
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
	// Validate and resolve path
	absPath, root, err := resolveAndValidatePath(path)
	if err != nil {
		return os.ErrPermission
	}

	// Don't allow deleting the root itself
	if absPath == root {
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
	// Validate source path
	sourceLocation, _, err := resolveAndValidatePath(sourcePath)
	if err != nil {
		return os.ErrPermission
	}

	// Validate destination path
	destinationDirLocation, _, err := resolveAndValidatePath(destinationDir)
	if err != nil {
		return os.ErrPermission
	}

	destinationFile := filepath.Join(destinationDirLocation, filepath.Base(sourcePath))

	// Check source file existence and is it not a directory?
	if f, err := os.Stat(sourceLocation); err != nil {
		return os.ErrNotExist
	} else if f.IsDir() {
		return os.ErrInvalid
	}

	// Check the destination file doesn't already exist
	if _, err := os.Stat(destinationFile); err == nil {
		return os.ErrExist
	}

	// Check destination is a directory
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

// Copy folder function
func CopyFolder(sourcePath string, destinationDir string) error {
	// Validate source path
	sourceLocation, _, err := resolveAndValidatePath(sourcePath)
	if err != nil {
		return os.ErrPermission
	}

	// Validate destination path
	destinationDirLocation, _, err := resolveAndValidatePath(destinationDir)
	if err != nil {
		return os.ErrPermission
	}

	destinationFolder := filepath.Join(destinationDirLocation, filepath.Base(sourcePath))

	// Check source existence and is it a directory?
	if f, err := os.Stat(sourceLocation); err != nil {
		return os.ErrNotExist
	} else if !f.IsDir() {
		return os.ErrInvalid
	}

	// Check the destination folder doesn't already exist
	if _, err := os.Stat(destinationFolder); err == nil {
		return os.ErrExist
	}

	// Check destination is a directory
	info, err := os.Stat(destinationDirLocation)
	if err != nil {
		return os.ErrNotExist
	}
	if !info.IsDir() {
		return os.ErrInvalid
	}

	// Perform the copy operation
	folderContents, err := os.ReadDir(sourceLocation)
	if err != nil {
		return err
	}

	err = os.Mkdir(destinationFolder, 0755)
	if err != nil {
		return err
	}

	for _, entry := range folderContents {
		srcPath := filepath.Join(sourceLocation, entry.Name())

		if entry.IsDir() {
			// Recursive call with absolute paths
			err = CopyFolder(srcPath, destinationFolder)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, destinationFolder)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Cut file function
func CutFile(sourcePath string, destinationDir string) error {
	// Validate source path
	sourceLocation, _, err := resolveAndValidatePath(sourcePath)
	if err != nil {
		return os.ErrPermission
	}

	// Validate destination path
	destinationDirLocation, _, err := resolveAndValidatePath(destinationDir)
	if err != nil {
		return os.ErrPermission
	}

	destinationFile := filepath.Join(destinationDirLocation, filepath.Base(sourcePath))

	// Check source file existence and is it not a directory?
	if f, err := os.Stat(sourceLocation); err != nil {
		return os.ErrNotExist
	} else if f.IsDir() {
		return os.ErrInvalid
	}

	// Check the destination file doesn't already exist
	if _, err := os.Stat(destinationFile); err == nil {
		return os.ErrExist
	}

	// Check destination is a directory
	info, err := os.Stat(destinationDirLocation)
	if err != nil {
		return os.ErrNotExist
	}
	if !info.IsDir() {
		return os.ErrInvalid
	}

	// Try direct rename first (faster, works on same filesystem)
	if err := os.Rename(sourceLocation, destinationFile); err == nil {
		return nil
	}

	// Fallback: cross-filesystem move (copy + delete)
	if err := CopyFile(sourcePath, destinationDir); err != nil {
		return err
	}

	return os.Remove(sourceLocation)
}

// Cut folder function
func CutFolder(sourcePath string, destinationDir string) error {
	// Validate source path
	sourceLocation, _, err := resolveAndValidatePath(sourcePath)
	if err != nil {
		return os.ErrPermission
	}

	// Validate destination path
	destinationDirLocation, _, err := resolveAndValidatePath(destinationDir)
	if err != nil {
		return os.ErrPermission
	}

	destinationFolder := filepath.Join(destinationDirLocation, filepath.Base(sourcePath))

	// Check source existence and is it a directory?
	if f, err := os.Stat(sourceLocation); err != nil {
		return os.ErrNotExist
	} else if !f.IsDir() {
		return os.ErrInvalid
	}

	// Check the destination folder doesn't already exist
	if _, err := os.Stat(destinationFolder); err == nil {
		return os.ErrExist
	}

	// Check destination is a directory
	info, err := os.Stat(destinationDirLocation)
	if err != nil {
		return os.ErrNotExist
	}
	if !info.IsDir() {
		return os.ErrInvalid
	}

	// Try direct rename first (faster, works on same filesystem)
	if err := os.Rename(sourceLocation, destinationFolder); err == nil {
		return nil
	}

	// Fallback: cross-filesystem move (copy + delete)
	if err := CopyFolder(sourceLocation, destinationDir); err != nil {
		return err
	}

	return os.RemoveAll(sourceLocation)
}

// Generic Delete function (handles both file and folder)
func Delete(path string) error {
	// Validate and resolve path
	absPath, _, err := resolveAndValidatePath(path)
	if err != nil {
		return os.ErrPermission
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return os.ErrNotExist
	}

	if info.IsDir() {
		return DeleteFolder(path)
	}
	return DeleteFile(path)
}

// Generic Copy function (handles both file and folder)
func Copy(sourcePath, destinationDir string) error {
	// Validate source path
	sourceLocation, _, err := resolveAndValidatePath(sourcePath)
	if err != nil {
		return os.ErrPermission
	}

	info, err := os.Stat(sourceLocation)
	if err != nil {
		return os.ErrNotExist
	}

	if info.IsDir() {
		return CopyFolder(sourcePath, destinationDir)
	}
	return CopyFile(sourcePath, destinationDir)
}

// Generic Cut function (handles both file and folder)
func Cut(sourcePath, destinationDir string) error {
	// Validate source path
	sourceLocation, _, err := resolveAndValidatePath(sourcePath)
	if err != nil {
		return os.ErrPermission
	}

	info, err := os.Stat(sourceLocation)
	if err != nil {
		return os.ErrNotExist
	}

	if info.IsDir() {
		return CutFolder(sourcePath, destinationDir)
	}
	return CutFile(sourcePath, destinationDir)
}
