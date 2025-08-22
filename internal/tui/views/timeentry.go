package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/b1tray3r/rmt/internal/tui/domain"
	"github.com/b1tray3r/rmt/internal/tui/messages"
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Common styles using Tokyo Night theme
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Primary).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1).
			Align(lipgloss.Center)

	fieldLabelStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Warning).
			Bold(true).
			Padding(0, 1)

	fieldValueStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Foreground).
			Padding(0, 1)

	focusedStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Success).
			Bold(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Muted).
			Italic(true).
			Padding(0, 1)

	emptyMessageStyle = lipgloss.NewStyle().
				Foreground(themes.TokyoNight.Muted).
				Italic(true).
				Padding(1, 2)

	loadingStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Info).
			Bold(true).
			Padding(1, 3).
			Align(lipgloss.Center)
)

// Activity represents a Redmine activity for time entries
type Activity struct {
	ID   int
	Name string
}

// TimeEntry represents a time entry to be submitted
type TimeEntry struct {
	IssueID     int
	ActivityID  int
	Hours       float64
	Description string
	Date        time.Time
}

// DatePicker handles date selection with calendar view
type DatePicker struct {
	selectedDate time.Time
	viewDate     time.Time // The month/year currently being viewed
	focused      bool
}

// NewDatePicker creates a new date picker with today's date selected
func NewDatePicker() *DatePicker {
	now := time.Now()
	return &DatePicker{
		selectedDate: now,
		viewDate:     now,
		focused:      false,
	}
}

// Update handles input for the date picker
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

// ensureDateInView adjusts viewDate if selectedDate is outside the current month
func (dp *DatePicker) ensureDateInView() {
	if dp.selectedDate.Year() != dp.viewDate.Year() || dp.selectedDate.Month() != dp.viewDate.Month() {
		dp.viewDate = time.Date(dp.selectedDate.Year(), dp.selectedDate.Month(), 1, 0, 0, 0, 0, dp.selectedDate.Location())
	}
}

// Focus enables input handling for the date picker
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
	// Calendar header with month/year
	header := fieldLabelStyle.Render(dp.viewDate.Format("January 2006"))

	// Day headers (starting with Monday)
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

	// Get last day of month
	lastDay := firstDay.AddDate(0, 1, -1).Day()

	// Build calendar grid
	var calendarRows []string
	currentRow := ""

	// Empty cells for days before the first day of month
	for i := 0; i < weekday; i++ {
		currentRow += fieldValueStyle.Width(4).Render("")
	}

	// Days of the month
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

	// Help text
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

// HoursSelector handles hours selection with predefined values
type HoursSelector struct {
	options       []float64
	selectedIndex int
	focused       bool
}

// NewHoursSelector creates a new hours selector
func NewHoursSelector() *HoursSelector {
	options := []float64{
		0.25, 0.50, 0.75, 1.00, 1.25, 1.50, 1.75, 2.00,
		2.25, 2.50, 2.75, 3.00, 3.25, 3.50, 3.75, 4.00,
		4.25, 4.50, 4.75, 5.00, 5.25, 5.50, 5.75, 6.00,
		6.25, 6.50, 6.75, 7.00, 7.25, 7.50, 7.75, 8.00,
	}

	return &HoursSelector{
		options:       options,
		selectedIndex: 3, // Default to 1 hour
		focused:       false,
	}
}

// Update handles input for the hours selector
func (hs *HoursSelector) Update(msg tea.KeyMsg) {
	if !hs.focused {
		return
	}

	switch msg.String() {
	case "left":
		if hs.selectedIndex > 0 {
			hs.selectedIndex--
		}
	case "right":
		if hs.selectedIndex < len(hs.options)-1 {
			hs.selectedIndex++
		}
	case "home":
		hs.selectedIndex = 0
	case "end":
		hs.selectedIndex = len(hs.options) - 1
	}
}

