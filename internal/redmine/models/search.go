package models

import "time"

// SearchParams defines the parameters for searching Redmine issues.
// SearchParams contains all the supported query parameters for the Redmine search API.
type SearchParams struct {
	Offset     int    `query:"offset"`      // Offset specifies the number of results to skip
	Limit      int    `query:"limit"`       // Limit specifies the maximum number of results to return
	Query      string `query:"q"`           // Query is the search term
	Scope      string `query:"scope"`       // Scope limits the search to specific areas (e.g., "issues", "projects")
	AllWords   bool   `query:"all_words"`   // AllWords indicates whether to match all words in the query
	TitlesOnly bool   `query:"titles_only"` // TitlesOnly limits search to titles only
	Issues     bool   `query:"issues"`      // Issues indicates whether to include issues in search results
	OpenIssues bool   `query:"open_issues"` // OpenIssues indicates whether to include only open issues
}

// SearchResult represents a single search result from Redmine.
// SearchResult contains the essential information about a found item in Redmine search.
type SearchResult struct {
	ID          int       `json:"id"`          // ID is the unique identifier of the search result
	Title       string    `json:"title"`       // Title is the title or subject of the result
	Type        string    `json:"type"`        // Type specifies the type of the result (issue, project, etc.)
	URL         string    `json:"url"`         // URL is the relative path to the result
	Description string    `json:"description"` // Description contains a brief description of the result
	DateTime    time.Time `json:"datetime"`    // DateTime indicates when the result was created or last modified
	Project     string    `json:"-"`           // Project name (not included in JSON response)
	ProjectID   string    `json:"-"`           // ProjectID is the project identifier (not included in JSON response)
}

// SearchResults represents the response from a Redmine search request.
// SearchResults contains paginated search results and metadata about the search.
type SearchResults struct {
	Results    []SearchResult `json:"results"`     // Results contains the array of search results
	TotalCount int            `json:"total_count"` // TotalCount is the total number of results available
	Offset     int            `json:"offset"`      // Offset is the number of results skipped
	Limit      int            `json:"limit"`       // Limit is the maximum number of results returned
}
