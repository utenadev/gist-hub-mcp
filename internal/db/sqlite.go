package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteRepository implements the Repository interface using SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dsn string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	repo := &SQLiteRepository{db: db}
	return repo, nil
}

// Init initializes the database schema
func (r *SQLiteRepository) Init(ctx context.Context) error {
	// Enable foreign key constraints
	_, err := r.db.ExecContext(ctx, "PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, GetSchema())
	return err
}


// SaveGist stores a new gist or updates an existing one
func (r *SQLiteRepository) SaveGist(ctx context.Context, gist *Gist, files []*File) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Upsert gist
	public := 0
	if gist.Public {
		public = 1
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO gists (id, path, description, public, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			path = excluded.path,
			description = excluded.description,
			public = excluded.public,
			updated_at = excluded.updated_at
	`, gist.ID, gist.Path, gist.Description, public, gist.CreatedAt.Unix(), gist.UpdatedAt.Unix())
	if err != nil {
		return err
	}

	// Delete existing files
	_, err = tx.ExecContext(ctx, `DELETE FROM files WHERE gist_id = ?`, gist.ID)
	if err != nil {
		return err
	}

	// Insert new files
	for _, file := range files {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO files (gist_id, original_path, gist_filename, size, language)
			VALUES (?, ?, ?, ?, ?)
		`, file.GistID, file.OriginalPath, file.GistFilename, file.Size, file.Language)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetGist retrieves a gist by ID
func (r *SQLiteRepository) GetGist(ctx context.Context, id string) (*Gist, error) {
	var gist Gist
	var public int
	var createdBy, updatedBy int64

	err := r.db.QueryRowContext(ctx, `
		SELECT id, path, description, public, created_at, updated_at
		FROM gists WHERE id = ?
	`, id).Scan(&gist.ID, &gist.Path, &gist.Description, &public, &createdBy, &updatedBy)

	if err != nil {
		return nil, err
	}

	gist.Public = public == 1
	gist.CreatedAt = time.Unix(createdBy, 0)
	gist.UpdatedAt = time.Unix(updatedBy, 0)

	return &gist, nil
}

// GetGistByPath retrieves a gist by wiki path
func (r *SQLiteRepository) GetGistByPath(ctx context.Context, path string) (*Gist, error) {
	var gist Gist
	var public int
	var createdBy, updatedBy int64

	err := r.db.QueryRowContext(ctx, `
		SELECT id, path, description, public, created_at, updated_at
		FROM gists WHERE path = ?
	`, path).Scan(&gist.ID, &gist.Path, &gist.Description, &public, &createdBy, &updatedBy)

	if err != nil {
		return nil, err
	}

	gist.Public = public == 1
	gist.CreatedAt = time.Unix(createdBy, 0)
	gist.UpdatedAt = time.Unix(updatedBy, 0)

	return &gist, nil
}

// ListGists returns all gists
func (r *SQLiteRepository) ListGists(ctx context.Context) ([]*Gist, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, path, description, public, created_at, updated_at
		FROM gists ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gists []*Gist
	for rows.Next() {
		var gist Gist
		var public int
		var createdBy, updatedBy int64

		err := rows.Scan(&gist.ID, &gist.Path, &gist.Description, &public, &createdBy, &updatedBy)
		if err != nil {
			return nil, err
		}

		gist.Public = public == 1
		gist.CreatedAt = time.Unix(createdBy, 0)
		gist.UpdatedAt = time.Unix(updatedBy, 0)

		gists = append(gists, &gist)
	}

	return gists, nil
}

// GetGistFiles returns all files for a given gist
func (r *SQLiteRepository) GetGistFiles(ctx context.Context, gistID string) ([]*File, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, gist_id, original_path, gist_filename, size, language
		FROM files WHERE gist_id = ?
		ORDER BY original_path
	`, gistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	for rows.Next() {
		var file File
		err := rows.Scan(&file.ID, &file.GistID, &file.OriginalPath, &file.GistFilename, &file.Size, &file.Language)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}

// DeleteGist removes a gist and all its files
func (r *SQLiteRepository) DeleteGist(ctx context.Context, gistID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM gists WHERE id = ?`, gistID)
	return err
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// GistJSON is a helper for JSON marshaling (optional for MCP)
type GistJSON struct {
	ID          string            `json:"id"`
	Path        string            `json:"path"`
	Description string            `json:"description"`
	Public      bool              `json:"public"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Files       []*File           `json:"files,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ToJSON converts a Gist with files to JSON-serializable format
func (r *SQLiteRepository) GetGistJSON(ctx context.Context, gistID string) (*GistJSON, error) {
	gist, err := r.GetGist(ctx, gistID)
	if err != nil {
		return nil, err
	}

	files, err := r.GetGistFiles(ctx, gistID)
	if err != nil {
		return nil, err
	}

	return &GistJSON{
		ID:          gist.ID,
		Path:        gist.Path,
		Description: gist.Description,
		Public:      gist.Public,
		CreatedAt:   gist.CreatedAt,
		UpdatedAt:   gist.UpdatedAt,
		Files:       files,
	}, nil
}

// ToJSONString returns JSON string representation
func (g *GistJSON) ToJSONString() (string, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
