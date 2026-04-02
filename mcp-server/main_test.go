package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// mockPassgenServer creates a test server that simulates the passgen API
func mockPassgenServer(t *testing.T, statusCode int, response interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/generate" {
			t.Errorf("expected path /api/generate, got %s", r.URL.Path)
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", contentType)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		var req PassgenRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Errorf("failed to unmarshal request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if response != nil {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf("failed to encode response: %v", err)
			}
		}
	}))
}

func TestHandleGeneratePasswordTool(t *testing.T) {
	tests := []struct {
		name           string
		arguments      map[string]interface{}
		mockStatusCode int
		mockResponse   *PassgenResponse
		wantErr        bool
		expectErrorMsg string
	}{
		{
			name: "successful password generation",
			arguments: map[string]interface{}{
				"length":    float64(12),
				"uppercase": true,
				"lowercase": true,
				"digits":    true,
				"symbols":   true,
				"count":     float64(2),
			},
			mockStatusCode: http.StatusOK,
			mockResponse: &PassgenResponse{
				Passwords: []string{"Abc123!@", "Xyz456#&"},
			},
			wantErr: false,
		},
		{
			name: "missing length parameter",
			arguments: map[string]interface{}{
				"uppercase": true,
			},
			mockStatusCode: http.StatusOK,
			mockResponse:   nil,
			wantErr:        true,
			expectErrorMsg: "length parameter is required",
		},
		{
			name: "no character sets enabled",
			arguments: map[string]interface{}{
				"length":    float64(10),
				"uppercase": false,
				"lowercase": false,
				"digits":    false,
				"symbols":   false,
			},
			mockStatusCode: http.StatusOK,
			mockResponse:   nil,
			wantErr:        true,
			expectErrorMsg: "at least one character set must be enabled",
		},
		{
			name: "passgen API returns error",
			arguments: map[string]interface{}{
				"length":    float64(8),
				"uppercase": true,
				"lowercase": true,
				"digits":    true,
				"symbols":   true,
			},
			mockStatusCode: http.StatusInternalServerError,
			mockResponse:   nil,
			wantErr:        true,
			expectErrorMsg: "passgen API returned error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockPassgenServer(t, tt.mockStatusCode, tt.mockResponse)
			defer server.Close()

			// Set environment variable to point to mock server
			originalEnv := os.Getenv("PASSGEN_URL")
			os.Setenv("PASSGEN_URL", server.URL)
			defer os.Setenv("PASSGEN_URL", originalEnv)

			ctx := context.Background()
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.arguments,
				},
			}
			result, err := handleGeneratePasswordTool(ctx, request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.expectErrorMsg != "" && !contains(err.Error(), tt.expectErrorMsg) {
					t.Errorf("error message mismatch: got %q, want containing %q", err.Error(), tt.expectErrorMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result == nil {
				t.Error("expected result but got nil")
				return
			}
			// Check that result contains the expected passwords
			if len(result.Content) != 1 {
				t.Errorf("expected 1 content item, got %d", len(result.Content))
			}
			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Errorf("expected TextContent, got %T", result.Content[0])
			}
			expectedText := "Generated passwords:\n1. Abc123!@\n2. Xyz456#&\n"
			if textContent.Text != expectedText {
				t.Errorf("result text mismatch:\ngot:\n%q\nwant:\n%q", textContent.Text, expectedText)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestNewPassgenMCPServer(t *testing.T) {
	server := NewPassgenMCPServer()
	if server == nil {
		t.Error("expected server not to be nil")
	}
}

func TestGetBaseURL(t *testing.T) {
	originalEnv := os.Getenv("PASSGEN_URL")
	defer os.Setenv("PASSGEN_URL", originalEnv)

	// Test default
	os.Setenv("PASSGEN_URL", "")
	if got := getBaseURL(); got != "http://localhost:8080" {
		t.Errorf("default base URL mismatch: got %q, want %q", got, "http://localhost:8080")
	}

	// Test custom
	os.Setenv("PASSGEN_URL", "http://example.com:3000")
	if got := getBaseURL(); got != "http://example.com:3000" {
		t.Errorf("custom base URL mismatch: got %q, want %q", got, "http://example.com:3000")
	}
}
