# Configuration Guide

This guide covers all configuration options and how to manage them in `institutionalized`.

## Configuration Overview

Institutionalized supports a YAML configuration file located at `~/.config/institutionalized/config.yaml`. The configuration file allows you to customize the behavior of the tool, including AI provider settings, emoji preferences, and timeout configurations.

## Setup and API Keys

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

## Available Configuration Options

### Core Settings

- **`use_emoji`**: Enable/disable emoji prefixes in commit messages (default: `false`)
  - When enabled, commit messages will include appropriate emoji based on the commit type
  - Can be overridden per-command using the `--emoji` flag

### Provider Settings

- **`providers.openai.enabled`**: Enable/disable OpenAI ChatGPT provider (default: `true`)
- **`providers.gemini.enabled`**: Enable/disable Google Gemini provider (default: `true`)
- **`providers.claude.enabled`**: Enable/disable Anthropic Claude provider (default: `true`)
- **`providers.priority`**: Which provider to try first when multiple are available (default: `"openai"`)
  - Valid values: `"openai"`, `"gemini"`, `"claude"`
- **`providers.delay_threshold`**: Maximum seconds to wait for a provider response before trying fallback (default: `10`, range: 1-300)

## Managing Configuration

### View Current Configuration

Display your current configuration settings:

```bash
institutionalized config show
```

This will show all current settings and the location of your configuration file.

### Create Default Configuration File

Create a default configuration file with all standard settings:

```bash
institutionalized config init
```

This creates `~/.config/institutionalized/config.yaml` with default values.

### Set Configuration Values

Use the `config set` command to update individual settings:

```bash
# Basic syntax
institutionalized config set <key> <value>
```

### Configuration Examples

#### Emoji Settings

```bash
# Enable emoji in commit messages
institutionalized config set use_emoji true

# Disable emoji in commit messages
institutionalized config set use_emoji false
```

#### Provider Management

```bash
# Enable/disable providers
institutionalized config set providers.openai.enabled true
institutionalized config set providers.gemini.enabled false
institutionalized config set providers.claude.enabled true

# Set provider priority
institutionalized config set providers.priority claude
institutionalized config set providers.priority gemini
institutionalized config set providers.priority openai

# Adjust timeout threshold
institutionalized config set providers.delay_threshold 20
```

#### Common Configuration Scenarios

**Use only OpenAI:**
```bash
institutionalized config set providers.openai.enabled true
institutionalized config set providers.gemini.enabled false
institutionalized config set providers.claude.enabled false
institutionalized config set providers.priority openai
```

**Use Claude as primary with Gemini fallback:**
```bash
institutionalized config set providers.claude.enabled true
institutionalized config set providers.gemini.enabled true
institutionalized config set providers.openai.enabled false
institutionalized config set providers.priority claude
```

**Enable all providers with Gemini priority:**
```bash
institutionalized config set providers.openai.enabled true
institutionalized config set providers.gemini.enabled true
institutionalized config set providers.claude.enabled true
institutionalized config set providers.priority gemini
```

## Emoji Support

When emoji support is enabled (`use_emoji: true`), commit messages will be prefixed with appropriate emoji based on the commit type:

| Commit Type | Emoji | Description |
|-------------|-------|-------------|
| `feat` | ✨ | New features |
| `fix` | 🐛 | Bug fixes |
| `docs` | 📚 | Documentation changes |
| `style` | 💄 | Code style changes |
| `refactor` | ♻️ | Code refactoring |
| `test` | ✅ | Adding or updating tests |
| `chore` | 🔧 | Build process or auxiliary tool changes |
| `perf` | ⚡ | Performance improvements |
| `ci` | 👷 | CI/CD changes |
| `build` | 🏗️ | Build system changes |
| `revert` | ⏪ | Reverting changes |

### Example Output

**With emoji enabled:**
```
✨ feat: add user authentication system

Implement JWT-based authentication with login and signup endpoints.
```

**With emoji disabled:**
```
feat: add user authentication system

Implement JWT-based authentication with login and signup endpoints.
```

### Override Emoji Setting Per Command

You can override the emoji setting for individual commands:

```bash
# Use emoji for this commit (overrides config)
institutionalized commit --emoji

# Don't use emoji for this commit (overrides config)
institutionalized commit --emoji=false
```

## Provider System

### How Provider Selection Works

1. **Primary Provider**: The tool first tries the provider specified in `providers.priority`
2. **Fallback Providers**: If the primary provider fails or times out, the tool tries other enabled providers
3. **Timeout Handling**: Each provider gets `providers.delay_threshold` seconds to respond before fallback kicks in

