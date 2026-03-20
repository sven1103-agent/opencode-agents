#!/bin/sh

set -eu

RELEASE_REPO_DEFAULT="sven1103-agent/opencode-agents"

RELEASE_REPO=${OPENCODE_HELPER_RELEASE_REPO:-$RELEASE_REPO_DEFAULT}
API_BASE=${OPENCODE_HELPER_API_BASE:-https://api.github.com/repos/$RELEASE_REPO/releases}
DOWNLOAD_BASE=${OPENCODE_HELPER_DOWNLOAD_BASE:-https://github.com/$RELEASE_REPO/releases/download}

err() {
  printf 'error: %s\n' "$*" >&2
}

usage() {
  cat <<EOF
install.sh - install opencode-helper

Usage:
  curl .../install.sh | sh
  curl .../install.sh | sh -s -- [options]

Options:
  --bin-dir PATH    Install directory (default: ~/.local/bin)
  --version TAG     Install a specific release tag (default: latest)
  -h, --help        Show this help

Environment:
  OPENCODE_HELPER_VERSION         Pin the release tag (overrides latest)
  OPENCODE_HELPER_RELEASE_REPO    Override release repo (default: sven1103-agent/opencode-agents)
EOF
}

fetch_url() {
  url=$1
  dest=$2

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL --connect-timeout 10 --max-time 60 "$url" -o "$dest"
    return $?
  fi
  if command -v wget >/dev/null 2>&1; then
    wget -q -O "$dest" --timeout=60 "$url"
    return $?
  fi

  err "missing downloader: require curl or wget"
  return 1
}

sha256_file() {
  file=$1
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$file" | awk '{print $1}'
  else
    shasum -a 256 "$file" | awk '{print $1}'
  fi
}

fetch_latest_tag() {
  out=$(mktemp)
  if ! fetch_url "${API_BASE}/latest" "$out"; then
    rm -f "$out"
    return 1
  fi
  tag=$(awk -F '"' '/"tag_name"[[:space:]]*:/ {print $4; exit}' "$out")
  rm -f "$out"
  [ -n "${tag:-}" ] || return 1
  printf '%s' "$tag"
}

verify_checksum() {
  file=$1
  checksums=$2
  name=$3

  [ -f "$checksums" ] || {
    err "missing checksums file: $checksums"
    return 1
  }

  expected=$(awk -v wanted="$name" '$2==wanted {print $1; found=1; exit} END {if (!found) exit 1}' "$checksums") || {
    err "checksums file missing entry: $name"
    return 1
  }

  actual=$(sha256_file "$file") || {
    err "failed to compute checksum for $name"
    return 1
  }

  if [ "$actual" != "$expected" ]; then
    err "checksum mismatch for $name"
    return 1
  fi
}

ASSUME_YES=0
BIN_DIR=
BIN_DIR_WAS_SET=0
REQUESTED_VERSION=

while [ "$#" -gt 0 ]; do
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
    --yes)
      ASSUME_YES=1
      ;;
    --bin-dir)
      shift
      [ "$#" -gt 0 ] || { err "--bin-dir requires an argument"; exit 1; }
      BIN_DIR=$1
      BIN_DIR_WAS_SET=1
      ;;
    --version)
      shift
      [ "$#" -gt 0 ] || { err "--version requires a tag"; exit 1; }
      REQUESTED_VERSION=$1
      ;;
    *)
      err "unknown option: $1"
      usage >&2
      exit 1
      ;;
  esac
  shift
done

TAG=
if [ -n "${REQUESTED_VERSION:-}" ]; then
  TAG=$REQUESTED_VERSION
elif [ -n "${OPENCODE_HELPER_VERSION:-}" ]; then
  TAG=$OPENCODE_HELPER_VERSION
else
  TAG=$(fetch_latest_tag) || {
    err "failed to determine latest release tag"
    exit 1
  }
fi

TEMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/opencode-helper-bootstrap.XXXXXX")
trap 'rm -rf "$TEMP_DIR"' EXIT INT TERM

INSTALLER_PATH="$TEMP_DIR/opencode-helper-install"
CHECKSUMS_PATH="$TEMP_DIR/checksums.txt"

installer_url="${DOWNLOAD_BASE}/${TAG}/opencode-helper-install"
checksums_url="${DOWNLOAD_BASE}/${TAG}/opencode-helper-${TAG}-checksums.txt"

fetch_url "$installer_url" "$INSTALLER_PATH" || {
  err "failed to download opencode-helper-install from $installer_url"
  exit 1
}

fetch_url "$checksums_url" "$CHECKSUMS_PATH" || {
  err "failed to download checksums from $checksums_url"
  exit 1
}

verify_checksum "$INSTALLER_PATH" "$CHECKSUMS_PATH" "opencode-helper-install" || {
  err "installer checksum verification failed"
  exit 1
}

printf 'ok: verified SHA-256 for opencode-helper-install\n'

chmod +x "$INSTALLER_PATH"

set -- "$INSTALLER_PATH" --yes --version "$TAG"
if [ "$BIN_DIR_WAS_SET" -eq 1 ]; then
  set -- "$@" --bin-dir "$BIN_DIR"
fi

OPENCODE_HELPER_RELEASE_API_BASE="$API_BASE" \
OPENCODE_HELPER_RELEASE_DOWNLOAD_BASE="$DOWNLOAD_BASE" \
exec "$@"
