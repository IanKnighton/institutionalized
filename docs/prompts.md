# Prompt Templates

This document describes the centralized prompt templates used by `institutionalized` for AI-generated commit messages and pull request content. All AI providers (OpenAI, Gemini, and Claude) use the same prompt templates, ensuring consistency across different providers.

## Overview

The prompt templates are centrally managed in `internal/llm/prompts.go` and are used by all three supported AI providers:

- **OpenAI ChatGPT**: Uses gpt-3.5-turbo model
- **Google Gemini**: Uses gemini-pro model  
- **Anthropic Claude**: Uses claude-3-haiku-20240307 model

## Commit Message Prompt Template

### Purpose
Generates conventional commit messages based on staged git changes.

### Template Function
```go
func CommitMessagePromptTemplate(diff string, useEmoji bool) string
```

### Base Template
```
Analyze the following git diff and generate a conventional commit message. 

The commit message should follow the Conventional Commits specification:
- Start with a type (feat, fix, docs, style, refactor, test, chore, etc.)
- Include a brief description in present tense
- Keep the first line under 50 characters if possible
- Add a body if the change is complex (separate with blank line)

Git diff:
{diff}

Return only the commit message, nothing else.
```

### With Emoji Support (useEmoji=true)
When emoji support is enabled, the following instruction is added:

```
- Add an appropriate emoji at the beginning of the commit type (✨ feat, 🐛 fix, 📚 docs, 💄 style, ♻️ refactor, ✅ test, 🔧 chore, ⚡ perf, 👷 ci, 🏗️ build, ⏪ revert)
```

### Example Usage
```bash
# Without emoji
institutionalized commit

# With emoji  
institutionalized commit --emoji
```

### Example Output
```
# Without emoji
feat: add user authentication system

# With emoji
✨ feat: add user authentication system
```

## Pull Request Content Prompt Template

### Purpose
Generates comprehensive pull request titles and descriptions based on commit history.

### Template Function
```go
func PRContentPromptTemplate(commits, currentBranch, defaultBranch, prTemplate string, useEmoji bool) string
```

### Base Template
```
Analyze the following git commits and generate a comprehensive pull request title and body.

The pull request merges branch '{currentBranch}' into '{defaultBranch}'.

Requirements:
- Generate a clear, concise PR title that summarizes the main purpose of the changes
- Create a detailed PR body with the following sections (unless overridden by template above):
  - ## Summary: Brief overview of what this PR accomplishes
  - ## Changes Made: Bullet points of key changes and improvements
  - ## Testing: Description of testing performed or needed
  - ## Additional Notes: Any important information for reviewers

Commits to analyze:
{commits}

Return the response in this exact format:
TITLE: [your generated title here]

BODY:
[your generated body here]
```

### With PR Template Support
When a repository has a pull request template (`.github/pull_request_template.md`), the following instruction is added:

```
IMPORTANT: This repository has a pull request template that you MUST follow. Please structure your response to match this template as closely as possible:

--- PR TEMPLATE START ---
{prTemplate}
--- PR TEMPLATE END ---

When generating the PR body, use the template structure above but fill it with content based on the commit analysis. Maintain the same sections and format from the template.
```

### With Emoji Support (useEmoji=true)
When emoji support is enabled, the following instruction is added:

```
- You may add appropriate emojis to make the PR more engaging if it fits naturally
```

### Example Usage
```bash
# Create PR with default settings
institutionalized pr

# Create PR with emoji support
institutionalized pr --emoji

# Preview PR content without creating
institutionalized pr --dry-run
```

### Example Complete Prompt (with emoji and PR template)

For a branch called `feature/auth` merging into `main` with emoji support enabled and a PR template present:

```
Analyze the following git commits and generate a comprehensive pull request title and body.

The pull request merges branch 'feature/auth' into 'main'.

IMPORTANT: This repository has a pull request template that you MUST follow. Please structure your response to match this template as closely as possible:

--- PR TEMPLATE START ---
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing
--- PR TEMPLATE END ---

When generating the PR body, use the template structure above but fill it with content based on the commit analysis. Maintain the same sections and format from the template.

Requirements:
- Generate a clear, concise PR title that summarizes the main purpose of the changes
- Create a detailed PR body with the following sections (unless overridden by template above):
  - ## Summary: Brief overview of what this PR accomplishes
  - ## Changes Made: Bullet points of key changes and improvements
  - ## Testing: Description of testing performed or needed
  - ## Additional Notes: Any important information for reviewers
- You may add appropriate emojis to make the PR more engaging if it fits naturally

Commits to analyze:
abc123 feat: add login endpoint
def456 test: add authentication tests
ghi789 docs: update API documentation

Return the response in this exact format:
TITLE: [your generated title here]

BODY:
[your generated body here]
```

## Modifying Prompts

To modify the prompt templates:

1. Edit the functions in `internal/llm/prompts.go`
2. The changes will automatically apply to all three AI providers
3. Test changes using the `--dry-run` flag to verify prompt behavior
4. Consider the impact on all providers when making changes

## Benefits of Centralized Prompts

- **Consistency**: All AI providers use identical prompts
- **Maintainability**: Single location to update prompts
- **Testing**: Easy to iterate and test prompt improvements
- **Provider Independence**: Prompts work across different AI models
- **Version Control**: Track prompt changes over time

## Provider-Specific Considerations

While the prompts are centralized, each provider may interpret them slightly differently:

- **OpenAI**: Generally follows instructions precisely
- **Gemini**: May be more creative with formatting
- **Claude**: Tends to be more conservative with responses

The centralized approach ensures consistent instructions while allowing each provider's strengths to shine through.