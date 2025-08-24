package views

import (
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingView represents a loading screen with a spinner.
type LoadingView struct {
	width, height int
	spinner       spinner.Model
	message       string
}

// NewLoadingView creates a new LoadingView with the given width and message.
func NewLoadingView(width int, message string) *LoadingView {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(themes.TokyoNight.Highlight)

	return &LoadingView{
		width:   width,
		spinner: s,
		message: message,
	}
}

// SetSize sets the dimensions of the LoadingView.
func (v *LoadingView) SetSize(width, height int) {
	v.width = width
	v.height = height - 3
}

// Init initializes the LoadingView and returns the spinner tick command.
func (v *LoadingView) Init() tea.Cmd {
	return v.spinner.Tick
}

// Update handles messages for the LoadingView and updates the spinner.
func (v *LoadingView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	v.spinner, cmd = v.spinner.Update(msg)
	return cmd
}

// Render renders the LoadingView as a string with centered spinner and message.
func (v *LoadingView) Render() string {
	spinnerView := lipgloss.NewStyle().
		Render(v.spinner.View())

	messageStyle := lipgloss.NewStyle().
		Background(themes.TokyoNight.Background).
		Foreground(themes.TokyoNight.Foreground).
		Bold(true)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		messageStyle.Render(v.message),
		"",
		lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Muted).
			Render(spinnerView),
	)

	// Center the content vertically and horizontally
	return lipgloss.Place(
		v.width,
		v.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
