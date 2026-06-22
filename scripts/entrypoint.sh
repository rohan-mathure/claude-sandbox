#!/bin/sh
set -e

# Install node if image doesn't have it (e.g. ubuntu:22.04)
if ! command -v npm >/dev/null 2>&1; then
    curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
    apt-get install -y nodejs
fi

npm install -g @anthropic-ai/claude-code --silent --prefer-offline 2>/dev/null || \
    npm install -g @anthropic-ai/claude-code

# Create non-root user
if ! id -u sandbox >/dev/null 2>&1; then
    useradd -m -s /bin/sh sandbox
fi

cd /workspace
chown -R sandbox:sandbox /workspace

# Run Claude as non-root user
exec su - sandbox -c "cd /workspace && claude --dangerously-skip-permissions"
