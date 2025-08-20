package main

import (
	"fmt"
	"os"

	"github.com/b1tray3r/rmt/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// run starts the TUI application
func run() error {
	app := tui.NewApplication()

	program := tea.NewProgram(app, tea.WithAltScreen())
	_, err := program.Run()
	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
