// Package main_test contains tests for the main package
package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRun_Success verifies that run function works correctly with a valid config file.
func TestRun_Success(t *testing.T) {
	// Create a temporary config file
	dir := t.TempDir()
	configDir := filepath.Join(dir, ".config", "rmt")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configContent := `
redmine:
  url: https://example.com
  token: secret
  activities:
    prefix:
      - dev
      - test
`
	configPath := filepath.Join(configDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Temporarily change the home directory for the test
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	defer os.Setenv("HOME", originalHome)

	// Test the run function
	err := run()
	if err != nil {
		t.Errorf("run() returned error: %v", err)
	}
}

// TestRun_ConfigFileNotFound verifies that run function returns an error when config file doesn't exist.
func TestRun_ConfigFileNotFound(t *testing.T) {
	// Use a temporary directory where no config file exists
	dir := t.TempDir()

	// Temporarily change the home directory for the test
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	defer os.Setenv("HOME", originalHome)

	// Test the run function
	err := run()
	if err == nil {
		t.Error("expected error when config file doesn't exist, got nil")
	}
}

// TestRun_InvalidConfig verifies that run function returns an error when config is invalid.
func TestRun_InvalidConfig(t *testing.T) {
	// Create a temporary config file with invalid content
	dir := t.TempDir()
	configDir := filepath.Join(dir, ".config", "rmt")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configContent := `
redmine:
  url: ""
  token: ""
`
	configPath := filepath.Join(configDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Temporarily change the home directory for the test
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	defer os.Setenv("HOME", originalHome)

	// Test the run function
	err := run()
	if err == nil {
		t.Error("expected error for invalid config, got nil")
	}
}
