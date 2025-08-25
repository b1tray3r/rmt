// Package views provides user interface components for the TUI application.
// This package contains various view components including time entry forms,
// search interfaces, and other interactive elements for the Redmine management tool.
package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/b1tray3r/rmt/internal/redmine/models"
	"github.com/b1tray3r/rmt/internal/tui/domain"
	"github.com/b1tray3r/rmt/internal/tui/messages"
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
			Foreground(themes.TokyoNight.Foreground).
			Italic(true).
			Padding(0, 1)

	emptyMessageStyle = lipgloss.NewStyle().
				Foreground(themes.TokyoNight.Foreground).
				Italic(true).
				Padding(1, 2)

	loadingStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Info).
			Bold(true).
			Padding(1, 3).
			Align(lipgloss.Center)
)

// Activity represents a Redmine activity that can be associated with time entries.
type Activity struct {
	ID   int
	Name string
}

// TimeEntry represents a time entry record with all necessary fields for logging work.
type TimeEntry struct {
	IssueID     int
	ActivityID  int
	Hours       float64
	Description string
	Date        time.Time
}

// TimeEntryIndex represents the field indices for navigation in the time entry form.
type TimeEntryIndex int

const (
	DateIndex TimeEntryIndex = iota
	HoursIndex
	DescriptionIndex
	ActivityIndex
	SubmitIndex
	TotalFields
)

// TimeEntryState represents the current state of the time entry view.
type TimeEntryState int

const (
	StateEditing TimeEntryState = iota
	StateSubmitting
	StateCompleted
	StateError
)

// TimeEntryView handles the time entry interface for logging work against issues.
type TimeEntryView struct {
	width  int
	height int

	timeLogService domain.TimeEntryCreator

	issue *domain.Issue

	activities       []Activity
	selectedActivity *Activity
	datePicker       *DatePicker
	hoursSelector    *HoursSelector
	descInput        textinput.Model
	focusIndex       TimeEntryIndex
	state            TimeEntryState
	errorMessage     string

	SearchInput *textinput.Model
}

// NewTimeEntryView creates a new time entry view instance with the specified dimensions and context.
// NewTimeEntryView initializes all necessary components including date picker, hours selector, and activity list.
// It returns an error if the issue or issueRepository parameters are nil, or if activities cannot be loaded.
func NewTimeEntryView(width, height int, issue *domain.Issue, issueRepository domain.IssueRepository, timeLogService domain.TimeEntryCreator) (*TimeEntryView, error) {
	if issue == nil {
		return nil, fmt.Errorf("issue cannot be nil")
	}

	if issueRepository == nil {
		return nil, fmt.Errorf("issueRepository cannot be nil")
	}

	activities, err := issueRepository.GetProjectActivities(issue.Project().ID())
	if err != nil {
		return nil, fmt.Errorf("failed to get project activities: %w", err)
	}

	var activitiesList []Activity
	for id, name := range activities {
		activity := Activity{
			ID:   id,
			Name: name,
		}
		activitiesList = append(activitiesList, activity)
	}

	descInput := textinput.New()
	descInput.Placeholder = "Data log: describe your digital work..."
	descInput.PlaceholderStyle = helpStyle
	descInput.CharLimit = 255
	descInput.Width = width - 2

	searchInput := textinput.New()
	searchInput.Placeholder = "Search (unused)"

	v := &TimeEntryView{
		width:          width,
		height:         height,
		timeLogService: timeLogService,
		issue:          issue,
		datePicker:     NewDatePicker(),
		hoursSelector:  NewHoursSelector(),
		descInput:      descInput,
		focusIndex:     DateIndex,
		activities:     activitiesList,
		state:          StateEditing,
		SearchInput:    &searchInput,
	}

	return v, nil
}

// Init implements the tea.Model interface and returns the initial command for the time entry view.
// Init sets up the text input blinking cursor animation.
func (v *TimeEntryView) Init() tea.Cmd {
	return textinput.Blink
}

// Update processes input messages and updates the time entry view state accordingly.
// Update handles keyboard input for navigation, field editing, and form submission.
// It returns tea commands for async operations like form submission or view transitions.
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

	if v.state == StateSubmitting {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if v.errorMessage != "" && v.state == StateEditing {
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
				return v.submitWithCommand(v.issue.ID())
			}
		case "up", "down":
			if v.focusIndex == ActivityIndex && len(v.activities) > 0 {
				v.handleActivitySelection(msg.String())
				return nil
			}
		}

		if v.focusIndex == DateIndex {
			v.datePicker.Update(msg)
			return nil
		}

		if v.focusIndex == HoursIndex {
			v.hoursSelector.Update(msg)
			return nil
		}

		if v.focusIndex == DescriptionIndex {
			var cmd tea.Cmd
			v.descInput, cmd = v.descInput.Update(msg)
			return cmd
		}
	}

	return nil
}

