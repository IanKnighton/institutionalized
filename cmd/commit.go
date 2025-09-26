package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate and commit changes using AI",
	Long:  `Analyze staged changes and generate a conventional commit message using ChatGPT, then prompt for confirmation.`,
	RunE:  runCommit,
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func runCommit(cmd *cobra.Command, args []string) error {
	// Get API key from flag or environment
	apiKey, _ := cmd.Flags().GetString("api-key")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		return fmt.Errorf("OpenAI API key is required. Set it via --api-key flag or OPENAI_API_KEY environment variable")
	}

	// Check if we're in a git repository
	if !isGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	// Get staged changes
	diff, err := getStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get staged changes: %w", err)
	}

	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf("no staged changes found. Use 'git add' to stage changes first")
	}

	fmt.Println("Analyzing staged changes...")

	// Generate commit message using ChatGPT
	commitMessage, err := generateCommitMessage(apiKey, diff)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Display the proposed commit message
	fmt.Printf("\nProposed commit message:\n%s\n\n", commitMessage)

	// Ask for user confirmation
	if !askForConfirmation("Do you want to commit with this message?") {
		fmt.Println("Commit cancelled.")
		return nil
	}

	// Commit the changes
	if err := commitChanges(commitMessage); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	fmt.Println("Changes committed successfully!")
	return nil
}

func isGitRepo() bool {
	_, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	return err == nil
}

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func generateCommitMessage(apiKey, diff string) (string, error) {
	prompt := fmt.Sprintf(`Analyze the following git diff and generate a conventional commit message. 

The commit message should follow the Conventional Commits specification:
- Start with a type (feat, fix, docs, style, refactor, test, chore, etc.)
- Include a brief description in present tense
- Keep the first line under 50 characters if possible
- Add a body if the change is complex (separate with blank line)

Git diff:
%s

Return only the commit message, nothing else.`, diff)

	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", err
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return strings.TrimSpace(openAIResp.Choices[0].Message.Content), nil
}

func askForConfirmation(question string) bool {
	fmt.Printf("%s (y/N): ", question)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func commitChanges(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}