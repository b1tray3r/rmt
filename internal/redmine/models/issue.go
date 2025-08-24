package models

import "time"

// Issue represents a Redmine issue.
// Issue contains the detailed information about a specific issue in Redmine.
type Issue struct {
	ID          int    `json:"id"`          // ID is the unique identifier of the issue
	Subject     string `json:"subject"`     // Subject is the title or summary of the issue
	Description string `json:"description"` // Description contains the detailed description of the issue
	Status      struct {
		Name string `json:"name"` // Name is the status name (e.g., "New", "In Progress", "Closed")
	} `json:"status"` // Status represents the current status of the issue
	Author struct {
		Name string `json:"name"` // Name is the full name of the issue author
	} `json:"author"` // Author represents the user who created the issue
	Project struct {
		Name string `json:"name"` // Name is the project name
		ID   int    `json:"id"`   // ID is the project identifier
	} `json:"project"` // Project represents the project this issue belongs to
	CreatedOn time.Time `json:"created_on"` // CreatedOn is the timestamp when the issue was created
	UpdatedOn time.Time `json:"updated_on"` // UpdatedOn is the timestamp when the issue was last updated
}
