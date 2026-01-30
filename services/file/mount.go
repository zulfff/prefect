package file

import (
	"bufio"
	"os"
	"strings"
)

type Mount struct {
	Device     string `json:"device"`
	MountPoint string `json:"mount_point"`
	FSType     string `json:"fs_type"`
}

func GetMountedDrives() ([]Mount, error) {
	// Directory where all the mounts data are located
	file, err := os.Open("/proc/mounts")

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mounts []Mount
	scanner := bufio.NewScanner(file)

	// Read each line from /proc/mounts
	// Then slice it into parts
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) >= 3 {
			mount := Mount{
				Device:     parts[0],
				MountPoint: parts[1],
				FSType:     parts[2],
			}
			mounts = append(mounts, mount)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mounts, nil
}
