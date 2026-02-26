#!/usr/bin/env sh
set -eu

if command -v tmux >/dev/null 2>&1; then
  tmux list-sessions 2>/dev/null || echo "no sessions"
else
  echo "tmux not installed"
fi
