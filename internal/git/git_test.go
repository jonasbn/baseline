package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonasbn/baseline/internal/types"
)

func TestNewGitOps(t *testing.T) {
	gitOps := NewGitOps(true)
	if gitOps == nil {
		t.Error("NewGitOps should not return nil")
	}

	if !gitOps.verbose {
		t.Error("GitOps verbose flag should be true")
	}
}

func TestRepositoryExists(t *testing.T) {
	gitOps := NewGitOps(false)

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	repo := types.Repository{
		Name:  "test-repo",
		Owner: "test-owner",
	}

	// Repository should not exist initially
	if gitOps.RepositoryExists(repo, tempDir) {
		t.Error("Repository should not exist initially")
	}

	// Create the repository directory
	repoPath := filepath.Join(tempDir, repo.Owner, repo.Name+".git")
	err := os.MkdirAll(repoPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test repository directory: %v", err)
	}

	// Repository should exist now
	if !gitOps.RepositoryExists(repo, tempDir) {
		t.Error("Repository should exist after creation")
	}
}

func TestSetReadOnlyPermissions(t *testing.T) {
	gitOps := NewGitOps(false)

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// Create a test file
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Set read-only permissions
	err = gitOps.setReadOnlyPermissions(tempDir)
	if err != nil {
		t.Fatalf("Failed to set read-only permissions: %v", err)
	}

	// Check that the file is read-only
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	mode := info.Mode()
	if mode&0200 != 0 { // Check if write permission is set for owner
		t.Error("File should be read-only")
	}

	// Cleanup: restore write permissions so temp dir can be cleaned up
	defer gitOps.setWritePermissions(tempDir)
}
