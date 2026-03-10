package db

import "time"

// Gist represents a GitHub Gist in the local cache
type Gist struct {
	ID          string    `db:"id"`
	Path        string    `db:"path"` // Wiki path with forward slashes
	Description string    `db:"description"`
	Public      bool      `db:"public"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// File represents a file within a Gist
type File struct {
	ID           int    `db:"id"`
	GistID       string `db:"gist_id"`
	OriginalPath string `db:"original_path"` // Original path with forward slashes
	GistFilename string `db:"gist_filename"` // Gist API filename with backslashes
	Size         int    `db:"size"`
	Language     string `db:"language"`
}
