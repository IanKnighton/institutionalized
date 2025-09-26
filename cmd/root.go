package cmd

import (
	"github.com/spf13/cobra"
)

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
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "OpenAI API key (can also be set via OPENAI_API_KEY environment variable)")
}
