package views

import (
	"fmt"

	"github.com/b1tray3r/rmt/internal/tui/domain"
	"github.com/b1tray3r/rmt/internal/tui/messages"
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type IssueView struct {
	width, height int
	Issue         *domain.Issue
	viewport      viewport.Model
}

func NewIssueView(width, height int, issue *domain.Issue) *IssueView {
	// Count the number of newline characters in the issue's full description.
	desc := issue.FullDescription()
	lineCount := 1
	for _, c := range desc {
		if c == '\n' {
			lineCount++
		}
	}

	maxheight := height - 12
	if lineCount <= maxheight {
		if lineCount < maxheight {
			for i := lineCount; i < maxheight; i++ {
				desc += "\n"
			}
		}
	}

	vp := viewport.New(width-4, maxheight)
	vp.SetContent(desc)

	return &IssueView{
		width:    width,
		Issue:    issue,
		viewport: vp,
	}
}

// Init initializes the IssueView and returns any initial command.
func (v *IssueView) Init() tea.Cmd {
	return nil
}

// Update updates the IssueView based on the incoming message.
func (v *IssueView) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			return func() tea.Msg {
				return messages.TimeEntryCreateMsg{Issue: v.Issue}
			}
		case "up", "k":
			v.viewport, cmd = v.viewport.Update(msg)
			return cmd
		case "down", "j":
			v.viewport, cmd = v.viewport.Update(msg)
			return cmd
		case "pgup":
			v.viewport, cmd = v.viewport.Update(msg)
			return cmd
		case "pgdown":
			v.viewport, cmd = v.viewport.Update(msg)
			return cmd
		case "home":
			v.viewport, cmd = v.viewport.Update(msg)
			return cmd
		case "end":
			v.viewport, cmd = v.viewport.Update(msg)
			return cmd
		}
	}

	// Update viewport with other messages
	v.viewport, cmd = v.viewport.Update(msg)
	return cmd
}

// Render renders the IssueView as a string.
func (v *IssueView) Render() string {
	style := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Secondary).
		Background(themes.TokyoNight.Background)

	pname := "--no-project--"
	p := v.Issue.Project()
	if p != nil {
		pname = p.Name()
	}

	projectInfo := lipgloss.JoinHorizontal(
		lipgloss.Left,
		style.
			Foreground(themes.TokyoNight.Primary).
			Italic(true).
			Render(pname),
	)

	titleInfo := lipgloss.JoinHorizontal(
		lipgloss.Left,
		style.Italic(true).Render(fmt.Sprintf("#%d", v.Issue.ID())),
		style.Render(" "),
		style.Bold(true).Render(v.Issue.FullTitle()),
		style.Foreground(themes.TokyoNight.Primary).Render(" by "),
		style.Italic(true).Foreground(themes.TokyoNight.Primary).Render(v.Issue.Author()),
	)

	helpText := "↑/↓/j/k: scroll • pgup/pgdown: page scroll • home/end: jump • t: log time • esc: back • ctrl+c: quit"
	help := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Muted).
		Background(themes.TokyoNight.Background).
		Padding(0, 1).
		Render(helpText)

	linkInfo := lipgloss.JoinHorizontal(
		lipgloss.Left,
		style.Padding(0, 0, 0, v.width-50).Foreground(themes.TokyoNight.Secondary).Render("  "),
		style.Foreground(themes.TokyoNight.Link).Render(v.Issue.Link()),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		projectInfo,
		style.Padding(1, 0, 1, 0).Render(linkInfo),
		style.Width(v.width).PaddingBottom(1).Render(titleInfo),
		style.PaddingBottom(1).Render(v.viewport.View()),
		help,
	)
}

func (v *IssueView) SetSize(width, height int) {
	v.width = width - 4
	v.height = height
}
