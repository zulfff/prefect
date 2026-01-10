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
	path := request.URL.Query().Get("path")

	filesData, err := parser.FileEntriesParser(path)
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

// DELETE
func DeleteAPI(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	if err := file.DeleteFile(path); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RENAME
func RenameAPI(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	newName := r.URL.Query().Get("name")

	if err := file.RenameFile(path, newName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// COPY
func CopyAPI(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	dst := r.URL.Query().Get("dst")

	if err := file.CopyFile(src, dst); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CUT
func CutAPI(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	dst := r.URL.Query().Get("dst")

	if err := file.CutFile(src, dst); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}