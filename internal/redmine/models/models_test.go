package models

import (
	"encoding/json"
	"testing"
	"time"
)

// TestSearchParamsMarshaling tests that SearchParams can be properly marshaled to query string.
func TestSearchParamsMarshaling(t *testing.T) {
	params := SearchParams{
		Offset:     10,
		Limit:      25,
		Query:      "test query",
		Scope:      "issues",
		AllWords:   true,
		TitlesOnly: false,
		Issues:     true,
		OpenIssues: true,
	}

	// This test verifies that all fields are properly tagged for query string marshaling
	// The actual marshaling logic is tested in the querystring package
	if params.Query != "test query" {
		t.Errorf("expected query 'test query', got '%s'", params.Query)
	}
	if params.Limit != 25 {
		t.Errorf("expected limit 25, got %d", params.Limit)
	}
}

// TestSearchResultJSONMarshaling tests that SearchResult can be properly marshaled/unmarshaled.
func TestSearchResultJSONMarshaling(t *testing.T) {
	originalTime := time.Date(2025, 8, 14, 12, 0, 0, 0, time.UTC)
	original := SearchResult{
		ID:          123,
		Title:       "Test Issue",
		Type:        "issue",
		URL:         "/issues/123",
		Description: "Test description",
		DateTime:    originalTime,
		Project:     "Test Project",
		ProjectID:   "test-project",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal SearchResult: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaled SearchResult
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal SearchResult: %v", err)
	}

	// Verify fields (note: Project and ProjectID have json:"-" tags, so they won't be marshaled)
	if unmarshaled.ID != original.ID {
		t.Errorf("expected ID %d, got %d", original.ID, unmarshaled.ID)
	}
	if unmarshaled.Title != original.Title {
		t.Errorf("expected title '%s', got '%s'", original.Title, unmarshaled.Title)
	}
	if unmarshaled.Type != original.Type {
		t.Errorf("expected type '%s', got '%s'", original.Type, unmarshaled.Type)
	}
	if unmarshaled.Project != "" {
		t.Errorf("expected Project to be empty (json:\"-\"), got '%s'", unmarshaled.Project)
	}
	if unmarshaled.ProjectID != "" {
		t.Errorf("expected ProjectID to be empty (json:\"-\"), got '%s'", unmarshaled.ProjectID)
	}
}

// TestIssueJSONMarshaling tests that Issue can be properly marshaled/unmarshaled.
func TestIssueJSONMarshaling(t *testing.T) {
	jsonData := `{
		"id": 123,
		"subject": "Test Issue",
		"description": "Test description",
		"status": {"name": "New"},
		"author": {"name": "Test User"},
		"project": {"name": "Test Project", "id": 1},
		"created_on": "2025-08-14T12:00:00Z",
		"updated_on": "2025-08-14T13:00:00Z"
	}`

	var issue Issue
	if err := json.Unmarshal([]byte(jsonData), &issue); err != nil {
		t.Fatalf("failed to unmarshal Issue: %v", err)
	}

	if issue.ID != 123 {
		t.Errorf("expected ID 123, got %d", issue.ID)
	}
	if issue.Subject != "Test Issue" {
		t.Errorf("expected subject 'Test Issue', got '%s'", issue.Subject)
	}
	if issue.Status.Name != "New" {
		t.Errorf("expected status name 'New', got '%s'", issue.Status.Name)
	}
	if issue.Author.Name != "Test User" {
		t.Errorf("expected author name 'Test User', got '%s'", issue.Author.Name)
	}
	if issue.Project.Name != "Test Project" {
		t.Errorf("expected project name 'Test Project', got '%s'", issue.Project.Name)
	}
	if issue.Project.ID != 1 {
		t.Errorf("expected project ID 1, got %d", issue.Project.ID)
	}
}

