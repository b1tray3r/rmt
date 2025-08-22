package redmine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/b1tray3r/rmt/internal/http/querystring"
	"github.com/b1tray3r/rmt/internal/redmine/models"
)

type RedmineSearcher interface {
	Search(params models.SearchParams) (*models.SearchResults, error)
}

type RedmineIssueGetter interface {
	GetIssue(issueID int) (*models.Issue, error)
}

type RedmineProjectGetter interface {
	GetProject(projectID int) (*models.Project, error)
	GetProjects() ([]models.Project, error)
}

type RedmineTimeEntryCreator interface {
	CreateTimeEntry(params models.CreateTimeEntryParams) (*models.TimeEntry, error)
}

type RedmineAPI interface {
	RedmineSearcher
	/*RedmineIssueGetter
	RedmineProjectGetter
	RedmineTimeEntryCreator*/
}

type RestClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewRestClient(baseURL, apiKey string) *RestClient {
	return &RestClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *RestClient) Login(ctx context.Context) error {
	// Test authentication by trying to get current user info
	req, err := c.newRequest(ctx, "GET", "/users/current.json", nil)
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication failed: invalid API key")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	return nil
}

// newRequest creates a new HTTP request with common headers and authentication.
// newRequest sets up the request with proper API key authentication and content type headers.
func (c *RestClient) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	fullURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	// Set authentication header
	req.Header.Set("X-Redmine-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// Search performs a search operation on the Redmine instance.
// Search queries the Redmine search API and returns paginated results based on the provided parameters.
func (c *RestClient) Search(params models.SearchParams) (*models.SearchResults, error) {
	ctx := context.Background()

	// Marshal search parameters using http/querystring for cleaner query construction.
	queryParams, err := querystring.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search parameters: %w", err)
	}

	// Construct the search endpoint URL
	searchPath := "/search.json"
	if len(queryParams) > 0 {
		searchPath += "?" + string(queryParams)
	}

	req, err := c.newRequest(ctx, "GET", searchPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read search response: %w", err)
	}

	var results models.SearchResults
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return &results, nil
}

// GetIssue retrieves a specific issue by ID from the Redmine instance.
// GetIssue is a stub implementation that will be implemented in future iterations.
func (c *RestClient) GetIssue(issueID int) (*models.Issue, error) {
	// TODO: Implement GetIssue method
	return nil, fmt.Errorf("GetIssue not implemented yet")
}

// GetProject retrieves a specific project by ID from the Redmine instance.
// GetProject is a stub implementation that will be implemented in future iterations.
func (c *RestClient) GetProject(projectID interface{}) (*models.Project, error) {
	// TODO: Implement GetProject method
	return nil, fmt.Errorf("GetProject not implemented yet")
}

// GetProjects retrieves all accessible projects from the Redmine instance.
// GetProjects is a stub implementation that will be implemented in future iterations.
func (c *RestClient) GetProjects() ([]models.Project, error) {
	// TODO: Implement GetProjects method
	return nil, fmt.Errorf("GetProjects not implemented yet")
}

// CreateTimeEntry creates a new time entry in the Redmine instance.
// CreateTimeEntry is a stub implementation that will be implemented in future iterations.
func (c *RestClient) CreateTimeEntry(params models.CreateTimeEntryParams) (*models.TimeEntry, error) {
	// TODO: Implement CreateTimeEntry method
	return nil, fmt.Errorf("CreateTimeEntry not implemented yet")
}
