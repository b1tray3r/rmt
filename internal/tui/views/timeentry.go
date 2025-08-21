package views

import (
	"github.com/b1tray3r/rmt/internal/tui/domain"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TimeLogView struct {
	width, height int
	SearchInput   *textinput.Model
}

func NewTimeLogView(width int, issue *domain.Issue) *TimeLogView {
	return &TimeLogView{
		width:  width,
		height: 0,
	}
}

func (v *TimeLogView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// Init initializes the TimeEntryView and returns the blinking cursor command.
func (v *TimeLogView) Init() tea.Cmd {
	return textinput.Blink
}

func (v *TimeLogView) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (v *TimeLogView) Render() string {
	return lipgloss.NewStyle().Width(v.width).Height(v.height).Render("Time Log")
}
