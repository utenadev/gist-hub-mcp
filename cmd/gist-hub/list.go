package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all gist-hub gists",
	Long:  `List all GitHub Gists that have the 'gist-hub:' prefix in their description.`,
	Run: func(cmd *cobra.Command, args []string) {
		listGists()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listGists() {
	gists, err := client.ListGists()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: Failed to list gists: %v\n", err)
		return
	}

	count := 0
	for _, gist := range gists {
		if strings.HasPrefix(gist.Description, GistHubPrefix) {
			fmt.Fprintf(os.Stdout, "%s - %s\n", gist.ID, gist.Description)
			for filename := range gist.Files {
				fmt.Fprintf(os.Stdout, "  %s\n", filename)
			}
			count++
		}
	}

	if count == 0 {
		fmt.Println("No gist-hub gists found.")
	}
}
