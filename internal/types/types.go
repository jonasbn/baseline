package types

import (
	"context"
	"time"
)

// Repository represents a Git repository with its metadata
type Repository struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	CloneURL    string    `json:"clone_url"`
	SSHURL      string    `json:"ssh_url"`
	HTTPSURL    string    `json:"https_url"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
	UpdatedAt   time.Time `json:"updated_at"`
	Language    string    `json:"language"`
	Owner       string    `json:"owner"`
}

// RepositorySource defines the interface for fetching repositories from different sources
type RepositorySource interface {
	// GetRepositories fetches all repositories for the given organization
	GetRepositories(ctx context.Context, organization string) ([]Repository, error)
	// GetName returns the name of the source (e.g., "github", "bitbucket")
	GetName() string
}

// CloneOptions contains configuration for cloning operations
type CloneOptions struct {
	Directory    string
	Organization string
	Threads      int
	Verbose      bool
	Source       string
}

// Credentials holds authentication information for API access
type Credentials struct {
	GitHubToken          string
	BitbucketUsername    string
	BitbucketAppPassword string
}

// SourceType represents the supported repository sources
type SourceType string

const (
	SourceGitHub    SourceType = "github"
	SourceBitbucket SourceType = "bitbucket"
)

// CloneResult represents the result of a clone operation
type CloneResult struct {
	Repository Repository
	Success    bool
	Error      error
	Duration   time.Duration
}

// UpdateResult represents the result of an update operation
type UpdateResult struct {
	Repository Repository
	Success    bool
	Error      error
	Duration   time.Duration
	Updated    bool // true if the repository was actually updated
}
