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
}

type RedmineTimeEntryCreator interface {
	CreateTimeEntry(params models.CreateTimeEntryParams) (*models.TimeEntry, error)
}

type RedmineBaseURLGetter interface {
	GetBaseURL() string
}

type RedmineAPI interface {
	RedmineSearcher
	RedmineProjectGetter
	RedmineIssueGetter
	RedmineBaseURLGetter
	RedmineTimeEntryCreator
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

func (c *RestClient) GetBaseURL() string {
	return c.baseURL
}

func (c *RestClient) Login(ctx context.Context) error {
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

func (c *RestClient) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	fullURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

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

func (c *RestClient) GetIssue(id int) (*models.Issue, error) {
	ctx := context.Background()

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/issues/%d.json", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetIssue request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetIssue request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetIssue failed with status: %d", resp.StatusCode)
	}

	var issueResponse struct {
		Issue models.Issue `json:"issue"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&issueResponse); err != nil {
		return nil, fmt.Errorf("failed to decode GetIssue response: %w", err)
	}

	return &issueResponse.Issue, nil
}

func (c *RestClient) GetProject(id int) (*models.Project, error) {
	ctx := context.Background()

	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/projects/%d.json?include=time_entry_activities", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetProject request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GetProject request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetProject failed with status: %d", resp.StatusCode)
	}

	var projectResponse struct {
		Project models.Project `json:"project"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&projectResponse); err != nil {
		return nil, fmt.Errorf("failed to decode GetProject response: %w", err)
	}

	return &projectResponse.Project, nil
}

func (c *RestClient) CreateTimeEntry(params models.CreateTimeEntryParams) (*models.TimeEntry, error) {
	ctx := context.Background()

	payload := struct {
		TimeEntry models.CreateTimeEntryParams `json:"time_entry"`
	}{
		TimeEntry: params,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CreateTimeEntry payload: %w", err)
	}

	req, err := c.newRequest(ctx, "POST", "/time_entries.json", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateTimeEntry request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute CreateTimeEntry request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CreateTimeEntry failed with status: %d", resp.StatusCode)
	}

	var timeEntryResponse struct {
		TimeEntry models.TimeEntry `json:"time_entry"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&timeEntryResponse); err != nil {
		return nil, fmt.Errorf("failed to decode CreateTimeEntry response: %w", err)
	}

	return &timeEntryResponse.TimeEntry, nil
}
