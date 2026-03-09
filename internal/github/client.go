package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Gist represents a GitHub Gist.
type Gist struct {
	ID          string              `json:"id"`
	Description string              `json:"description,omitempty"`
	Public      bool                `json:"public,omitempty"`
	Files       map[string]GistFile `json:"files,omitempty"`
	HTMLURL     string              `json:"html_url,omitempty"`
	UpdatedAt   time.Time           `json:"updated_at,omitempty"`
	CreatedAt   time.Time           `json:"created_at,omitempty"`
}

// GistFile represents a single file within a Gist.
type GistFile struct {
	Filename string `json:"filename,omitempty"`
	Type     string `json:"type,omitempty"`
	Language string `json:"language,omitempty"`
	RawURL   string `json:"raw_url,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Content  string `json:"content,omitempty"`
}

// APIError represents an error response from GitHub API.
type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("GitHub API error: %s (status: %d)", e.Message, e.Status)
}

// Client provides methods to interact with GitHub Gist API.
type Client struct {
	client    *http.Client
	token     string
	baseURL   string
	userAgent string
}

// NewClient creates a new GitHub Gist API client.
func NewClient(token string) *Client {
	return &Client{
		client:    &http.Client{Timeout: 30 * time.Second},
		token:     token,
		baseURL:   "https://api.github.com",
		userAgent: "gist-hub-mcp/0.1.0",
	}
}

// ListGists retrieves all gists for the authenticated user.
func (c *Client) ListGists() ([]Gist, error) {
	url := fmt.Sprintf("%s/gists", c.baseURL)
	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var gists []Gist
	if err := c.do(req, &gists); err != nil {
		return nil, err
	}

	return gists, nil
}

// GetGist retrieves a single gist by ID.
func (c *Client) GetGist(gistID string) (*Gist, error) {
	url := fmt.Sprintf("%s/gists/%s", c.baseURL, gistID)
	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var gist Gist
	if err := c.do(req, &gist); err != nil {
		return nil, err
	}

	return &gist, nil
}

// CreateGist creates a new gist.
func (c *Client) CreateGist(gist *Gist) (*Gist, error) {
	url := fmt.Sprintf("%s/gists", c.baseURL)
	body, err := json.Marshal(gist)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var createdGist Gist
	if err := c.do(req, &createdGist); err != nil {
		return nil, err
	}

	return &createdGist, nil
}

// UpdateGist updates an existing gist.
func (c *Client) UpdateGist(gistID string, gist *Gist) (*Gist, error) {
	url := fmt.Sprintf("%s/gists/%s", c.baseURL, gistID)
	updateReq := struct {
		Description *string             `json:"description,omitempty"`
		Files       map[string]GistFile `json:"files,omitempty"`
	}{
		Description: &gist.Description,
		Files:       gist.Files,
	}
	body, err := json.Marshal(updateReq)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var updatedGist Gist
	if err := c.do(req, &updatedGist); err != nil {
		return nil, err
	}

	return &updatedGist, nil
}

// DeleteGist deletes a gist by ID.
func (c *Client) DeleteGist(gistID string) error {
	url := fmt.Sprintf("%s/gists/%s", c.baseURL, gistID)
	req, err := c.newRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// newRequest creates an HTTP request with authentication headers.
func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do sends the request and decodes the response.
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
			apiErr.Status = resp.StatusCode
			return &apiErr
		}
		return &APIError{
			Message: string(body),
			Status:  resp.StatusCode,
		}
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
