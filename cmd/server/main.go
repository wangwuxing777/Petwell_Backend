package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vf0429/Petwell_Backend/internal/config"
	"github.com/vf0429/Petwell_Backend/internal/handlers"
	"github.com/vf0429/Petwell_Backend/internal/services/chat"
	"github.com/vf0429/Petwell_Backend/internal/services/places"
	"github.com/vf0429/Petwell_Backend/internal/services/rag"
)

const port = "8000"

func main() {
	// Load Configuration
	cfg := config.LoadConfig()
	ragClient := rag.NewClient(cfg)

	// Initialize chat session store (30-minute TTL)
	sessionStore := chat.NewSessionStore(30 * time.Minute)

	// Core handlers
	http.HandleFunc("/vaccines", handlers.VaccinesHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/posts", handlers.PostsHandler)
	http.HandleFunc("/clinics", handlers.NewClinicsHandler(cfg))
	http.HandleFunc("/emergency-clinics", handlers.NewEmergencyClinicsHandler(cfg))

	// Insurance handlers
	http.HandleFunc("/insurance-companies", handlers.InsuranceCompaniesHandler)
	http.HandleFunc("/insurance-products", handlers.InsuranceProductsHandler)
	http.HandleFunc("/coverage-list", handlers.CoverageListHandler)
	http.HandleFunc("/coverage-limits", handlers.CoverageLimitsHandler)
	http.HandleFunc("/sub-coverage-limits", handlers.SubCoverageLimitsHandler)

	// Legacy handlers
	http.HandleFunc("/insurance-providers", handlers.InsuranceProvidersHandler)
	http.HandleFunc("/service-subcategories", handlers.ServiceSubcategoriesHandler)

	// Legacy RAG chat handler (now with session support)
	http.HandleFunc("/api/chat", handlers.NewRAGHandler(ragClient, sessionStore))

	// New chat session endpoints
	http.HandleFunc("/api/chat/session", handlers.NewChatSessionHandler(sessionStore))
	http.HandleFunc("/api/chat/session/", handlers.NewChatSelectProviderHandler(sessionStore)) // matches /api/chat/session/{id}/provider
	http.HandleFunc("/api/chat/providers", handlers.NewChatProvidersHandler(ragClient))
	http.HandleFunc("/api/chat/ask", handlers.NewChatAskHandler(sessionStore, ragClient))

	// Vets handler
	placesClient := places.NewClient(cfg.MapsAPIKey)
	http.HandleFunc("/api/vets", handlers.NewVetsHandler(placesClient))

	fmt.Printf("PetWell Backend running at http://localhost:%s\n", port)
	fmt.Println("Chat endpoints:")
	fmt.Println("  POST /api/chat/session          - Create session")
	fmt.Println("  POST /api/chat/session/{id}/provider - Select provider")
	fmt.Println("  GET  /api/chat/providers         - List providers")
	fmt.Println("  POST /api/chat/ask               - Ask with context")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
