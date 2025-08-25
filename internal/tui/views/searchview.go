package views

import (
	"github.com/b1tray3r/rmt/internal/tui/messages"
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	SearchInput = iota
	Favorites
)

type SearchView struct {
	width, height int
	focusedIndex  int
	views         map[int]any
}

func NewSearchView(width int) *SearchView {
	ti := textinput.New()
	ti.Placeholder = "Search issues... (e.g., #123, 'bug fix', etc.)"

	// Apply Tokyo Night theme styling to placeholder
	ti.PlaceholderStyle = ti.PlaceholderStyle.
		Background(themes.TokyoNight.Background).
		Foreground(themes.TokyoNight.Foreground)

	// Also style the text input itself to match the theme
	ti.TextStyle = ti.TextStyle.
		Background(themes.TokyoNight.Background).
		Foreground(themes.TokyoNight.Foreground)

	ti.Focus()
	ti.CharLimit = 256
	ti.Width = width
	ti.Prompt = "󰅬 "

	list := list.New(
		[]list.Item{},
		NewIssueDelegate(width),
		0,
		0,
	)
	list.SetShowTitle(false)
	list.SetShowHelp(false) // Disable default help to render our own

	// Apply Tokyo Night theme to list styles (help text, status bar, etc.)
	list.Styles.HelpStyle = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Foreground).
		Background(themes.TokyoNight.Background)

	list.Styles.StatusBar = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Foreground).
		Background(themes.TokyoNight.Background)

	list.Styles.FilterPrompt = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Primary).
		Background(themes.TokyoNight.Background)

	list.Styles.FilterCursor = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Highlight).
		Background(themes.TokyoNight.Background)

	return &SearchView{
		width:        width,
		focusedIndex: SearchInput,
		views: map[int]any{
			SearchInput: ti,
			Favorites:   list,
		},
	}
}

func (v *SearchView) SetSize(width, height int) {
	v.width = width
	v.height = height - 1
}

// Init initializes the SearchView and returns the blinking cursor command.
func (v *SearchView) Init() tea.Cmd {
	return textinput.Blink
}

func (v *SearchView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if v.focusedIndex == SearchInput {
				v.focusedIndex = Favorites
			} else {
				v.focusedIndex = SearchInput
			}
		case "enter":
			if v.focusedIndex == SearchInput {
				view := v.views[SearchInput].(textinput.Model)
				if view.Value() != "" {
					query := view.Value()
					return func() tea.Msg {
						return messages.SearchSubmittedMsg{Query: query}
					}
				}
			}
		case "esc":
			if v.focusedIndex == SearchInput {
				view := v.views[SearchInput].(textinput.Model)
				view.SetValue("")
				v.views[SearchInput] = view
			}
		}
	}

	var cmd tea.Cmd
	if v.focusedIndex == SearchInput {
		view := v.views[SearchInput].(textinput.Model)
		view, cmd = view.Update(msg)
		v.views[SearchInput] = view
	}
	return cmd
}

func (v *SearchView) Render() string {
	headline := lipgloss.NewStyle().
		Bold(true).
		Foreground(themes.TokyoNight.Secondary).
		Padding(1, 0).
		Render("󰡦 Search issues:")

	inputView := lipgloss.NewStyle().
		Width(v.width).
		Padding(0, 2).
		Render(v.views[SearchInput].(textinput.Model).View())

	listView := lipgloss.NewStyle().
		Width(v.width).
		Height(v.height - 7).
		Render(v.views[Favorites].(list.Model).View())

	hint := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Foreground).
		Render("Press 'Enter' to search, 'Esc' to clear, 'Ctrl+c' to quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		headline,
		inputView,
		listView,
		hint,
	)
}
