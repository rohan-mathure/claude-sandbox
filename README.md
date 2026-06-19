# claude-sandbox

Run Claude Code with `--dangerously-skip-permissions` in an isolated smolvm container. Skip all prompts, zero local trust required.

## Prerequisites

Install [smolvm](https://github.com/smol-machines/smolvm):

```bash
curl -sSL https://smolmachines.com/install.sh | bash
```

## Usage

```bash
make sandbox                                    # node:lts image, current directory
make sandbox IMAGE=ubuntu:22.04                 # custom image
make sandbox REPO=/path/to/other/project        # external repo
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

## Limitations

- Non-Debian images must have `/bin/sh` and network access
- Alpine, musl-based images need manual Node install (or pre-built image)
- smolvm macOS: requires OS 11+, Linux: requires KVM
