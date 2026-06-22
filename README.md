# claude-sandbox

Run Claude Code with `--dangerously-skip-permissions` in an isolated smolvm container. Skip all prompts, zero local trust required.

## Installation

### Pre-built Binary (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/rohan-mathure/claude-sandbox/releases) and add to your PATH:

**macOS:**
```bash
# Apple Silicon (M1/M2/M3)
curl -L https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-darwin-arm64 \
  -o /usr/local/bin/claude-sandbox && chmod +x /usr/local/bin/claude-sandbox

# Intel Mac
curl -L https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-darwin-amd64 \
  -o /usr/local/bin/claude-sandbox && chmod +x /usr/local/bin/claude-sandbox
```

**Linux:**
```bash
# x86_64
curl -L https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-linux-amd64 \
  -o ~/.local/bin/claude-sandbox && chmod +x ~/.local/bin/claude-sandbox

# ARM64
curl -L https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-linux-arm64 \
  -o ~/.local/bin/claude-sandbox && chmod +x ~/.local/bin/claude-sandbox

# Ensure ~/.local/bin is in PATH
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc && source ~/.bashrc
```

**Windows:**
```powershell
# Download to a directory in PATH (e.g., C:\Users\YourName\AppData\Local\Programs\bin)
curl -L -o "$env:USERPROFILE\AppData\Local\Programs\bin\claude-sandbox.exe" `
  https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-windows-amd64.exe
```

Then use from anywhere:
```bash
claude-sandbox
```

### Build from Source

Requires [Go 1.21+](https://golang.org/dl/).

```bash
git clone https://github.com/rohan-mathure/claude-sandbox.git
cd claude-sandbox
go build -o claude-sandbox .

# Add to PATH
sudo mv claude-sandbox /usr/local/bin/  # macOS/Linux
# or
move-item claude-sandbox $env:USERPROFILE\AppData\Local\Programs\bin\  # Windows
```

## Prerequisites

Install [smolvm](https://github.com/smol-machines/smolvm):

```bash
curl -sSL https://smolmachines.com/install.sh | bash
```

## Usage

```bash
claude-sandbox                                   # node:lts image, current directory
claude-sandbox -image ubuntu:22.04               # custom image
claude-sandbox -repo /path/to/other/project     # external repo
```

Launches smolvm VM, mounts repo at `/workspace`, installs claude-code, runs with permissions skipped. Exit with `Ctrl+D` or `exit`.

## How It Works

1. **claude-sandbox binary** — Go wrapper with embedded scripts:
   - Accepts `-image` flag (default: `node:lts`) for OCI image
   - Accepts `-repo` flag (default: current dir) to mount at `/workspace`
   - Extracts scripts to temp directory
   - Executes `smolvm machine run` with proper mounts and entrypoint

2. **scripts/entrypoint.sh** — runs inside VM:
   - Auto-installs Node.js if missing (debian-based images only)
   - Installs `@anthropic-ai/claude-code` globally
   - Runs Claude as non-root user with `--dangerously-skip-permissions`

3. **API Key** — provide via env inside VM or when claude prompts (no passthrough needed)

## Examples

```bash
# Default: isolated Claude Code against current repo
claude-sandbox

# Use Alpine Linux
claude-sandbox -image alpine:latest

# Analyze external project without trusting your machine
claude-sandbox -repo /path/to/untrusted/code

# Use custom Node image
claude-sandbox -image node:18
```

## Release Process

Releases are automated via GitHub Actions. To create a new release:

```bash
git tag v0.1.0  # Use semantic versioning
git push origin v0.1.0
```

This triggers the release workflow which:
1. Builds binaries for macOS (amd64/arm64), Linux (amd64/arm64), and Windows (amd64)
2. Creates a GitHub Release with all binaries as assets
3. Verifies each binary is executable

## Limitations

- Non-Debian images must have `/bin/sh` and network access
- Alpine, musl-based images need manual Node install (or pre-built image)
- smolvm macOS: requires OS 11+, Linux: requires KVM
