package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// PassgenRequest matches the API request structure
type PassgenRequest struct {
	Length    int  `json:"length"`
	Uppercase bool `json:"uppercase"`
	Lowercase bool `json:"lowercase"`
	Digits    bool `json:"digits"`
	Symbols   bool `json:"symbols"`
	Count     int  `json:"count"`
}

// PassgenResponse matches the API response structure
type PassgenResponse struct {
	Passwords []string `json:"passwords"`
}

// getBaseURL returns the base URL for the passgen API
func getBaseURL() string {
	url := os.Getenv("PASSGEN_URL")
	if url == "" {
		return "http://localhost:8080"
	}
	return url
}

// handleGeneratePasswordTool handles the generate_password tool
func handleGeneratePasswordTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	// Get parameters with defaults
	length := request.GetInt("length", 0)
	if length == 0 {
		return nil, fmt.Errorf("length parameter is required")
	}

	uppercase := request.GetBool("uppercase", true)
	lowercase := request.GetBool("lowercase", true)
	digits := request.GetBool("digits", true)
	symbols := request.GetBool("symbols", true)
	count := request.GetInt("count", 1)

	// Validate at least one character set is enabled
	if !uppercase && !lowercase && !digits && !symbols {
		return nil, fmt.Errorf("at least one character set must be enabled")
	}

	// Prepare request
	reqBody := PassgenRequest{
		Length:    length,
		Uppercase: uppercase,
		Lowercase: lowercase,
		Digits:    digits,
		Symbols:   symbols,
		Count:     count,
	}

	// Call passgen API
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := getBaseURL() + "/api/generate"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call passgen API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("passgen API returned error: %s, body: %s", resp.Status, string(body))
	}

	var passgenResp PassgenResponse
	if err := json.NewDecoder(resp.Body).Decode(&passgenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Format result
	resultText := "Generated passwords:\n"
	for i, pwd := range passgenResp.Passwords {
		resultText += fmt.Sprintf("%d. %s\n", i+1, pwd)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

// NewPassgenMCPServer creates and configures the MCP server
func NewPassgenMCPServer() *server.MCPServer {
	mcpServer := server.NewMCPServer(
		"passgen-mcp",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	// Add the generate_password tool
	mcpServer.AddTool(mcp.NewTool("generate_password",
		mcp.WithDescription("Generate random passwords using the passgen API"),
		mcp.WithNumber("length",
			mcp.Description("Length of the password (required)"),
			mcp.Required(),
			mcp.Min(8),
			mcp.Max(128),
		),
		mcp.WithBoolean("uppercase",
			mcp.Description("Include uppercase letters"),
			mcp.DefaultBool(true),
		),
		mcp.WithBoolean("lowercase",
			mcp.Description("Include lowercase letters"),
			mcp.DefaultBool(true),
		),
		mcp.WithBoolean("digits",
			mcp.Description("Include digits"),
			mcp.DefaultBool(true),
		),
		mcp.WithBoolean("symbols",
			mcp.Description("Include symbols"),
			mcp.DefaultBool(true),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of passwords to generate"),
			mcp.DefaultNumber(1),
			mcp.Min(1),
			mcp.Max(100),
		),
	), handleGeneratePasswordTool)

	return mcpServer
}

func main() {
	// Parse command line flags
	httpFlag := flag.Bool("http", false, "Start HTTP/SSE server instead of stdio")
	portFlag := flag.String("port", "8080", "Port for HTTP server (default: 8080)")
	flag.Parse()

	mcpServer := NewPassgenMCPServer()

	if *httpFlag {
		// Start SSE server (compatible with Open WebUI and most MCP clients)
		log.Printf("Starting MCP SSE server on port %s", *portFlag)
		sseServer := server.NewSSEServer(mcpServer)
		addr := ":" + *portFlag
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("SSE server error: %v", err)
		}
	} else {
		// Start server with stdio transport (default)
		log.Println("Starting MCP server with stdio transport")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
