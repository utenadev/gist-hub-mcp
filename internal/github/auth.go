package github

import (
	"fmt"

	ghauth "github.com/cli/go-gh/v2/pkg/auth"
)

// GetToken retrieves the GitHub access token from gh CLI configuration.
// It returns the token for the specified host (default: github.com).
func GetToken(host string) (string, error) {
	if host == "" {
		host = "github.com"
	}
	token, _ := ghauth.TokenForHost(host)
	if token == "" {
		return "", fmt.Errorf("no GitHub token found for host: %s", host)
	}
	return token, nil
}
