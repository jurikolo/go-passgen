package handler

import (
	_ "embed"
	"encoding/json"
	"io"
	"net/http"

	"github.com/jurikolo/go-passgen/generator"
)

// GenerateRequest represents the JSON request body for password generation
type GenerateRequest struct {
	Length    int  `json:"length"`
	Uppercase bool `json:"uppercase"`
	Lowercase bool `json:"lowercase"`
	Digits    bool `json:"digits"`
	Symbols   bool `json:"symbols"`
	Count     int  `json:"count"`
}

// GenerateResponse represents the JSON response with generated passwords
type GenerateResponse struct {
	Passwords []string `json:"passwords"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

//go:embed templates/index.html
var indexHTMLContent string

// GeneratePassword handles POST /api/generate
func GeneratePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req GenerateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Convert to generator config
	config := generator.Config{
		Length:    req.Length,
		Uppercase: req.Uppercase,
		Lowercase: req.Lowercase,
		Digits:    req.Digits,
		Symbols:   req.Symbols,
		Count:     req.Count,
	}

	// Generate passwords
	passwords, err := generator.GenerateMultiple(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send response
	resp := GenerateResponse{Passwords: passwords}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := HealthResponse{Status: "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// ServeIndex serves the HTML UI at GET /
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve the HTML page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTMLContent))
}
