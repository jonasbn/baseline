package bitbucket

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jonasbn/baseline/internal/types"
)

// BitbucketClient implements the RepositorySource interface for Bitbucket
type BitbucketClient struct {
	username    string
	appPassword string
	httpClient  *http.Client
	baseURL     string
}

// BitbucketRepository represents a Bitbucket repository response
type BitbucketRepository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Links    struct {
		Clone []struct {
			Name string `json:"name"`
			Href string `json:"href"`
		} `json:"clone"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	UpdatedOn   time.Time `json:"updated_on"`
	Language    string    `json:"language"`
	Owner       struct {
		Username string `json:"username"`
	} `json:"owner"`
}

// BitbucketResponse represents the paginated response from Bitbucket
type BitbucketResponse struct {
	Values []BitbucketRepository `json:"values"`
	Next   string                `json:"next"`
}

// NewBitbucketClient creates a new Bitbucket client
func NewBitbucketClient(username, appPassword string) *BitbucketClient {
	return &BitbucketClient{
		username:    username,
		appPassword: appPassword,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.bitbucket.org/2.0",
	}
}

// GetName returns the source name
func (b *BitbucketClient) GetName() string {
	return "bitbucket"
}

// GetRepositories fetches all repositories for the given organization
func (b *BitbucketClient) GetRepositories(ctx context.Context, organization string) ([]types.Repository, error) {
	var allRepos []types.Repository
	url := fmt.Sprintf("%s/repositories/%s?pagelen=100&sort=-updated_on", b.baseURL, organization)

	for url != "" {
		repos, nextURL, err := b.getRepositoriesPage(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repositories: %w", err)
		}

		allRepos = append(allRepos, repos...)
		url = nextURL
	}

	return allRepos, nil
}

func (b *BitbucketClient) getRepositoriesPage(ctx context.Context, url string) ([]types.Repository, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add basic auth if credentials are provided
	if b.username != "" && b.appPassword != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(b.username + ":" + b.appPassword))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("Bitbucket API returned status %d", resp.StatusCode)
	}

	var response BitbucketResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	repos := make([]types.Repository, len(response.Values))
	for i, repo := range response.Values {
		repos[i] = b.convertToRepository(repo)
	}

	return repos, response.Next, nil
}

func (b *BitbucketClient) convertToRepository(repo BitbucketRepository) types.Repository {
	// Extract clone URLs
	var cloneURL, sshURL string
	for _, link := range repo.Links.Clone {
		if link.Name == "https" {
			cloneURL = link.Href
		} else if link.Name == "ssh" {
			sshURL = link.Href
		}
	}

	return types.Repository{
		Name:        repo.Name,
		FullName:    repo.FullName,
		CloneURL:    cloneURL,
		SSHURL:      sshURL,
		HTTPSURL:    repo.Links.HTML.Href,
		Description: repo.Description,
		Private:     repo.IsPrivate,
		UpdatedAt:   repo.UpdatedOn,
		Language:    repo.Language,
		Owner:       repo.Owner.Username,
	}
}
