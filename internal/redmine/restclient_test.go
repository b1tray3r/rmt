package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/b1tray3r/rmt/internal/redmine/models"
)

// TestNewRestClient tests the creation of a new RestClient instance.
func TestNewRestClient(t *testing.T) {
	client := NewRestClient("https://example.com/", "test-api-key")

	if client.baseURL != "https://example.com" {
		t.Errorf("expected baseURL 'https://example.com', got '%s'", client.baseURL)
	}
	if client.apiKey != "test-api-key" {
		t.Errorf("expected apiKey 'test-api-key', got '%s'", client.apiKey)
	}
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestRestClient_GetBaseURL(t *testing.T) {
	client := NewRestClient("https://example.com/", "test-api-key")
	if client.GetBaseURL() != "https://example.com" {
		t.Errorf("expected baseURL 'https://example.com', got '%s'", client.GetBaseURL())
	}
}

// TestRestClient_Login tests the login functionality with mock server.
func TestRestClient_Login(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError bool
		errorContains string
	}{
		{
			name:          "successful login",
			statusCode:    http.StatusOK,
			responseBody:  `{"user": {"id": 1, "login": "admin"}}`,
			expectedError: false,
		},
		{
			name:          "unauthorized",
			statusCode:    http.StatusUnauthorized,
			responseBody:  `{"errors": ["Invalid credentials"]}`,
			expectedError: true,
			errorContains: "authentication failed",
		},
		{
			name:          "server error",
			statusCode:    http.StatusInternalServerError,
			responseBody:  `{"errors": ["Internal server error"]}`,
			expectedError: true,
			errorContains: "login failed with status: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request headers
				if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
					t.Errorf("expected API key 'test-api-key', got '%s'", r.Header.Get("X-Redmine-API-Key"))
				}
				if r.URL.Path != "/users/current.json" {
					t.Errorf("expected path '/users/current.json', got '%s'", r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewRestClient(server.URL, "test-api-key")
			err := client.Login(context.Background())

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

// TestRestClient_Search tests the search functionality with mock server.
func TestRestClient_Search(t *testing.T) {
	mockSearchResults := models.SearchResults{
		Results: []models.SearchResult{
			{
				ID:          1,
				Title:       "Test Issue",
				Type:        "issue",
				URL:         "/issues/1",
				Description: "Test description",
				DateTime:    time.Date(2025, 8, 22, 12, 0, 0, 0, time.UTC),
			},
		},
		TotalCount: 1,
		Offset:     0,
		Limit:      25,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("expected API key 'test-api-key', got '%s'", r.Header.Get("X-Redmine-API-Key"))
		}

		// Verify query parameters
		query := r.URL.Query()
		if query.Get("q") != "test query" {
			t.Errorf("expected query 'test query', got '%s'", query.Get("q"))
		}
		if query.Get("limit") != "10" {
			t.Errorf("expected limit '10', got '%s'", query.Get("limit"))
		}
		if query.Get("scope") != "issues" {
			t.Errorf("expected scope 'issues', got '%s'", query.Get("scope"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockSearchResults)
	}))
	defer server.Close()

	client := NewRestClient(server.URL, "test-api-key")

	searchParams := models.SearchParams{
		Query:  "test query",
		Limit:  10,
		Scope:  "issues",
		Issues: true,
	}

	results, err := client.Search(searchParams)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if results.TotalCount != 1 {
		t.Errorf("expected total count 1, got %d", results.TotalCount)
	}
	if len(results.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results.Results))
	}
	if results.Results[0].Title != "Test Issue" {
		t.Errorf("expected title 'Test Issue', got '%s'", results.Results[0].Title)
	}
}

// TestRestClient_SearchError tests search error handling.
func TestRestClient_SearchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors": ["Bad request"]}`))
	}))
	defer server.Close()

	client := NewRestClient(server.URL, "test-api-key")

	searchParams := models.SearchParams{
		Query: "test",
	}

	_, err := client.Search(searchParams)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !contains(err.Error(), "search failed with status: 400") {
		t.Errorf("expected error to contain 'search failed with status: 400', got '%s'", err.Error())
	}
}

// contains is a helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && s[:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		containsSubstring(s, substr))
}

// containsSubstring checks if s contains substr as a substring.
func containsSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
