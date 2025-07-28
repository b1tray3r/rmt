// Package redmine provides an user friendly interface which can be used to read or write data to a redmine instance.
package redmine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/b1tray3r/rmt/internal/http/querystring"
	rm "github.com/mattn/go-redmine"
)

type SearchParams struct {
	Offset     int    `query:"offset"`
	Limit      int    `query:"limit"`
	Query      string `query:"q"`
	Scope      string `query:"scope"`
	AllWords   bool   `query:"all_words"`
	TitlesOnly bool   `query:"titles_only"`
	Issues     bool   `query:"issues"`
	OpenIssues bool   `query:"open_issues"`
}

type SearchResult struct {
	ID          int
	Title       string
	Type        string
	URL         string
	Description string
	DateTime    time.Time
}

type SearchResults struct {
	Results    []SearchResult `json:"results"`
	TotalCount int            `json:"total_count"`
	Offset     int            `json:"offset"`
	Limit      int            `json:"limit"`
}

type Redmine interface {
	Search(SearchParams) ([]SearchResult, error)
}

type RedmineAPI struct {
	endpoint string
	key      string
	client   rm.Client
}

func NewRedmineAPI(endpoint, apikey string) *RedmineAPI {
	return &RedmineAPI{
		endpoint: endpoint,
		key:      apikey,
		client:   *rm.NewClient(endpoint, apikey),
	}
}

func (r *RedmineAPI) GetClient() *rm.Client {
	return &r.client
}

func (r *RedmineAPI) Search(params SearchParams) (*SearchResults, error) {
	// If limit is 0 or less, user wants all results
	fetchAll := params.Limit <= 0
	pageLimit := params.Limit
	if fetchAll {
		pageLimit = 100 // default page size when fetching all
	}

	allResults := &SearchResults{
		Results: []SearchResult{},
	}

	offset := 0
	for {
		params.Limit = pageLimit
		params.Offset = offset

		q, err := querystring.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal query params: %w", err)
		}

		url := r.endpoint + "/search.json?" + string(q)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Add("X-Redmine-API-Key", r.key)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed with status: %s", resp.Status)
		}

		var page SearchResults
		if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		if offset == 0 {
			allResults.TotalCount = page.TotalCount
			allResults.Offset = 0
		}

		allResults.Results = append(allResults.Results, page.Results...)

		// If not fetching all, return after the first page
		if !fetchAll {
			break
		}

		offset += page.Limit
		if offset >= page.TotalCount {
			break
		}
	}

	return allResults, nil
}
