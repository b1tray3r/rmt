package models

// IssueFilter defines the parameters for filtering Redmine issues with advanced criteria.
// IssueFilter supports custom field filtering and other advanced Redmine query parameters.
type IssueFilter struct {
	Offset       int            `query:"offset"`         // Offset specifies the number of results to skip
	Limit        int            `query:"limit"`          // Limit specifies the maximum number of results to return
	StatusID     []int          `query:"status_id"`      // StatusID filters by status IDs (e.g., 1=New, 2=In Progress)
	TrackerID    []int          `query:"tracker_id"`     // TrackerID filters by tracker IDs (e.g., 1=Bug, 2=Feature)
	ProjectID    []int          `query:"project_id"`     // ProjectID filters by project IDs
	AssignedTo   string         `query:"assigned_to_id"` // AssignedTo filters by assignee ("me" for current user)
	AuthorID     []int          `query:"author_id"`      // AuthorID filters by author IDs
	Subject      string         `query:"subject"`        // Subject searches in issue subjects
	CustomFields map[int]string `query:"-"`              // CustomFields maps custom field ID to value
}

// IssueResults represents the response from a Redmine issues request.
// IssueResults contains paginated issue results and metadata.
type IssueResults struct {
	Issues     []Issue `json:"issues"`      // Issues contains the array of issues
	TotalCount int     `json:"total_count"` // TotalCount is the total number of issues available
	Offset     int     `json:"offset"`      // Offset is the number of results skipped
	Limit      int     `json:"limit"`       // Limit is the maximum number of results returned
}
