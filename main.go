package main

import (
	"fmt"
	"os"

	"github.com/b1tray3r/rmt/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	app := tui.NewApplication()

	program := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