// Focus enables input handling
func (hs *HoursSelector) Focus() {
	hs.focused = true
}

// Blur disables input handling
func (hs *HoursSelector) Blur() {
	hs.focused = false
}

// SelectedHours returns the currently selected hours
func (hs *HoursSelector) SelectedHours() float64 {
	return hs.options[hs.selectedIndex]
}

// Render returns the horizontal hours selector view
func (hs *HoursSelector) Render() string {
	var items []string

	// Show 7 items at a time for better visibility
	start := hs.selectedIndex - 3
	end := hs.selectedIndex + 4

	if start < 0 {
		start = 0
		end = 7
	}
	if end > len(hs.options) {
		end = len(hs.options)
		start = end - 7
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end && i < len(hs.options); i++ {
		hourStr := fmt.Sprintf("%.2f", hs.options[i])

		var style lipgloss.Style
		if i == hs.selectedIndex {
			if hs.focused {
				style = focusedStyle.Padding(0, 1)
			} else {
				style = fieldValueStyle.Bold(true).Foreground(themes.TokyoNight.Info).Padding(0, 1)
			}
		} else {
			style = fieldValueStyle.Padding(0, 1)
		}

		items = append(items, style.Render(hourStr))
	}

	// Add navigation indicators
	leftArrow := ""
	rightArrow := ""
	if start > 0 {
		leftArrow = fieldValueStyle.Foreground(themes.TokyoNight.Info).Render("< ")
	} else {
		leftArrow = "  "
	}
	if end < len(hs.options) {
		rightArrow = fieldValueStyle.Foreground(themes.TokyoNight.Info).Render(" >")
	} else {
		rightArrow = "  "
	}

	content := leftArrow + strings.Join(items, " ") + rightArrow

	// Help text when focused
	helpText := ""
	if hs.focused {
		helpText = "\n" + helpStyle.Render("Left/Right: select hours | Home/End: first/last")
	}

	return content + helpText
}

type TimeEntryIndex int

const (
	DateIndex TimeEntryIndex = iota
	HoursIndex
	DescriptionIndex
	ActivityIndex
	SubmitIndex
	TotalFields // This should always be last
)

// TimeEntryState represents the current state of the time entry view
type TimeEntryState int

const (
	StateEditing TimeEntryState = iota
	StateSubmitting
	StateCompleted
	StateError
)

// TimeEntryView handles the time entry interface
type TimeEntryView struct {
	activities       []Activity
	selectedActivity *Activity
	datePicker       *DatePicker
	hoursSelector    *HoursSelector
	descInput        textinput.Model
	focusIndex       TimeEntryIndex
	state            TimeEntryState
	errorMessage     string
}

// NewTimeEntryView creates a new time entry view
func NewTimeEntryView() *TimeEntryView {
	descInput := textinput.New()
	descInput.Placeholder = "Data log: describe your digital work..."
	descInput.CharLimit = 255
	descInput.Width = 48 // Set a reasonable width

	v := &TimeEntryView{
		datePicker:    NewDatePicker(),
		hoursSelector: NewHoursSelector(),
		descInput:     descInput,
		focusIndex:    DateIndex,
		activities:    []Activity{},
		state:         StateEditing,
	}

	// Set initial focus
	v.updateFocus()

	return v
}

// Init implements tea.Model interface
func (v *TimeEntryView) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles input
func (v *TimeEntryView) Update(msg tea.Msg) tea.Cmd {
	// Handle submission result messages
	switch msg := msg.(type) {
	case TimeEntrySubmissionSuccess:
		v.state = StateCompleted
		return nil
	case TimeEntrySubmissionError:
		v.state = StateError
		v.errorMessage = msg.Error.Error()
		return nil
	}

	// Handle input in error or completed state
	if v.state == StateError || v.state == StateCompleted {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":
				if v.state == StateError {
					// Reset to editing state to try again
					v.state = StateEditing
					v.errorMessage = ""
					return nil
				}
				// For completed state, return to issue view
				return func() tea.Msg { return messages.ReturnToIssueMsg{} }
			default:
				if v.state == StateCompleted {
					// Any key press after completion should return to issue view
					return func() tea.Msg { return messages.ReturnToIssueMsg{} }
				}
			}
		}
		return nil
	}

	// Don't handle input during submission
	if v.state == StateSubmitting {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Clear validation errors when user starts interacting
		if v.errorMessage != "" && v.state == StateEditing {
			// Don't clear on enter (submit) to allow validation error to show
			if msg.String() != "enter" {
				v.errorMessage = ""
			}
		}

		switch msg.String() {
		case "tab":
			v.nextField()
			return nil
		case "shift+tab":
			v.prevField()
			return nil
		case "enter":
			if v.focusIndex == SubmitIndex {
				return v.submitWithCommand()
			}
		case "up", "down":
			if v.focusIndex == ActivityIndex && len(v.activities) > 0 {
				v.handleActivitySelection(msg.String())
				return nil
			}
		}

		// Handle date picker input when focused
		if v.focusIndex == DateIndex {
			v.datePicker.Update(msg)
			return nil
		}

		// Handle hours selector input when focused
		if v.focusIndex == HoursIndex {
			v.hoursSelector.Update(msg)
			return nil
		}

		// Handle description input when focused
		if v.focusIndex == DescriptionIndex {
			var cmd tea.Cmd
			v.descInput, cmd = v.descInput.Update(msg)
			return cmd
		}
	}

	return nil
}

