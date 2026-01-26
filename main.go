package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	port     = "8000"
	jsonFile = "vaccines.json"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func vaccinesHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for testing if needed, though SimpleHTTPRequestHandler usually doesn't by default.
	// We'll stick to the Python behavior: just serve the file.
	
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the directory of the executable to find the json file reliably
	// equivalent to: script_dir = os.path.dirname(os.path.realpath(__file__))
	ex, err := os.Executable()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	exPath := filepath.Dir(ex)
	
	// However, usually when running 'go run .', the CWD is the project root.
	// Let's try to open local file first.
	file, err := os.Open(jsonFile)
	if err != nil {
		// Fallback to executable path if strictly needed, but for 'go run' local is better.
		file, err = os.Open(filepath.Join(exPath, jsonFile))
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")
	// Copy file content to response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/vaccines", vaccinesHandler)

	fmt.Printf("Serving at http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
