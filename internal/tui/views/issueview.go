package views

import (
	"fmt"

	"github.com/b1tray3r/rmt/internal/tui/domain"
	"github.com/b1tray3r/rmt/internal/tui/messages"
	"github.com/b1tray3r/rmt/internal/tui/themes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type IssueView struct {
	width, height int
	Issue         *domain.Issue
}

func NewIssueView(width int, issue *domain.Issue) *IssueView {
	return &IssueView{
		width: width,
		Issue: issue,
	}
}

// Init initializes the IssueView and returns any initial command.
func (v *IssueView) Init() tea.Cmd {
	return nil
}

// Update updates the IssueView based on the incoming message.
func (v *IssueView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			return func() tea.Msg {
				return messages.TimeEntryCreateMsg{Issue: v.Issue}
			}
		}
	}

	return nil
}

// Render renders the IssueView as a string.
func (v *IssueView) Render() string {
	style := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Secondary).
		Background(themes.TokyoNight.Background)

	titleInfo := lipgloss.JoinHorizontal(
		lipgloss.Left,
		style.Italic(true).Render(fmt.Sprintf("#%d", v.Issue.ID())),
		style.Render(" "),
		style.Bold(true).Render(v.Issue.Title()),
		style.Foreground(themes.TokyoNight.Primary).Render(" by "),
		style.Italic(true).Foreground(themes.TokyoNight.Primary).Render(v.Issue.Author()),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		style.Width(v.width).Render(titleInfo),
		style.Padding(1, 0).Foreground(themes.TokyoNight.Link).Render(v.Issue.Link()),
		style.Height(v.height-2).Render(v.Issue.Description()),
	)
}

func (v *IssueView) SetSize(width, height int) {
	v.width = width
	v.height = height
}
