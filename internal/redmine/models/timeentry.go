package models

import "time"

// TimeEntry represents a time entry in Redmine.
// TimeEntry contains information about logged time for an issue or project.
type TimeEntry struct {
	ID       int     `json:"id"`       // ID is the unique identifier of the time entry
	Hours    float64 `json:"hours"`    // Hours is the amount of time logged
	Comments string  `json:"comments"` // Comments contain additional notes about the time entry
	SpentOn  string  `json:"spent_on"` // SpentOn is the date when the work was performed (YYYY-MM-DD format)
	Issue    struct {
		ID int `json:"id"` // ID is the issue identifier this time entry belongs to
	} `json:"issue"` // Issue represents the issue this time entry is associated with
	Activity struct {
		ID   int    `json:"id"`   // ID is the activity identifier
		Name string `json:"name"` // Name is the activity name
	} `json:"activity"` // Activity represents the type of work performed
	User struct {
		ID   int    `json:"id"`   // ID is the user identifier
		Name string `json:"name"` // Name is the user's full name
	} `json:"user"` // User represents the user who logged this time
	CreatedOn time.Time `json:"created_on"` // CreatedOn is the timestamp when the time entry was created
	UpdatedOn time.Time `json:"updated_on"` // UpdatedOn is the timestamp when the time entry was last updated
}

// CreateTimeEntryParams represents the request payload for creating a time entry.
// CreateTimeEntryParams contains all the required and optional fields for creating a new time entry.
type CreateTimeEntryParams struct {
	IssueID    int     `json:"issue_id"`    // IssueID is the ID of the issue to log time against
	Hours      float64 `json:"hours"`       // Hours is the amount of time to log
	ActivityID int     `json:"activity_id"` // ActivityID is the ID of the time tracking activity
	Comments   string  `json:"comments"`    // Comments contain additional notes about the work performed
	SpentOn    string  `json:"spent_on"`    // SpentOn is the date when the work was performed (YYYY-MM-DD format)
}