// submitWithCommand validates the form data and handles the time entry submission process.
// submitWithCommand performs field validation, creates the time entry, and manages state transitions.
// It returns a tea command for async submission or nil if validation fails.
func (v *TimeEntryView) submitWithCommand(issueID int) tea.Cmd {
	if !v.HasValidEntry() {
		v.errorMessage = "Please fill in all required fields"
		return nil
	}

	v.errorMessage = ""

	v.state = StateSubmitting

	activityID := 0
	if v.selectedActivity != nil {
		activityID = v.selectedActivity.ID
	}

	_, err := v.timeLogService.CreateTimeEntry(models.CreateTimeEntryParams{
		IssueID:    issueID,
		ActivityID: activityID,
		Hours:      v.hoursSelector.SelectedHours(),
		Comments:   strings.TrimSpace(v.descInput.Value()),
		SpentOn:    v.datePicker.SelectedDate().Format("2006-01-02"),
	})
	if err != nil {
		v.state = StateError
		v.errorMessage = err.Error()
		return nil
	}

	return func() tea.Msg { return TimeEntrySubmissionSuccess{} }
}

// Render returns the time entry view using internal state and implements the View interface.
// Render delegates to RenderWithParams using the view's internal state values.
func (v *TimeEntryView) Render() string {
	title := titleStyle.Width(v.width).Render("TIME ENTRY")

	if v.issue == nil {
		return title + "\n\n" + emptyMessageStyle.Render("No issue selected")
	}

	switch v.state {
	case StateSubmitting:
		return v.renderSubmittingState(title, v.issue, v.width, v.height)
	case StateCompleted:
		return v.renderCompletedState(title, v.issue, v.width, v.height)
	case StateError:
		return v.renderErrorState(title, v.issue, v.width, v.height)
	default:
		return v.renderEditingState(title, v.issue, v.activities, v.width, v.height)
	}
}

// renderEditingState renders the normal editing form with all input fields and validation.
// renderEditingState creates the complete form interface including date picker, hours selector, and activity selection.
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

	descFocused := v.focusIndex == DescriptionIndex
	if descFocused {
		v.descInput.Focus()
	} else {
		v.descInput.Blur()
	}
	descSection := v.renderInput("Description:", v.descInput, descFocused)

	submitSection := v.renderSubmitButton()

	var errorSection string
	if v.errorMessage != "" && v.state == StateEditing {
		errorSection = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Error).
			Bold(true).
			Padding(0, 1).
			Render("⚠ " + v.errorMessage)
	}

	helpText := helpStyle.Render("Tab/Shift+Tab: next/previous field | Esc: back")

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
		lipgloss.NewStyle().Height(height - 35).Render(activitySection),
		"",
		submitSection,
	}

	if errorSection != "" {
		sections = append(sections, "", errorSection)
	}

	sections = append(sections, "", helpText)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderSubmittingState renders the loading state during time entry submission.
