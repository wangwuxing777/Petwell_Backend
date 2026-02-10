package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vf0429/Petwell_Backend/internal/services/rag"
)

func NewRAGHandler(client *rag.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		EnableCors(&w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req rag.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if client == nil {
			http.Error(w, "RAG client not initialized", http.StatusInternalServerError)
			return
		}

		resp, err := client.Ask(req.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
