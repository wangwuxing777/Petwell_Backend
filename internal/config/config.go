package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RAGServiceURL string
	MapsAPIKey    string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	ragURL := os.Getenv("RAG_SERVICE_URL")
	if ragURL == "" {
		ragURL = "http://localhost:8001"
	}

	mapsKey := os.Getenv("MAPS_API_KEY")

	return &Config{
		RAGServiceURL: ragURL,
		MapsAPIKey:    mapsKey,
	}
}
