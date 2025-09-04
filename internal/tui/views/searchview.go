package views

import (
	"fmt"

	"github.com/b1tray3r/rmt/internal/config"
	"github.com/b1tray3r/rmt/internal/tui/domain"
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
	config        *config.Config
}

func NewSearchView(width int, cfg *config.Config) *SearchView {
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
		NewFavoriteDelegate(width),
		0,
		0,
	)
	list.SetShowTitle(false)
	list.SetShowHelp(false) // Disable default help to render our own

	// Apply Tokyo Night theme to list styles (consistent with listview)
	list.Styles.HelpStyle = lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Foreground).
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

	return &SearchView{
		width:        width,
		focusedIndex: SearchInput,
		views: map[int]any{
			SearchInput: ti,
			Favorites:   list,
		},
		config: cfg,
	}
}

// InitializeFavorites sets up the initial favorites list
func (v *SearchView) InitializeFavorites() {
	favoritesList := v.views[Favorites].(list.Model)

	followUpFieldID := v.config.Redmine.FollowUpFieldID

	favoritesList.SetItems([]list.Item{})
	favorites := []list.Item{
		domain.NewFavorite(1, "Follow-up: diese Woche", fmt.Sprintf("f%%5B%%5D=status_id&op%%5Bstatus_id%%5D=o&f%%5B%%5D=assigned_to_id&op%%5Bassigned_to_id%%5D=%%3D&v%%5Bassigned_to_id%%5D%%5B%%5D=me&f%%5B%%5D=cf_%d&op%%5Bcf_%d%%5D=w", followUpFieldID, followUpFieldID)),
		domain.NewFavorite(2, "Meine offenen Tickets", "f%5B%5D=status_id&op%5Bstatus_id%5D=o&f%5B%5D=assigned_to_id&op%5Bassigned_to_id%5D=%3D&v%5Bassigned_to_id%5D%5B%5D=me"),
	}

	favoritesList.SetItems(favorites)
	v.views[Favorites] = favoritesList
}

func (v *SearchView) SetSize(width, height int) {
	v.width = width
	v.height = height - 1

	// Update the favorites list size (consistent with listview)
	if favoritesList, ok := v.views[Favorites].(list.Model); ok {
		favoritesList.SetSize(width-8, height-15)
		v.views[Favorites] = favoritesList
	}
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
				// Blur the text input when switching to favorites
				textInput := v.views[SearchInput].(textinput.Model)
				textInput.Blur()
				v.views[SearchInput] = textInput
			} else {
				v.focusedIndex = SearchInput
				// Focus the text input when switching back to search
				textInput := v.views[SearchInput].(textinput.Model)
				textInput.Focus()
				v.views[SearchInput] = textInput
			}
		case "enter":
			switch v.focusedIndex {
			case SearchInput:
				view := v.views[SearchInput].(textinput.Model)
				if view.Value() != "" {
					query := view.Value()
					return func() tea.Msg {
						return messages.SearchSubmittedMsg{Query: query}
					}
				}
			case Favorites:
				// Handle favorite selection
				favoritesList := v.views[Favorites].(list.Model)
				if selectedItem := favoritesList.SelectedItem(); selectedItem != nil {
					if favorite, ok := selectedItem.(*domain.Favorite); ok {
						query := favorite.Config()
						return func() tea.Msg {
							return messages.SearchSubmittedMsg{Query: query}
						}
					}
				}
			}
			// Note: ESC is intentionally not handled here - it should be ignored in SearchView
		}
	}

	var cmd tea.Cmd
	switch v.focusedIndex {
	case SearchInput:
		view := v.views[SearchInput].(textinput.Model)
		view, cmd = view.Update(msg)
		v.views[SearchInput] = view
	case Favorites:
		view := v.views[Favorites].(list.Model)
		view, cmd = view.Update(msg)
		v.views[Favorites] = view
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

	var hint string
	if v.focusedIndex == SearchInput {
		hint = "Press 'Enter' to search, 'Tab' to switch to favorites, 'Esc' to clear, 'Ctrl+c' to quit"
	} else {
		hint = "Press 'Enter' to select favorite, 'Tab' to switch to search, 'Ctrl+c' to quit"
	}

	hintView := lipgloss.NewStyle().
		Foreground(themes.TokyoNight.Foreground).
		Render(hint)

	return lipgloss.JoinVertical(lipgloss.Left,
		headline,
		inputView,
		listView,
		hintView,
	)
}
