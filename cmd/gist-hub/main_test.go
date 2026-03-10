package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	testFile1 := filepath.Join(tmpDir, "test1.txt")
	testFile2 := filepath.Join(tmpDir, "test2.md")

	if err := os.WriteFile(testFile1, []byte("content 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("content 2"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	if _, ok := files["test1.txt"]; !ok {
		t.Error("Expected test1.txt to be in files")
	}

	if _, ok := files["test2.md"]; !ok {
		t.Error("Expected test2.md to be in files")
	}

	if files["test1.txt"].Content != "content 1" {
		t.Errorf("Expected content 'content 1', got '%s'", files["test1.txt"].Content)
	}

	if files["test2.md"].Content != "content 2" {
		t.Errorf("Expected content 'content 2', got '%s'", files["test2.md"].Content)
	}
}

func TestScanDirectoryWithSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()

	testFile1 := filepath.Join(tmpDir, "root.txt")
	if err := os.WriteFile(testFile1, []byte("root content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	subDir := filepath.Join(tmpDir, "report")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	testFile2 := filepath.Join(subDir, "2026-03-09.md")
	if err := os.WriteFile(testFile2, []byte("report content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files (root + subdirectory), got %d", len(files))
	}

	if _, ok := files["root.txt"]; !ok {
		t.Error("Expected root.txt to be in files")
	}

	if _, ok := files["report\\2026-03-09.md"]; !ok {
		t.Error("Expected report\\2026-03-09.md to be in files")
	}

	if files["report\\2026-03-09.md"].Content != "report content" {
		t.Errorf("Expected content 'report content', got '%s'", files["report\\2026-03-09.md"].Content)
	}
}
func TestScanDirectoryEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(files))
	}
}

func TestGistHubPrefix(t *testing.T) {
	expected := "gist-hub:"
	if GistHubPrefix != expected {
		t.Errorf("Expected prefix '%s', got '%s'", expected, GistHubPrefix)
	}
}

func TestVersion(t *testing.T) {
	expected := "0.1.0"
	if Version != expected {
		t.Errorf("Expected version '%s', got '%s'", expected, Version)
	}
}

func TestRootCmdExists(t *testing.T) {
	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}

	if rootCmd.Use != "gist-hub" {
		t.Errorf("Expected Use to be 'gist-hub', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short != "GitHub Gist CLI with gist-hub prefix support" {
		t.Errorf("Expected Short description, got '%s'", rootCmd.Short)
	}
}

func TestSubCommandsExist(t *testing.T) {
	subcommandNames := []string{"list", "get", "create", "edit"}

	for _, name := range subcommandNames {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand '%s' not found in rootCmd", name)
		}
	}
}

