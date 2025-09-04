package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/b1tray3r/rmt/internal/config"
	"github.com/b1tray3r/rmt/internal/redmine"
	"github.com/b1tray3r/rmt/internal/tui"
	"github.com/b1tray3r/rmt/internal/tui/domain"
	tea "github.com/charmbracelet/bubbletea"
)

func loadConfig() (*config.Config, error) {
	configPath := "./"
	if _, err := os.Stat(configPath + "config.yml"); err != nil {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %w", err)
			}
			xdgConfigHome = filepath.Join(homeDir, ".config")
		}

		configPath = filepath.Join(xdgConfigHome, "rmt") + "/"
	}

	configFile := configPath + "config.yml"

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed loading config: %w", err)
	}

	client := redmine.NewRestClient(cfg.Redmine.URL, cfg.Redmine.Token)

	issueService := domain.NewRedmineIssueRepository(client)

	program := tea.NewProgram(
		tui.NewApplication(issueService, cfg),
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
