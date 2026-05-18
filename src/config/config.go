package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Title       string `yaml:"site_title"`
	Desc        string `yaml:"site_desc"`
	BaseURL     string `yaml:"base_url"`
	Attribution bool   `yaml:"attribution"`
	FooterNote  string `yaml:"footer_note"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file - %w", err)
	}

	if len(bytes.TrimSpace(data)) == 0 {
		return nil, fmt.Errorf("config file at '%s' is empty", path)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config file - %w", err)
	}

	return &cfg, nil
}

func Default(host string, port int) *Config {
	return &Config{
		Title:       "Snipraw",
		Desc:        "Code snippets server",
		BaseURL:     fmt.Sprintf("http://%s:%d", host, port),
		Attribution: true,
	}
}

func Write(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not serialize config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
