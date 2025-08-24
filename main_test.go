// Package main_test contains tests for the main package.
package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/b1tray3r/rmt/internal/config"
)

// TestLoadConfig_Success verifies that loadConfig successfully loads a valid configuration file from current directory
func TestLoadConfig_Success(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd) // Restore original directory
	}()

	// Create a valid config file in current directory
	configContent := `
redmine:
  url: https://redmine.example.com
  token: test-token-123
  activities:
    prefix:
      - dev
      - test
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Test loadConfig
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() returned error: %v", err)
	}
	if cfg == nil {
		t.Fatal("loadConfig() returned nil config")
	}

	// Verify config contents
	if cfg.Redmine.URL != "https://redmine.example.com" {
		t.Errorf("Redmine.URL = %q, want %q", cfg.Redmine.URL, "https://redmine.example.com")
	}
	if cfg.Redmine.Token != "test-token-123" {
		t.Errorf("Redmine.Token = %q, want %q", cfg.Redmine.Token, "test-token-123")
	}
}

// TestLoadConfig_XDGConfigHome verifies that loadConfig loads from XDG_CONFIG_HOME when local config doesn't exist
func TestLoadConfig_XDGConfigHome(t *testing.T) {
	// Create temporary directories
	tempDir := t.TempDir()
	xdgConfigDir := filepath.Join(tempDir, "config")
	rmtConfigDir := filepath.Join(xdgConfigDir, "rmt")

	if err := os.MkdirAll(rmtConfigDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Set XDG_CONFIG_HOME environment variable
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", xdgConfigDir)
	defer func() {
		if originalXDG == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		}
	}()

	// Change to a directory without local config
	workDir := filepath.Join(tempDir, "work")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		t.Fatalf("failed to create work directory: %v", err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Create config in XDG location
	configContent := `
redmine:
  url: https://xdg.example.com
  token: xdg-token-456
`
	configPath := filepath.Join(rmtConfigDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Test loadConfig
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() returned error: %v", err)
	}
	if cfg == nil {
		t.Fatal("loadConfig() returned nil config")
	}

	// Verify config was loaded from XDG location
	if cfg.Redmine.URL != "https://xdg.example.com" {
		t.Errorf("Redmine.URL = %q, want %q", cfg.Redmine.URL, "https://xdg.example.com")
	}
}

// TestLoadConfig_HomeConfig verifies fallback to ~/.config when XDG_CONFIG_HOME is not set
func TestLoadConfig_HomeConfig(t *testing.T) {
	// Skip this test in environments where we can't control HOME
	if testing.Short() {
		t.Skip("Skipping HOME directory test in short mode")
	}

	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	configDir := filepath.Join(homeDir, ".config", "rmt")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Unset XDG_CONFIG_HOME and set HOME
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Setenv("HOME", homeDir)

	defer func() {
		if originalXDG == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		}
		if originalHome == "" {
			os.Unsetenv("HOME")
		} else {
			os.Setenv("HOME", originalHome)
		}
	}()

	// Change to a directory without local config
	workDir := filepath.Join(tempDir, "work")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		t.Fatalf("failed to create work directory: %v", err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Create config in home location
	configContent := `
redmine:
  url: https://home.example.com
  token: home-token-789
