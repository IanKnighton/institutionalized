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
	Long:  `Set a configuration value. Available keys: use_emoji (true/false)`,
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
	default:
		return fmt.Errorf("unknown config key: %s (available: use_emoji)", key)
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
