# Contributing to Institutionalized

Thank you for your interest in contributing to Institutionalized! We welcome contributions from the community to help make this tool better for everyone.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Coding Standards](#coding-standards)
- [Issue Guidelines](#issue-guidelines)

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct:

- Be respectful and inclusive
- Use welcoming and inclusive language
- Be collaborative
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Go 1.24+ installed
- Git
- A GitHub account
- (Optional) OpenAI API key and/or Google Gemini API key for testing

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/institutionalized.git
   cd institutionalized
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/IanKnighton/institutionalized.git
   ```

## Development Setup

### Building the Project

```bash
# Clean build (recommended first step)
make clean

# Build the binary
make build
```

### Installing Dependencies

Dependencies are managed by Go modules and will be automatically downloaded during the first build.

```bash
# Ensure dependencies are up to date
go mod tidy
```

### Environment Variables

For testing the tool's functionality, set up API keys:

```bash
# For OpenAI (ChatGPT)
export OPENAI_API_KEY="your-openai-key-here"

# For Google Gemini
export GEMINI_API_KEY="your-gemini-key-here"
```

## Making Changes

### Creating a Branch

Always create a new branch for your changes:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

### Branch Naming Conventions

- `feature/description` - for new features
- `fix/description` - for bug fixes
- `docs/description` - for documentation changes
- `refactor/description` - for code refactoring
- `test/description` - for adding or updating tests

### Making Commits

We follow the Conventional Commits specification. Use the tool itself to generate commit messages when possible:

```bash
# Stage your changes
git add .

# Use institutionalized to create the commit message
./institutionalized commit
```

If you're making changes to the commit generation functionality, you may need to create manual commits:

```bash
git commit -m "feat: add new AI provider support"
git commit -m "fix: resolve API timeout issue"
git commit -m "docs: update installation instructions"
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Format code and run tests
make check
```

### Testing Your Changes

Before submitting changes, ensure they work correctly:

```bash
# Test basic functionality
./institutionalized --help
./institutionalized version

# Test with a real git repository (requires API key)
# 1. Make some changes in a test directory
# 2. Stage the changes: git add .
# 3. Test commit functionality: ./institutionalized commit --dry-run
# 4. Test full functionality: ./institutionalized commit
```

### Manual Testing Scenarios

When testing changes, consider these scenarios:

- No staged changes (should show error)
- Various types of code changes (new features, bug fixes, documentation)
- Network connectivity issues
- Invalid API keys
- Different git repository states

## Submitting Changes

### Before Submitting

1. Ensure your code builds successfully: `make build`
2. Run the test suite: `make check`
3. Test the functionality manually
4. Update documentation if necessary
5. Ensure your commits follow conventional commit format

### Pull Request Process

1. Push your branch to your fork:

   ```bash
   git push origin your-branch-name
   ```

2. Create a pull request on GitHub with:

   - A clear title describing the change
   - A detailed description of what changed and why
   - Steps to test the changes
   - Screenshots or examples if applicable

3. Address any feedback from code reviews

4. Once approved, your changes will be merged

### Pull Request Guidelines

- Keep changes focused and atomic
- Include tests for new functionality
- Update documentation for user-facing changes
- Ensure backwards compatibility when possible
- Reference related issues using keywords (e.g., "Fixes #123")

## Coding Standards

### Go Style Guide

We follow standard Go conventions:

- Use `gofmt` to format code (run `make fmt`)
- Follow effective Go practices
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Handle errors appropriately

### Code Organization

- Follow the existing package structure
- Keep functions focused and single-purpose
- Use dependency injection where appropriate
- Separate concerns (API calls, Git operations, user interaction)

### Error Handling

- Always handle errors explicitly
- Provide meaningful error messages to users
- Use wrapped errors for better debugging: `fmt.Errorf("context: %w", err)`
- Log errors at appropriate levels

### Configuration

- Use the existing configuration system in `internal/config/`
- Add new configuration options through the config package
- Provide sensible defaults
- Document configuration options

## Issue Guidelines

### Reporting Bugs

When reporting bugs, please use the bug report template and include:

- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment information (OS, Go version, etc.)
- Relevant logs or error messages

### Requesting Features

For feature requests, please use the feature request template and include:

- Clear description of the desired feature
- Use case and motivation
- Proposed implementation approach
- Consider backwards compatibility

### General Issues

- Search existing issues before creating new ones
- Use clear, descriptive titles
- Provide as much relevant information as possible
- Tag issues appropriately
- Be patient and respectful in discussions

## Development Tips

### Useful Make Targets

- `make build` - Build the binary with version info
- `make clean` - Remove build artifacts and clean Go cache
- `make test` - Run tests
- `make check` - Format, tidy dependencies, and test
- `make install` - Install binary to GOPATH/bin
- `make lint` - Run golangci-lint (if installed)

### Debugging

- Use `--dry-run` flag for testing without API calls
- Set `--verbose` for detailed output
- Check logs in `~/.config/institutionalized/` directory
- Use Go's built-in debugger or print statements for development

### Working with APIs

When working with AI provider integrations:

- Test with both providers when possible
- Handle rate limiting and timeouts gracefully
- Mock API responses for unit tests
- Respect API usage guidelines

## Getting Help

- Check existing issues and discussions
- Ask questions in GitHub issues with the "question" label
- Review the codebase and existing patterns
- Refer to the project documentation

Thank you for contributing to Institutionalized! 🎉
