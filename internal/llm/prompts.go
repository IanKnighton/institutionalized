package llm

import "fmt"

// CommitMessagePromptTemplate generates the prompt for commit message generation
func CommitMessagePromptTemplate(diff string, useEmoji bool, userContext string) string {
	emojiInstruction := ""
	if useEmoji {
		emojiInstruction = "\n- Add an appropriate emoji at the beginning of the commit type (✨ feat, 🐛 fix, 📚 docs, 💄 style, ♻️ refactor, ✅ test, 🔧 chore, ⚡ perf, 👷 ci, 🏗️ build, ⏪ revert)"
	}

	contextSection := ""
	if userContext != "" {
		contextSection = fmt.Sprintf("\n\nAdditional context from the developer:\n%s", userContext)
	}

	return fmt.Sprintf(`Analyze the following git diff and generate a conventional commit message. 

The commit message should follow the Conventional Commits specification:
- Start with a type (feat, fix, docs, style, refactor, test, chore, etc.)
- Include a brief description in present tense
- Keep the first line under 50 characters if possible
- Add a body if the change is complex (separate with blank line)%s

Git diff:
%s%s

Return only the commit message, nothing else.`, emojiInstruction, diff, contextSection)
}

// PRContentPromptTemplate generates the prompt for PR content generation
func PRContentPromptTemplate(commits, currentBranch, defaultBranch, prTemplate string, useEmoji bool, userContext string) string {
	emojiInstruction := ""
	if useEmoji {
		emojiInstruction = "\n- You may add appropriate emojis to make the PR more engaging if it fits naturally"
	}

	templateInstruction := ""
	if prTemplate != "" {
		templateInstruction = fmt.Sprintf(`

IMPORTANT: This repository has a pull request template that you MUST follow. Please structure your response to match this template as closely as possible:

--- PR TEMPLATE START ---
%s
--- PR TEMPLATE END ---

When generating the PR body, use the template structure above but fill it with content based on the commit analysis. Maintain the same sections and format from the template.`, prTemplate)
	}

	contextSection := ""
	if userContext != "" {
		contextSection = fmt.Sprintf("\n\nAdditional context from the developer:\n%s", userContext)
	}

	return fmt.Sprintf(`Analyze the following git commits and generate a comprehensive pull request title and body.

The pull request merges branch '%s' into '%s'.%s

Requirements:
- Generate a clear, concise PR title that summarizes the main purpose of the changes
- Create a detailed PR body with the following sections (unless overridden by template above):
  - ## Summary: Brief overview of what this PR accomplishes
  - ## Changes Made: Bullet points of key changes and improvements
  - ## Testing: Description of testing performed or needed
  - ## Additional Notes: Any important information for reviewers%s

Commits to analyze:
%s%s

Return the response in this exact format:
TITLE: [your generated title here]

BODY:
[your generated body here]`, currentBranch, defaultBranch, templateInstruction, emojiInstruction, commits, contextSection)
}
