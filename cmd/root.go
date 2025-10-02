package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	directory      string
	githubToken    string
	bitbucketUser  string
	bitbucketToken string
	organization   string
	verbose        bool
	source         string
	threads        int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "baseline",
	Short: "A tool for creating a baseline of Git repositories for easy searching",
	Long: `baseline is a Go program that creates a baseline of Git repositories for easy searching.

It integrates with GitHub and Bitbucket to clone repositories into a specified directory,
setting permissions to disallow write access, making them suitable for searching rather 
than active development.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "./baseline", "Target directory for the baseline")
	rootCmd.PersistentFlags().StringVarP(&githubToken, "github-token", "g", "", "GitHub token for accessing private repositories")
	rootCmd.PersistentFlags().StringVarP(&bitbucketUser, "bitbucket-username", "u", "", "Bitbucket username or email for API authentication")
	rootCmd.PersistentFlags().StringVarP(&bitbucketToken, "bitbucket-token", "b", "", "Bitbucket API token (repository, project, or workspace access token)")
	rootCmd.PersistentFlags().StringVarP(&organization, "organization", "o", "jonasbn", "Organization to fetch repositories from")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output for debugging")
	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "github", "Source platform (github or bitbucket)")

	// Flags specific to clone and update commands
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 4, "Number of concurrent threads for cloning/updating repositories")
}
