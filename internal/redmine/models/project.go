package models

import "time"

// Project represents a Redmine project.
// Project contains detailed information about a specific project in Redmine.
type Project struct {
	ID                  int                 `json:"id"`                    // ID is the unique identifier of the project
	Name                string              `json:"name"`                  // Name is the display name of the project
	Identifier          string              `json:"identifier"`            // Identifier is the unique string identifier of the project
	Description         string              `json:"description"`           // Description contains the project description
	Status              int                 `json:"status"`                // Status indicates the project status (1=active, 5=closed, etc.)
	IsPublic            bool                `json:"is_public"`             // IsPublic indicates whether the project is publicly visible
	CreatedOn           time.Time           `json:"created_on"`            // CreatedOn is the timestamp when the project was created
	UpdatedOn           time.Time           `json:"updated_on"`            // UpdatedOn is the timestamp when the project was last updated
	TimeEntryActivities []TimeEntryActivity `json:"time_entry_activities"` // TimeEntryActivities contains available time tracking activities
}

// TimeEntryActivity represents an activity that can be used for time tracking.
// TimeEntryActivity defines the types of work activities available for logging time.
type TimeEntryActivity struct {
	ID        int    `json:"id"`         // ID is the unique identifier of the activity
	Name      string `json:"name"`       // Name is the display name of the activity
	IsDefault bool   `json:"is_default"` // IsDefault indicates whether this is the default activity
	Active    bool   `json:"active"`     // Active indicates whether this activity is currently available for use
}
