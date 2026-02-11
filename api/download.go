package api

import (
	"net/http"
	"os"
	"prefect/services/file"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Extract the source path from the request
	src := r.URL.Query().Get("path")
	if src == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}

	// Validate and resolve path
	absPath, _, err := file.ResolveAndValidatePath(src)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	if d, err := os.Stat(absPath); err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	} else if d.IsDir() {
		http.Error(w, "Path is a directory", http.StatusBadRequest)
		return
	}

	info, _ := os.Stat(absPath)

	// Force download instead of inline display
	w.Header().Set("Content-Disposition", `attachment; filename="`+info.Name()+`"`)
	http.ServeFile(w, r, absPath)
}
