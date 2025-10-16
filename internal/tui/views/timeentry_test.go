package views

import (
	"testing"

	"github.com/b1tray3r/rmt/internal/redmine/models"
	"github.com/b1tray3r/rmt/internal/tui/domain"
)

// Mock implementations for testing

// mockIssueRepository implements domain.IssueRepository for testing
type mockIssueRepository struct {
	activities map[int]string
	err        error
}

// GetBaseURL returns a mock base URL
func (m *mockIssueRepository) GetBaseURL() string {
	return "http://example.com"
}

// GetIssue returns a mock issue
func (m *mockIssueRepository) GetIssue(id int) (*domain.Issue, error) {
	return createTestIssue(), nil
}

// GetProjectActivities returns mock activities or an error
func (m *mockIssueRepository) GetProjectActivities(projectID int, activityPatterns []string) (map[int]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.activities != nil {
		return m.activities, nil
	}
	// Default activities
	return map[int]string{
		1: "Development",
		2: "Testing",
		3: "Documentation",
	}, nil
}

// Search is a mock implementation - not used in TimeEntryView tests
func (m *mockIssueRepository) Search(query string) ([]*domain.Issue, error) {
	return nil, nil
}

// SearchWithFilter searches issues using filters (mock implementation)
func (m *mockIssueRepository) SearchWithFilter(query string) ([]*domain.Issue, error) {
	return nil, nil
}

// CreateTimeEntry creates a mock time entry or returns an error
func (m *mockIssueRepository) CreateTimeEntry(params models.CreateTimeEntryParams) (*models.TimeEntry, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &models.TimeEntry{
		ID:       1,
		Hours:    params.Hours,
		Comments: params.Comments,
		SpentOn:  params.SpentOn,
		Issue: struct {
			ID int `json:"id"`
		}{ID: params.IssueID},
		Activity: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{ID: params.ActivityID, Name: "Development"},
	}, nil
}

// Helper functions for creating test data

// createTestIssue creates a test issue for testing
func createTestIssue() *domain.Issue {
	// We need to create a project properly. Since Project has unexported fields,
	// we'll create it the same way the production code does
	project := &domain.Project{} // This creates a zero-value Project
	return domain.NewIssue(123, "http://example.com/issues/123", "test.user", "Test Issue", "Test Description", project)
}

// Test NewTimeEntryView

// TestNewTimeEntryView_Success verifies that NewTimeEntryView creates a valid instance with proper initialization
func TestNewTimeEntryView_Success(t *testing.T) {
	issue := createTestIssue()
	repo := &mockIssueRepository{}

	view, err := NewTimeEntryView(80, 24, []string{"Development"}, issue, repo, repo)

	if err != nil {
		t.Fatalf("NewTimeEntryView returned error: %v", err)
	}
	if view == nil {
		t.Fatal("NewTimeEntryView returned nil view")
	}
	if view.width != 80 {
		t.Errorf("width = %d, want 80", view.width)
	}
	if view.height != 24 {
		t.Errorf("height = %d, want 24", view.height)
	}
	if view.issue != issue {
		t.Error("issue not set correctly")
	}
	if len(view.activities) != 3 {
		t.Errorf("activities length = %d, want 3", len(view.activities))
	}
	if view.state != StateEditing {
		t.Errorf("initial state = %v, want StateEditing", view.state)
	}
	if view.focusIndex != DateIndex {
		t.Errorf("initial focus = %v, want DateIndex", view.focusIndex)
	}
}
