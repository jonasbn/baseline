package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jonasbn/baseline/internal/types"
)

// GitHubClient implements the RepositorySource interface for GitHub
type GitHubClient struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

// GitHubRepository represents a GitHub repository response
type GitHubRepository struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	CloneURL    string    `json:"clone_url"`
	SSHURL      string    `json:"ssh_url"`
	HTTPSURL    string    `json:"html_url"`
	Description *string   `json:"description"`
	Private     bool      `json:"private"`
	UpdatedAt   time.Time `json:"updated_at"`
	Language    *string   `json:"language"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.github.com",
	}
}

// GetName returns the source name
func (g *GitHubClient) GetName() string {
	return "github"
}

// GetRepositories fetches all repositories for the given organization
func (g *GitHubClient) GetRepositories(ctx context.Context, organization string) ([]types.Repository, error) {
	var allRepos []types.Repository
	page := 1
	perPage := 100

	for {
		repos, hasMore, err := g.getRepositoriesPage(ctx, organization, page, perPage)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repositories page %d: %w", page, err)
		}

		allRepos = append(allRepos, repos...)

		if !hasMore {
			break
		}
		page++
	}

	return allRepos, nil
}

func (g *GitHubClient) getRepositoriesPage(ctx context.Context, organization string, page, perPage int) ([]types.Repository, bool, error) {
	// Try organization endpoint first, fall back to user endpoint if 404
	url := fmt.Sprintf("%s/orgs/%s/repos?page=%d&per_page=%d&sort=updated", g.baseURL, organization, page, perPage)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create request: %w", err)
	}

	if g.token != "" {
		req.Header.Set("Authorization", "token "+g.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Try user endpoint if organization endpoint returns 404
		userURL := fmt.Sprintf("%s/users/%s/repos?page=%d&per_page=%d&sort=updated", g.baseURL, organization, page, perPage)

		req, err = http.NewRequestWithContext(ctx, "GET", userURL, nil)
		if err != nil {
			return nil, false, fmt.Errorf("failed to create user request: %w", err)
		}

		if g.token != "" {
			req.Header.Set("Authorization", "token "+g.token)
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err = g.httpClient.Do(req)
		if err != nil {
			return nil, false, fmt.Errorf("failed to make user request: %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var githubRepos []GitHubRepository
	if err := json.NewDecoder(resp.Body).Decode(&githubRepos); err != nil {
		return nil, false, fmt.Errorf("failed to decode response: %w", err)
	}

	repos := make([]types.Repository, len(githubRepos))
	for i, repo := range githubRepos {
		repos[i] = g.convertToRepository(repo)
	}

	// Check if there are more pages
	hasMore := len(githubRepos) == perPage

	return repos, hasMore, nil
}

func (g *GitHubClient) convertToRepository(repo GitHubRepository) types.Repository {
	description := ""
	if repo.Description != nil {
		description = *repo.Description
	}

	language := ""
	if repo.Language != nil {
		language = *repo.Language
	}

	return types.Repository{
		Name:        repo.Name,
		FullName:    repo.FullName,
		CloneURL:    repo.CloneURL,
		SSHURL:      repo.SSHURL,
		HTTPSURL:    repo.HTTPSURL,
		Description: description,
		Private:     repo.Private,
		UpdatedAt:   repo.UpdatedAt,
		Language:    language,
		Owner:       repo.Owner.Login,
	}
}
