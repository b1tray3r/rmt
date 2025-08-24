package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/b1tray3r/rmt/internal/tui/themes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DatePicker struct {
	selectedDate time.Time
	viewDate     time.Time
	focused      bool
}

func NewDatePicker() *DatePicker {
	now := time.Now()
	return &DatePicker{
		selectedDate: now,
		viewDate:     now,
		focused:      false,
	}
}

func (dp *DatePicker) Update(msg tea.KeyMsg) {
	if !dp.focused {
		return
	}

	switch msg.String() {
	case "left":
		dp.selectedDate = dp.selectedDate.AddDate(0, 0, -1)
		dp.ensureDateInView()
	case "right":
		dp.selectedDate = dp.selectedDate.AddDate(0, 0, 1)
		dp.ensureDateInView()
	case "up":
		dp.selectedDate = dp.selectedDate.AddDate(0, 0, -7)
		dp.ensureDateInView()
	case "down":
		dp.selectedDate = dp.selectedDate.AddDate(0, 0, 7)
		dp.ensureDateInView()
	case "shift+left":
		dp.viewDate = dp.viewDate.AddDate(0, -1, 0)
	case "shift+right":
		dp.viewDate = dp.viewDate.AddDate(0, 1, 0)
	case "home":
		dp.selectedDate = time.Now()
		dp.viewDate = dp.selectedDate
	}
}

func (dp *DatePicker) ensureDateInView() {
	if dp.selectedDate.Year() != dp.viewDate.Year() || dp.selectedDate.Month() != dp.viewDate.Month() {
		dp.viewDate = time.Date(dp.selectedDate.Year(), dp.selectedDate.Month(), 1, 0, 0, 0, 0, dp.selectedDate.Location())
	}
}

func (dp *DatePicker) Focus() {
	dp.focused = true
}

// Blur disables input handling for the date picker
func (dp *DatePicker) Blur() {
	dp.focused = false
}

// SelectedDate returns the currently selected date
func (dp *DatePicker) SelectedDate() time.Time {
	return dp.selectedDate
}

// Render returns the calendar view
func (dp *DatePicker) Render() string {
	header := fieldLabelStyle.Render(dp.viewDate.Format("January 2006"))
	dayHeaders := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	headerRow := ""
	for _, day := range dayHeaders {
		headerRow += fieldValueStyle.Width(4).Align(lipgloss.Center).Render(day)
	}

	// Calculate first day of month and adjust for Monday start
	firstDay := time.Date(dp.viewDate.Year(), dp.viewDate.Month(), 1, 0, 0, 0, 0, dp.viewDate.Location())
	weekday := int(firstDay.Weekday())
	if weekday == 0 { // Sunday = 0, but we want Monday = 0
		weekday = 6
	} else {
		weekday--
	}

	lastDay := firstDay.AddDate(0, 1, -1).Day()

	// Build calendar grid
	var calendarRows []string
	currentRow := ""

	// Empty cells for days before the first day of month
	for i := 0; i < weekday; i++ {
		currentRow += fieldValueStyle.Width(4).Render("")
	}

	for day := 1; day <= lastDay; day++ {
		dayDate := time.Date(dp.viewDate.Year(), dp.viewDate.Month(), day, 0, 0, 0, 0, dp.viewDate.Location())
		dayStr := fmt.Sprintf("%2d", day)

		var style lipgloss.Style
		if dp.focused && dayDate.Year() == dp.selectedDate.Year() &&
			dayDate.Month() == dp.selectedDate.Month() &&
			dayDate.Day() == dp.selectedDate.Day() {
			// Selected date
			style = focusedStyle.Width(4).Align(lipgloss.Center)
		} else if dayDate.Year() == time.Now().Year() &&
			dayDate.Month() == time.Now().Month() &&
			dayDate.Day() == time.Now().Day() {
			// Today
			style = fieldValueStyle.Width(4).Align(lipgloss.Center).Bold(true).Foreground(themes.TokyoNight.Info)
		} else {
			// Regular day
			style = fieldValueStyle.Width(4).Align(lipgloss.Center)
		}

		currentRow += style.Render(dayStr)

		// Start new row after Sunday (7 days)
		if (weekday+day)%7 == 0 {
			calendarRows = append(calendarRows, currentRow)
			currentRow = ""
		}
	}

	// Add remaining row if not complete
	if currentRow != "" {
		calendarRows = append(calendarRows, currentRow)
	}

	helpText := ""
	if dp.focused {
		helpText = helpStyle.Render("Left/Right: navigate days | Up/Down: navigate weeks | Shift+Left/Right: change month | Home: today")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		headerRow,
		strings.Join(calendarRows, "\n"),
		helpText,
	)
}
