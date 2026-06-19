#!/bin/sh
set -e

# Install node if image doesn't have it (e.g. ubuntu:22.04)
if ! command -v npm >/dev/null 2>&1; then
    curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
    apt-get install -y nodejs
fi

npm install -g @anthropic-ai/claude-code --silent --prefer-offline 2>/dev/null || \
    npm install -g @anthropic-ai/claude-code

cd /workspace
exec claude --dangerously-skip-permissions
