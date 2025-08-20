package views

import (
	"github.com/b1tray3r/rmt/internal/tui/messages"
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SearchView struct {
	width, height int
	SearchInput   *textinput.Model
}

func NewSearchView(width int) *SearchView {
	ti := textinput.New()
	ti.Placeholder = "Search issues... (e.g., #123, 'bug fix', etc.)"

	// Apply Tokyo Night theme styling to placeholder
	ti.PlaceholderStyle = ti.PlaceholderStyle.
		Background(themes.TokyoNight.Background).
		Foreground(themes.TokyoNight.Muted)

	// Also style the text input itself to match the theme
	ti.TextStyle = ti.TextStyle.
		Background(themes.TokyoNight.Background).
		Foreground(themes.TokyoNight.Foreground)

	ti.Focus()
	ti.CharLimit = 256
	ti.Width = width
	ti.Prompt = "> "

	return &SearchView{
		SearchInput: &ti,
		width:       width,
	}
}

func (v *SearchView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

func (v *SearchView) Init() tea.Cmd {
	return nil
}

func (v *SearchView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if v.SearchInput.Value() != "" {
				query := v.SearchInput.Value()
				return func() tea.Msg {
					return messages.SearchSubmittedMsg{Query: query}
				}
			}
		case "esc":
			v.SearchInput.SetValue("")
			v.SearchInput.Blur()
			return nil
		}
	}

	var cmd tea.Cmd
	*v.SearchInput, cmd = v.SearchInput.Update(msg)
	return cmd
}

func (v *SearchView) Render() string {
	return v.SearchInput.View()
}
