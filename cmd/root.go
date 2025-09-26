package cmd

import (
	"github.com/IanKnighton/institutionalized/internal/config"
	"github.com/spf13/cobra"
)

var appConfig *config.Config

var rootCmd = &cobra.Command{
	Use:   "institutionalized",
	Short: "A simple tool that uses an LLM to create commit and PR messages based on git status",
	Long: `institutionalized is a CLI tool that analyzes your staged git changes
and uses ChatGPT to generate conventional commit messages, then prompts
you to confirm before committing the changes.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Load configuration
	var err error
	appConfig, err = config.LoadConfig()
	if err != nil {
		// Use default config if loading fails
		appConfig = config.DefaultConfig()
	}

	rootCmd.PersistentFlags().StringP("api-key", "k", "", "OpenAI API key (can also be set via OPENAI_API_KEY environment variable)")
	rootCmd.PersistentFlags().Bool("emoji", appConfig.UseEmoji, "Use emoji in commit messages (overrides config file setting)")
}
