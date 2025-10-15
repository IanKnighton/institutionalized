# Copilot Instructions for Institutionalized

## Repository Overview

**institutionalized** is a Go CLI tool that analyzes staged git changes and uses AI providers (OpenAI ChatGPT, Google Gemini, or Anthropic Claude) to generate conventional commit messages, then prompts users to confirm before committing the changes. It can also create comprehensive pull requests using GitHub CLI. The tool follows the Conventional Commits specification and helps developers create consistent, meaningful commit messages and PRs automatically.

## High-Level Repository Information

- **Size**: Small, focused CLI application (~15 files, ~2500 lines total)
- **Type**: Go CLI application using Cobra framework
- **Language**: Go 1.24.7+
- **Key Dependencies**: 
  - `github.com/spf13/cobra` v1.10.1 (CLI framework)
  - `github.com/spf13/pflag` v1.0.9 (command-line flag parsing)
  - `gopkg.in/yaml.v3` v3.0.1 (YAML configuration parsing)
- **Target Runtime**: Cross-platform CLI binary
- **External APIs**: OpenAI ChatGPT API, Google Gemini API, and Anthropic Claude API with configurable provider priority

## Build and Validation Instructions

### Prerequisites
- Go 1.24.7+ installed
- Git repository (tool must be run within a git repo)
- OpenAI API key, Google Gemini API key, and/or Anthropic Claude API key for full functionality
- GitHub CLI (`gh`) installed and authenticated (for PR creation)
- **Always run commands from the repository root directory**

### Build Process
```bash
# Clean build (recommended first step)
make clean

# Build the binary
make build
```

**Build Notes**:
- Build time: ~10-30 seconds depending on network speed for dependency downloads
- Creates `institutionalized` binary in root directory
- Uses ldflags to inject version information during build
- Binary name is controlled by `BINARY_NAME` variable in Makefile

### Testing
```bash
# Run all tests
make test

# Format code and run tests (recommended)
make check
```

**Testing Notes**:
- Test files exist in `cmd/pr_test.go` covering PR template functionality
- `make check` runs: `go fmt`, `go mod tidy`, then `go test`
- All commands should complete successfully with exit code 0

### Available Make Targets
- `make build` - Build the binary with version info
- `make clean` - Remove build artifacts and clean Go cache
- `make test` - Run tests (currently PR functionality tests)
- `make install` - Install binary to GOPATH/bin
- `make fmt` - Format Go code
- `make lint` - Run golangci-lint (requires golangci-lint to be installed)
- `make tidy` - Clean up go.mod dependencies
- `make check` - Run fmt, tidy, and test in sequence
- `make all` - Clean then build (default target)

### Runtime Testing
```bash
# Test basic functionality (no API key required)
./institutionalized --help
./institutionalized version

# Test configuration commands
./institutionalized config show
./institutionalized config init  # Creates default config file

# Test with staged changes (requires API key)
git add <files>
export OPENAI_API_KEY="your-openai-key"
# OR/AND
export GEMINI_API_KEY="your-gemini-key"
# OR/AND
export CLAUDE_API_KEY="your-claude-key"

./institutionalized commit --dry-run  # Safe test without API call

# Test PR functionality (requires GitHub CLI)
./institutionalized pr --dry-run  # Preview PR without creating

# Full functionality test
./institutionalized commit  # Requires valid API key
./institutionalized pr      # Requires gh CLI authentication
```

**Runtime Notes**:
- Tool must be run from within a git repository
- Requires staged changes (`git add` files first) for commit command
- Without staged changes, returns error: "no staged changes found. Use 'git add' to stage changes first"
- `--dry-run` flag shows staged changes without calling API or committing
- Multiple AI providers supported with configurable priority and fallback
- PR functionality requires GitHub CLI (`gh`) to be installed and authenticated

## Project Layout and Architecture

### Core Architecture
The application follows a standard Cobra CLI pattern with command-based architecture:

