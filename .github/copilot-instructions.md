# Copilot Instructions for Institutionalized

## Repository Overview

**institutionalized** is a Go CLI tool that analyzes staged git changes and uses OpenAI's ChatGPT to generate conventional commit messages, then prompts users to confirm before committing the changes. The tool follows the Conventional Commits specification and helps developers create consistent, meaningful commit messages automatically.

## High-Level Repository Information

- **Size**: Small, focused CLI application (~10 files)
- **Type**: Go CLI application using Cobra framework
- **Language**: Go 1.24+
- **Key Dependencies**: 
  - `github.com/spf13/cobra` v1.10.1 (CLI framework)
  - `github.com/spf13/pflag` v1.0.9 (command-line flag parsing)
- **Target Runtime**: Cross-platform CLI binary
- **External APIs**: OpenAI ChatGPT API (gpt-3.5-turbo model)

## Build and Validation Instructions

### Prerequisites
- Go 1.24+ installed
- Git repository (tool must be run within a git repo)
- OpenAI API key for full functionality
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
- Currently no test files exist (shows `[no test files]` message)
- `make check` runs: `go fmt`, `go mod tidy`, then `go test`
- All commands should complete successfully with exit code 0

### Available Make Targets
- `make build` - Build the binary with version info
- `make clean` - Remove build artifacts and clean Go cache
- `make test` - Run tests (currently none exist)
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

# Test with staged changes (requires OpenAI API key)
git add <files>
export OPENAI_API_KEY="your-key"
./institutionalized commit --dry-run  # Safe test without API call

# Full functionality test
./institutionalized commit  # Requires valid API key
```

**Runtime Notes**:
- Tool must be run from within a git repository
- Requires staged changes (`git add` files first)
- Without staged changes, returns error: "no staged changes found"
- `--dry-run` flag shows staged changes without calling API or committing

## Project Layout and Architecture

### Core Architecture
The application follows a standard Cobra CLI pattern with command-based architecture:

```
/
├── main.go                 # Entry point, calls cmd.Execute()
├── cmd/                    # Command implementations
│   ├── root.go            # Root command setup, persistent flags
│   ├── commit.go          # Main commit command logic
│   └── version.go         # Version command
├── Makefile               # Build automation
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── README.md              # User documentation
├── LICENSE                # MIT license
└── .gitignore            # Git ignore rules
```

### Key Source Files

**main.go** (15 lines):
- Simple entry point that calls `cmd.Execute()`
- Handles error output and exit codes

**cmd/root.go** (21 lines):
- Defines root command with description
- Sets up persistent `--api-key` flag
- Uses Cobra's Execute() pattern

**cmd/commit.go** (202 lines):
- Main functionality: analyzes git diff and generates commit messages
- Handles OpenAI API integration
- Implements user confirmation flow
- Key functions:
  - `runCommit()` - Main command logic
  - `getStagedDiff()` - Gets git staged changes
  - `generateCommitMessage()` - Calls OpenAI API
  - `askForConfirmation()` - User interaction
  - `commitChanges()` - Executes git commit

**cmd/version.go** (22 lines):
- Simple version command
- Version variable set via ldflags during build

### Configuration Files

**Makefile**:
- Defines build targets and automation
- Sets version to 1.0.0
- No configuration for linting (golangci-lint not installed by default)

**go.mod**:
- Go 1.24.7 requirement
- Minimal dependencies (only Cobra framework)

**.gitignore**:
- Standard Go ignore patterns
- Excludes built binary (`institutionalized`)
- Includes IDE and OS-specific ignores

### Dependencies and External Integrations

**External APIs**:
- OpenAI ChatGPT API (https://api.openai.com/v1/chat/completions)
- Uses gpt-3.5-turbo model
- Requires Bearer token authentication

**Git Integration**:
- Uses `git` command directly via `os/exec`
- Checks for git repository with `.git` directory
- Gets staged changes via `git diff --cached`
- Commits changes via `git commit -m`

### Development Workflow

**No CI/CD Pipeline**: Currently no GitHub Actions or other automation exists.

**Manual Validation Steps**:
1. Run `make check` to format, tidy, and test
2. Build with `make build`
3. Test basic functionality with `./institutionalized --help`
4. Test with actual git repository and staged changes
5. Verify conventional commit message format

**Common Development Patterns**:
- Uses standard Go error handling with wrapped errors
- HTTP client for API calls with proper error handling
- Command-line interaction with `bufio.Scanner`
- JSON marshaling/unmarshaling for API communication

### Files in Repository Root
```
.git/          # Git repository data
.gitignore     # Git ignore rules (35 lines)
LICENSE        # MIT license (21 lines)
Makefile       # Build automation (40 lines)
README.md      # Documentation (131 lines)
cmd/           # Command implementations directory
go.mod         # Go module file (10 lines)
go.sum         # Dependency checksums (10 lines)
main.go        # Application entry point (15 lines)
```

## Important Implementation Notes

**API Key Handling**:
- Accepts API key via `--api-key` flag or `OPENAI_API_KEY` environment variable
- Environment variable takes precedence over flag
- No API key validation until actual API call

**Error Conditions**:
- Must be run in git repository
- Must have staged changes
- Requires valid OpenAI API key for full functionality
- Network connectivity required for API calls

**Conventional Commits**:
- Generates messages following Conventional Commits specification
- Common types: feat, fix, docs, style, refactor, test, chore
- Format: `<type>: <description>` with optional body

## Agent Guidelines

**Trust these instructions** - they have been validated through comprehensive testing. Only search for additional information if:
- Instructions are incomplete for your specific task
- Instructions appear to be incorrect based on your findings
- You need details about code not covered in this overview

**Always validate changes** by running `make check` and testing the built binary before considering the work complete.