// submitWithCommand handles submission and returns the command
func (v *TimeEntryView) submitWithCommand() tea.Cmd {
	// Validate required fields (this handles activity requirement properly)
	if !v.HasValidEntry() {
		// Stay in editing state but show validation error inline
		v.errorMessage = "Please fill in all required fields"
		return nil
	}

	// Clear any previous error messages
	v.errorMessage = ""

	// Set loading state
	v.state = StateSubmitting

	// Create time entry data
	activityID := 0
	if v.selectedActivity != nil {
		activityID = v.selectedActivity.ID
	}

	entry := TimeEntry{
		IssueID:     0, // Will be set by parent when issue is available
		ActivityID:  activityID,
		Hours:       v.hoursSelector.SelectedHours(),
		Description: strings.TrimSpace(v.descInput.Value()),
		Date:        v.datePicker.SelectedDate(),
	}

	// Submit through mock service
	service := &MockTimeEntryService{}
	return service.SubmitTimeEntry(entry)
}

// Render returns the view string
func (v *TimeEntryView) Render(issue *domain.Issue, activities []Activity, width, height int) string {
	title := titleStyle.Width(width).Render("TIME ENTRY")

	if issue == nil {
		return title + "\n\n" + emptyMessageStyle.Render("No issue selected")
	}

	// Show different content based on state
	switch v.state {
	case StateSubmitting:
		return v.renderSubmittingState(title, issue, width, height)
	case StateCompleted:
		return v.renderCompletedState(title, issue, width, height)
	case StateError:
		return v.renderErrorState(title, issue, width, height)
	default:
		return v.renderEditingState(title, issue, activities, width, height)
	}
}

