package file

import (
	"os"
	"log"
	"path/filepath"
)

func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

func DefaultDirectories() {
	var defaultDirs = []string{"Downloads", "Documents", "Media"}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not find home directory:", err)
	}

	for _, dir := range defaultDirs {
		// Make folder inside home directory
		dirsPath := filepath.Join(home, dir)
		
		err := os.MkdirAll(dirsPath, 0755)
		if err != nil {
			log.Printf("Failed to create %s: %v", dirsPath, err)
		}
	}
}
