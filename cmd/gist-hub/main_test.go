package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create test files
	testFile1 := filepath.Join(tmpDir, "test1.txt")
	testFile2 := filepath.Join(tmpDir, "test2.md")

	if err := os.WriteFile(testFile1, []byte("content 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("content 2"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test scanning
	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	// Verify results
	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	if _, ok := files["test1.txt"]; !ok {
		t.Error("Expected test1.txt to be in files")
	}

	if _, ok := files["test2.md"]; !ok {
		t.Error("Expected test2.md to be in files")
	}

	// Verify content
	if files["test1.txt"].Content != "content 1" {
		t.Errorf("Expected content 'content 1', got '%s'", files["test1.txt"].Content)
	}

	if files["test2.md"].Content != "content 2" {
		t.Errorf("Expected content 'content 2', got '%s'", files["test2.md"].Content)
	}
}

func TestScanDirectoryWithSubdirectory(t *testing.T) {
	// Create a temporary directory with subdirectory
	tmpDir := t.TempDir()

	// Create files in root
	testFile1 := filepath.Join(tmpDir, "root.txt")
	if err := os.WriteFile(testFile1, []byte("root content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create subdirectory (should now be included)
	subDir := filepath.Join(tmpDir, "report")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create file in subdirectory
	testFile2 := filepath.Join(subDir, "2026-03-09.md")
	if err := os.WriteFile(testFile2, []byte("report content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test scanning
	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	// Verify both files are included with hierarchy
	if len(files) != 2 {
		t.Errorf("Expected 2 files (root + subdirectory), got %d", len(files))
	}

	if _, ok := files["root.txt"]; !ok {
		t.Error("Expected root.txt to be in files")
	}

	if _, ok := files["report/2026-03-09.md"]; !ok {
		t.Error("Expected report/2026-03-09.md to be in files")
	}

	// Verify content
	if files["report/2026-03-09.md"].Content != "report content" {
		t.Errorf("Expected content 'report content', got '%s'", files["report/2026-03-09.md"].Content)
	}
}


func TestScanDirectoryEmpty(t *testing.T) {
	// Create empty directory
	tmpDir := t.TempDir()

	// Test scanning
	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	// Verify directory scan
	if len(files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(files))
	}
}
