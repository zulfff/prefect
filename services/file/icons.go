package file

import (
	"strings"
)

// GetIconName returns the icon name for a given file name and directory status
func GetIconName(fileName string, isDir bool) string {
	if fileName == "" {
		if isDir {
			return "folder"
		}
		return "file"
	}

	name := strings.ToLower(fileName)

	if isDir {
		// Check original name (case-insensitive)
		if icon, ok := folderIcons[name]; ok {
			return icon
		}
		// Check without leading dot if applicable
		if strings.HasPrefix(name, ".") {
			withoutDot := name[1:]
			if icon, ok := folderIcons[withoutDot]; ok {
				return icon
			}
		}
		return "folder"
	}

	// Check exact filename
	if icon, ok := fileIconsByFilename[name]; ok {
		return icon
	}

	// Check extensions (including multi-part like .js.map)
	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		for i := 1; i < len(parts); i++ {
			ext := strings.Join(parts[i:], ".")
			if icon, ok := fileIconsByExtension[ext]; ok {
				return icon
			}
		}
	}

	return "file"
}
