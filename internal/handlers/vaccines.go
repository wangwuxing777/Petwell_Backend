package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const jsonFile = "assets/vaccines.json"

func VaccinesHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ex, err := os.Executable()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	exPath := filepath.Dir(ex)
	file, err := os.Open(jsonFile)
	if err != nil {
		file, err = os.Open(filepath.Join(exPath, jsonFile))
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	}
	defer file.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, file)
}
