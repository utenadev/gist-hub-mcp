package github

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		verify func(*testing.T, *Client)
	}{
		{
			name:  "Create client with valid token",
			token: "test-token",
			verify: func(t *testing.T, c *Client) {
				if c == nil {
					t.Fatal("Expected non-nil client")
				}
				if c.token != "test-token" {
					t.Errorf("Expected token 'test-token', got '%s'", c.token)
				}
				if c.baseURL != "https://api.github.com" {
					t.Errorf("Expected base URL 'https://api.github.com', got '%s'", c.baseURL)
				}
				if c.userAgent != "gist-hub-mcp/0.1.0" {
					t.Errorf("Expected user agent 'gist-hub-mcp/0.1.0', got '%s'", c.userAgent)
				}
			},
		},
		{
			name:  "Create client with empty token",
			token: "",
			verify: func(t *testing.T, c *Client) {
				if c == nil {
					t.Fatal("Expected non-nil client")
				}
				if c.token != "" {
					t.Errorf("Expected empty token, got '%s'", c.token)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.token)
			tt.verify(t, client)
		})
	}
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name:     "APIError message",
			err:      &APIError{Message: "Bad request", Status: 400},
			expected: "GitHub API error: Bad request (status: 400)",
		},
		{
			name:     "APIError with empty message",
			err:      &APIError{Message: "", Status: 404},
			expected: "GitHub API error:  (status: 404)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Integration tests - require actual GitHub token
// These tests are skipped if no token is available

// func TestClientIntegration(t *testing.T) {
// 	token, err := GetToken("github.com")
// 	if err != nil {
// 		t.Skipf("No GitHub token available for integration tests: %v", err)
// 	}

// 	client := NewClient(token)

// 	t.Run("ListGists", func(t *testing.T) {
// 		gists, err := client.ListGists()
// 		if err != nil {
// 			t.Fatalf("ListGists() error = %v", err)
// 		}
// 		// We expect at least zero gists (even if user has none)
// 		if gists == nil {
// 			t.Error("Expected non-nil gists slice")
// 		}
// 	})

// 	t.Run("GetGist with invalid ID", func(t *testing.T) {
// 		_, err := client.GetGist("invalid-gist-id")
// 		if err == nil {
// 			t.Error("Expected error for invalid gist ID, got nil")
// 		}
// 	})

// 	t.Run("Create Get Update Delete Gist", func(t *testing.T) {
// 		// Create
// 		gist := &Gist{
// 			Description: "Test gist from gist-hub-mcp",
// 			Public:      false,
// 			Files: map[string]GistFile{
// 				"test.txt": {
// 					Content: "Test content",
// 				},
// 			},
// 		}

// 		created, err := client.CreateGist(gist)
// 		if err != nil {
// 			t.Fatalf("CreateGist() error = %v", err)
// 		}
// 		if created.ID == "" {
// 			t.Error("Expected non-empty gist ID")
// 		}

// 		// Get
// 		got, err := client.GetGist(created.ID)
// 		if err != nil {
// 			t.Fatalf("GetGist() error = %v", err)
// 		}
// 		if got.Description != gist.Description {
// 			t.Errorf("Expected description %q, got %q", gist.Description, got.Description)
// 		}

// 		// Update
// 		newDescription := "Updated test gist"
// 		created.Description = newDescription
// 		updated, err := client.UpdateGist(created.ID, created)
// 		if err != nil {
// 			t.Fatalf("UpdateGist() error = %v", err)
// 		}
// 		if updated.Description != newDescription {
// 			t.Errorf("Expected description %q, got %q", newDescription, updated.Description)
// 		}

// 		// Delete
// 		err = client.DeleteGist(created.ID)
// 		if err != nil {
// 			t.Fatalf("DeleteGist() error = %v", err)
// 		}

// 		// Verify deletion
// 		_, err = client.GetGist(created.ID)
// 		if err == nil {
// 			t.Error("Expected error when getting deleted gist, got nil")
// 		}
// 	})
// }
