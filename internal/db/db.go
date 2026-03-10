package db

import "context"

// Repository defines the interface for storing and retrieving Gist data
type Repository interface {
	// Initialize the database schema
	Init(ctx context.Context) error

	// SaveGist stores a new gist or updates an existing one
	SaveGist(ctx context.Context, gist *Gist, files []*File) error

	// GetGist retrieves a gist by ID
	GetGist(ctx context.Context, id string) (*Gist, error)

	// GetGistByPath retrieves a gist by wiki path
	GetGistByPath(ctx context.Context, path string) (*Gist, error)

	// ListGists returns all gists
	ListGists(ctx context.Context) ([]*Gist, error)

	// GetGistFiles returns all files for a given gist
	GetGistFiles(ctx context.Context, gistID string) ([]*File, error)

	// DeleteGist removes a gist and all its files
	DeleteGist(ctx context.Context, gistID string) error

	// Close closes the database connection
	Close() error
}
