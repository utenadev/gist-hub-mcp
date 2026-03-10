package db

const schema = `
-- gists table: stores GitHub Gist metadata with wiki path
CREATE TABLE IF NOT EXISTS gists (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL UNIQUE,
    description TEXT,
    public INTEGER DEFAULT 0,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- files table: stores individual files within gists
CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gist_id TEXT NOT NULL,
    original_path TEXT NOT NULL,
    gist_filename TEXT NOT NULL,
    size INTEGER DEFAULT 0,
    language TEXT,
    FOREIGN KEY (gist_id) REFERENCES gists(id) ON DELETE CASCADE
);

-- Indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_gists_path ON gists(path);
CREATE INDEX IF NOT EXISTS idx_files_gist_id ON files(gist_id);
CREATE INDEX IF NOT EXISTS idx_files_original_path ON files(original_path);
`

// GetSchema returns the SQL schema for initialization
func GetSchema() string {
	return schema
}
