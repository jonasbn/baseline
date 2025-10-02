package bitbucket

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/jonasbn/baseline/internal/types"
)

// BitbucketClient implements the RepositorySource interface for Bitbucket
type BitbucketClient struct {
	username   string
	apiToken   string
	debug      bool
	httpClient *http.Client
	baseURL    string
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
// username should be your Bitbucket username or email
// apiToken should be a repository, project, or workspace access token
func NewBitbucketClient(username, apiToken string, debug bool) *BitbucketClient {
	return &BitbucketClient{
		username: username,
		apiToken: apiToken,
		debug:    debug,
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
	url := fmt.Sprintf("%s/repositories/%s?role=member&pagelen=100", b.baseURL, organization)

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
	if b.username != "" && b.apiToken != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(b.username + ":" + b.apiToken))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	// Debug logging: print the entire request
	if b.debug {
		fmt.Println("\n=== DEBUG: Bitbucket HTTP Request ===")
		fmt.Printf("Method: %s\n", req.Method)
		fmt.Printf("URL: %s\n", req.URL.String())
		fmt.Printf("Headers:\n")
		for name, values := range req.Header {
			for _, value := range values {
				if name == "Authorization" {
					// Mask the token for security, but show the format
					fmt.Printf("  %s: %s\n", name, "Basic [BASE64_ENCODED_USERNAME:TOKEN]")
				} else {
					fmt.Printf("  %s: %s\n", name, value)
				}
			}
		}

		// Show equivalent curl command
		fmt.Printf("\nEquivalent curl command:\n")
		fmt.Printf("curl -v -u '%s:[TOKEN]' \\\n", b.username)
		fmt.Printf("  -H 'Accept: application/json' \\\n")
		fmt.Printf("  '%s'\n", req.URL.String())
		fmt.Println("========================================\n")

		// Also dump the raw request if needed
		if reqDump, err := httputil.DumpRequestOut(req, false); err == nil {
			fmt.Printf("Raw HTTP Request:\n%s\n", string(reqDump))
		}
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Debug logging: print response details
	if b.debug {
		fmt.Println("=== DEBUG: Bitbucket HTTP Response ===")
		fmt.Printf("Status: %s (%d)\n", resp.Status, resp.StatusCode)
		fmt.Printf("Response Headers:\n")
		for name, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("  %s: %s\n", name, value)
			}
		}
		fmt.Println("======================================\n")
	}

	if resp.StatusCode != http.StatusOK {
		if b.debug {
			log.Printf("DEBUG: Request failed with status %d", resp.StatusCode)
		}
		return nil, "", fmt.Errorf("bitbucket API returned status %d", resp.StatusCode)
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
		switch link.Name {
		case "https":
			cloneURL = link.Href
		case "ssh":
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
