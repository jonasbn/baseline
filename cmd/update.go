package cmd

import (
	"context"
	"fmt"

	"github.com/jonasbn/baseline/internal/sources/bitbucket"
	"github.com/jonasbn/baseline/internal/sources/github"
	"github.com/jonasbn/baseline/internal/types"
	"github.com/jonasbn/baseline/internal/worker"
	"github.com/spf13/cobra"
)

var (
	updateUseSSH bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update repositories in the target directory from the specified source",
	Long: `Update repositories in the target directory from the specified source (GitHub or Bitbucket).

This command fetches the latest changes for existing repositories in the baseline directory.
Only repositories that already exist locally will be updated.

Use the --ssh flag to update using SSH URLs instead of HTTPS URLs, which is useful 
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
			fmt.Printf("Updating repositories from %s for organization: %s\n", source, organization)
			fmt.Printf("Target directory: %s\n", directory)
			fmt.Printf("Using %d concurrent threads\n", threads)
			if updateUseSSH {
				fmt.Printf("Using SSH URLs for updating\n")
			} else {
				fmt.Printf("Using HTTPS URLs for updating\n")
			}
		}

		// Fetch repositories
		repositories, err := sourceClient.GetRepositories(ctx, organization)
		if err != nil {
			return fmt.Errorf("failed to fetch repositories: %w", err)
		}

		// Convert to SSH URLs if --ssh flag is set
		if updateUseSSH {
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

		fmt.Printf("Found %d repositories to check for updates\n", len(repositories))

		// Create worker pool and start updating
		wp := worker.NewWorkerPool(threads, verbose)
		resultChan := wp.UpdateRepositories(ctx, repositories, directory)

		// Process results
		var successful, failed, skipped, updated int
		for result := range resultChan {
			if result.Error != nil {
				if verbose {
					fmt.Printf("❌ %s: %v\n", result.Repository.FullName, result.Error)
				}
				failed++
			} else if !result.Success {
				if verbose {
					fmt.Printf("⏭️  %s: does not exist locally\n", result.Repository.FullName)
				}
				skipped++
			} else {
				if result.Updated {
					if verbose {
						fmt.Printf("🔄 %s: updated (%.2fs)\n", result.Repository.FullName, result.Duration.Seconds())
					}
					updated++
				} else {
					if verbose {
						fmt.Printf("✅ %s: up to date (%.2fs)\n", result.Repository.FullName, result.Duration.Seconds())
					}
				}
				successful++
			}
		}

		// Print summary
		fmt.Printf("\nUpdate Summary:\n")
		fmt.Printf("  Successful: %d\n", successful)
		fmt.Printf("  Updated:    %d\n", updated)
		fmt.Printf("  Skipped:    %d (not found locally)\n", skipped)
		fmt.Printf("  Failed:     %d\n", failed)

		if failed > 0 {
			return fmt.Errorf("some repositories failed to update")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&updateUseSSH, "ssh", false, "Use SSH URLs for updating instead of HTTPS")
}