// renderEditingState renders the normal editing form
func (v *TimeEntryView) renderEditingState(title string, issue *domain.Issue, activities []Activity, width, height int) string {
	issueInfo := fieldLabelStyle.Render(fmt.Sprintf("Issue: #%d %s", issue.ID(), issue.Title()))

	var activitySection string
	if len(activities) > 0 {
		activitySection = v.renderActivitySelection(activities)
	} else {
		activitySection = loadingStyle.Render("Loading activities...")
	}

	dateSection := v.renderDatePicker()
	hoursSection := v.renderHoursSelector()
	descSection := v.renderInput("Description:", v.descInput, v.focusIndex == DescriptionIndex)
	submitSection := v.renderSubmitButton()

	// Show validation error if any (for inline feedback)
	var errorSection string
	if v.errorMessage != "" && v.state == StateEditing {
		errorSection = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Error).
			Bold(true).
			Padding(0, 1).
			Render("⚠ " + v.errorMessage)
	}

	helpText := helpStyle.Render("Tab: next field | Shift+Tab: previous | Up/Down: select activity/navigate date | Left/Right: navigate day/hours | Shift+Left/Right: change month | Home/End: first/last hours | Enter: submit | Esc: back")

	sections := []string{
		title,
		"",
		issueInfo,
		"",
		dateSection,
		"",
		hoursSection,
		"",
		descSection,
		"",
		activitySection,
		"",
		submitSection,
	}

	if errorSection != "" {
		sections = append(sections, "", errorSection)
	}

	sections = append(sections, "", helpText)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderSubmittingState renders the loading state during submission
func (v *TimeEntryView) renderSubmittingState(title string, issue *domain.Issue, width, height int) string {
	issueInfo := fieldLabelStyle.Render(fmt.Sprintf("Issue: #%d %s", issue.ID(), issue.Title()))

	loadingMessage := loadingStyle.Render("Submitting time entry...")

	// Create a simple spinner animation
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerChar := spinner[int(time.Now().Unix())%len(spinner)]
	spinnerStyle := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Info).
		Bold(true)

	loadingWithSpinner := lipgloss.JoinHorizontal(lipgloss.Left,
		spinnerStyle.Render(spinnerChar+" "),
		loadingMessage,
	)

	helpText := helpStyle.Render("Please wait while your time entry is being submitted...")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		issueInfo,
		"",
		"",
		"",
		loadingWithSpinner,
		"",
		"",
		helpText,
	)
}

// renderCompletedState renders the success state after submission
func (v *TimeEntryView) renderCompletedState(title string, issue *domain.Issue, width, height int) string {
	issueInfo := fieldLabelStyle.Render(fmt.Sprintf("Issue: #%d %s", issue.ID(), issue.Title()))

	successMessage := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Success).
		Bold(true).
		Padding(1, 3).
		Align(lipgloss.Center).
		Render("✓ Time entry submitted successfully!")

	helpText := helpStyle.Render("Press any key to return to issue view...")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		issueInfo,
		"",
		"",
		"",
		successMessage,
		"",
		"",
		helpText,
	)
}

// renderErrorState renders the error state when submission fails
func (v *TimeEntryView) renderErrorState(title string, issue *domain.Issue, width, height int) string {
	issueInfo := fieldLabelStyle.Render(fmt.Sprintf("Issue: #%d %s", issue.ID(), issue.Title()))

	errorMessage := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Error).
		Bold(true).
		Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(themes.TokyoNight.Error).
		Render("✗ Error: " + v.errorMessage)

	helpText := helpStyle.Render("Press Enter to try again or Esc to return to issue view")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		issueInfo,
		"",
		"",
		"",
		errorMessage,
		"",
		"",
		helpText,
	)
}

// SetActivities sets the available activities
func (v *TimeEntryView) SetActivities(activities []Activity) {
	v.activities = activities
	if len(activities) > 0 && v.selectedActivity == nil {
		v.selectedActivity = &activities[0]
	}
}

// GetState returns the current state of the time entry view
func (v *TimeEntryView) GetState() TimeEntryState {
	return v.state
}

// GetErrorMessage returns the current error message
func (v *TimeEntryView) GetErrorMessage() string {
	return v.errorMessage
}

// Reset resets the time entry view to initial state
func (v *TimeEntryView) Reset() {
	v.state = StateEditing
	v.errorMessage = ""
	v.descInput.SetValue("")
	v.focusIndex = DateIndex
	v.updateFocus()
}

