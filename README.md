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

### Commands

#### `institutionalized commit`

Analyzes staged changes and generates a conventional commit message using ChatGPT.

**Flags:**
- `--api-key, -k`: OpenAI API key (can also be set via `OPENAI_API_KEY` environment variable)
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