package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SlackToken   string
	WorkspaceURL string
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "slack-tui", "config.json"), nil
}

func LoadConfig() (Config, error) {
	var cfg Config
	path, err := ConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func SaveConfig(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
