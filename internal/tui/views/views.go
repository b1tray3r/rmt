package views

import (
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	Render() string
	SetSize(width, height int)
}

// NewIssueDelegate creates a new delegate for issue items with Tokyo Night theme
func NewIssueDelegate(maxWidth int) list.ItemDelegate {
	d := list.NewDefaultDelegate()

	d.SetHeight(3)
	d.SetSpacing(1)

	d.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Secondary).
		Background(themes.TokyoNight.Background).
		Width(maxWidth).
		Bold(true).
		Padding(0, 1)

	d.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Muted).
		Background(themes.TokyoNight.Background).
		Width(maxWidth).
		Italic(true).
		Padding(0, 1)

	d.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Primary).
		Width(maxWidth).
		Bold(true)

	d.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Muted).
		Width(maxWidth).
		Italic(true)

	d.Styles.DimmedTitle = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Muted).
		Width(maxWidth)

	d.Styles.DimmedDesc = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Muted).
		Width(maxWidth).
		Italic(true)

	d.Styles.FilterMatch = lipgloss.NewStyle().
		Background(themes.TokyoNight.Warning).
		Foreground(themes.TokyoNight.Background).
		Width(maxWidth).
		Bold(true)

	return d
}
