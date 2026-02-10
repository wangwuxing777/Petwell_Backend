package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vf0429/Petwell_Backend/internal/config"
)

type Client struct {
	Config *config.Config
}

type Request struct {
	Query string `json:"query"`
}

type Response struct {
	Answer  string   `json:"answer"`
	Sources []string `json:"sources"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		Config: cfg,
	}
}

func (c *Client) Ask(query string) (*Response, error) {
	reqBody := Request{Query: query}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(c.Config.RAGServiceURL+"/ask", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call RAG service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("RAG service returned status %d: %s", resp.StatusCode, string(body))
	}

	var ragResp Response
	if err := json.NewDecoder(resp.Body).Decode(&ragResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ragResp, nil
}
