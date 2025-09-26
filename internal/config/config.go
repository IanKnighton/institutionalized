package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	UseEmoji bool `yaml:"use_emoji"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		UseEmoji: false,
	}
}

// LoadConfig loads configuration from the config file or returns default config
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), nil
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), nil
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "institutionalized", "config.yaml"), nil
}

// GetEmojiForCommitType returns the appropriate emoji for a commit type
func GetEmojiForCommitType(commitType string) string {
	emojiMap := map[string]string{
		"feat":     "✨",
		"fix":      "🐛",
		"docs":     "📚",
		"style":    "💄",
		"refactor": "♻️",
		"test":     "✅",
		"chore":    "🔧",
		"perf":     "⚡",
		"ci":       "👷",
		"build":    "🏗️",
		"revert":   "⏪",
	}

	if emoji, exists := emojiMap[commitType]; exists {
		return emoji + " "
	}
	return ""
}
