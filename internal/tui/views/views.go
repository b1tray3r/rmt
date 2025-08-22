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

// customIssueDelegate implements list.ItemDelegate with custom filtering behavior
type customIssueDelegate struct {
	maxWidth int
}

func (d customIssueDelegate) Height() int                             { return 3 }
func (d customIssueDelegate) Spacing() int                            { return 1 }
func (d customIssueDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d customIssueDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
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

	// Apply highlighting if filtering is active
	if isFiltered && m.FilterValue() != "" {
		title = d.highlightMatches(title, m.FilterValue())
		description = d.highlightMatches(description, m.FilterValue())
	}

	// Choose styles based on selection state
	var titleStyle, descStyle lipgloss.Style
	if isSelected {
		titleStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Warning).
			Background(themes.TokyoNight.Background).
			Bold(true).
			Width(d.maxWidth)
		descStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Muted).
			Background(themes.TokyoNight.Background).
			Italic(true).
			Width(d.maxWidth)
	} else {
		titleStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Background).
			Bold(true).
			Width(d.maxWidth)
		descStyle = lipgloss.NewStyle().
			Foreground(themes.TokyoNight.Border).
			Italic(true).
			Width(d.maxWidth)
	}

	// Render the item
	fmt.Fprint(w, titleStyle.Render(title))
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, descStyle.Render(description))
}

func (d customIssueDelegate) highlightMatches(text, filter string) string {
	if filter == "" {
		return text
	}

	filterLower := strings.ToLower(filter)
	textLower := strings.ToLower(text)

	// Find all matches
	var result strings.Builder
	lastEnd := 0

	for i := 0; i < len(textLower); {
		matchIndex := strings.Index(textLower[i:], filterLower)
		if matchIndex == -1 {
			// No more matches, append the rest
			result.WriteString(text[lastEnd:])
			break
		}

		// Adjust match index to be relative to the full string
		matchIndex += i
		matchEnd := matchIndex + len(filter)

		// Append text before match
		result.WriteString(text[lastEnd:matchIndex])

		// Append highlighted match
		matchStyle := lipgloss.NewStyle().
			Background(themes.TokyoNight.Warning).
			Foreground(themes.TokyoNight.Background).
			Bold(true)
		result.WriteString(matchStyle.Render(text[matchIndex:matchEnd]))

		// Update positions
		lastEnd = matchEnd
		i = matchEnd
	}

	return result.String()
}

// NewIssueDelegate creates a new delegate for issue items with Tokyo Night theme
func NewIssueDelegate(maxWidth int) list.ItemDelegate {
	return customIssueDelegate{maxWidth: maxWidth}
}
