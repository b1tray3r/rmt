package views

import (
	"github.com/b1tray3r/rmt/internal/tui/domain"
	"github.com/b1tray3r/rmt/internal/tui/messages"
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListView struct {
	list list.Model
}

func NewListView(maxWidth int) *ListView {
	list := list.New(
		[]list.Item{},
		NewIssueDelegate(maxWidth),
		0,
		0,
	)
	list.SetShowTitle(false)
	list.SetShowHelp(false) // Disable default help to render our own

	// Apply Tokyo Night theme to list styles (help text, status bar, etc.)
	list.Styles.HelpStyle = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Muted).
		Background(themes.TokyoNight.Background)

	list.Styles.StatusBar = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Foreground).
		Background(themes.TokyoNight.Background).
		MarginBottom(1)

	list.Styles.FilterPrompt = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Primary).
		Background(themes.TokyoNight.Background)

	list.Styles.FilterCursor = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Highlight).
		Background(themes.TokyoNight.Background)

	return &ListView{
		list: list,
	}
}

func (v *ListView) SetSize(width, height int) {
	v.list.SetSize(width-8, height-4)
}

func (v *ListView) SetItems(items []*domain.Issue) {
	i := make([]list.Item, 0, len(items))
	for _, issue := range items {
		i = append(i, issue)
	}
	v.list.SetItems(i)
}

// Init initializes the ListView and returns any initial command.
func (v *ListView) Init() tea.Cmd {
	return nil
}

// Update updates the ListView based on the incoming message.
func (v *ListView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit
		case "/":
			if v.list.FilterState() == list.Filtering {
				v.list.ResetFilter()
				v.list.SetFilterState(list.Unfiltered)
				v.list.SetFilteringEnabled(false)
			} else {
				v.list.ResetFilter()
				v.list.SetFilterState(list.Filtering)
				v.list.SetFilteringEnabled(true)
			}

			return nil
		case "t":
			// ListView handles the "t" key only when the filter input is not active.
			if v.list.FilterState() != list.Filtering {
				if selectedItem := v.list.SelectedItem(); selectedItem != nil {
					if issueItem, ok := selectedItem.(*domain.Issue); ok {
						return func() tea.Msg {
							return messages.TimeEntryCreateMsg{Issue: issueItem}
						}
					}
				}
			}
		case "enter":
			if selectedItem := v.list.SelectedItem(); selectedItem != nil {
				if issueItem, ok := selectedItem.(*domain.Issue); ok {
					return func() tea.Msg {
						return messages.IssueSelectedMsg{Issue: issueItem}
					}
				}
			}
		}
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return cmd
}

// Render renders the ListView to a string.
func (v *ListView) Render() string {
	listView := v.list.View()

	helpText := "↑/↓ • enter: select • /: filter • t: log time • esc: back • ctrl+c: quit"
	helpStyle := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Muted).
		Background(themes.TokyoNight.Background).
		Padding(0, 1).
		Render(helpText)

	return lipgloss.JoinVertical(lipgloss.Left, listView, helpStyle)
}
