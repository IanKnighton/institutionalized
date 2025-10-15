# institutionalized

A simple tool that uses an LLM to create commit and PR messages based on git status.

## Overview

Institutionalized is a Go CLI tool that analyzes your staged git changes and uses AI providers (OpenAI ChatGPT, Google Gemini, or Anthropic Claude) to generate conventional commit messages, then prompts you to confirm before committing the changes. It can also create comprehensive pull requests using GitHub CLI.

## Features

- 🤖 **AI-powered commit messages**: Uses OpenAI's ChatGPT, Google Gemini, or Anthropic Claude to generate meaningful commit messages
- 📝 **Conventional Commits**: Follows the Conventional Commits specification by default
- 🔍 **Smart analysis**: Analyzes your staged git changes to understand the context
- 🛡️ **User confirmation**: Always asks for confirmation before committing or creating PRs
- 🔧 **Flexible configuration**: Support for multiple AI providers with fallback capability
- 🚀 **Pull Request creation**: Creates comprehensive PRs with GitHub CLI integration
- 📋 **Draft PR support**: Option to create draft pull requests
- 🔍 **Dry-run mode**: Preview PR content without creating actual PRs
- ⚡ **Provider fallback**: Automatically switches to backup provider if primary fails or times out
- 😊 **Emoji support**: Optional emoji prefixes for commit types (✨ feat, 🐛 fix, etc.)
- 🚫 **Skip confirmation**: Optional flag to bypass confirmation prompts for automation

## Installation

### Quick Install

#### macOS (Homebrew - Recommended)

The fastest way to install on macOS:

```bash
brew install ianknighton/tap/institutionalized
```

#### Other Platforms

Install using Go (works on all platforms):

```bash
go install github.com/IanKnighton/institutionalized@latest
```

This installs the binary to your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your `$PATH` to use `institutionalized` from anywhere.

### Other Installation Methods

For detailed installation instructions including manual installation, PATH setup, and troubleshooting, see our comprehensive [Installation Guide](docs/installation.md).

**Available methods:**

- Install via Homebrew (macOS - recommended)
- Install via Go (all platforms)
- Build from source with global installation
- Download pre-built binaries from GitHub releases
- Manual binary placement and PATH configuration
- Platform-specific instructions for Linux, macOS, and Windows

## Usage

### Setup

You need an API key from one or more of the supported providers:

