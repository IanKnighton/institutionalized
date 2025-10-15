# Homebrew Tap Setup Instructions

This document explains how to set up the Homebrew tap repository for `institutionalized`.

## One-Time Setup

### 1. Create the Homebrew Tap Repository

1. Create a new GitHub repository named `homebrew-tap` in your account
2. Repository URL should be: `https://github.com/IanKnighton/homebrew-tap`
3. Make it public
4. Initialize with a README (optional)

### 2. Set Up GitHub Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Give it a name like "GoReleaser Homebrew Tap"
4. Select scopes:
   - `public_repo` (Recommended: access to public repositories only)
   - **Alternatively, use a [fine-grained personal access token](https://github.com/settings/tokens?type=beta) limited to the `homebrew-tap` repository with "Contents: Read and write" permission.**
5. Generate and copy the token

### 3. Add Secret to Repository

1. Go to the `institutionalized` repository settings
2. Navigate to Secrets and variables → Actions
3. Click "New repository secret"
4. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
5. Value: Paste the token from step 2
6. Click "Add secret"

## How It Works

When a new release is created:

1. The CI/CD workflow runs GoReleaser
2. GoReleaser builds binaries for multiple platforms
3. GoReleaser creates/updates the Homebrew formula in `homebrew-tap` repository
4. The formula is automatically committed to `homebrew-tap/Formula/institutionalized.rb`

## Testing the Tap

After the first release with Homebrew support:

```bash
# Tap the repository
brew tap ianknighton/tap

# Install
brew install ianknighton/tap/institutionalized

# Verify
institutionalized version
```

## Manual Formula Creation (Optional)

If you need to manually create or update the formula, create a file at `Formula/institutionalized.rb` in the `homebrew-tap` repository:

> **Important:** In the example below, replace **every** occurrence of `VERSION` with the actual release tag (e.g., `1.2.3`), and **every** occurrence of `CHECKSUM` with the computed SHA256 for each binary. Failing to update all placeholders will cause the formula to fail.
```ruby
class Institutionalized < Formula
  desc "CLI tool that uses LLMs to create commit and PR messages based on git status"
  homepage "https://github.com/IanKnighton/institutionalized"
  version "VERSION"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/IanKnighton/institutionalized/releases/download/vVERSION/institutionalized-VERSION-macos-arm64.tar.gz"
      sha256 "CHECKSUM"
    end
    if Hardware::CPU.intel?
      url "https://github.com/IanKnighton/institutionalized/releases/download/vVERSION/institutionalized-VERSION-macos-amd64.tar.gz"
      sha256 "CHECKSUM"
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/IanKnighton/institutionalized/releases/download/vVERSION/institutionalized-VERSION-linux-arm64.tar.gz"
      sha256 "CHECKSUM"
    end
    if Hardware::CPU.intel?
      url "https://github.com/IanKnighton/institutionalized/releases/download/vVERSION/institutionalized-VERSION-linux-amd64.tar.gz"
      sha256 "CHECKSUM"
    end
  end

  def install
    bin.install "institutionalized"
  end

  test do
    system "#{bin}/institutionalized", "version"
  end
end
```

However, GoReleaser will handle this automatically.

## Updating the Formula

GoReleaser automatically updates the formula on each release. No manual intervention needed!

## Troubleshooting

### Token Issues
- Make sure the token has `repo` scope
- Verify the secret name is exactly `HOMEBREW_TAP_GITHUB_TOKEN`
- Check the token hasn't expired

### Formula Not Updating
- Check the GitHub Actions logs for errors
- Verify the `homebrew-tap` repository exists and is public
- Ensure GoReleaser has write access to the tap repository

### Installation Issues
Users should tap the repository first:
```bash
brew tap ianknighton/tap
brew install ianknighton/tap/institutionalized
```

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GoReleaser Homebrew Documentation](https://goreleaser.com/customization/homebrew/)
- [Creating a Tap](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
