package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jurikolo/go-passgen/handler"
)

func main() {
	port := getPort()

	// Setup handlers
	http.HandleFunc("/", handler.ServeIndex)
	http.HandleFunc("/api/generate", handler.GeneratePassword)
	http.HandleFunc("/health", handler.HealthCheck)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Validate port is numeric
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Invalid PORT environment variable: %s", port)
	}
	return port
}