func (v *TimeEntryView) renderActivitySelection(activities []Activity) string {
	if len(activities) == 0 {
		return fieldLabelStyle.Render("Activity: ") + "No activities available"
	}

	var content strings.Builder
	content.WriteString(fieldLabelStyle.Render("Activity:") + "\n")

	for _, activity := range activities {
		prefix := "  "
		style := fieldValueStyle
		if v.selectedActivity != nil && activity.ID == v.selectedActivity.ID {
			prefix = "> "
			if v.focusIndex == ActivityIndex {
				style = focusedStyle
			}
		}
		content.WriteString(prefix + style.Render(activity.Name) + "\n")
	}

	return content.String()
}

func (v *TimeEntryView) renderInput(label string, input textinput.Model, focused bool) string {
	style := fieldValueStyle
	if focused {
		style = focusedStyle
		input.Focus()
	} else {
		input.Blur()
	}

	return fieldLabelStyle.Render(label) + "\n" +
		style.Width(50).Padding(0, 1).Render(input.View())
}

func (v *TimeEntryView) renderDatePicker() string {
	focused := v.focusIndex == DateIndex
	if focused {
		v.datePicker.Focus()
	} else {
		v.datePicker.Blur()
	}

	return fieldLabelStyle.Render("Date:") + "\n" + v.datePicker.Render()
}

func (v *TimeEntryView) renderHoursSelector() string {
	focused := v.focusIndex == HoursIndex
	if focused {
		v.hoursSelector.Focus()
	} else {
		v.hoursSelector.Blur()
	}

	return fieldLabelStyle.Render("Hours:") + "\n" + v.hoursSelector.Render()
}

func (v *TimeEntryView) renderSubmitButton() string {
	buttonText := "Submit Time Entry"
	focused := v.focusIndex == SubmitIndex

	var buttonStyle lipgloss.Style
	if focused {
		buttonStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Background).
			Background(themes.TokyoNight.Info).
			Padding(0, 2).
			Bold(true)
	} else {
		buttonStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Foreground).
			Background(themes.TokyoNight.BackgroundAlt).
			Padding(0, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(themes.TokyoNight.Border)
	}

	return buttonStyle.Render(buttonText)
}

func (v *TimeEntryView) nextField() {
	v.focusIndex = (v.focusIndex + 1) % TotalFields
	v.updateFocus()
}

func (v *TimeEntryView) prevField() {
	v.focusIndex = (v.focusIndex - 1 + TotalFields) % TotalFields
	v.updateFocus()
}

// updateFocus manages the focus state of all components based on the current focus index
func (v *TimeEntryView) updateFocus() {
	// Reset all focus states
	v.datePicker.Blur()
	v.hoursSelector.Blur()
	v.descInput.Blur()

	// Set focus on the currently selected field
	switch v.focusIndex {
	case DateIndex:
		v.datePicker.Focus()
	case HoursIndex:
		v.hoursSelector.Focus()
	case DescriptionIndex:
		v.descInput.Focus()
	case ActivityIndex:
		// Activity selection doesn't need explicit focus, handled by visual styling
	case SubmitIndex:
		// Submit button doesn't need explicit focus, handled by visual styling
	}
}

func (v *TimeEntryView) handleActivitySelection(direction string) {
	if len(v.activities) == 0 {
		return
	}

	currentIndex := 0
	if v.selectedActivity != nil {
		for i, activity := range v.activities {
			if activity.ID == v.selectedActivity.ID {
				currentIndex = i
				break
			}
		}
	}

	switch direction {
	case "up":
		if currentIndex > 0 {
			v.selectedActivity = &v.activities[currentIndex-1]
		}
	case "down":
		if currentIndex < len(v.activities)-1 {
			v.selectedActivity = &v.activities[currentIndex+1]
		}
	}
}

// TimeEntryService interface for submitting time entries
type TimeEntryService interface {
	SubmitTimeEntry(entry TimeEntry) tea.Cmd
}

// MockTimeEntryService provides a fake implementation for testing
type MockTimeEntryService struct{}

