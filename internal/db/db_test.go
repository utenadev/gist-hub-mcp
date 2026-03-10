package db

import (
	"context"
	"testing"
	"time"
)

// setupTestDB creates an in-memory database for testing
func setupTestDB(t *testing.T) Repository {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	ctx := context.Background()
	if err := repo.Init(ctx); err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
		t.Cleanup(func() { repo.Close() })
	}

	t.Cleanup(func() { repo.Close() })
	return repo
}

func TestSaveAndGetGist(t *testing.T) {
	ctx := context.Background()
	repo := setupTestDB(t)

	// Create test gist
	now := time.Now()
	gist := &Gist{
		ID:          "abc123",
		Path:        "docs/arch/spec.md",
		Description: "Architecture Specification",
		Public:      false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	files := []*File{
		{
			GistID:       "abc123",
			OriginalPath: "docs/arch/spec.md",
			GistFilename: "docs\\arch\\spec.md",
			Size:         1024,
			Language:     "markdown",
		},
	}

	// Save gist
	err := repo.SaveGist(ctx, gist, files)
	if err != nil {
		t.Fatalf("SaveGist failed: %v", err)
	}

	// Get gist by ID
	retrieved, err := repo.GetGist(ctx, "abc123")
	if err != nil {
		t.Fatalf("GetGist failed: %v", err)
	}

	if retrieved == nil {
		t.Error("GetGist returned nil")
		return
	}

	if retrieved.ID != gist.ID {
		t.Errorf("Expected ID %s, got %s", gist.ID, retrieved.ID)
	}

	if retrieved.Path != gist.Path {
		t.Errorf("Expected Path %s, got %s", gist.Path, retrieved.Path)
	}

	if retrieved.Description != gist.Description {
		t.Errorf("Expected Description %s, got %s", gist.Description, retrieved.Description)
	}

	if retrieved.Public != gist.Public {
		t.Errorf("Expected Public %v, got %v", gist.Public, retrieved.Public)
	}
}

func TestGetGistByPath(t *testing.T) {
	ctx := context.Background()
	repo := setupTestDB(t)

	now := time.Now()
	gist := &Gist{
		ID:          "test123",
		Path:        "docs/guide/tutorial.md",
		Description: "Tutorial",
		Public:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	files := []*File{
		{
			GistID:       "test123",
			OriginalPath: "docs/guide/tutorial.md",
			GistFilename: "docs\\guide\\tutorial.md",
			Size:         512,
			Language:     "markdown",
		},
	}

	err := repo.SaveGist(ctx, gist, files)
	if err != nil {
		t.Fatalf("SaveGist failed: %v", err)
	}

	// Get gist by path
	retrieved, err := repo.GetGistByPath(ctx, "docs/guide/tutorial.md")
	if err != nil {
		t.Fatalf("GetGistByPath failed: %v", err)
	}

	if retrieved == nil {
		t.Error("GetGistByPath returned nil")
		return
	}

	if retrieved.ID != gist.ID {
		t.Errorf("Expected ID %s, got %s", gist.ID, retrieved.ID)
	}
}

func TestListGists(t *testing.T) {
	ctx := context.Background()
	repo := setupTestDB(t)

	now := time.Now()
	gists := []*Gist{
		{ID: "1", Path: "docs/a.md", Description: "A", Public: true, CreatedAt: now, UpdatedAt: now},
		{ID: "2", Path: "docs/b.md", Description: "B", Public: false, CreatedAt: now, UpdatedAt: now},
		{ID: "3", Path: "docs/c.md", Description: "C", Public: true, CreatedAt: now, UpdatedAt: now},
	}

	for _, gist := range gists {
		files := []*File{
			{
				GistID:       gist.ID,
				OriginalPath: gist.Path,
				GistFilename: "docs\\" + gist.ID + ".md",
				Size:         100,
				Language:     "markdown",
			},
		}
		err := repo.SaveGist(ctx, gist, files)
		if err != nil {
			t.Fatalf("SaveGist failed: %v", err)
		}
	}

	// List gists
	list, err := repo.ListGists(ctx)
	if err != nil {
		t.Fatalf("ListGists failed: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 gists, got %d", len(list))
	}
}

func TestGetGistFiles(t *testing.T) {
	ctx := context.Background()
	repo := setupTestDB(t)

	now := time.Now()
	gist := &Gist{
		ID:          "filetest",
		Path:        "docs/comprehensive.md",
		Description: "Test files",
		Public:      false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	files := []*File{
		{GistID: "filetest", OriginalPath: "docs/comprehensive.md", GistFilename: "docs\\comprehensive.md", Size: 1000, Language: "markdown"},
		{GistID: "filetest", OriginalPath: "docs/appendix.md", GistFilename: "docs\\appendix.md", Size: 500, Language: "markdown"},
		{GistID: "filetest", OriginalPath: "docs/index.md", GistFilename: "docs\\index.md", Size: 200, Language: "markdown"},
	}

	err := repo.SaveGist(ctx, gist, files)
	if err != nil {
		t.Fatalf("SaveGist failed: %v", err)
	}

	// Get gist files
	retrievedFiles, err := repo.GetGistFiles(ctx, "filetest")
	if err != nil {
		t.Fatalf("GetGistFiles failed: %v", err)
	}

	if len(retrievedFiles) != 3 {
		t.Errorf("Expected 3 files, got %d", len(retrievedFiles))
	}

	// Verify original_path vs gist_filename
	for _, f := range retrievedFiles {
		if f.GistID != "filetest" {
			t.Errorf("Expected gist_id filetest, got %s", f.GistID)
		}
	}
}

func TestUpdateGist(t *testing.T) {
	ctx := context.Background()
	repo := setupTestDB(t)

	now := time.Now()
	gist := &Gist{
		ID:          "update123",
		Path:        "docs/original.md",
		Description: "Original",
		Public:      false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	files := []*File{
		{
			GistID:       "update123",
			OriginalPath: "docs/original.md",
			GistFilename: "docs\\original.md",
			Size:         300,
			Language:     "markdown",
		},
	}

	err := repo.SaveGist(ctx, gist, files)
	if err != nil {
		t.Fatalf("SaveGist failed: %v", err)
	}

	// Update gist
	updatedNow := time.Now()
	gist.Path = "docs/updated.md"
	gist.Description = "Updated"
	gist.Public = true
	gist.UpdatedAt = updatedNow

	files[0].OriginalPath = "docs/updated.md"
	files[0].GistFilename = "docs\\updated.md"
	files[0].Size = 400

	err = repo.SaveGist(ctx, gist, files)
	if err != nil {
		t.Fatalf("Update SaveGist failed: %v", err)
	}

	// Verify update
	retrieved, err := repo.GetGist(ctx, "update123")
	if err != nil {
		t.Fatalf("GetGist failed: %v", err)
	}

	if retrieved.Path != "docs/updated.md" {
		t.Errorf("Expected updated path docs/updated.md, got %s", retrieved.Path)
	}

	if retrieved.Description != "Updated" {
		t.Errorf("Expected updated description 'Updated', got '%s'", retrieved.Description)
	}

	if !retrieved.Public {
		t.Error("Expected Public to be true after update")
	}

	retrievedFiles, err := repo.GetGistFiles(ctx, "update123")
	if err != nil {
		t.Fatalf("GetGistFiles failed: %v", err)
	}

	if len(retrievedFiles) != 1 {
		t.Errorf("Expected 1 file after update, got %d", len(retrievedFiles))
	}

	if retrievedFiles[0].OriginalPath != "docs/updated.md" {
		t.Errorf("Expected updated file path docs/updated.md, got %s", retrievedFiles[0].OriginalPath)
	}

	if retrievedFiles[0].Size != 400 {
		t.Errorf("Expected updated file size 400, got %d", retrievedFiles[0].Size)
	}
}

func TestDeleteGist(t *testing.T) {
	ctx := context.Background()
	repo := setupTestDB(t)

	now := time.Now()
	gist := &Gist{
		ID:          "delete123",
		Path:        "docs/delete.md",
		Description: "Delete me",
		Public:      false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	files := []*File{
		{
			GistID:       "delete123",
			OriginalPath: "docs/delete.md",
			GistFilename: "docs\\delete.md",
			Size:         100,
			Language:     "markdown",
		},
	}

	err := repo.SaveGist(ctx, gist, files)
	if err != nil {
		t.Fatalf("SaveGist failed: %v", err)
	}

	// Verify file exists
	_, err = repo.GetGistFiles(ctx, "delete123")
	if err != nil {
		t.Fatalf("GetGistFiles failed before delete: %v", err)
	}

	// Delete gist
	err = repo.DeleteGist(ctx, "delete123")
	if err != nil {
		t.Fatalf("DeleteGist failed: %v", err)
	}

	// Verify deletion
	_, err = repo.GetGist(ctx, "delete123")
	if err == nil {
		t.Error("Expected error when getting deleted gist, got nil")
	}

	// Verify files are also deleted
	files, err = repo.GetGistFiles(ctx, "delete123")
	if err != nil {
		t.Fatalf("GetGistFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files after deleting gist, got %d", len(files))
	}
}

func TestGetGistJSON(t *testing.T) {
	ctx := context.Background()
	repo := setupTestDB(t)

	now := time.Now()
	gist := &Gist{
		ID:          "json123",
		Path:        "docs/api/spec.md",
		Description: "API Specification",
		Public:      false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	files := []*File{
		{
			GistID:       "json123",
			OriginalPath: "docs/api/spec.md",
			GistFilename: "docs\\api\\spec.md",
			Size:         2048,
			Language:     "markdown",
		},
	}

	err := repo.SaveGist(ctx, gist, files)
	if err != nil {
		t.Fatalf("SaveGist failed: %v", err)
	}

	// Get gist JSON
	gistJSON, err := repo.(*SQLiteRepository).GetGistJSON(ctx, "json123")
	if err != nil {
		t.Fatalf("GetGistJSON failed: %v", err)
	}

	if gistJSON == nil {
		t.Error("GetGistJSON returned nil")
		return
	}

	if gistJSON.ID != gist.ID {
		t.Errorf("Expected ID %s, got %s", gist.ID, gistJSON.ID)
	}

	if len(gistJSON.Files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(gistJSON.Files))
	}

	if gistJSON.Files[0].OriginalPath != "docs/api/spec.md" {
		t.Errorf("Expected OriginalPath docs/api/spec.md, got %s", gistJSON.Files[0].OriginalPath)
	}

	// Test ToJSONString
	jsonStr, err := gistJSON.ToJSONString()
	if err != nil {
		t.Fatalf("ToJSONString failed: %v", err)
	}

	if jsonStr == "" {
		t.Error("ToJSONString returned empty string")
	}
}

func TestSchemaInitialization(t *testing.T) {
	ctx := context.Background()
	repo := setupTestDB(t)

	// Verify tables exist
	var tableCount int
	err := repo.(*SQLiteRepository).db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='gists'
	`).Scan(&tableCount)
	if err != nil {
		t.Fatalf("Failed to query gists table: %v", err)
	}

	if tableCount != 1 {
		t.Errorf("Expected 1 gists table, got %d", tableCount)
	}

	err = repo.(*SQLiteRepository).db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='files'
	`).Scan(&tableCount)
	if err != nil {
		t.Fatalf("Failed to query files table: %v", err)
	}

	if tableCount != 1 {
		t.Errorf("Expected 1 files table, got %d", tableCount)
	}

	// Verify indexes
	var indexCount int
	err = repo.(*SQLiteRepository).db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND tbl_name='files'
	`).Scan(&indexCount)
	if err != nil {
		t.Fatalf("Failed to query files indexes: %v", err)
	}

	if indexCount < 2 { // Should have at least gist_id and original_path indexes
		t.Errorf("Expected at least 2 indexes on files table, got %d", indexCount)
	}
}

func TestClose(t *testing.T) {
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = repo.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Double close should not panic
	err = repo.Close()
	if err != nil {
		t.Fatalf("Double close failed: %v", err)
	}
}
