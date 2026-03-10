package db

import (
	"context"
	"testing"
	"time"
)

func TestSQLiteRepository(t *testing.T) {
	ctx := context.Background()

	// Use in-memory database for testing
	repo, err := NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Initialize schema
	if err := repo.Init(ctx); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Test data
	gistID := "test-gist-id"
	now := time.Now().Truncate(time.Second) // SQLite INTEGER stores Unix timestamp
	gist := &Gist{
		ID:          gistID,
		Path:        "docs/test",
		Description: "gist-hub: Test Gist",
		Public:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	files := []*File{
		{
			GistID:       gistID,
			OriginalPath: "docs/test/readme.md",
			GistFilename: "docs\\test\\readme.md",
			Size:         100,
			Language:     "Markdown",
		},
		{
			GistID:       gistID,
			OriginalPath: "docs/test/config.json",
			GistFilename: "docs\\test\\config.json",
			Size:         200,
			Language:     "JSON",
		},
	}

	t.Run("Save and Get Gist", func(t *testing.T) {
		if err := repo.SaveGist(ctx, gist, files); err != nil {
			t.Fatalf("SaveGist failed: %v", err)
		}

		got, err := repo.GetGist(ctx, gistID)
		if err != nil {
			t.Fatalf("GetGist failed: %v", err)
		}

		if got.ID != gist.ID || got.Path != gist.Path || got.Description != gist.Description {
			t.Errorf("Gist metadata mismatch. Got %+v", got)
		}

		// Check timestamp (Unix seconds)
		if got.UpdatedAt.Unix() != gist.UpdatedAt.Unix() {
			t.Errorf("Timestamp mismatch. Expected %v, got %v", gist.UpdatedAt.Unix(), got.UpdatedAt.Unix())
		}
	})

	t.Run("Get Gist by Path", func(t *testing.T) {
		got, err := repo.GetGistByPath(ctx, "docs/test")
		if err != nil {
			t.Fatalf("GetGistByPath failed: %v", err)
		}
		if got.ID != gistID {
			t.Errorf("Expected gist ID %s, got %s", gistID, got.ID)
		}
	})

	t.Run("Get Gist Files", func(t *testing.T) {
		gotFiles, err := repo.GetGistFiles(ctx, gistID)
		if err != nil {
			t.Fatalf("GetGistFiles failed: %v", err)
		}

		if len(gotFiles) != 2 {
			t.Fatalf("Expected 2 files, got %d", len(gotFiles))
		}

		// Verify backslash name
		found := false
		for _, f := range gotFiles {
			if f.GistFilename == "docs\\test\\readme.md" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Backslash filename not found in retrieved files")
		}
	})

	t.Run("Update Gist (Upsert)", func(t *testing.T) {
		gist.Description = "Updated Description"
		newFiles := []*File{
			{
				GistID:       gistID,
				OriginalPath: "docs/test/readme.md",
				GistFilename: "docs\\test\\readme.md",
				Size:         150,
			},
		}

		if err := repo.SaveGist(ctx, gist, newFiles); err != nil {
			t.Fatalf("Update SaveGist failed: %v", err)
		}

		got, _ := repo.GetGist(ctx, gistID)
		if got.Description != "Updated Description" {
			t.Errorf("Expected updated description, got %s", got.Description)
		}

		gotFiles, _ := repo.GetGistFiles(ctx, gistID)
		if len(gotFiles) != 1 {
			t.Errorf("Expected 1 file after update, got %d", len(gotFiles))
		}
	})

	t.Run("Delete Gist", func(t *testing.T) {
		if err := repo.DeleteGist(ctx, gistID); err != nil {
			t.Fatalf("DeleteGist failed: %v", err)
		}

		_, err := repo.GetGist(ctx, gistID)
		if err == nil {
			t.Error("Expected error getting deleted gist, got nil")
		}

		gotFiles, _ := repo.GetGistFiles(ctx, gistID)
		if len(gotFiles) != 0 {
			t.Errorf("Expected 0 files after deletion, got %d", len(gotFiles))
		}
	})
}
