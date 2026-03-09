package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yklcs/gist-hub-mcp/internal/github"
)

const (
	GistHubPrefix = "gist-hub:"
)

var (
	token   string
	client  *github.Client
	Version = "0.1.0"
)

func init() {
	var err error
	token, err = github.GetToken("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get GitHub token: %v\n", err)
		os.Exit(1)
	}
	client = github.NewClient(token)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "list":
		listGists()
	case "get":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: gist ID is required\n")
			printUsage()
			os.Exit(1)
		}
		getGist(os.Args[2])
	case "create":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: directory path is required\n")
			printUsage()
			os.Exit(1)
		}
		createGist(os.Args[2])
	case "edit":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "Error: gist ID and directory path are required\n")
			printUsage()
			os.Exit(1)
		}
		editGist(os.Args[2], os.Args[3])
	case "--help", "-h":
		printUsage()
	case "--version":
		fmt.Println("gist-hub version", Version)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func listGists() {
	gists, err := client.ListGists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to list gists: %v\n", err)
		os.Exit(1)
	}

	count := 0
	for _, gist := range gists {
		if strings.HasPrefix(gist.Description, GistHubPrefix) {
			fmt.Printf("%s - %s\n", gist.ID, gist.Description)
			for filename := range gist.Files {
				fmt.Printf("  %s\n", filename)
			}
			count++
		}
	}

	if count == 0 {
		fmt.Println("No gist-hub gists found.")
	}
}

func getGist(gistID string) {
	gist, err := client.GetGist(gistID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get gist: %v\n", err)
		os.Exit(1)
	}

	if !strings.HasPrefix(gist.Description, GistHubPrefix) {
		fmt.Fprintf(os.Stderr, "Error: Gist does not have the required prefix '%s'\n", GistHubPrefix)
		os.Exit(1)
	}

	fmt.Printf("=== %s ===\n", gist.ID)
	fmt.Printf("Description: %s\n", gist.Description)
	fmt.Printf("URL: %s\n", gist.HTMLURL)
	fmt.Printf("Files:\n")

	for filename, file := range gist.Files {
		fmt.Printf("\n--- %s ---\n", filename)
		fmt.Printf("Language: %s\n", file.Language)
		fmt.Printf("Size: %d bytes\n", file.Size)
		fmt.Printf("Content:\n%s\n", file.Content)
	}
}

func createGist(dirPath string) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}

	files, err := scanDirectory(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to scan directory: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No files found in directory\n")
		os.Exit(1)
	}

	gist := &github.Gist{
		Description: fmt.Sprintf("%s From %s", GistHubPrefix, filepath.Base(absPath)),
		Public:      false,
		Files:       files,
	}

	created, err := client.CreateGist(gist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create gist: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Gist created successfully!\n")
	fmt.Printf("ID: %s\n", created.ID)
	fmt.Printf("URL: %s\n", created.HTMLURL)
}

func editGist(gistID, dirPath string) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}

	files, err := scanDirectory(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to scan directory: %v\n", err)
		os.Exit(1)
	}

	gist := &github.Gist{
		Description: fmt.Sprintf("%s From %s", GistHubPrefix, filepath.Base(absPath)),
		Files:       files,
	}

	updated, err := client.UpdateGist(gistID, gist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to update gist: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Gist updated successfully!\n")
	fmt.Printf("ID: %s\n", updated.ID)
	fmt.Printf("URL: %s\n", updated.HTMLURL)
}

func scanDirectory(dirPath string) (map[string]github.GistFile, error) {
	files := make(map[string]github.GistFile)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != dirPath {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			files[relPath] = github.GistFile{
				Content: string(content),
			}
		}

		return nil
	})

	return files, err
}

func printUsage() {
	fmt.Println("gist-hub - GitHub Gist CLI with gist-hub prefix support")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gist-hub list                    List all gist-hub gists")
	fmt.Println("  gist-hub get <id>               Get a gist-hub gist by ID")
	fmt.Println("  gist-hub create <dir>           Create a new gist from directory")
	fmt.Println("  gist-hub edit <id> <dir>        Update an existing gist from directory")
	fmt.Println("  gist-hub --help                 Show this help message")
	fmt.Println("  gist-hub --version              Show version information")
	fmt.Println()
	fmt.Println("Notes:")
	fmt.Println("  - All gists must have the prefix 'gist-hub:' in their description")
	fmt.Println("  - Requires gh CLI to be authenticated")
}