```
/
├── main.go                 # Entry point, calls cmd.Execute()
├── cmd/                    # Command implementations
│   ├── root.go            # Root command setup, persistent flags
│   ├── commit.go          # Main commit command logic
│   ├── config.go          # Configuration management commands
│   ├── pr.go              # Pull request creation functionality
│   ├── pr_test.go         # Tests for PR template functionality
│   └── version.go         # Version command
├── docs/                   # User documentation
│   ├── configuration.md   # Detailed configuration guide
│   └── installation.md    # Installation instructions
├── internal/               # Internal packages
│   ├── config/            # Configuration management
│   │   └── config.go      # Config structure and file handling
│   └── llm/               # LLM provider implementations
│       └── providers.go   # OpenAI and Gemini provider implementations
├── Makefile               # Build automation
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── CONTRIBUTING.md        # Development guidelines and contribution process
├── README.md              # User documentation
├── LICENSE                # MIT license
└── .gitignore            # Git ignore rules
```

### Key Source Files

**main.go** (15 lines):
- Simple entry point that calls `cmd.Execute()`
- Handles error output and exit codes

**cmd/root.go** (33 lines):
- Defines root command with description
- Sets up persistent `--api-key` flag (deprecated) and `--emoji` flag
- Uses Cobra's Execute() pattern
- Loads configuration on initialization

**cmd/commit.go** (215 lines):
- Main functionality: analyzes git diff and generates commit messages
- Handles multiple AI provider integration with fallback support
- Implements user confirmation flow
- Key functions:
  - `runCommit()` - Main command logic
  - `getStagedDiff()` - Gets git staged changes
  - `generateCommitMessage()` - Calls AI providers with fallback
  - `askForConfirmation()` - User interaction
  - `commitChanges()` - Executes git commit

**cmd/config.go** (157 lines):
- Configuration management commands (show, set, init)
- Handles YAML configuration file operations
- Supports setting provider preferences, emoji settings, etc.

**cmd/pr.go** (301 lines):
- Pull request creation using GitHub CLI
- Analyzes commit history and generates comprehensive PR descriptions
- Supports draft PRs and dry-run mode
- Requires GitHub CLI authentication

**cmd/version.go** (22 lines):
- Simple version command
- Version variable set via ldflags during build

**internal/config/config.go** (131 lines):
- Configuration structure and file handling
- Supports provider configuration, emoji settings, API timeouts
- Default config creation and YAML serialization

**internal/llm/providers.go** (513 lines):
- OpenAI, Gemini, and Claude provider implementations
- Provider interface for consistent API across providers
- Handles API requests, response parsing, and error handling
- Supports both commit message and PR content generation

### Configuration Files

**Makefile**:
- Defines build targets and automation
- Sets version to 1.0.0
- No configuration for linting (golangci-lint not installed by default)

**go.mod**:
- Go 1.24.7 requirement
- Dependencies: Cobra framework, YAML parsing, standard library

**.gitignore**:
- Standard Go ignore patterns
- Excludes built binary (`institutionalized`)
- Includes IDE and OS-specific ignores

**Configuration File** (~/.config/institutionalized/config.yaml):
- YAML-based configuration for provider settings
- Controls emoji usage, provider priority, API timeouts
- Created automatically with `institutionalized config init`

### Dependencies and External Integrations