// SubmitTimeEntry simulates submitting a time entry with a delay
func (s *MockTimeEntryService) SubmitTimeEntry(entry TimeEntry) tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		// Simulate success most of the time, occasional failure for testing
		if time.Now().Unix()%10 == 0 { // 10% failure rate
			return TimeEntrySubmissionError{Error: fmt.Errorf("network error: failed to submit time entry")}
		}
		return TimeEntrySubmissionSuccess{Entry: entry}
	})
}

// TimeEntrySubmissionSuccess indicates successful submission
type TimeEntrySubmissionSuccess struct {
	Entry TimeEntry
}

// TimeEntrySubmissionError indicates submission failure
type TimeEntrySubmissionError struct {
	Error error
}

// HasValidEntry checks if all required fields are filled for submission
func (v *TimeEntryView) HasValidEntry() bool {
	// Check basic required fields
	hasValidHours := v.hoursSelector.SelectedHours() > 0
	hasValidDescription := strings.TrimSpace(v.descInput.Value()) != ""

	// Activity is only required if activities are available
	hasValidActivity := true
	if len(v.activities) > 0 {
		hasValidActivity = v.selectedActivity != nil
	}

	return hasValidActivity && hasValidHours && hasValidDescription
}

// TimeLogView maintains backward compatibility with the original structure
type TimeLogView struct {
	width, height int
	SearchInput   *textinput.Model
	timeEntryView *TimeEntryView
	issue         *domain.Issue
	activities    []Activity
}

// NewTimeLogView creates a new time log view for backward compatibility
func NewTimeLogView(width int, issue *domain.Issue) *TimeLogView {
	return &TimeLogView{
		width:         width,
		height:        0,
		timeEntryView: NewTimeEntryView(),
		issue:         issue,
		activities:    []Activity{}, // Initialize with empty activities, will be loaded later
	}
}

// SetSize sets the dimensions of the time log view
func (v *TimeLogView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// SetIssue sets the current issue for the time entry
func (v *TimeLogView) SetIssue(issue *domain.Issue) {
	v.issue = issue
}

// SetActivities sets the available activities for time entries
func (v *TimeLogView) SetActivities(activities []Activity) {
	v.activities = activities
	if v.timeEntryView != nil {
		v.timeEntryView.SetActivities(activities)
	}
}

// GetTimeEntryState returns the current time entry state
func (v *TimeLogView) GetTimeEntryState() TimeEntryState {
	if v.timeEntryView != nil {
		return v.timeEntryView.GetState()
	}
	return StateEditing
}

// GetTimeEntryError returns any time entry error message
func (v *TimeLogView) GetTimeEntryError() string {
	if v.timeEntryView != nil {
		return v.timeEntryView.GetErrorMessage()
	}
	return ""
}

// ResetTimeEntry resets the time entry form
func (v *TimeLogView) ResetTimeEntry() {
	if v.timeEntryView != nil {
		v.timeEntryView.Reset()
	}
}

// HasValidTimeEntry checks if the time entry form has valid data
func (v *TimeLogView) HasValidTimeEntry() bool {
	if v.timeEntryView != nil {
		return v.timeEntryView.HasValidEntry()
	}
	return false
}

// Init initializes the TimeLogView and returns the blinking cursor command
func (v *TimeLogView) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles input for the time log view
func (v *TimeLogView) Update(msg tea.Msg) tea.Cmd {
	// Handle window resize messages
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		v.SetSize(msg.Width, msg.Height)
	}

	// Forward all other messages to the time entry view
	if v.timeEntryView != nil {
		return v.timeEntryView.Update(msg)
	}
	return nil
}

// Render returns the time log view string
func (v *TimeLogView) Render() string {
	if v.timeEntryView == nil {
		return emptyMessageStyle.Render("Time entry view not initialized")
	}

	// Use the TimeEntryView's render method with current issue and activities
	return v.timeEntryView.Render(v.issue, v.activities, v.width, v.height)
}