### Provider Models Used

- **OpenAI**: Uses `gpt-3.5-turbo` model
- **Gemini**: Uses `gemini-1.5-flash` model
- **Claude**: Uses `claude-3-haiku-20240307` model

### Error Handling

When no providers are available or configured, you'll see:
```
Error: no LLM providers available. Please set OPENAI_API_KEY, GEMINI_API_KEY, or CLAUDE_API_KEY environment variable
```

## Configuration File Location

The configuration file is stored at:
- **Linux/macOS**: `~/.config/institutionalized/config.yaml`
- **Windows**: `%USERPROFILE%\.config\institutionalized\config.yaml`

### Sample Configuration File

```yaml
use_emoji: false
providers:
  openai:
    enabled: true
  gemini:
    enabled: true
  claude:
    enabled: true
  priority: openai
  delay_threshold: 10
```

## Advanced Configuration

### Environment Variable Priority

Environment variables take precedence over configuration file settings for API keys:
- `OPENAI_API_KEY` - OpenAI API key
- `GEMINI_API_KEY` - Gemini API key
- `CLAUDE_API_KEY` - Claude API key

### Command-Line Overrides

Some settings can be overridden via command-line flags:
- `--emoji` / `--emoji=false` - Override emoji setting for the current command
- `--api-key` (deprecated) - Override OpenAI API key (use environment variable instead)

### Timeout Configuration

The `delay_threshold` setting controls how long to wait for each provider:
- **Minimum**: 1 second
- **Maximum**: 300 seconds (5 minutes)
- **Default**: 10 seconds
- **Recommended**: 10-30 seconds for most use cases

Higher values give providers more time to respond but may slow down the fallback process.

## Troubleshooting Configuration

### Configuration Not Found

If you see "using defaults" when running `institutionalized config show`:
```bash
# Create the default configuration file
institutionalized config init
```

### Invalid Configuration Values

The tool validates configuration values. Common errors:

**Invalid provider priority:**
```
Error: invalid value for providers.priority: invalid (expected openai/gemini/claude)
```

**Invalid timeout:**
```
Error: invalid value for providers.delay_threshold: 400 (expected 1-300 seconds)
```

**Invalid boolean:**
```
Error: invalid value for use_emoji: maybe (expected true/false)
```

### Configuration File Permissions

If you encounter permission errors:
```bash
# Fix configuration directory permissions
chmod 755 ~/.config/institutionalized/
chmod 644 ~/.config/institutionalized/config.yaml
```

### Reset Configuration

To reset to defaults:
```bash
# Remove existing configuration
rm ~/.config/institutionalized/config.yaml

# Create new default configuration
institutionalized config init
```

## Best Practices

### Recommended Settings

**For reliability (multiple providers with failover):**
```bash
institutionalized config set providers.openai.enabled true
institutionalized config set providers.gemini.enabled true
institutionalized config set providers.claude.enabled true
institutionalized config set providers.priority openai
institutionalized config set providers.delay_threshold 15
```

**For speed (single provider, shorter timeout):**
```bash
institutionalized config set providers.claude.enabled true
institutionalized config set providers.openai.enabled false
institutionalized config set providers.gemini.enabled false
institutionalized config set providers.delay_threshold 5
```

**For consistent style (with emoji):**
```bash
institutionalized config set use_emoji true
institutionalized config set providers.priority gemini  # Often good at creative tasks
```

### Security Considerations

- **Never commit API keys** to your repository
- **Use environment variables** instead of hardcoding keys
- **Rotate API keys regularly** as per your organization's security policy
- **Monitor API usage** to detect unauthorized access

### Performance Tips

- **Enable fewer providers** for faster response times
- **Set appropriate timeout values** based on your network and provider reliability
- **Choose the right primary provider** based on your use case (Claude for analysis, Gemini for creativity, OpenAI for reliability)

## Migration Guide

### From Version 1.x (OpenAI + Gemini only)

If you're upgrading from a version that only supported OpenAI and Gemini:

1. **Your existing configuration will continue to work**
2. **Claude will be automatically enabled** with default settings
3. **Update your configuration** to take advantage of Claude:

```bash
# Check current settings
institutionalized config show

# Optionally set Claude as primary
institutionalized config set providers.priority claude

# Or disable Claude if not needed
institutionalized config set providers.claude.enabled false
```

Your existing `providers.priority` values (`openai` or `gemini`) will continue to work unchanged.