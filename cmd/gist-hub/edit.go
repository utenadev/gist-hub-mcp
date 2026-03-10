package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yklcs/gist-hub-mcp/internal/github"
)

var editCmd = &cobra.Command{
	Use:   "edit <id> <dir>",
	Short: "Update an existing gist from directory",
	Long:  `Update an existing GitHub Gist with files from the specified directory.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		editGist(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func editGist(gistID, dirPath string) {
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

	gist := &github.Gist{
		Description: fmt.Sprintf("%s From %s", GistHubPrefix, filepath.Base(absPath)),
		Files:       files,
	}

	updated, err := client.UpdateGist(gistID, gist)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Failed to update gist: %v\n", err)
		return
	}

	fmt.Fprintf(os.Stdout, "✓ Gist updated successfully!\n")
	fmt.Fprintf(os.Stdout, "ID: %s\n", updated.ID)
	fmt.Fprintf(os.Stdout, "URL: %s\n", updated.HTMLURL)
}
