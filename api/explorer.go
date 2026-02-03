package api

import (
	"net/http"
	"prefect/services/parser"
	"prefect/services/file"
	"encoding/json"
)

type FileExplorer struct {
	Files   []file.FilesList   `json:"files"`
	Sidebar []parser.SidebarData `json:"sidebar"`
	Drives  []parser.DrivesData  `json:"drives"`
}

func FileExplorerAPI(response http.ResponseWriter, request *http.Request) {
	filesData, err := parser.FileEntriesParser()
	if err != nil {
		http.Error(response, "Failed to get file entries", http.StatusInternalServerError)
		return
	}

	sidebarData, err := parser.SidebarDataParser()
	if err != nil {
		http.Error(response, "Failed to get sidebar data", http.StatusInternalServerError)
		return
	}

	drivesData, err := parser.DrivesDataParser()
	if err != nil {
		http.Error(response, "Failed to get drives data", http.StatusInternalServerError)
		return
	}

	fileExplorerResponse := FileExplorer{
		Files:   filesData,
		Sidebar: sidebarData,
		Drives:  drivesData,
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(fileExplorerResponse)
}
