package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jurikolo/go-passgen/handler"
)

const version = "1.0.0"

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Print version and exit")
	vFlag := flag.Bool("v", false, "Print version and exit (shorthand)")
	flag.Parse()

	// Handle version flags
	if *versionFlag || *vFlag {
		fmt.Printf("go-passgen version %s\n", version)
		os.Exit(0)
	}

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
