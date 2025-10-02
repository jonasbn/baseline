package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jonasbn/baseline/internal/types"
)

// GitOps provides Git operations for baseline
type GitOps struct {
	verbose bool
}

// NewGitOps creates a new GitOps instance
func NewGitOps(verbose bool) *GitOps {
	return &GitOps{
		verbose: verbose,
	}
}

// CloneRepository clones a repository as bare to the specified directory
func (g *GitOps) CloneRepository(repo types.Repository, targetDir string) types.CloneResult {
	start := time.Now()
	result := types.CloneResult{
		Repository: repo,
		Success:    false,
		Duration:   0,
	}

	// Create the repository directory path
	repoPath := filepath.Join(targetDir, repo.Owner, repo.Name)

	// Check if repository already exists
	if _, err := os.Stat(repoPath); err == nil {
		result.Error = fmt.Errorf("repository already exists at %s", repoPath)
		result.Duration = time.Since(start)
		return result
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(repoPath), 0755); err != nil {
		result.Error = fmt.Errorf("failed to create parent directory: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", repo.CloneURL, repoPath)
	if g.verbose {
		fmt.Printf("Cloning %s to %s\n", repo.FullName, repoPath)
	}

	if err := cmd.Run(); err != nil {
		result.Error = fmt.Errorf("failed to clone repository %s: %w", repo.FullName, err)
		result.Duration = time.Since(start)
		return result
	}

	// Set permissions to read-only
	if err := g.setReadOnlyPermissions(repoPath); err != nil {
		result.Error = fmt.Errorf("failed to set read-only permissions for %s: %w", repoPath, err)
		result.Duration = time.Since(start)
		return result
	}

	result.Success = true
	result.Duration = time.Since(start)
	return result
}

// UpdateRepository updates an existing bare repository
func (g *GitOps) UpdateRepository(repo types.Repository, targetDir string) types.UpdateResult {
	start := time.Now()
	result := types.UpdateResult{
		Repository: repo,
		Success:    false,
		Updated:    false,
		Duration:   0,
	}

	repoPath := filepath.Join(targetDir, repo.Owner, repo.Name)

	// Check if repository exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		result.Error = fmt.Errorf("repository does not exist at %s", repoPath)
		result.Duration = time.Since(start)
		return result
	}

	// Temporarily set write permissions
	if err := g.setWritePermissions(repoPath); err != nil {
		result.Error = fmt.Errorf("failed to set write permissions for %s: %w", repoPath, err)
		result.Duration = time.Since(start)
		return result
	}

	// Get current HEAD before update
	oldHead, err := g.getCurrentHead(repoPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to get current HEAD for %s: %w", repoPath, err)
		result.Duration = time.Since(start)
		return result
	}

	// Fetch updates
	cmd := exec.Command("git", "-C", repoPath, "fetch", "origin")
	if g.verbose {
		fmt.Printf("Updating %s at %s\n", repo.FullName, repoPath)
	}

	if err := cmd.Run(); err != nil {
		// Restore read-only permissions even if update fails
		g.setReadOnlyPermissions(repoPath)
		result.Error = fmt.Errorf("failed to update repository %s: %w", repo.FullName, err)
		result.Duration = time.Since(start)
		return result
	}

	// Get new HEAD after update
	newHead, err := g.getCurrentHead(repoPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to get new HEAD for %s: %w", repoPath, err)
		result.Duration = time.Since(start)
		return result
	}

	// Check if repository was actually updated
	result.Updated = oldHead != newHead

	// Restore read-only permissions
	if err := g.setReadOnlyPermissions(repoPath); err != nil {
		result.Error = fmt.Errorf("failed to restore read-only permissions for %s: %w", repoPath, err)
		result.Duration = time.Since(start)
		return result
	}

	result.Success = true
	result.Duration = time.Since(start)
	return result
}

// RepositoryExists checks if a repository already exists in the target directory
func (g *GitOps) RepositoryExists(repo types.Repository, targetDir string) bool {
	repoPath := filepath.Join(targetDir, repo.Owner, repo.Name)
	_, err := os.Stat(repoPath)
	return err == nil
}

// setReadOnlyPermissions sets read-only permissions recursively
func (g *GitOps) setReadOnlyPermissions(path string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Set read-only permissions (owner: read+execute, group: read, others: read)
		var mode os.FileMode
		if info.IsDir() {
			mode = 0555 // read + execute for directories
		} else {
			mode = 0444 // read-only for files
		}

		return os.Chmod(filePath, mode)
	})
}

// setWritePermissions sets write permissions recursively
func (g *GitOps) setWritePermissions(path string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Set write permissions (owner: read+write+execute, group: read, others: read)
		var mode os.FileMode
		if info.IsDir() {
			mode = 0755 // read + write + execute for directories
		} else {
			mode = 0644 // read + write for files
		}

		return os.Chmod(filePath, mode)
	})
}

// getCurrentHead gets the current HEAD commit hash
func (g *GitOps) getCurrentHead(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