`
	configPath := filepath.Join(configDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Test loadConfig
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() returned error: %v", err)
	}
	if cfg == nil {
		t.Fatal("loadConfig() returned nil config")
	}

	// Verify config was loaded from home location
	if cfg.Redmine.URL != "https://home.example.com" {
		t.Errorf("Redmine.URL = %q, want %q", cfg.Redmine.URL, "https://home.example.com")
	}
}

// TestLoadConfig_FileNotFound verifies that loadConfig returns an error when config file doesn't exist
func TestLoadConfig_FileNotFound(t *testing.T) {
	// Create empty temporary directory
	tempDir := t.TempDir()

	// Unset environment variables to avoid finding config elsewhere
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, "nonexistent"))
	os.Setenv("HOME", filepath.Join(tempDir, "nonexistent"))

	defer func() {
		if originalXDG == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		}
		if originalHome == "" {
			os.Unsetenv("HOME")
		} else {
			os.Setenv("HOME", originalHome)
		}
	}()

	// Change to temp directory (no config file)
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Test loadConfig should fail
	_, err = loadConfig()
	if err == nil {
		t.Fatal("loadConfig() should return error when config file doesn't exist")
	}
	if !strings.Contains(err.Error(), "failed to load config") {
		t.Errorf("error should contain 'failed to load config', got: %v", err)
	}
}

// TestLoadConfig_InvalidYAML verifies that loadConfig returns an error for invalid YAML
func TestLoadConfig_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Create invalid YAML config
	invalidConfig := `
redmine:
  url: https://example.com
  token: test-token
  activities:
    prefix:
      - dev
    - invalid: yaml structure
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(invalidConfig), 0600); err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}

	// Test loadConfig should fail
	_, err = loadConfig()
	if err == nil {
		t.Fatal("loadConfig() should return error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "failed to load config") {
		t.Errorf("error should contain 'failed to load config', got: %v", err)
	}
}

// TestRun_Success verifies that run function works with valid configuration
func TestRun_Success(t *testing.T) {
	// This test is tricky because run() starts a Bubble Tea program
	// We can test up to the point where the TUI would start

	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Create a valid config file
	configContent := `
redmine:
  url: https://redmine.example.com
  token: test-token-123
  activities:
    prefix:
      - dev
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// We can't easily test the full run() function without mocking the TUI
	// Instead, test the components that run() uses
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() failed: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("config validation failed: %v", err)
	}

	// If we reach here, run() would succeed up to starting the TUI
	t.Log("run() components (config loading and validation) work correctly")
}

// TestRun_InvalidConfig verifies that run returns an error for invalid configuration
func TestRun_InvalidConfig(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Create invalid config (missing required fields)
	invalidConfig := `
redmine:
  # Missing url and token
  activities:
    prefix:
      - dev
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(invalidConfig), 0600); err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}

	// Test that run returns an error
	err = run()
	if err == nil {
		t.Fatal("run() should return error for invalid config")
	}
	var missingFieldErr *config.MissingFieldError
	if !errors.As(err, &missingFieldErr) {
		t.Errorf("error should be of type MissingFieldError, got: %v", err)
	}
}

// TestRun_MissingConfig verifies that run returns an error when no config file exists
func TestRun_MissingConfig(t *testing.T) {
	tempDir := t.TempDir()

	// Set environment to point to non-existent config locations
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, "nonexistent"))
	os.Setenv("HOME", filepath.Join(tempDir, "nonexistent"))

	defer func() {
		if originalXDG == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		}
		if originalHome == "" {
			os.Unsetenv("HOME")
		} else {
			os.Setenv("HOME", originalHome)
		}
	}()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Test that run returns an error
	err = run()
	if err == nil {
		t.Fatal("run() should return error when config file doesn't exist")
	}
	if !strings.Contains(err.Error(), "failed loading config") {
		t.Errorf("error should contain 'failed loading config', got: %v", err)
	}
}

func TestRun_InValidConfig(t *testing.T) {
	// This test checks that run() returns an error when config is invalid
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		os.Chdir(originalWd)
	}()

	// Create an invalid config file (missing required fields)
	invalidConfig := `
redmine:
  # Missing url and token
  activities:
    prefix:
      - dev
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(invalidConfig), 0600); err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}

	// Test that run returns an error
	err = run()
	if err == nil {
		t.Fatal("run() should return error for invalid config")
	}

	var invalidConfigErr *config.MissingFieldError
	if !errors.As(err, &invalidConfigErr) {
		t.Errorf("error should be of type MissingFieldError, got: %v", err)
	}
}
