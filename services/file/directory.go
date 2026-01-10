package file

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func GetHomeDirectory() (string, error) {
	// If not running under sudo, normal home directory
	if os.Geteuid() != 0 {
		return os.UserHomeDir()
	}

	// Running as root → check sudo
	sudoUID := os.Getenv("SUDO_UID")
	if sudoUID == "" {
		// root without sudo
		return "/root", nil
	}

	uid, err := strconv.Atoi(sudoUID)
	if err != nil {
		return "", err
	}

	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return "", err
	}

	return u.HomeDir, nil
}

func DefaultDirectories() {
	var defaultDirs = []string{"Downloads", "Documents", "Media"}

	home, err := GetHomeDirectory()
	if err != nil {
		log.Printf("Could not find home directory: %v", err)
		return
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
