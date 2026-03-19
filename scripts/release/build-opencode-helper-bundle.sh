#!/bin/sh

set -eu

usage() {
  cat <<EOF
Usage: build-opencode-helper-bundle.sh --tag <tag> --commit-sha <sha> [--output-dir <dir>]
EOF
}

sha256_file() {
  file=$1
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$file" | awk '{print $1}'
  else
    shasum -a 256 "$file" | awk '{print $1}'
  fi
}

file_mode() {
  file=$1
  if stat --version >/dev/null 2>&1; then
    stat -c '%a' "$file"
  else
    stat -f '%Lp' "$file"
  fi
}

TAG=
COMMIT_SHA=
OUTPUT_DIR=${PWD}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --tag)
      shift
      TAG=${1:-}
      ;;
    --commit-sha)
      shift
      COMMIT_SHA=${1:-}
      ;;
    --output-dir)
      shift
      OUTPUT_DIR=${1:-}
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      printf 'error: unknown option: %s\n' "$1" >&2
      usage >&2
      exit 1
      ;;
  esac
  shift
done

[ -n "$TAG" ] || { printf 'error: --tag is required\n' >&2; exit 1; }
[ -n "$COMMIT_SHA" ] || { printf 'error: --commit-sha is required\n' >&2; exit 1; }

SCRIPT_DIR=$(CDPATH= cd "$(dirname "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd "$SCRIPT_DIR/../.." && pwd)
OUTPUT_DIR=$(mkdir -p "$OUTPUT_DIR" && CDPATH= cd "$OUTPUT_DIR" && pwd)

BUNDLE_ROOT="opencode-helper-$TAG"
TARBALL_NAME="$BUNDLE_ROOT.tar.gz"
MANIFEST_NAME="$BUNDLE_ROOT-manifest.json"
CHECKSUMS_NAME="$BUNDLE_ROOT-checksums.txt"

TMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/opencode-helper-bundle.XXXXXX")
trap 'rm -rf "$TMP_DIR"' EXIT INT TERM

STAGE="$TMP_DIR/stage/$BUNDLE_ROOT"
mkdir -p "$STAGE/scripts" "$STAGE/.opencode/schemas"

cp "$REPO_ROOT/scripts/opencode-helper" "$STAGE/scripts/opencode-helper"
cp "$REPO_ROOT/scripts/opencode-helper-install" "$STAGE/scripts/opencode-helper-install"
cp "$REPO_ROOT/opencode.openai.json" "$STAGE/opencode.openai.json"
cp "$REPO_ROOT/opencode.mixed.json" "$STAGE/opencode.mixed.json"
cp "$REPO_ROOT/.opencode/schemas/handoff.schema.json" "$STAGE/.opencode/schemas/handoff.schema.json"
cp "$REPO_ROOT/.opencode/schemas/result.schema.json" "$STAGE/.opencode/schemas/result.schema.json"

chmod 755 "$STAGE/scripts/opencode-helper" "$STAGE/scripts/opencode-helper-install"
chmod 644 "$STAGE/opencode.openai.json" "$STAGE/opencode.mixed.json"
chmod 644 "$STAGE/.opencode/schemas/handoff.schema.json" "$STAGE/.opencode/schemas/result.schema.json"

FILES='scripts/opencode-helper
scripts/opencode-helper-install
opencode.openai.json
opencode.mixed.json
.opencode/schemas/handoff.schema.json
.opencode/schemas/result.schema.json'

CONTENTS_FILE="$TMP_DIR/contents.jsonl"
: > "$CONTENTS_FILE"
IFS='
'
for rel in $FILES; do
  abs="$STAGE/$rel"
  size=$(wc -c < "$abs" | tr -d ' ')
  mode=$(file_mode "$abs")
  sha=$(sha256_file "$abs")
  printf '{"path":"%s","size":%s,"mode":"%s","sha256":"%s"}\n' "$rel" "$size" "$mode" "$sha" >> "$CONTENTS_FILE"
done
unset IFS

MANIFEST_INNER="$STAGE/release-manifest.json"
{
  printf '{\n'
  printf '  "bundle_format_version": 1,\n'
  printf '  "release_tag": "%s",\n' "$TAG"
  printf '  "commit_sha": "%s",\n' "$COMMIT_SHA"
  printf '  "bundle_root": "%s",\n' "$BUNDLE_ROOT"
  printf '  "contents": [\n'
  awk 'NR==1{printf "    %s",$0} NR>1{printf ",\n    %s",$0} END{printf "\n"}' "$CONTENTS_FILE"
  printf '  ]\n'
  printf '}\n'
} > "$MANIFEST_INNER"
chmod 644 "$MANIFEST_INNER"

cp "$MANIFEST_INNER" "$OUTPUT_DIR/$MANIFEST_NAME"

touch -t 198001010000 "$STAGE" "$STAGE/scripts" "$STAGE/.opencode" "$STAGE/.opencode/schemas"
touch -t 198001010000 "$STAGE/scripts/opencode-helper" "$STAGE/scripts/opencode-helper-install"
touch -t 198001010000 "$STAGE/opencode.openai.json" "$STAGE/opencode.mixed.json"
touch -t 198001010000 "$STAGE/.opencode/schemas/handoff.schema.json" "$STAGE/.opencode/schemas/result.schema.json"
touch -t 198001010000 "$MANIFEST_INNER"

LIST_FILE="$TMP_DIR/tar.list"
cat > "$LIST_FILE" <<EOF
$BUNDLE_ROOT
$BUNDLE_ROOT/scripts
$BUNDLE_ROOT/scripts/opencode-helper
$BUNDLE_ROOT/scripts/opencode-helper-install
$BUNDLE_ROOT/opencode.openai.json
$BUNDLE_ROOT/opencode.mixed.json
$BUNDLE_ROOT/.opencode
$BUNDLE_ROOT/.opencode/schemas
$BUNDLE_ROOT/.opencode/schemas/handoff.schema.json
$BUNDLE_ROOT/.opencode/schemas/result.schema.json
$BUNDLE_ROOT/release-manifest.json
EOF

TAR_PATH="$OUTPUT_DIR/$TARBALL_NAME"
if tar --help 2>/dev/null | grep -q -- '--uid'; then
  tar --format=ustar --uid 0 --gid 0 --uname root --gname root -cf "$TMP_DIR/bundle.tar" -C "$TMP_DIR/stage" -T "$LIST_FILE"
elif tar --help 2>/dev/null | grep -q -- '--owner'; then
  tar --format=ustar --owner 0 --group 0 --numeric-owner -cf "$TMP_DIR/bundle.tar" -C "$TMP_DIR/stage" -T "$LIST_FILE"
else
  tar --format=ustar -cf "$TMP_DIR/bundle.tar" -C "$TMP_DIR/stage" -T "$LIST_FILE"
fi
gzip -n -c "$TMP_DIR/bundle.tar" > "$TAR_PATH"

TAR_SHA=$(sha256_file "$TAR_PATH")
MANIFEST_SHA=$(sha256_file "$OUTPUT_DIR/$MANIFEST_NAME")
{
  printf '%s  %s\n' "$TAR_SHA" "$TARBALL_NAME"
  printf '%s  %s\n' "$MANIFEST_SHA" "$MANIFEST_NAME"
} > "$OUTPUT_DIR/$CHECKSUMS_NAME"

printf '%s\n' "$OUTPUT_DIR/$TARBALL_NAME"
printf '%s\n' "$OUTPUT_DIR/$MANIFEST_NAME"
printf '%s\n' "$OUTPUT_DIR/$CHECKSUMS_NAME"
