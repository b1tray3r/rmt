package redmine

import (
	"fmt"
	"time"

	"github.com/b1tray3r/rmt/internal/redmine/models"
)

// WeeklyFollowUpFilter creates a filter for issues marked for this week in the follow-up custom field.
// This function assumes the follow-up custom field contains date values and filters for the current week.
func WeeklyFollowUpFilter(followUpFieldID int) models.IssueFilter {
	// Calculate the start and end of the current week (Monday to Sunday)
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}

	// Find Monday of this week
	monday := now.AddDate(0, 0, -(weekday - 1))
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())

	// Find Sunday of this week
	sunday := monday.AddDate(0, 0, 6)
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 0, sunday.Location())

	// Format dates for Redmine (YYYY-MM-DD format)
	startDate := monday.Format("2006-01-02")
	endDate := sunday.Format("2006-01-02")

	// Create date range query for the custom field
	dateRange := fmt.Sprintf(">=%s&cf_%d=<=%s", startDate, followUpFieldID, endDate)

	return models.IssueFilter{
		Limit: 100, // Default limit
		CustomFields: map[int]string{
			followUpFieldID: dateRange,
		},
	}
}

// ThisWeekFollowUpFilter creates a simple filter for issues where the follow-up custom field
// contains "this week" as text (case-insensitive approach).
func ThisWeekFollowUpFilter(followUpFieldID int) models.IssueFilter {
	return models.IssueFilter{
		Limit: 100,
		CustomFields: map[int]string{
			followUpFieldID: "*this week*", // Wildcard search for "this week"
		},
	}
}

// CustomFieldContainsFilter creates a filter for issues where a custom field contains specific text.
func CustomFieldContainsFilter(fieldID int, searchText string) models.IssueFilter {
	return models.IssueFilter{
		Limit: 100,
		CustomFields: map[int]string{
			fieldID: fmt.Sprintf("*%s*", searchText),
		},
	}
}
