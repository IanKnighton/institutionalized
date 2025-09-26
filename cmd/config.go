package cmd

import (
	"fmt"
	"os"

	"github.com/IanKnighton/institutionalized/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Manage institutionalized configuration settings. You can view, set, or create a default configuration file.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration values from the config file or defaults.`,
	RunE:  runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `Set a configuration value. Available keys: use_emoji (true/false), providers.openai.enabled (true/false), providers.gemini.enabled (true/false), providers.priority (openai/gemini), providers.delay_threshold (seconds)`,
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default configuration file",
	Long:  `Create a default configuration file at ~/.config/institutionalized/config.yaml`,
	RunE:  runConfigInit,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configInitCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Current configuration:")
	fmt.Printf("  use_emoji: %t\n", cfg.UseEmoji)
	fmt.Printf("  providers:\n")
	fmt.Printf("    openai:\n")
	fmt.Printf("      enabled: %t\n", cfg.Providers.OpenAI.Enabled)
	fmt.Printf("    gemini:\n")
	fmt.Printf("      enabled: %t\n", cfg.Providers.Gemini.Enabled)
	fmt.Printf("    priority: %s\n", cfg.Providers.Priority)
	fmt.Printf("    delay_threshold: %d seconds\n", cfg.Providers.DelayThreshold)

	// Show config file location
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := fmt.Sprintf("%s/.config/institutionalized/config.yaml", homeDir)
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("\nConfig file: %s\n", configPath)
		} else {
			fmt.Printf("\nConfig file: %s (not found - using defaults)\n", configPath)
		}
	}

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "use_emoji":
		switch value {
		case "true", "1", "yes", "on":
			cfg.UseEmoji = true
		case "false", "0", "no", "off":
			cfg.UseEmoji = false
		default:
			return fmt.Errorf("invalid value for use_emoji: %s (expected true/false)", value)
		}
	case "providers.openai.enabled":
		switch value {
		case "true", "1", "yes", "on":
			cfg.Providers.OpenAI.Enabled = true
		case "false", "0", "no", "off":
			cfg.Providers.OpenAI.Enabled = false
		default:
			return fmt.Errorf("invalid value for providers.openai.enabled: %s (expected true/false)", value)
		}
	case "providers.gemini.enabled":
		switch value {
		case "true", "1", "yes", "on":
			cfg.Providers.Gemini.Enabled = true
		case "false", "0", "no", "off":
			cfg.Providers.Gemini.Enabled = false
		default:
			return fmt.Errorf("invalid value for providers.gemini.enabled: %s (expected true/false)", value)
		}
	case "providers.priority":
		if value == "openai" || value == "gemini" {
			cfg.Providers.Priority = value
		} else {
			return fmt.Errorf("invalid value for providers.priority: %s (expected openai/gemini)", value)
		}
	case "providers.delay_threshold":
		// Parse the value as integer
		var delayThreshold int
		if _, err := fmt.Sscanf(value, "%d", &delayThreshold); err != nil {
			return fmt.Errorf("invalid value for providers.delay_threshold: %s (expected number of seconds)", value)
		}
		if delayThreshold < 1 || delayThreshold > 300 {
			return fmt.Errorf("invalid value for providers.delay_threshold: %d (expected 1-300 seconds)", delayThreshold)
		}
		cfg.Providers.DelayThreshold = delayThreshold
	default:
		return fmt.Errorf("unknown config key: %s (available: use_emoji, providers.openai.enabled, providers.gemini.enabled, providers.priority, providers.delay_threshold)", key)
	}

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Configuration updated: %s = %s\n", key, value)
	return nil
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	cfg := config.DefaultConfig()

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Default configuration file created successfully!")
	} else {
		configPath := fmt.Sprintf("%s/.config/institutionalized/config.yaml", homeDir)
		fmt.Printf("Default configuration file created at: %s\n", configPath)
	}

	return nil
}
