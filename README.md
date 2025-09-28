# institutionalized

A simple tool that uses an LLM to create commit and PR messages based on git status.

## Overview

Institutionalized is a Go CLI tool that analyzes your staged git changes and uses AI providers (OpenAI ChatGPT or Google Gemini) to generate conventional commit messages, then prompts you to confirm before committing the changes. It can also create comprehensive pull requests using GitHub CLI.

## Features

- 🤖 **AI-powered commit messages**: Uses OpenAI's ChatGPT or Google Gemini to generate meaningful commit messages
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

The fastest way to install and use `institutionalized` from anywhere:

```bash
go install github.com/IanKnighton/institutionalized@latest
```

This installs the binary to your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your `$PATH` to use `institutionalized` from anywhere.

### Other Installation Methods

For detailed installation instructions including manual installation, PATH setup, and troubleshooting, see our comprehensive [Installation Guide](docs/installation.md).

**Available methods:**
- Install via Go (recommended)
- Build from source with global installation
- Manual binary placement and PATH configuration
- Platform-specific instructions for Linux, macOS, and Windows

## Usage

### Setup

You need an API key from one or both of the supported providers:

**OpenAI**: Get your API key from [OpenAI's platform](https://platform.openai.com/api-keys)
**Google Gemini**: Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)

Set your API key(s) as environment variables:

```bash
# For OpenAI (ChatGPT)
export OPENAI_API_KEY="your-openai-key-here"

# For Google Gemini
export GEMINI_API_KEY="your-gemini-key-here"

# You can set both - the tool will use them based on your configuration
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

Institutionalized supports a configuration file located at `~/.config/institutionalized/config.yaml`. The configuration file allows you to customize the behavior of the tool.

#### Available Configuration Options

- `use_emoji`: Enable/disable emoji prefixes in commit messages (default: `false`)
- `providers.openai.enabled`: Enable/disable OpenAI ChatGPT provider (default: `true`)
- `providers.gemini.enabled`: Enable/disable Google Gemini provider (default: `true`)
- `providers.priority`: Which provider to try first when both are available - `openai` or `gemini` (default: `openai`)
- `providers.delay_threshold`: Maximum seconds to wait for a provider response before trying fallback (default: `10`, range: 1-300)

#### Managing Configuration

**View current configuration:**
```bash
institutionalized config show
```

**Create a default configuration file:**
```bash
institutionalized config init
```

**Set configuration values:**
```bash
# Enable emoji in commit messages
institutionalized config set use_emoji true

# Disable emoji in commit messages
institutionalized config set use_emoji false

# Set Gemini as the primary provider
institutionalized config set providers.priority gemini

# Disable OpenAI provider (use only Gemini)
institutionalized config set providers.openai.enabled false

# Set delay threshold to 20 seconds
institutionalized config set providers.delay_threshold 20
```

#### Emoji Support

When emoji support is enabled (`use_emoji: true`), commit messages will be prefixed with appropriate emoji:

- ✨ `feat`: New features
- 🐛 `fix`: Bug fixes
- 📚 `docs`: Documentation changes
- 💄 `style`: Code style changes
- ♻️ `refactor`: Code refactoring
- ✅ `test`: Adding or updating tests
- 🔧 `chore`: Build process or auxiliary tool changes
- ⚡ `perf`: Performance improvements
- 👷 `ci`: CI/CD changes
- 🏗️ `build`: Build system changes
- ⏪ `revert`: Reverting changes

**Example with emoji enabled:**
```
✨ feat: add user authentication system

Implement JWT-based authentication with login and signup endpoints.
```

**Override emoji setting per command:**
```bash
# Use emoji for this commit (overrides config)
institutionalized commit --emoji

# Don't use emoji for this commit (overrides config)
institutionalized commit --emoji=false
```

### Commands

#### `institutionalized commit`

Analyzes staged changes and generates a conventional commit message using available AI providers (OpenAI ChatGPT or Google Gemini).

**Flags:**
- `--api-key, -k`: OpenAI API key (deprecated: use `OPENAI_API_KEY` environment variable)
- `--emoji`: Use emoji in commit messages (overrides config file setting)
- `--dry-run`: Show staged changes without calling API or committing (useful for testing)

**Examples:**

```bash
# Basic usage with environment variables
export OPENAI_API_KEY="your-openai-key"
export GEMINI_API_KEY="your-gemini-key"
institutionalized commit

# Using only OpenAI (if you have both keys but want to use only OpenAI)
institutionalized config set providers.gemini.enabled false
institutionalized commit

# Dry run to see what changes would be analyzed
institutionalized commit --dry-run

# Use emoji for this specific commit
institutionalized commit --emoji
```

#### `institutionalized config`

Manage configuration settings for institutionalized.

**Subcommands:**
- `show`: Display current configuration values
- `set <key> <value>`: Set a configuration value
- `init`: Create a default configuration file

**Available configuration keys:**
- `use_emoji`: Enable/disable emoji support (true/false)
- `providers.openai.enabled`: Enable/disable OpenAI provider (true/false)
- `providers.gemini.enabled`: Enable/disable Gemini provider (true/false)
- `providers.priority`: Set provider priority (openai/gemini)
- `providers.delay_threshold`: Set timeout in seconds (1-300)

**Examples:**

```bash
# View current configuration
institutionalized config show

# Enable emoji support
institutionalized config set use_emoji true

# Set Gemini as primary provider
institutionalized config set providers.priority gemini

# Create default config file
institutionalized config init
```

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
- OpenAI API key or Google Gemini API key (for commit message generation)
- GitHub CLI (`gh`) installed and authenticated (for PR creation)
- Staged changes (use `git add` to stage files before running commit command)

## License

MIT License - see [LICENSE](LICENSE) file for details.# Test comment
