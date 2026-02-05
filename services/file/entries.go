package file

import (
	"os"
	"path/filepath"
)

type FilesList struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Path       string `json:"path"`
	IsDir      bool   `json:"is_dir"`
	Extensions string `json:"extensions"`
	Icon       string `json:"icon"`
}

func getExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	if len(ext) > 0 {
		return ext[1:] // Remove the dot
	}
	return ""
}

func determineExtension(entry os.DirEntry) string {
	if entry.IsDir() {
		return "directory"
	}
	return getExtension(entry.Name())
}

func FileEntries(path string) ([]FilesList, error) {
	var currentPath string

	// Default to HOME if path is empty
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		currentPath = home
	} else {
		currentPath = path
	}

	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return nil, err
	}

	var filesList []FilesList
	for _, entry := range entries {
		// entry.Info() gives the file size
		info, err := entry.Info()
		if err != nil {
			continue
		}

		filesList = append(filesList, FilesList{
			Name:       entry.Name(),
			Size:       info.Size(),
			Path:       filepath.Join(currentPath, entry.Name()),
			IsDir:      entry.IsDir(),
			Extensions: determineExtension(entry),
			Icon:       GetIconName(entry.Name(), entry.IsDir()),
		})
	}

	return filesList, nil
}
