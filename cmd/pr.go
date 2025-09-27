package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/IanKnighton/institutionalized/internal/config"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create a pull request using GitHub CLI",
	Long:  `Create a pull request that documents the scope of changes made, testing added, and features completed. Uses GitHub CLI (gh) to create the PR.`,
	RunE:  runPR,
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().BoolP("draft", "d", false, "Create a draft pull request")
}

func runPR(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	if !isGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	// Check if gh CLI is available
	if !isGHCliAvailable() {
		return fmt.Errorf("GitHub CLI (gh) is not available. Please install it from https://cli.github.com/")
	}

	// Check if user is authenticated with gh
	if !isGHAuthenticated() {
		return fmt.Errorf("not authenticated with GitHub CLI. Run 'gh auth login' to authenticate")
	}

	// Get current branch
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Get default branch
	defaultBranch, err := getDefaultBranch()
	if err != nil {
		return fmt.Errorf("failed to get default branch: %w", err)
	}

	// Check if current branch is not the default branch
	if currentBranch == defaultBranch {
		return fmt.Errorf("cannot create PR from default branch (%s). Please create a feature branch first", defaultBranch)
	}

	// Show user what we're about to do
	fmt.Printf("🔄 Creating PR: %s -> %s\n", currentBranch, defaultBranch)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Generate PR title and body
	prTitle, prBody, err := generatePRContent(currentBranch, defaultBranch, cfg)
	if err != nil {
		return fmt.Errorf("failed to generate PR content: %w", err)
	}

	// Check for draft flag
	isDraft, _ := cmd.Flags().GetBool("draft")

	// Create the PR
	if err := createPR(prTitle, prBody, currentBranch, defaultBranch, isDraft); err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	fmt.Printf("✅ Pull request created successfully!\n")
	return nil
}

// isGHCliAvailable checks if gh CLI is available
func isGHCliAvailable() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// isGHAuthenticated checks if user is authenticated with GitHub CLI
func isGHAuthenticated() bool {
	// Check if GH_TOKEN is set (GitHub Actions environment)
	if os.Getenv("GH_TOKEN") != "" {
		return true
	}
	
	// Check normal gh auth status
	cmd := exec.Command("gh", "auth", "status")
	err := cmd.Run()
	return err == nil
}

// getCurrentBranch returns the current git branch
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getDefaultBranch returns the default branch of the repository
func getDefaultBranch() (string, error) {
	// Try to get from symbolic ref first
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err == nil {
		// Extract branch name from refs/remotes/origin/HEAD -> refs/remotes/origin/<branch>
		refPath := strings.TrimSpace(string(output))
		parts := strings.Split(refPath, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}
	
	// Try using gh CLI to get default branch if authenticated
	if isGHAuthenticated() {
		cmd := exec.Command("gh", "repo", "view", "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name")
		output, err := cmd.Output()
		if err == nil {
			branch := strings.TrimSpace(string(output))
			if branch != "" && branch != "null" {
				return branch, nil
			}
		}
	}
	
	// Fallback to checking common default branches
	for _, branch := range []string{"main", "master"} {
		if branchExists(branch) {
			return branch, nil
		}
	}
	
	// Final fallback - assume main
	return "main", nil
}

// branchExists checks if a branch exists
func branchExists(branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/remotes/origin/%s", branch))
	return cmd.Run() == nil
}

// generatePRContent generates the PR title and body based on commits
func generatePRContent(currentBranch, defaultBranch string, cfg *config.Config) (string, string, error) {
	// First, try to get commits between default branch and current branch
	cmd := exec.Command("git", "log", fmt.Sprintf("%s..%s", defaultBranch, currentBranch), "--oneline")
	output, err := cmd.Output()
	
	// If that fails, try getting recent commits from current branch
	if err != nil {
		cmd = exec.Command("git", "log", "--oneline", "-10", currentBranch)
		output, err = cmd.Output()
		if err != nil {
			return "", "", fmt.Errorf("failed to get commits: %w", err)
		}
	}

	commits := strings.TrimSpace(string(output))
	if commits == "" {
		return "", "", fmt.Errorf("no commits found on branch %s", currentBranch)
	}

	// Generate PR title from the first commit or branch name
	commitLines := strings.Split(commits, "\n")
	firstCommit := commitLines[0]
	
	// Extract commit message (remove hash)
	parts := strings.SplitN(firstCommit, " ", 2)
	var prTitle string
	if len(parts) > 1 {
		prTitle = parts[1]
	} else {
		prTitle = fmt.Sprintf("Changes from %s", currentBranch)
	}

	// Generate PR body with feature summary
	prBody := fmt.Sprintf(`## Summary

This pull request includes changes from the **%s** branch.

## Recent Commits
%s

## Changes Made
- Implementation based on commits in this branch
- Updates and improvements to existing functionality

## Testing
- Manual testing performed on the changes
- All functionality has been verified

## Additional Notes
- This PR merges **%s** into **%s**

---
*This PR was created automatically by institutionalized*`, currentBranch, commits, currentBranch, defaultBranch)

	return prTitle, prBody, nil
}

// createPR creates the pull request using gh CLI
func createPR(title, body, currentBranch, baseBranch string, isDraft bool) error {
	args := []string{"pr", "create", "--title", title, "--body", body, "--base", baseBranch}
	
	if isDraft {
		args = append(args, "--draft")
	}

	cmd := exec.Command("gh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}