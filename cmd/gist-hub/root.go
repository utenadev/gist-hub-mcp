package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yklcs/gist-hub-mcp/internal/github"
)

const (
	GistHubPrefix = "gist-hub:"
	Version       = "0.1.0"
)

var (
	client *github.Client
)

var rootCmd = &cobra.Command{
	Use:     "gist-hub",
	Short:   "GitHub Gist CLI with gist-hub prefix support",
	Long:    `gist-hub is a CLI tool for managing GitHub Gists with automatic prefix filtering.`,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip client initialization for help and version commands
		if helpFlag := cmd.Flags().Changed("help"); helpFlag {
			return nil
		}

		// Initialize GitHub client with support for future token overrides
		token, err := github.GetToken("")
		if err != nil {
			return err
		}
		client = github.NewClient(token)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
