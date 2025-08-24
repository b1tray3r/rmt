package tui

import (
	"github.com/b1tray3r/rmt/internal/tui/domain"
	"github.com/b1tray3r/rmt/internal/tui/messages"
	"github.com/b1tray3r/rmt/internal/tui/themes"
	"github.com/b1tray3r/rmt/internal/tui/views"
	tea "github.com/charmbracelet/bubbletea"
	lippgloss "github.com/charmbracelet/lipgloss"
)

const (
	SearchView = iota
	LoadingView
	ListView
	IssueView
	TimeLogView
)

type Application struct {
	width, height int

	currentView int
	views       map[int]views.View

	issueService *domain.RedmineIssueRepository
}

// NewApplication creates and returns a new Application instance.
func NewApplication(issueService *domain.RedmineIssueRepository) *Application {
	return &Application{
		width:  75,
		height: 0,
		views: map[int]views.View{
			SearchView: views.NewSearchView(75),
		},
		issueService: issueService,
	}
}

// Init initializes the Application and returns the initial command.
func (a *Application) Init() tea.Cmd {
	return tea.Batch(
		a.views[a.currentView].Init(),
	)
}

func (a *Application) searchIssues(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := a.issueService.Search(query)
		if err != nil {
			return messages.SearchCompletedMsg{Error: err}
		}
		return messages.SearchCompletedMsg{Results: results}
	}
}

// Update handles incoming messages and updates the Application's state.
func (a *Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = 75
		a.height = 25

		if msg.Width > a.width {
			a.width = msg.Width - 2
		}
		if msg.Height > a.height {
			a.height = msg.Height - 4
		}

		for _, view := range a.views {
			view.SetSize(a.width, a.height)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			newIndex := a.currentView - 1
			if newIndex == LoadingView {
				newIndex = SearchView
			}
			// Skip the loading view
			if _, ok := a.views[newIndex]; !ok {
				newIndex = SearchView
			}
			a.currentView = newIndex

			var cmd tea.Cmd
			if newIndex == SearchView {
				cmd = a.views[SearchView].Update(msg)
			}
			return a, cmd
		case "alt+f":
			a.currentView = SearchView

			var cmd tea.Cmd
			if a.currentView == SearchView {
				cmd = a.views[SearchView].Update(msg)
			}
			return a, cmd
		case "ctrl+c":
			return a, tea.Quit
		}
	case messages.SearchSubmittedMsg:
		// Switch to loading view
		a.currentView = LoadingView
		lv := views.NewLoadingView(a.width, "Searching issues")
		lv.SetSize(a.width, a.height)
		a.views[LoadingView] = lv

		// Start the search operation
		return a, tea.Batch(
			lv.Init(),
			a.searchIssues(msg.Query),
		)

	case messages.TimeEntryCreateMsg:
		a.currentView = TimeLogView
		iv := views.NewIssueView(a.width, a.height, msg.Issue)
		iv.SetSize(a.width, a.height)
		a.views[IssueView] = iv

		tv, err := views.NewTimeEntryView(a.width, a.height, msg.Issue, a.issueService, a.issueService)
		if err != nil {
			return a, nil
		}
		a.views[TimeLogView] = tv
		return a, tv.Init()

	case messages.SearchCompletedMsg:
		if msg.Error != nil {
			a.currentView = SearchView
			return a, nil
		}

		// Switch to list view with results
		a.currentView = ListView
		lv := views.NewListView(a.width)
		lv.SetItems(msg.Results)
		lv.SetSize(a.width, a.height)
		a.views[ListView] = lv

		return a, nil

	case messages.IssueSelectedMsg:
		a.currentView = IssueView
		iv := views.NewIssueView(a.width, a.height, msg.Issue)
		iv.SetSize(a.width, a.height)
		a.views[IssueView] = iv
		return a, iv.Init()

	case messages.ReturnToIssueMsg:
		a.currentView = IssueView
		return a, nil
	}

	cmd := a.views[a.currentView].Update(msg)
	return a, cmd
}

// View renders the Application's UI as a string.
func (a *Application) View() string {
	// View sets up the style using colors defined in the colors.go file.
	style := lippgloss.NewStyle().
		Padding(1, 2).
		Width(a.width).
		Border(lippgloss.RoundedBorder()).
		BorderForeground(themes.TokyoNight.Border).
		Background(themes.TokyoNight.Background).
		Foreground(themes.TokyoNight.Foreground)

	title := lippgloss.NewStyle().
		Bold(true).
		Margin(0, 0, 1, 0).
		Foreground(themes.TokyoNight.Highlight).
		Render("RMT - Redmine Management Tool")

	// Application.View renders the main application UI with a title and the current view.
	return lippgloss.JoinVertical(
		lippgloss.Top,
		style.Render(
			title,
			"\n",
			a.views[a.currentView].Render(),
		),
	)
}
