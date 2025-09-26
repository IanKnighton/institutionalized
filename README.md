# institutionalized

A simple tool that uses an LLM to create commit and PR messages based on git status.

## Overview

Institutionalized is a Go CLI tool that analyzes your staged git changes and uses ChatGPT to generate conventional commit messages, then prompts you to confirm before committing the changes.

## Features

- 🤖 **AI-powered commit messages**: Uses OpenAI's ChatGPT to generate meaningful commit messages
- 📝 **Conventional Commits**: Follows the Conventional Commits specification by default
- 🔍 **Smart analysis**: Analyzes your staged git changes to understand the context
- 🛡️ **User confirmation**: Always asks for confirmation before committing
- 🔧 **Flexible configuration**: Supports API key via environment variable or command flag
- 😊 **Emoji support**: Optional emoji prefixes for commit types (✨ feat, 🐛 fix, etc.)

## Installation

### Build from source

```bash
git clone https://github.com/IanKnighton/institutionalized.git
cd institutionalized
go build -o institutionalized .
```

## Usage

### Setup

First, you need an OpenAI API key. You can get one from [OpenAI's platform](https://platform.openai.com/api-keys).

Set your API key either as an environment variable:

```bash
export OPENAI_API_KEY="your-api-key-here"
```

Or pass it as a flag:

```bash
institutionalized commit --api-key "your-api-key-here"
```

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

Analyzes staged changes and generates a conventional commit message using ChatGPT.

**Flags:**
- `--api-key, -k`: OpenAI API key (can also be set via `OPENAI_API_KEY` environment variable)
- `--emoji`: Use emoji in commit messages (overrides config file setting)
- `--dry-run`: Show staged changes without calling API or committing (useful for testing)

**Examples:**

```bash
# Basic usage with environment variable
export OPENAI_API_KEY="your-key"
institutionalized commit

# Using API key flag
institutionalized commit --api-key "your-key"

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

**Examples:**

```bash
# View current configuration
institutionalized config show

# Enable emoji support
institutionalized config set use_emoji true

# Create default config file
institutionalized config init
```

### Example Workflow

```bash
# Make some changes
echo "console.log('Hello, World!');" > hello.js

# Stage the changes
git add hello.js

# Let AI generate and commit with a conventional message
institutionalized commit
```

The tool will analyze your changes and might generate a commit message like:
```
feat: add hello world JavaScript example

Add a simple JavaScript file that logs "Hello, World!" to the console.
```

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
- OpenAI API key
- Staged changes (use `git add` to stage files before running)

## License

MIT License - see [LICENSE](LICENSE) file for details.