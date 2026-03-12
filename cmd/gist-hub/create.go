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
	Use:   "create <path>",
	Short: "Create a new gist from directory or file",
	Long:  `Create a new GitHub Gist from files in the specified directory or a single file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createGist(args[0])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createGist(targetPath string) {
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Failed to get absolute path: %v\n", err)
		return
	}

	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Failed to stat target: %v\n", err)
		return
	}

	files := make(map[string]github.GistFile)

	if info.IsDir() {
		files, err = scanDirectory(absPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to scan directory: %v\n", err)
			return
		}
	} else {
		// Single file mode
		content, err := os.ReadFile(absPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: Failed to read file: %v\n", err)
			return
		}

		if IsEncryptionEnabled() {
			encrypted, err := crypto.Encrypt(content, GetPassphrase())
			if err != nil {
				fmt.Fprintf(os.Stdout, "Error: Failed to encrypt file: %v\n", err)
				return
			}
			content = []byte(encrypted)
		}

		files[filepath.Base(absPath)] = github.GistFile{
			Content: string(content),
		}
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stdout, "Error: No files found to upload\n")
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

			// GistPad uses backslashes for directory separators.
			// However, some API environments might reject it in initial creation.
			// Let's stick to the user's "backslash" requirement but ensure it's sanitized.
			fileName := strings.ReplaceAll(relPath, "/", "\\")
			
			files[fileName] = github.GistFile{
				Content: string(content),
			}
		}

		return nil
	})

	return files, err
}
