package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/b1tray3r/rmt/internal/tui/domain"
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

type RMTIssueDelegate struct {
	maxWidth int
}

func (d RMTIssueDelegate) Height() int                             { return 3 }
func (d RMTIssueDelegate) Spacing() int                            { return 0 }
func (d RMTIssueDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d RMTIssueDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	issue, ok := item.(*domain.Issue)
	if !ok {
		return
	}

	var (
		title       = issue.Title()
		description = issue.Description()
		isSelected  = index == m.Index()
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if isFiltered && m.FilterValue() != "" {
		title = d.highlightMatches(title, m.FilterValue())
		description = d.highlightMatches(description, m.FilterValue())
	}

	prefix := "- "
	var titleStyle, descStyle lipgloss.Style
	if isSelected {
		prefix = " "
		titleStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Warning).
			Background(themes.TokyoNight.Background).
			Bold(true).
			Width(d.maxWidth)
		descStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Foreground).
			Background(themes.TokyoNight.Background).
			Italic(true).
			Width(d.maxWidth)
	} else {
		titleStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Foreground).
			Background(themes.TokyoNight.Background).
			Bold(true).
			Width(d.maxWidth)
		descStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Foreground).
			Background(themes.TokyoNight.Background).
			Italic(true).
			Width(d.maxWidth)
	}

	fmt.Fprint(w, titleStyle.Render(prefix+title))
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, descStyle.Padding(0, 1).PaddingBottom(1).Render(description))
}

func (d RMTIssueDelegate) highlightMatches(text, filter string) string {
	if filter == "" {
		return text
	}

	filterLower := strings.ToLower(filter)
	textLower := strings.ToLower(text)

	var result strings.Builder
	lastEnd := 0

	for i := 0; i < len(textLower); {
		matchIndex := strings.Index(textLower[i:], filterLower)
		if matchIndex == -1 {
			result.WriteString(text[lastEnd:])
			break
		}

		matchIndex += i
		result.WriteString(text[lastEnd:matchIndex])

		matchStyle := lipgloss.NewStyle().
			Background(themes.TokyoNight.Warning).
			Foreground(themes.TokyoNight.Background).
			Bold(true)
		matchEnd := matchIndex + len(filter)
		result.WriteString(matchStyle.Render(text[matchIndex:matchEnd]))

		lastEnd = matchEnd
		i = matchEnd
	}

	return result.String()
}

func NewIssueDelegate(maxWidth int) list.ItemDelegate {
	return RMTIssueDelegate{maxWidth: maxWidth}
}
