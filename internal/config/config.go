// Package config provides configuration management for the application.
package config

import (
	"fmt"
	"os"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

// Config holds the application configuration.
type Config struct {
	Redmine RedmineConfig `yaml:"redmine"`
}

// RedmineConfig holds Redmine-specific configuration.
type RedmineConfig struct {
	URL        string `yaml:"url"`
	Token      string `yaml:"token"`
	Activities struct {
		Prefix []string `yaml:"prefix"`
	} `yaml:"activities"`
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(file string) (*Config, error) {
	var config Config
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// MissingFieldError represents an error for a missing required configuration field.
type MissingFieldError struct {
	Field string
}

// Error returns the error message for MissingFieldError.
func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("missing required field: %s", e.Field)
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Redmine.URL == "" {
		return &MissingFieldError{Field: "redmine.url"}
	}
	if c.Redmine.Token == "" {
		return &MissingFieldError{Field: "redmine.token"}
	}
	return nil
}
