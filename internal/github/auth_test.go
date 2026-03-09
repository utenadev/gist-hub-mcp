package github

import (
	"testing"
)

func TestGetToken(t *testing.T) {
	// This test requires actual gh CLI configuration.
	// It's an integration-style test that verifies the function works with real credentials.

	t.Run("GetToken with GitHub.com", func(t *testing.T) {
		token, err := GetToken("github.com")
		if err != nil {
			t.Skipf("No GitHub token configured: %v", err)
		}
		if token == "" {
			t.Error("Expected non-empty token, got empty string")
		}
	})

	t.Run("GetToken with empty host (should default to github.com)", func(t *testing.T) {
		token, err := GetToken("")
		if err != nil {
			t.Skipf("No GitHub token configured: %v", err)
		}
		if token == "" {
			t.Error("Expected non-empty token for empty host, got empty string")
		}
	})

	t.Run("GetToken with unknown host", func(t *testing.T) {
		token, err := GetToken("unknown.example.com")
		if err == nil {
			t.Error("Expected error for unknown host, got nil")
		}
		if token != "" {
			t.Error("Expected empty token for unknown host")
		}
	})
}
