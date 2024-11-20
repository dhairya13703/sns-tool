# Development and Deployment Guide

## Setting Up Development Environment

### Prerequisites
- Go 1.21 or later
- Git

### Installing GoReleaser

Choose one of these methods to install GoReleaser:

1. Using Go:
```bash
go install github.com/goreleaser/goreleaser@latest
```

2. Using Homebrew (macOS/Linux):
```bash
brew install goreleaser
```

3. Direct download (manually):
   - Visit [GoReleaser Releases](https://github.com/goreleaser/goreleaser/releases)
   - Download the appropriate version for your OS/architecture
   - Extract and move to your PATH

4. Using curl (Linux/macOS):
```bash
# For macOS (Intel)
curl -L https://github.com/goreleaser/goreleaser/releases/download/v1.24.0/goreleaser_Darwin_x86_64.tar.gz | tar -xz -C /usr/local/bin/ goreleaser

# For macOS (M1/M2)
curl -L https://github.com/goreleaser/goreleaser/releases/download/v1.24.0/goreleaser_Darwin_arm64.tar.gz | tar -xz -C /usr/local/bin/ goreleaser

# For Linux (x86_64)
curl -L https://github.com/goreleaser/goreleaser/releases/download/v1.24.0/goreleaser_Linux_x86_64.tar.gz | tar -xz -C /usr/local/bin/ goreleaser
```

Verify installation:
```bash
goreleaser --version
```

## Deployment Steps

### 1. GitHub Token Setup

1. Create a GitHub Personal Access Token (PAT):
   - Go to GitHub → Settings → Developer Settings → Personal Access Tokens → Tokens (classic)
   - Click "Generate new token"
   - Name: "GORELEASER_TOKEN"
   - Select scopes:
     - `repo` (Full control of private repositories)
     - `write:packages` (Upload packages)
   - Copy the generated token

2. Set up the token locally:
```bash
# For current session
export GITHUB_TOKEN="your_token_here"

# For permanent storage (choose based on your shell):
# For bash
echo 'export GITHUB_TOKEN="your_token_here"' >> ~/.bashrc
source ~/.bashrc

# For zsh
echo 'export GITHUB_TOKEN="your_token_here"' >> ~/.zshrc
source ~/.zshrc
```

### 2. Release Process

1. Make sure all changes are committed:
```bash
git status
git add .
git commit -m "Prepare for release x.y.z"
```

2. Create a new tag:
```bash
git tag -a vx.y.z -m "Release vx.y.z"
```

3. Test the release locally:
```bash
# Dry run without publishing
goreleaser release --snapshot --clean --skip=publish
```

4. Make the actual release:
```bash
goreleaser release --clean
```

### Troubleshooting Release Issues

If you encounter tag-related issues:

1. Remove problematic tags:
```bash
# Remove local tag
git tag -d vx.y.z

# Remove remote tag (if pushed)
git push origin :refs/tags/vx.y.z
```

2. Re-create tag:
```bash
git tag -a vx.y.z -m "Release vx.y.z"
```

3. Verify tag placement:
```bash
git show vx.y.z
```
