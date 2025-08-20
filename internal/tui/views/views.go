package views

import tea "github.com/charmbracelet/bubbletea"

type View interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	Render() string
	SetSize(width, height int)
}
