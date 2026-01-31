package parser

import (
	"encoding/json"
	"log"
	"os"
	"prefect/services/file"
	"prefect/services/sys"
	"strings"
	"time"
)

type SysData struct {
	// CPU
	CPUCores   int `json:"cpu_cores"`
	CPUThreads int `json:"cpu_threads"`
	CPUUsage   int `json:"cpu_usage"`
	CPUTemp    int `json:"cpu_temp"`
	CPUPower   int `json:"cpu_power"`

	// RAM
	RAMTotal int `json:"ram_total"`
	RAMUsed  int `json:"ram_used"`
	RAMUsage int `json:"ram_usage"`

	// Disk
	DiskTotal uint64 `json:"disk_total"`
	DiskUsed  uint64 `json:"disk_used"`
	DiskUsage int    `json:"disk_usage"`
}

type DrivesData struct {
	Device     string `json:"device"`
	MountPoint string `json:"mount_point"`
	FSType     string `json:"fs_type"`
}

type SidebarData struct {
	DirectoryName string `json:"directory_name"`
	DirectoryPath string `json:"directory_path"`
}

func SysDataParser() SysData {
	RAMTotal, RAMUsed, RAMUsage := sys.RAM()
	DiskTotal, DiskUsed, DiskUsage := sys.Disk()

	// Data Structures
	return SysData{
		CPUCores:   sys.CPUCores(),
		CPUThreads: sys.CPUThreads(),
		CPUUsage:   sys.CPUUsage(),
		CPUTemp:    sys.CPUTemp(),
		CPUPower:   sys.CPUPower(1 * time.Second),
		RAMTotal:   RAMTotal,
		RAMUsed:    RAMUsed,
		RAMUsage:   RAMUsage,
		DiskUsed:   DiskUsed,
		DiskTotal:  DiskTotal,
		DiskUsage:  DiskUsage,
	}
}

func DrivesDataParser() {
	allMounts, err := file.GetMountedDrives()

	if err != nil {
		log.Println("Error fetching mounted drives:", err)
		return
	}

	var filteredMounts []DrivesData

	for _, mount := range allMounts {
		// Filter these filesystems
		if mount.FSType == "tmpfs" || mount.FSType == "devtmpfs" || mount.FSType == "proc" || mount.FSType == "sysfs" || mount.FSType == "cgroup" || mount.FSType == "overlay" || mount.FSType == "rootfs" || mount.FSType == "cgroup2" || mount.FSType == "debugfs" || mount.FSType == "tracefs" || mount.FSType == "configfs" || mount.FSType == "binfmt_misc" || mount.FSType == "fusectl" || mount.FSType == "hugetlbfs" || mount.FSType == "mqueue" || mount.FSType == "pstore" || mount.FSType == "securityfs" || mount.FSType == "efivarfs" || mount.FSType == "bpf" {
			continue
		}

		// Exclude system mount points
		if strings.HasPrefix(mount.MountPoint, "/proc") || strings.HasPrefix(mount.MountPoint, "/sys") || strings.HasPrefix(mount.MountPoint, "/sys/fs") || strings.HasPrefix(mount.MountPoint, "/sys/kernel") || strings.HasPrefix(mount.MountPoint, "/dev") || strings.HasPrefix(mount.MountPoint, "/dev/pts") || strings.HasPrefix(mount.MountPoint, "/dev/mqueue") || strings.HasPrefix(mount.MountPoint, "/dev/hugepages") || strings.HasPrefix(mount.MountPoint, "/run") || strings.HasPrefix(mount.MountPoint, "/run/user") || strings.HasPrefix(mount.MountPoint, "/usr") || strings.HasPrefix(mount.MountPoint, "/usr/lib") || strings.HasPrefix(mount.MountPoint, "/usr/lib/wsl") || strings.HasPrefix(mount.MountPoint, "/usr/lib/modules") || strings.HasPrefix(mount.MountPoint, "/var/run") || strings.HasPrefix(mount.MountPoint, "/var/lib/docker") || strings.HasPrefix(mount.MountPoint, "/var/lib/containerd") || strings.HasPrefix(mount.MountPoint, "/var/lib/kubelet") || strings.HasPrefix(mount.MountPoint, "/init") {
			continue
		}

		// Only include root, internal and external drives
		if mount.MountPoint == "/" || strings.HasPrefix(mount.MountPoint, "/mnt/") || strings.HasPrefix(mount.MountPoint, "/media") {
			filteredMounts = append(filteredMounts, DrivesData{
				Device:     mount.Device,
				MountPoint: mount.MountPoint,
				FSType:     mount.FSType,
			})
		}
	}

	mounts := filteredMounts

	// Write to drives.json
	jsoner, jerr := os.OpenFile("drives.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

	if jerr != nil {
		log.Println("Error writing JSON:", jerr)
		return
	}
	defer jsoner.Close()

	encoder := json.NewEncoder(jsoner)

	if err := encoder.Encode(mounts); err != nil {
		log.Println("Error encoding mounted drives:", err)
	}
}

func SidebarDataParser() {
	homeDir, err := file.GetHomeDirectory()
	if err != nil {
		log.Println("Error fetching home directory:", err)
		return
	}

	var sidebar []SidebarData

	// Add Home directory to sidebar
	sidebar = append(sidebar, SidebarData{
		DirectoryName: "Home",
		DirectoryPath: homeDir,
	})

	// Add contents inside $HOME to sidebar
	entries, err := os.ReadDir(homeDir)
	if err != nil {
		log.Println("Error reading home directory contents:", err)
		return
	}

	// Initialize default directories in case they don't exist
	file.DefaultDirectories()

	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "Downloads" || entry.Name() == "Documents" || entry.Name() == "Media" {
				sidebar = append(sidebar, SidebarData{
					DirectoryName: entry.Name(),
					DirectoryPath: homeDir + "/" + entry.Name(),
				})
			}
		}
	}

	jsoner, jerr := os.OpenFile("sidebar.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

	if jerr != nil {
		log.Println("Error writing JSON:", jerr)
		return
	}
	defer jsoner.Close()

	encoder := json.NewEncoder(jsoner)

	if err := encoder.Encode(sidebar); err != nil {
		log.Println("Error encoding sidebar data:", err)
	}
}