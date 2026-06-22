# claude-sandbox

Run Claude Code with `--dangerously-skip-permissions` in an isolated smolvm container. Skip all prompts, zero local trust required.

## Installation

### Pre-built Binary (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/rohan-mathure/claude-sandbox/releases):

**macOS:**
```bash
curl -L -o claude-sandbox https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-darwin-arm64
# For Intel Macs:
# curl -L -o claude-sandbox https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-darwin-amd64
chmod +x claude-sandbox
./claude-sandbox
```

**Linux:**
```bash
curl -L -o claude-sandbox https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-linux-amd64
# For ARM64:
# curl -L -o claude-sandbox https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-linux-arm64
chmod +x claude-sandbox
./claude-sandbox
```

**Windows:**
```powershell
curl -L -o claude-sandbox.exe https://github.com/rohan-mathure/claude-sandbox/releases/download/latest/claude-sandbox-windows-amd64.exe
.\claude-sandbox.exe
```

### Build from Source

Requires [Go 1.21+](https://golang.org/dl/).

```bash
git clone https://github.com/rohan-mathure/claude-sandbox.git
cd claude-sandbox
go build -o claude-sandbox .
./claude-sandbox
```

## Prerequisites

Install [smolvm](https://github.com/smol-machines/smolvm):

```bash
curl -sSL https://smolmachines.com/install.sh | bash
```

## Usage

```bash
./claude-sandbox                                 # node:lts image, current directory
./claude-sandbox -image ubuntu:22.04             # custom image
./claude-sandbox -repo /path/to/other/project   # external repo
```

Launches smolvm VM, mounts repo at `/workspace`, installs claude-code, runs with permissions skipped. Exit with `Ctrl+D` or `exit`.

## How It Works

1. **Makefile** — entry point with two optional vars:
   - `IMAGE` (default: `node:lts`) — OCI image to run
   - `REPO` (default: current dir) — repo path to mount at `/workspace`

2. **scripts/entrypoint.sh** — runs inside VM:
   - Auto-installs Node.js if missing (debian-based images only)
   - Installs `@anthropic-ai/claude-code` globally
   - Starts claude with `--dangerously-skip-permissions`

3. **API Key** — provide via env inside VM or when claude prompts (no passthrough needed)

## Examples

```bash
# Default: isolated Claude Code against current repo
make sandbox

# Use Alpine Linux
make sandbox IMAGE=alpine

# Analyze external project without trusting your machine
cd /tmp && make -f ~/Projects/claude-sandbox/Makefile sandbox REPO=/path/to/untrusted/code
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