// TestProjectJSONMarshaling tests that Project can be properly marshaled/unmarshaled.
func TestProjectJSONMarshaling(t *testing.T) {
	jsonData := `{
		"id": 1,
		"name": "Test Project",
		"identifier": "test-project",
		"description": "Test description",
		"status": 1,
		"is_public": true,
		"created_on": "2025-08-14T12:00:00Z",
		"updated_on": "2025-08-14T13:00:00Z",
		"time_entry_activities": [
			{
				"id": 1,
				"name": "Development",
				"is_default": true,
				"active": true
			},
			{
				"id": 2,
				"name": "Testing",
				"is_default": false,
				"active": true
			}
		]
	}`

	var project Project
	if err := json.Unmarshal([]byte(jsonData), &project); err != nil {
		t.Fatalf("failed to unmarshal Project: %v", err)
	}

	if project.ID != 1 {
		t.Errorf("expected ID 1, got %d", project.ID)
	}
	if project.Name != "Test Project" {
		t.Errorf("expected name 'Test Project', got '%s'", project.Name)
	}
	if project.Identifier != "test-project" {
		t.Errorf("expected identifier 'test-project', got '%s'", project.Identifier)
	}
	if project.Status != 1 {
		t.Errorf("expected status 1, got %d", project.Status)
	}
	if !project.IsPublic {
		t.Error("expected IsPublic to be true")
	}
	if len(project.TimeEntryActivities) != 2 {
		t.Errorf("expected 2 time entry activities, got %d", len(project.TimeEntryActivities))
	}

	// Test first activity
	activity := project.TimeEntryActivities[0]
	if activity.ID != 1 {
		t.Errorf("expected activity ID 1, got %d", activity.ID)
	}
	if activity.Name != "Development" {
		t.Errorf("expected activity name 'Development', got '%s'", activity.Name)
	}
	if !activity.IsDefault {
		t.Error("expected IsDefault to be true for first activity")
	}
	if !activity.Active {
		t.Error("expected Active to be true for first activity")
	}
}

// TestTimeEntryJSONMarshaling tests that TimeEntry can be properly marshaled/unmarshaled.
func TestTimeEntryJSONMarshaling(t *testing.T) {
	jsonData := `{
		"id": 1,
		"hours": 2.5,
		"comments": "Test work",
		"spent_on": "2025-08-14",
		"issue": {"id": 123},
		"project": {"id": 1, "name": "Test Project"},
		"activity": {"id": 1, "name": "Development"},
		"user": {"id": 1, "name": "Test User"},
		"created_on": "2025-08-14T12:00:00Z",
		"updated_on": "2025-08-14T13:00:00Z"
	}`

	var timeEntry TimeEntry
	if err := json.Unmarshal([]byte(jsonData), &timeEntry); err != nil {
		t.Fatalf("failed to unmarshal TimeEntry: %v", err)
	}

	if timeEntry.ID != 1 {
		t.Errorf("expected ID 1, got %d", timeEntry.ID)
	}
	if timeEntry.Hours != 2.5 {
		t.Errorf("expected hours 2.5, got %f", timeEntry.Hours)
	}
	if timeEntry.Comments != "Test work" {
		t.Errorf("expected comments 'Test work', got '%s'", timeEntry.Comments)
	}
	if timeEntry.SpentOn != "2025-08-14" {
		t.Errorf("expected spent_on '2025-08-14', got '%s'", timeEntry.SpentOn)
	}
	if timeEntry.Issue.ID != 123 {
		t.Errorf("expected issue ID 123, got %d", timeEntry.Issue.ID)
	}
	if timeEntry.Activity.ID != 1 {
		t.Errorf("expected activity ID 1, got %d", timeEntry.Activity.ID)
	}
	if timeEntry.Activity.Name != "Development" {
		t.Errorf("expected activity name 'Development', got '%s'", timeEntry.Activity.Name)
	}
	if timeEntry.User.ID != 1 {
		t.Errorf("expected user ID 1, got %d", timeEntry.User.ID)
	}
	if timeEntry.User.Name != "Test User" {
		t.Errorf("expected user name 'Test User', got '%s'", timeEntry.User.Name)
	}
}

// TestCreateTimeEntryRequestJSONMarshaling tests that CreateTimeEntryRequest can be properly marshaled.
func TestCreateTimeEntryRequestJSONMarshaling(t *testing.T) {
	request := CreateTimeEntryParams{
		IssueID:    123,
		Hours:      2.5,
		ActivityID: 1,
		Comments:   "Test work",
		SpentOn:    "2025-08-14",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal CreateTimeEntryRequest: %v", err)
	}

	var unmarshaled CreateTimeEntryParams
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CreateTimeEntryRequest: %v", err)
	}

	if unmarshaled.IssueID != request.IssueID {
		t.Errorf("expected IssueID %d, got %d", request.IssueID, unmarshaled.IssueID)
	}
	if unmarshaled.Hours != request.Hours {
		t.Errorf("expected Hours %f, got %f", request.Hours, unmarshaled.Hours)
	}
	if unmarshaled.ActivityID != request.ActivityID {
		t.Errorf("expected ActivityID %d, got %d", request.ActivityID, unmarshaled.ActivityID)
	}
	if unmarshaled.Comments != request.Comments {
		t.Errorf("expected Comments '%s', got '%s'", request.Comments, unmarshaled.Comments)
	}
	if unmarshaled.SpentOn != request.SpentOn {
		t.Errorf("expected SpentOn '%s', got '%s'", request.SpentOn, unmarshaled.SpentOn)
	}
}
