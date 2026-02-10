package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/vf0429/Petwell_Backend/internal/models"
)

var (
	Posts = []models.BlogPost{
		{ID: "1", AuthorName: "System", Title: "Welcome to PetWell Blog", Content: "This is the start of our community.", Likes: 10, ImageColor: "green", Timestamp: time.Now()},
	}
	postsMutex sync.RWMutex
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method == http.MethodGet {
		postsMutex.RLock()
		defer postsMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Posts)
		return
	}
	if r.Method == http.MethodPost {
		var post models.BlogPost
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		postsMutex.Lock()
		Posts = append(Posts, post)
		postsMutex.Unlock()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}
