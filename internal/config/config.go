package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RAGServiceURL string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	ragURL := os.Getenv("RAG_SERVICE_URL")
	if ragURL == "" {
		ragURL = "http://localhost:8001"
	}

	return &Config{
		RAGServiceURL: ragURL,
	}
}