**OpenAI**: Get your API key from [OpenAI's platform](https://platform.openai.com/api-keys)
**Google Gemini**: Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
**Anthropic Claude**: Get your API key from [Anthropic Console](https://console.anthropic.com/)

Set your API key(s) as environment variables:

```bash
# For OpenAI (ChatGPT)
export OPENAI_API_KEY="your-openai-key-here"

# For Google Gemini
export GEMINI_API_KEY="your-gemini-key-here"

# For Anthropic Claude
export CLAUDE_API_KEY="your-claude-key-here"

# You can set multiple keys - the tool will use them based on your configuration
```

The tool will automatically detect which API keys are available and use them according to your configuration preferences.

### Basic Usage

1. Make your changes and stage them:

   ```bash
   git add .
   ```

2. Run institutionalized to generate and commit:

   ```bash
   institutionalized commit
   ```

3. Review the proposed commit message and confirm or cancel.

### Configuration

Institutionalized supports extensive configuration options to customize AI provider settings, emoji preferences, timeouts, and more. Configuration is managed through a YAML file and CLI commands.

**Quick configuration:**

```bash
# View current settings
institutionalized config show

# Create default configuration
institutionalized config init

# Set a provider as primary
institutionalized config set providers.priority claude
```

For comprehensive configuration documentation including all available options, provider management, emoji settings, and troubleshooting, see our detailed [Configuration Guide](docs/configuration.md).

**Key features:**

- Multiple AI provider support (OpenAI, Gemini, Claude)
- Configurable provider priority and fallback
- Emoji support for commit types
- Timeout and performance tuning
- Environment variable integration

### Commands

#### `institutionalized commit`

Analyzes staged changes and generates a conventional commit message using available AI providers (OpenAI ChatGPT, Google Gemini, or Anthropic Claude).

**Flags:**

- `--api-key, -k`: OpenAI API key (deprecated: use `OPENAI_API_KEY` environment variable)
- `--emoji`: Use emoji in commit messages (overrides config file setting)
- `--dry-run`: Show staged changes without calling API or committing (useful for testing)

**Examples:**

```bash
# Basic usage with environment variables
export OPENAI_API_KEY="your-openai-key"
export GEMINI_API_KEY="your-gemini-key"
export CLAUDE_API_KEY="your-claude-key"
institutionalized commit

# Using only OpenAI (if you have multiple keys but want to use only OpenAI)
institutionalized config set providers.gemini.enabled false
institutionalized commit

# Dry run to see what changes would be analyzed
institutionalized commit --dry-run

# Use emoji for this specific commit
institutionalized commit --emoji
```

#### `institutionalized config`

Manage configuration settings for institutionalized. Supports provider management, emoji preferences, timeout settings, and more.

**Subcommands:**

- `show`: Display current configuration values
- `set <key> <value>`: Set a configuration value
- `init`: Create a default configuration file

**Quick Examples:**

```bash
# View current configuration
institutionalized config show

# Set Claude as primary provider
institutionalized config set providers.priority claude

# Enable emoji support
institutionalized config set use_emoji true
```

For complete configuration documentation including all available options and advanced settings, see the [Configuration Guide](docs/configuration.md).

#### `institutionalized pr`

Create a pull request using GitHub CLI that documents the scope of changes made, testing added, and features completed.

**Requirements:**

- GitHub CLI (`gh`) installed and authenticated
- Must be on a feature branch (not the default branch)
- Must be in a git repository

**Flags:**

- `--draft, -d`: Create a draft pull request
- `--dry-run`: Show what would be done without creating the PR (doesn't require authentication)
- `--yes, -y`: Skip confirmation prompt and create PR immediately

**Examples:**

```bash
# Create a standard pull request (prompts for confirmation)
institutionalized pr

# Create a draft pull request (prompts for confirmation)
institutionalized pr --draft

# Skip confirmation and create PR immediately
institutionalized pr --yes

# Create draft PR without confirmation
institutionalized pr --draft --yes

# Preview what the PR would look like (no authentication required)
institutionalized pr --dry-run

# Preview a draft PR
institutionalized pr --draft --dry-run
```

**Pull Request Templates:**

The tool automatically detects and respects pull request templates in your repository. It looks for templates in the following locations (in order of priority):

- `.github/pull_request_template.md`
- `.github/PULL_REQUEST_TEMPLATE.md`
- `.github/PULL_REQUEST_TEMPLATE/pull_request_template.md`
- `docs/pull_request_template.md`

When a template is found, the LLM is instructed to follow the template structure while generating the PR description based on your commit history.

**How it works:**

1. Verifies `gh` CLI is installed and user is authenticated
2. Checks that current branch is not the default branch
3. Analyzes commit history between current branch and default branch
4. Detects and reads pull request template (if available)
5. Generates PR title from the most recent commit message
6. Creates comprehensive PR description following the template structure (if available) with commit summary and structured content
7. Uses `gh pr create` to create the pull request

### Example Workflow

```bash
# Create a feature branch
git checkout -b feature/new-functionality

# Make some changes
echo "console.log('Hello, World!');" > hello.js

# Stage the changes
git add hello.js

# Let AI generate and commit with a conventional message
institutionalized commit

# Make additional changes and commits as needed
# ... more development work ...

# When ready, create a pull request
institutionalized pr

# Or create a draft PR to share work-in-progress
institutionalized pr --draft

# Preview the PR without creating it
institutionalized pr --dry-run
```

The tool will analyze your changes and might generate a commit message like:

```
feat: add hello world JavaScript example

Add a simple JavaScript file that logs "Hello, World!" to the console.
```

When creating a PR, it will generate a comprehensive description including:

- Summary of all commits in the branch
- Structured sections for changes made and testing
- Automatic formatting for easy review

## Conventional Commits

The tool follows the [Conventional Commits](https://www.conventionalcommits.org/) specification, generating messages in the format:

```
<type>: <description>

[optional body]
```

Common types include:

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Build process or auxiliary tool changes

## Requirements

- Go 1.24+ for building from source
- Git repository (the tool must be run within a git repository)
- OpenAI API key, Google Gemini API key, or Anthropic Claude API key (for commit message generation)
- GitHub CLI (`gh`) installed and authenticated (for PR creation)
- Staged changes (use `git add` to stage files before running commit command)

## License

MIT License - see [LICENSE](LICENSE) file for details.