**External APIs**:
- OpenAI ChatGPT API (https://api.openai.com/v1/chat/completions)
- Google Gemini API (https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent)
- Anthropic Claude API (https://api.anthropic.com/v1/messages)
- Configurable provider priority with automatic fallback
- Uses gpt-3.5-turbo for OpenAI, gemini-1.5-flash for Gemini, and claude-3-5-sonnet for Claude
- Requires Bearer token authentication for OpenAI and Gemini, API key for Claude

**Git Integration**:
- Uses `git` command directly via `os/exec`
- Checks for git repository with `.git` directory
- Gets staged changes via `git diff --cached`
- Commits changes via `git commit -m`

**GitHub CLI Integration**:
- Uses `gh` command for PR creation
- Requires authentication via `gh auth login`
- Supports draft PRs and dry-run preview
- Generates comprehensive PR descriptions with commit analysis

### Development Workflow

**CI/CD Pipeline**: GitHub Actions workflow in `.github/workflows/ci.yml` with three jobs:
1. **test**: Runs `make check` and `make build` on Ubuntu
2. **release-tag**: Creates release tags on main branch pushes  
3. **create-release**: Uses GoReleaser to build multi-platform binaries and publish to Homebrew tap

**Manual Validation Steps**:
1. Run `make check` to format, tidy, and test
2. Build with `make build`
3. Test basic functionality with `./institutionalized --help`
4. Test configuration with `./institutionalized config show`
5. Test with actual git repository and staged changes
6. Test PR functionality with `./institutionalized pr --dry-run`
7. Verify conventional commit message format and emoji support

**Common Development Patterns**:
- Uses standard Go error handling with wrapped errors
- HTTP client for API calls with proper error handling
- Command-line interaction with `bufio.Scanner`
- JSON marshaling/unmarshaling for API communication

### Files in Repository Root
```
.git/               # Git repository data
.github/            # GitHub configuration
├── copilot-instructions.md  # This file
.gitignore          # Git ignore rules (35 lines)
LICENSE             # MIT license (21 lines)
Makefile            # Build automation (39 lines)
README.md           # Documentation (311 lines)
CONTRIBUTING.md     # Development guidelines (293 lines)
cmd/                # Command implementations directory
├── commit.go       # Commit command (215 lines)
├── config.go       # Config command (157 lines)
├── pr.go          # PR command (301 lines)
├── pr_test.go     # Tests for PR functionality (77 lines)
├── root.go        # Root command (33 lines)
└── version.go     # Version command (22 lines)
docs/               # User documentation
├── configuration.md # Configuration guide (360 lines)
└── installation.md  # Installation instructions
internal/           # Internal packages
├── config/        # Configuration management
│   └── config.go  # Config handling (131 lines)
└── llm/           # LLM providers
    └── providers.go # Provider implementations (513 lines)
go.mod              # Go module file (13 lines)
go.sum              # Dependency checksums (10 lines)
main.go             # Application entry point (15 lines)
```

## Important Implementation Notes

**API Key Handling**:
- Accepts API keys via environment variables: `OPENAI_API_KEY` and/or `GEMINI_API_KEY` and/or `CLAUDE_API_KEY`
- Legacy `--api-key` flag still supported but deprecated for OpenAI only
- Environment variables take precedence over flags
- Multiple providers can be enabled simultaneously with configurable priority
- Automatic fallback to secondary provider if primary fails or times out

**Configuration System**:
- YAML configuration file at `~/.config/institutionalized/config.yaml`
- Supports provider enable/disable, priority setting, emoji preferences
- Default configuration created with `institutionalized config init`
- Runtime flag overrides (e.g., `--emoji`) override config file settings

**Command Structure**:
- `commit`: Generate and commit changes using AI providers
- `config`: Manage configuration (show, set, init subcommands)
- `pr`: Create pull requests using GitHub CLI
- `version`: Display version information
- Global flags: `--emoji`, `--api-key` (deprecated)

**Error Conditions**:
- Must be run in git repository
- Must have staged changes for commit command
- Requires valid API key(s) for AI functionality
- Requires GitHub CLI authentication for PR functionality
- Network connectivity required for API calls

**Conventional Commits**:
- Generates messages following Conventional Commits specification
- Common types: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert
- Format: `<type>: <description>` with optional body
- Optional emoji prefix based on configuration (✨ feat, 🐛 fix, etc.)

## Agent Guidelines

**Trust these instructions** - they have been validated through comprehensive testing. Only search for additional information if:
- Instructions are incomplete for your specific task
- Instructions appear to be incorrect based on your findings
- You need details about code not covered in this overview

**Always validate changes** by running `make check` and testing the built binary before considering the work complete.