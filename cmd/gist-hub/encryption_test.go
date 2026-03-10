package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yklcs/gist-hub-mcp/internal/crypto"
)

// Test encryption/decryption with scanDirectory
func TestScanDirectoryWithEncryption(t *testing.T) {
	// Set test passphrase
	passphrase = "test-secure-passphrase-123"

	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "secret.txt")
	originalContent := "This is a secret message"
	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Scan directory (should encrypt content)
	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	// Verify content is encrypted
	if _, ok := files["secret.txt"]; !ok {
		t.Error("Expected secret.txt to be in files")
	}

	encryptedContent := files["secret.txt"].Content
	if encryptedContent == originalContent {
		t.Error("Expected content to be encrypted, but it matches original")
	}

	// Verify encryption by trying to decrypt
	decrypted, err := crypto.Decrypt(encryptedContent, passphrase)
	if err != nil {
		t.Fatalf("Failed to decrypt content: %v", err)
	}

	if string(decrypted) != originalContent {
		t.Errorf("Expected decrypted content '%s', got '%s'", originalContent, string(decrypted))
	}
}

// Test scanDirectory without encryption
func TestScanDirectoryWithoutEncryption(t *testing.T) {
	// Clear passphrase (no encryption)
	passphrase = ""

	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "public.txt")
	originalContent := "This is public content"
	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	// Verify content is NOT encrypted
	content := files["public.txt"].Content
	if content != originalContent {
		t.Errorf("Expected unencrypted content '%s', got '%s'", originalContent, content)
	}
}

// Test encryption with empty passphrase
func TestEncryptionWithEmptyPassphrase(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.md")
	originalContent := "# Test Document"
	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test with empty passphrase (should not encrypt)
	passphrase = ""
	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	content := files["test.md"].Content
	if content != originalContent {
		t.Errorf("Expected content to remain unchanged when passphrase is empty, got '%s'", content)
	}

	// Test with passphrase (should encrypt)
	passphrase = "test-key"
	files, err = scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	content = files["test.md"].Content
	if content == originalContent {
		t.Error("Expected content to be encrypted with passphrase, but it matches original")
	}

	// Verify decryption
	decrypted, err := crypto.Decrypt(content, passphrase)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != originalContent {
		t.Errorf("Decrypted content doesn't match original: got '%s'", string(decrypted))
	}
}

// Test multi-file encryption
func TestMultiFileEncryption(t *testing.T) {
	passphrase = "multi-file-pass"

	tmpDir := t.TempDir()

	files := map[string]string{
		"doc1.md":  "Content of document 1",
		"doc2.txt": "Content of document 2",
		"doc3.dat": "Binary content for doc3",
	}

	for filename, content := range files {
		testFile := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	scannedFiles, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	if len(scannedFiles) != len(files) {
		t.Fatalf("Expected %d files, got %d", len(files), len(scannedFiles))
	}

	// Verify all files are encrypted and can be decrypted
	for filename, expectedContent := range files {
		file, ok := scannedFiles[filename]
		if !ok {
			t.Errorf("Expected file %s not found", filename)
			continue
		}

		if file.Content == expectedContent {
			t.Errorf("File %s was not encrypted", filename)
			continue
		}

		decrypted, err := crypto.Decrypt(file.Content, passphrase)
		if err != nil {
			t.Errorf("Failed to decrypt file %s: %v", filename, err)
			continue
		}

		if string(decrypted) != expectedContent {
			t.Errorf("Decrypted content for file %s doesn't match: got '%s'", filename, string(decrypted))
		}
	}
}

// Test encryption with subdirectories
func TestEncryptionWithSubdirectories(t *testing.T) {
	passphrase = "nested-pass"

	tmpDir := t.TempDir()

	// Create subdirectories
	subDir := filepath.Join(tmpDir, "docs", "guide")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create files in different directories
	rootFile := filepath.Join(tmpDir, "readme.md")
	subFile := filepath.Join(subDir, "tutorial.md")

	rootContent := "# README"
	subContent := "Tutorial content"

	if err := os.WriteFile(rootFile, []byte(rootContent), 0644); err != nil {
		t.Fatalf("Failed to create root file: %v", err)
	}
	if err := os.WriteFile(subFile, []byte(subContent), 0644); err != nil {
		t.Fatalf("Failed to create sub file: %v", err)
	}

	files, err := scanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("scanDirectory() error = %v", err)
	}

	// Verify both files are encrypted
	if len(files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(files))
	}

	// Check root file
	if _, ok := files["readme.md"]; !ok {
		t.Error("Expected readme.md")
	} else {
		decrypted, err := crypto.Decrypt(files["readme.md"].Content, passphrase)
		if err != nil {
			t.Errorf("Failed to decrypt readme.md: %v", err)
		} else if string(decrypted) != rootContent {
			t.Errorf("Decrypted content mismatch for readme.md")
		}
	}

	// Check subdirectory file
	if _, ok := files["docs\\guide\\tutorial.md"]; !ok {
		t.Error("Expected docs\\guide\\tutorial.md")
	} else {
		decrypted, err := crypto.Decrypt(files["docs\\guide\\tutorial.md"].Content, passphrase)
		if err != nil {
			t.Errorf("Failed to decrypt tutorial.md: %v", err)
		} else if string(decrypted) != subContent {
			t.Errorf("Decrypted content mismatch for tutorial.md")
		}
	}
}

// Test encryption helper functions
func TestEncryptionHelpers(t *testing.T) {
	// Test IsEncryptionEnabled
	passphrase = ""
	if IsEncryptionEnabled() {
		t.Error("IsEncryptionEnabled returned true for empty passphrase")
	}

	passphrase = "test"
	if !IsEncryptionEnabled() {
		t.Error("IsEncryptionEnabled returned false for non-empty passphrase")
	}

	// Test GetPassphrase
	passphrase = "my-secret-pass"
	if GetPassphrase() != "my-secret-pass" {
		t.Errorf("GetPassphrase returned '%s', expected 'my-secret-pass'", GetPassphrase())
	}
}

// Clean up passphrase after tests
func TestCleanup(t *testing.T) {
	passphrase = ""
}
