package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonasbn/baseline/internal/sources/bitbucket"
	"github.com/jonasbn/baseline/internal/sources/github"
	"github.com/jonasbn/baseline/internal/types"
	"github.com/jonasbn/baseline/internal/worker"
	"github.com/spf13/cobra"
)

var (
	useSSH bool
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories from the specified source into the target directory",
	Long: `Clone repositories from the specified source (GitHub or Bitbucket) into the target 
directory and update existing ones.

This command fetches all repositories from the organization and clones them as bare 
repositories, setting read-only permissions for searching purposes.

Use the --ssh flag to clone using SSH URLs instead of HTTPS URLs, which is useful 
when you have SSH keys configured and want to avoid HTTPS authentication issues.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Create the appropriate source client
		var sourceClient types.RepositorySource

		switch source {
		case "github":
			sourceClient = github.NewGitHubClient(githubToken)
		case "bitbucket":
			sourceClient = bitbucket.NewBitbucketClient(bitbucketUser, bitbucketToken, verbose)
		default:
			return fmt.Errorf("unsupported source: %s (supported: github, bitbucket)", source)
		}

		if verbose {
			fmt.Printf("Cloning repositories from %s for organization: %s\n", source, organization)
			fmt.Printf("Target directory: %s\n", directory)
			fmt.Printf("Using %d concurrent threads\n", threads)
			if useSSH {
				fmt.Printf("Using SSH URLs for cloning\n")
			} else {
				fmt.Printf("Using HTTPS URLs for cloning\n")
			}
		}

		// Create target directory if it doesn't exist
		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("failed to create target directory %s: %w", directory, err)
		}

		// Fetch repositories
		repositories, err := sourceClient.GetRepositories(ctx, organization)
		if err != nil {
			return fmt.Errorf("failed to fetch repositories: %w", err)
		}

		// Convert to SSH URLs if --ssh flag is set
		if useSSH {
			for i := range repositories {
				if repositories[i].SSHURL != "" {
					repositories[i].CloneURL = repositories[i].SSHURL
					if verbose {
						fmt.Printf("Using SSH URL for %s: %s\n", repositories[i].FullName, repositories[i].SSHURL)
					}
				} else if verbose {
					fmt.Printf("Warning: No SSH URL available for %s, using HTTPS\n", repositories[i].FullName)
				}
			}
		}

		fmt.Printf("Found %d repositories to clone\n", len(repositories))

		// Create worker pool and start cloning
		wp := worker.NewWorkerPool(threads, verbose)
		resultChan := wp.CloneRepositories(ctx, repositories, directory)

		// Process results
		var successful, failed, skipped int
		for result := range resultChan {
			if result.Error != nil {
				if verbose {
					fmt.Printf("❌ %s: %v\n", result.Repository.FullName, result.Error)
				}
				if result.Error.Error() == fmt.Sprintf("repository already exists at %s",
					filepath.Join(directory, result.Repository.Owner, result.Repository.Name)) {
					skipped++
				} else {
					failed++
				}
			} else if result.Success {
				if verbose {
					fmt.Printf("✅ %s (%.2fs)\n", result.Repository.FullName, result.Duration.Seconds())
				}
				successful++
			}
		}

		// Print summary
		fmt.Printf("\nClone Summary:\n")
		fmt.Printf("  Successful: %d\n", successful)
		fmt.Printf("  Skipped:    %d (already exists)\n", skipped)
		fmt.Printf("  Failed:     %d\n", failed)

		if failed > 0 {
			return fmt.Errorf("some repositories failed to clone")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().BoolVar(&useSSH, "ssh", false, "Use SSH URLs for cloning instead of HTTPS")
}
