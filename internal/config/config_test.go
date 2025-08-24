// Package config_test contains tests for the config package
package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadConfig_Success verifies that LoadConfig successfully loads a valid configuration file in YAML format.
func TestLoadConfig_Success(t *testing.T) {
	dir := t.TempDir()
	configContent := `
redmine:
  url: https://example.com
  token: secret
  activities:
    prefix:
      - dev
      - test
`
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig returned nil config")
	}
	if cfg.Redmine.URL != "https://example.com" {
		t.Errorf("Redmine.URL = %q, want %q", cfg.Redmine.URL, "https://example.com")
	}
	if cfg.Redmine.Token != "secret" {
		t.Errorf("Redmine.Token = %q, want %q", cfg.Redmine.Token, "secret")
	}
	if cfg.Redmine.Activities.Prefix == nil {
		t.Errorf("Redmine.Activities.Prefix is nil, want [dev test]")
	} else if len(cfg.Redmine.Activities.Prefix) != 2 || cfg.Redmine.Activities.Prefix[0] != "dev" || cfg.Redmine.Activities.Prefix[1] != "test" {
		t.Errorf("Redmine.Activities.Prefix = %v, want [dev test]", cfg.Redmine.Activities.Prefix)
	}
}

// TestConfig_Validate verifies that Validate returns errors for missing required fields.
func TestConfig_Validate(t *testing.T) {
	// Missing URL
	cfg := &Config{
		Redmine: RedmineConfig{
			URL:   "",
			Token: "token",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing URL, got nil")
	}
	if mfe, ok := err.(*MissingFieldError); !ok || mfe.Field != "redmine.url" {
		t.Errorf("expected MissingFieldError for redmine.url, got %v", err)
	}

	// Missing Token
	cfg = &Config{
		Redmine: RedmineConfig{
			URL:   "https://example.com",
			Token: "",
		},
	}
	err = cfg.Validate()
	if err == nil {
		t.Error("expected error for missing Token, got nil")
	}
	if mfe, ok := err.(*MissingFieldError); !ok || mfe.Field != "redmine.token" {
		t.Errorf("expected MissingFieldError for redmine.token, got %v", err)
	}

	// All fields present
	cfg = &Config{
		Redmine: RedmineConfig{
			URL:   "https://example.com",
			Token: "token",
		},
	}
	err = cfg.Validate()
	if err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}

// TestMissingFieldError_Error verifies the error message of MissingFieldError.
func TestMissingFieldError_Error(t *testing.T) {
	e := &MissingFieldError{Field: "redmine.url"}
	want := "missing required field: redmine.url"
	if e.Error() != want {
		t.Errorf("MissingFieldError.Error() = %q, want %q", e.Error(), want)
	}
}

// TestLoadConfig_FileNotFound verifies that LoadConfig returns an error if the config file does not exist.
func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error when config file does not exist, got nil")
	}
	// Check that the error message contains information about opening the file
	errStr := err.Error()
	if !strings.Contains(errStr, "failed to open config file") {
		t.Errorf("expected error message to contain 'failed to open config file', got: %v", err)
	}
}

// TestLoadConfig_InvalidYAML verifies that LoadConfig returns an error for invalid YAML.
func TestLoadConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("invalid: [unclosed"), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

// TestLoadConfig_ValidationFailure verifies that LoadConfig returns an error when validation fails.
func TestLoadConfig_ValidationFailure(t *testing.T) {
	dir := t.TempDir()
	configContent := `
redmine:
  url: ""
  token: ""
`
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("expected error for invalid config, got nil")
	}

	var missingFieldErr *MissingFieldError
	if !errors.As(err, &missingFieldErr) {
		t.Errorf("expected error to be of type MissingFieldError, got: %v", err)
	}
}
