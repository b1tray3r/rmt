package main

import (
	"fmt"
	"os"

	"github.com/b1tray3r/rmt/internal/config"
	"github.com/b1tray3r/rmt/internal/redmine"
	"github.com/b1tray3r/rmt/internal/tui"
	"github.com/b1tray3r/rmt/internal/tui/domain"
	tea "github.com/charmbracelet/bubbletea"
)

// run starts the TUI application
func run() error {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		xdgConfigHome = homeDir + "/.config"
	}
	configFile := xdgConfigHome + "/rmt/config.yml"

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	client := redmine.NewRestClient(cfg.Redmine.URL, cfg.Redmine.Token)

	issueService := domain.NewIssueService(client)

	program := tea.NewProgram(
		tui.NewApplication(issueService),
		tea.WithAltScreen(),
	)
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
