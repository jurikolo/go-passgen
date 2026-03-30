package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HealthCheck returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	body := rr.Body.String()
	// json.Encoder adds a newline, trim it
	body = strings.TrimSpace(body)
	if body != expected {
		t.Errorf("HealthCheck returned unexpected body: got %v want %v", body, expected)
	}

	// Test wrong method
	req = httptest.NewRequest("POST", "/health", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("HealthCheck with POST should return 405, got %v", rr.Code)
	}
}

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "valid request",
			method: "POST",
			body: GenerateRequest{
				Length:    16,
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
				Symbols:   true,
				Count:     3,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var resp GenerateResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if len(resp.Passwords) != 3 {
					t.Errorf("Expected 3 passwords, got %d", len(resp.Passwords))
				}
				for _, pwd := range resp.Passwords {
					if len(pwd) != 16 {
						t.Errorf("Expected password length 16, got %d", len(pwd))
					}
				}
			},
		},
		{
			name:           "invalid JSON",
			method:         "POST",
			body:           "{invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "length too small",
			method: "POST",
			body: GenerateRequest{
				Length:    5,
				Uppercase: true,
				Count:     1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "count too large",
			method: "POST",
			body: GenerateRequest{
				Length:    12,
				Uppercase: true,
				Count:     200,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "wrong method GET",
			method:         "GET",
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "no character sets",
			method: "POST",
			body: GenerateRequest{
				Length: 12,
				Count:  1,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			switch b := tt.body.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				var err error
				bodyBytes, err = json.Marshal(b)
				if err != nil {
					t.Fatalf("Failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/api/generate", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GeneratePassword)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("GeneratePassword returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}
		})
	}
}

func TestServeIndex(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ServeIndex returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("ServeIndex returned wrong content type: got %v want text/html; charset=utf-8", contentType)
	}

	if len(rr.Body.String()) == 0 {
		t.Error("ServeIndex returned empty body")
	}

	// Should contain some expected HTML
	body := rr.Body.String()
	if !bytes.Contains([]byte(body), []byte("PassGen")) {
		t.Error("ServeIndex body does not contain 'PassGen'")
	}

	// Test wrong method
	req = httptest.NewRequest("POST", "/", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("ServeIndex with POST should return 405, got %v", rr.Code)
	}
}
