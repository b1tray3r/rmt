package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/b1tray3r/rmt/internal/config"
)

func run() error {
	// Initialize configuration settings
	// Get the user's home directory and construct the config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	configFile := filepath.Join(homeDir, ".config", "rmt", "config.yml")
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	_ = cfg

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
