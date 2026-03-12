package api

import (
	"archive/zip"
	"fmt"
	"net/http"
	"os"
	"prefect/services/file"
	"path/filepath"
	"io"
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

	info, err := os.Stat(absPath)

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if info.IsDir() {
		http.Error(w, "Path is a directory", http.StatusBadRequest)
		return
	}

	// Force download instead of inline display
	w.Header().Set("Content-Disposition", `attachment; filename="`+info.Name()+`"`)
	http.ServeFile(w, r, absPath)
}

func DownloadFolder(w http.ResponseWriter, r *http.Request) {
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

	info, err := os.Stat(absPath)

	if err != nil {
		http.Error(w, "Folder not found", http.StatusNotFound)
		return
	}

	if !info.IsDir() {
		http.Error(w, "Path is not a directory", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", info.Name()))

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			http.Error(w, "Failed to create archive", http.StatusInternalServerError)
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(absPath, path)
		if err != nil {
			return err
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		io.Copy(zipFile, file)
		file.Close()

		return nil
	})
}
