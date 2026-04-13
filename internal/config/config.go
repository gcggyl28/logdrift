package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Service defines a single log source to tail.
type Service struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
	Color   string `yaml:"color"`
}

// Config is the top-level logdrift configuration.
type Config struct {
	Services  []Service `yaml:"services"`
	DiffMode  string    `yaml:"diff_mode"`  // "unified" | "inline"
	BufferSize int      `yaml:"buffer_size"` // lines to keep per service
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		DiffMode:   "unified",
		BufferSize: 200,
	}
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parsing yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks required fields and constraints.
func (c *Config) validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("at least one service must be defined")
	}
	for i, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("service[%d]: name is required", i)
		}
		if svc.Command == "" {
			return fmt.Errorf("service %q: command is required", svc.Name)
		}
	}
	if c.DiffMode != "unified" && c.DiffMode != "inline" {
		return fmt.Errorf("diff_mode must be \"unified\" or \"inline\", got %q", c.DiffMode)
	}
	if c.BufferSize <= 0 {
		return fmt.Errorf("buffer_size must be positive")
	}
	return nil
}
