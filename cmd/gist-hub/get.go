package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yklcs/gist-hub-mcp/internal/crypto"
)

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a gist-hub gist by ID",
	Long:  `Retrieve a GitHub Gist by its ID and display its details.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getGist(args[0])
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func getGist(gistID string) {
	gist, err := client.GetGist(gistID)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Failed to get gist: %v\n", err)
		return
	}

	if !strings.HasPrefix(gist.Description, GistHubPrefix) {
		fmt.Fprintf(os.Stdout, "Error: Gist does not have the required prefix '%s'\n", GistHubPrefix)
		return
	}

	fmt.Fprintf(os.Stdout, "=== %s ===\n", gist.ID)
	fmt.Fprintf(os.Stdout, "Description: %s\n", gist.Description)
	fmt.Fprintf(os.Stdout, "URL: %s\n", gist.HTMLURL)
	fmt.Fprintf(os.Stdout, "Files:\n")

	for filename, file := range gist.Files {
		content := file.Content

		// Decrypt content if passphrase is provided
		if IsEncryptionEnabled() {
			decrypted, err := crypto.Decrypt(content, GetPassphrase())
			if err != nil {
				fmt.Fprintf(os.Stdout, "\n--- %s ---\n", filename)
				fmt.Fprintf(os.Stdout, "Language: %s\n", file.Language)
				fmt.Fprintf(os.Stdout, "Size: %d bytes\n", file.Size)
				fmt.Fprintf(os.Stdout, "Content: [ERROR: Decryption failed - invalid passphrase or corrupted data]\n")
				continue
			}
			content = string(decrypted)
		}

		fmt.Fprintf(os.Stdout, "\n--- %s ---\n", filename)
		fmt.Fprintf(os.Stdout, "Language: %s\n", file.Language)
		fmt.Fprintf(os.Stdout, "Size: %d bytes\n", file.Size)
		fmt.Fprintf(os.Stdout, "Content:\n%s\n", content)
	}
}
