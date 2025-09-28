# Installation Guide

This guide covers all the ways to install `institutionalized` so you can use it from anywhere on your system.

## Quick Install (Recommended)

### Option 1: Install from Source with Go

If you have Go installed, this is the fastest way to get started:

```bash
go install github.com/IanKnighton/institutionalized@latest
```

This will:
- Download and build the latest version automatically
- Install the binary to your `$GOPATH/bin` directory
- Make `institutionalized` available from anywhere (if `$GOPATH/bin` is in your `$PATH`)

### Option 2: Download Pre-built Binary

*Note: Pre-built binaries are not currently available. Please use one of the other installation methods.*

## Manual Installation Methods

### Build from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/IanKnighton/institutionalized.git
   cd institutionalized
   ```

2. **Build the binary:**
   ```bash
   make build
   ```
   
   Or manually with Go:
   ```bash
   go build -o institutionalized .
   ```

3. **Install globally using one of these methods:**

#### Method A: Use Make Install (Recommended)
```bash
make install
```

This installs the binary to your `$GOPATH/bin` directory.

#### Method B: Copy to System PATH

**On Linux/macOS:**
```bash
# Copy to /usr/local/bin (requires sudo)
sudo cp institutionalized /usr/local/bin/

# Or copy to ~/bin (user-only, no sudo required)
mkdir -p ~/bin
cp institutionalized ~/bin/

# Make sure ~/bin is in your PATH by adding this to ~/.bashrc or ~/.zshrc:
export PATH="$HOME/bin:$PATH"
```

**On Windows:**
```cmd
# Copy to a directory in your PATH, or create a new directory
mkdir C:\tools\institutionalized
copy institutionalized.exe C:\tools\institutionalized\

# Add C:\tools\institutionalized to your PATH environment variable
# via System Properties > Environment Variables
```

#### Method C: Create a Symlink

**On Linux/macOS:**
```bash
# Create a symlink in /usr/local/bin
sudo ln -s /path/to/institutionalized/institutionalized /usr/local/bin/institutionalized

# Or in ~/bin
mkdir -p ~/bin
ln -s /path/to/institutionalized/institutionalized ~/bin/institutionalized
export PATH="$HOME/bin:$PATH"
```

## Setting up your PATH

After installation, you may need to ensure the binary is in your system's PATH:

### Check if it's working

Test that the installation worked:
```bash
institutionalized --help
institutionalized version
```

### If the command is not found

#### Linux/macOS (Bash/Zsh)

1. **Check your current PATH:**
   ```bash
   echo $PATH
   ```

2. **Add the binary location to your PATH:**
   
   Add one of these lines to your shell profile (`~/.bashrc`, `~/.zshrc`, or `~/.profile`):
   
   ```bash
   # If you used 'make install' or 'go install'
   export PATH="$GOPATH/bin:$PATH"
   
   # If you copied to ~/bin
   export PATH="$HOME/bin:$PATH"
   
   # If you installed to /usr/local/bin, it should already be in PATH
   ```

3. **Reload your shell:**
   ```bash
   source ~/.bashrc  # or ~/.zshrc
   # Or restart your terminal
   ```

#### Windows

1. **Open System Properties:**
   - Press `Win + R`, type `sysdm.cpl`, and press Enter
   - Go to the "Advanced" tab
   - Click "Environment Variables"

2. **Edit the PATH variable:**
   - In the "System Variables" section, find and select "Path"
   - Click "Edit"
   - Click "New" and add the directory containing `institutionalized.exe`
   - Click "OK" to save

3. **Restart your command prompt or PowerShell**

#### Verify GOPATH (for Go installations)

If you used `go install` or `make install`, make sure your GOPATH is set correctly:

```bash
# Check GOPATH
go env GOPATH

# If not set, add to your shell profile:
export GOPATH=$HOME/go
export PATH="$GOPATH/bin:$PATH"
```

## Installation Verification

After installation, verify everything is working:

1. **Check the command is available:**
   ```bash
   which institutionalized  # Should show the path to the binary
   institutionalized --help # Should show help output
   ```

2. **Check the version:**
   ```bash
   institutionalized version
   ```

3. **Test in a git repository:**
   ```bash
   cd /path/to/any/git/repo
   institutionalized config show  # Should show configuration
   ```

## Troubleshooting

### "command not found" error

- **Check PATH**: Make sure the directory containing the binary is in your `$PATH`
- **Verify installation**: Use `which institutionalized` or `where institutionalized` (Windows) to locate the binary
- **Permissions**: Ensure the binary has execute permissions (`chmod +x institutionalized` on Linux/macOS)

### Permission denied errors

- **Linux/macOS**: Use `chmod +x institutionalized` to make the binary executable
- **Windows**: Run your terminal as Administrator if installing to system directories

### Go-related issues

- **GOPATH not set**: Make sure `$GOPATH/bin` is in your PATH if you used Go to install
- **Go version**: Ensure you have Go 1.24+ installed (`go version`)
- **Module issues**: Try `go clean -modcache` and reinstall

### PATH issues

- **Shell profile**: Make sure you added the PATH export to the correct shell profile file
- **Reload**: Restart your terminal or run `source ~/.bashrc` (or equivalent)
- **Verify**: Use `echo $PATH` to confirm the directory is included

## Next Steps

After successful installation:

1. **Set up API keys**: Follow the [Usage section in the README](../README.md#usage) to configure your OpenAI and/or Gemini API keys
2. **Initialize configuration**: Run `institutionalized config init` to create a default configuration file
3. **Test with a repository**: Try `institutionalized commit --dry-run` in a git repository with staged changes

## Uninstalling

To remove institutionalized:

### If installed via Go
```bash
# Remove the binary
rm $(go env GOPATH)/bin/institutionalized

# Remove configuration (optional)
rm -rf ~/.config/institutionalized/
```

### If installed manually
```bash
# Remove the binary (adjust path as needed)
sudo rm /usr/local/bin/institutionalized
# or
rm ~/bin/institutionalized

# Remove configuration (optional)
rm -rf ~/.config/institutionalized/
```

### Windows
```cmd
# Remove the binary from wherever you installed it
del C:\tools\institutionalized\institutionalized.exe

# Remove configuration (optional)
rmdir /s "%USERPROFILE%\.config\institutionalized"
```