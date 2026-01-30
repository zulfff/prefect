package parser

import (
	"prefect/services/file"
	"prefect/services/sys"
	// "prefect/services/file"
	"encoding/json"
	"log"
	"os"
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
	mounts, err := file.GetMountedDrives()

	if err != nil {
		log.Println("Error fetching mounted drives:", err)
		return
	}

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