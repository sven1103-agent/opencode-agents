#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO_ROOT=$(cd "$SCRIPT_DIR/.." && pwd)
TMP_DIR=/tmp/oc-demo-recording
WORK_DIR="$TMP_DIR/work"
BIN_DIR="$TMP_DIR/bin"
HOME_DIR="$TMP_DIR/home"
CONFIG_DIR="$TMP_DIR/config"
TAPE_PATH="$REPO_ROOT/docs/demo.tape"
GIF_PATH="$REPO_ROOT/docs/demo.gif"

cleanup() {
  rm -rf "$TMP_DIR"
}

trap cleanup EXIT

rm -rf "$TMP_DIR"
mkdir -p "$WORK_DIR" "$BIN_DIR" "$HOME_DIR" "$CONFIG_DIR"

command -v vhs >/dev/null 2>&1 || {
  printf 'error: vhs is required. Install it with Homebrew: brew install vhs\n' >&2
  exit 1
}

go build -o "$BIN_DIR/oc" "$REPO_ROOT"

env \
  HOME="$HOME_DIR" \
  XDG_CONFIG_HOME="$CONFIG_DIR" \
  PATH="$BIN_DIR:$PATH" \
  TERM="xterm-256color" \
  vhs "$TAPE_PATH"

printf 'written: %s\n' "$GIF_PATH"
