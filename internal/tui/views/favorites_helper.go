package views

import (
	"fmt"
	"time"

	"github.com/b1tray3r/rmt/internal/tui/domain"
)

// CreateWeeklyFollowUpFavorite creates a favorite for issues marked for this week in the follow-up custom field.
// The followUpFieldID should be the ID of the custom field used for follow-up dates in your Redmine instance.
func CreateWeeklyFollowUpFavorite(followUpFieldID int) *domain.Favorite {
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

	// Create query that filters for issues where the follow-up date is within this week
	query := fmt.Sprintf("cf_%d=>=%s&cf_%d=<=%s", followUpFieldID, startDate, followUpFieldID, endDate)

	return domain.NewFavorite(1, "Diese Woche (Follow-up)", query)
}

// CreateFollowUpContainsFavorite creates a favorite for issues where follow-up field contains specific text.
func CreateFollowUpContainsFavorite(followUpFieldID int, searchText string) *domain.Favorite {
	query := fmt.Sprintf("cf_%d=*%s*", followUpFieldID, searchText)
	return domain.NewFavorite(2, fmt.Sprintf("Follow-up: %s", searchText), query)
}