// renderSubmittingState displays a spinner animation and loading message while the form is being processed.
func (v *TimeEntryView) renderSubmittingState(title string, issue *domain.Issue, width, height int) string {
	issueInfo := fieldLabelStyle.Render(fmt.Sprintf("Issue: #%d %s", issue.ID(), issue.Title()))

	loadingMessage := loadingStyle.Render("Submitting time entry...")

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

// renderCompletedState renders the success state after successful time entry submission.
// renderCompletedState displays a success message and instructions for returning to the issue view.
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

// renderErrorState renders the error state when time entry submission fails.
// renderErrorState displays the error message with options to retry or return to the issue view.
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

// SetIssue updates the current issue context for the time entry view.
// SetIssue allows changing the issue that time will be logged against.
func (v *TimeEntryView) SetIssue(issue *domain.Issue) {
	v.issue = issue
}

// SetSize updates the dimensions of the time entry view and adjusts child components accordingly.
// SetSize resizes the view and updates the description input field width to match.
func (v *TimeEntryView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.descInput.Width = width - 2
}

// GetState returns the current operational state of the time entry view.
// GetState provides access to the current state for external components to react accordingly.
func (v *TimeEntryView) GetState() TimeEntryState {
	return v.state
}

// GetErrorMessage returns the current error message if the view is in an error state.
// GetErrorMessage provides access to error details for display or logging purposes.
func (v *TimeEntryView) GetErrorMessage() string {
	return v.errorMessage
}

// Reset resets the time entry view to its initial state, clearing all form data and errors.
// Reset restores default values, clears error messages, and returns focus to the first field.
func (v *TimeEntryView) Reset() {
	v.state = StateEditing
	v.errorMessage = ""
	v.descInput.SetValue("")
	v.focusIndex = DateIndex
}

// renderActivitySelection creates the visual representation of available activities with selection highlighting.
// renderActivitySelection displays the list of activities and highlights the currently selected one.
func (v *TimeEntryView) renderActivitySelection(activities []Activity) string {
	if len(activities) == 0 {
		return fieldLabelStyle.Render("Activity: ") + "No activities available"
	}

	if v.selectedActivity == nil && len(activities) > 0 {
		v.selectedActivity = &activities[0]
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

	if v.focusIndex == ActivityIndex {
		helpText := helpStyle.Render("Use Up/Down arrows to select activity")
		content.WriteString("\n" + helpText)
	}

	return content.String()
}

// renderInput renders a labeled input field for the TimeEntryView.
// If the label is "Description:", it displays a multi-line, textarea-like input with word wrapping,
// a visible cursor when focused, and a character count indicator. For other fields, it renders a
// single-line input. The appearance of the input adapts based on focus state, applying different
// styles and border colors accordingly.
//
// Parameters:
//   - label: the label to display above the input field.
//   - input: the textinput.Model representing the current input state.
//   - focused: whether the input field is currently focused.
func (v *TimeEntryView) renderInput(label string, input textinput.Model, focused bool) string {
	style := fieldValueStyle
	if focused {
		style = focusedStyle
	}

	var result string

	if label == "Description:" {
		// For description field, create a text area-like display
		inputText := input.Value()
		availableWidth := v.width - 10

		var wrappedLines []string
		if len(inputText) == 0 {
			wrappedLines = []string{""}
		} else {
			// Split text into lines that fit the available width
			for len(inputText) > 0 {
				if len(inputText) <= availableWidth {
					wrappedLines = append(wrappedLines, inputText)
					break
				}
				breakPoint := availableWidth
				for i := availableWidth - 1; i >= 0; i-- {
					if inputText[i] == ' ' {
						breakPoint = i + 1
						break
					}
				}
				wrappedLines = append(wrappedLines, inputText[:breakPoint])
				inputText = inputText[breakPoint:]
			}
		}

		if focused && len(wrappedLines) > 0 {
			lastIndex := len(wrappedLines) - 1
			wrappedLines[lastIndex] = wrappedLines[lastIndex] + "│"
		}

		wrappedText := strings.Join(wrappedLines, "\n")

		textAreaStyle := style.
			Width(availableWidth + 4).
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(themes.TokyoNight.Border)

		if focused {
			textAreaStyle = textAreaStyle.BorderForeground(themes.TokyoNight.Info)
		}

		result = fieldLabelStyle.Render(label) + "\n" +
			textAreaStyle.Render(wrappedText)
	} else {
		// For other fields, use the original single-line input
		result = fieldLabelStyle.Render(label) + "\n" +
			style.Width(v.width-4).Padding(0, 1).Render(input.View())
	}

	if focused && label == "Description:" {
		currentLength := len(input.Value())
		maxLength := 255

		charCountText := fmt.Sprintf("%d/%d characters", currentLength, maxLength)
		helpText := helpStyle.Render(charCountText)
		result += "\n" + helpText
	}

	return result
}

// renderDatePicker creates the date selection component with appropriate focus handling.
// renderDatePicker manages the visual state and focus of the date picker component.
func (v *TimeEntryView) renderDatePicker() string {
	focused := v.focusIndex == DateIndex
	if focused {
		v.datePicker.Focus()
	} else {
		v.datePicker.Blur()
	}

	return fieldLabelStyle.Width(v.width-4).Render("Date:") + "\n" + v.datePicker.Render()
}

// renderHoursSelector creates the hours selection component with appropriate focus handling.
// renderHoursSelector manages the visual state and focus of the hours selector component.
func (v *TimeEntryView) renderHoursSelector() string {
	focused := v.focusIndex == HoursIndex
	if focused {
		v.hoursSelector.Focus()
	} else {
		v.hoursSelector.Blur()
	}

	return fieldLabelStyle.Render("Hours:") + "\n" + v.hoursSelector.Render()
}

// renderSubmitButton creates the submit button with appropriate styling based on focus state.
// renderSubmitButton handles the visual representation of the form submission button.
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

// nextField advances the focus to the next input field in the form.
// nextField implements circular navigation through all form fields.
func (v *TimeEntryView) nextField() {
	v.focusIndex = (v.focusIndex + 1) % TotalFields
}

// prevField moves the focus to the previous input field in the form.
// prevField implements reverse circular navigation through all form fields.
func (v *TimeEntryView) prevField() {
	v.focusIndex = (v.focusIndex - 1 + TotalFields) % TotalFields
}

// handleActivitySelection processes up/down navigation within the activity selection list.
// handleActivitySelection updates the selected activity based on keyboard navigation input.
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

// TimeEntrySubmissionSuccess represents a successful time entry submission event.
type TimeEntrySubmissionSuccess struct{}

// TimeEntrySubmissionError represents a failed time entry submission event with error details.
type TimeEntrySubmissionError struct {
	Error error
}

// HasValidEntry validates that all required fields are properly filled for submission.
// HasValidEntry checks hours, description, and activity selection to ensure form completeness.
func (v *TimeEntryView) HasValidEntry() bool {
	hasValidHours := v.hoursSelector.SelectedHours() > 0
	hasValidDescription := strings.TrimSpace(v.descInput.Value()) != ""

	hasValidActivity := true
	if len(v.activities) > 0 {
		hasValidActivity = v.selectedActivity != nil
	}

	return hasValidActivity && hasValidHours && hasValidDescription
}
