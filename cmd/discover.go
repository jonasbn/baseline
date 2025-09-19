package cmd

import (
	"context"
	"fmt"

	"github.com/jonasbn/baseline/internal/sources/bitbucket"
	"github.com/jonasbn/baseline/internal/sources/github"
	"github.com/jonasbn/baseline/internal/types"
	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "List repositories available in the specified source",
	Long: `Discover and list all repositories available in the specified source (GitHub or Bitbucket)
for the given organization.

This command helps you see what repositories are available before cloning them.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Create the appropriate source client
		var sourceClient types.RepositorySource
		var err error

		switch source {
		case "github":
			sourceClient = github.NewGitHubClient(githubToken)
		case "bitbucket":
			sourceClient = bitbucket.NewBitbucketClient(bitbucketUser, bitbucketPass)
		default:
			return fmt.Errorf("unsupported source: %s (supported: github, bitbucket)", source)
		}

		if verbose {
			fmt.Printf("Discovering repositories from %s for organization/user: %s\n", source, organization)
		}

		// Fetch repositories
		repositories, err := sourceClient.GetRepositories(ctx, organization)
		if err != nil {
			return fmt.Errorf("failed to fetch repositories: %w", err)
		}

		// Display results
		fmt.Printf("Found %d repositories in %s/%s:\n", len(repositories), source, organization)
		fmt.Println()

		for _, repo := range repositories {
			fmt.Printf("  %-30s %s\n", repo.Name, repo.Description)
			if verbose {
				fmt.Printf("    Full name: %s\n", repo.FullName)
				fmt.Printf("    Clone URL: %s\n", repo.CloneURL)
				fmt.Printf("    Language:  %s\n", repo.Language)
				fmt.Printf("    Private:   %t\n", repo.Private)
				fmt.Printf("    Updated:   %s\n", repo.UpdatedAt.Format("2006-01-02 15:04:05"))
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}
