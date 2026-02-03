package file

import (
	"path/filepath"
	"log"
	"os"
)

type FilesList struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Path       string `json:"path"`
	IsDir      bool   `json:"is_dir"`
	Extensions string `json:"extensions"`
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

func FileEntries() ([]FilesList, error) {
	prototypeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	currentPath := prototypeDirectory

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
			IsDir:      entry.Type().IsDir(),
			Extensions: determineExtension(entry),
		})
	}

	return filesList, err
}
