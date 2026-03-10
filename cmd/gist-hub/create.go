package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yklcs/gist-hub-mcp/internal/crypto"
	"github.com/yklcs/gist-hub-mcp/internal/github"
)

var createCmd = &cobra.Command{
	Use:   "create <dir>",
	Short: "Create a new gist from directory",
	Long:  `Create a new GitHub Gist from files in the specified directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createGist(args[0])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createGist(dirPath string) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Failed to get absolute path: %v\n", err)
		return
	}

	files, err := scanDirectory(absPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Failed to scan directory: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stdout, "Error: No files found in directory\n")
		return
	}

	gist := &github.Gist{
		Description: fmt.Sprintf("%s From %s", GistHubPrefix, filepath.Base(absPath)),
		Public:      false,
		Files:       files,
	}

	created, err := client.CreateGist(gist)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Failed to create gist: %v\n", err)
		return
	}

	fmt.Fprintf(os.Stdout, "✓ Gist created successfully!\n")
	fmt.Fprintf(os.Stdout, "ID: %s\n", created.ID)
	fmt.Fprintf(os.Stdout, "URL: %s\n", created.HTMLURL)

	if IsEncryptionEnabled() {
		fmt.Fprintf(os.Stdout, "Note: Encrypted with passphrase\n")
	}
}

func scanDirectory(dirPath string) (map[string]github.GistFile, error) {
	files := make(map[string]github.GistFile)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			// Encrypt content if passphrase is provided
			if IsEncryptionEnabled() {
				encrypted, err := crypto.Encrypt(content, GetPassphrase())
				if err != nil {
					return fmt.Errorf("failed to encrypt file %s: %w", path, err)
				}
				content = []byte(encrypted)
			}

			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			// Ensure forward slashes for Gist path compatibility
			relPath = filepath.ToSlash(relPath)

			// Convert forward slashes to backslashes for GistPad compatibility
			// GistPad interprets backslashes as directory separators
			fileName := strings.ReplaceAll(relPath, "/", "\\")
			files[fileName] = github.GistFile{
				Content: string(content),
			}
		}

		return nil
	})

	return files, err
}
