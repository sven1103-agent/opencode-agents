# Installation

This guide covers how to install the OpenCode Config CLI (`oc`) on macOS and Linux.

## Prerequisites

- A Unix-like operating system (macOS or Linux)
- Go 1.21+ (for `go install` method)
- `curl` and `tar` (for manual download)

## Install Methods

### Method 1: Go Install (Recommended)

If you have Go installed:

```sh
go install github.com/sven1103-agent/opencode-config-cli@latest
```

This installs the `oc` binary to `$GOPATH/bin` (or `$HOME/go/bin`). Make sure this directory is in your PATH.

**Note:** `@latest` follows Go module version selection. It does not mean "latest GitHub prerelease".

**Install a specific version:**

```sh
go install github.com/sven1103-agent/opencode-config-cli@v1.0.0-alpha.5
```

### Method 2: Manual Download

Download the correct tarball for your platform from [GitHub Releases](https://github.com/sven1103-agent/opencode-config-cli/releases):

```sh
# Example: macOS ARM64
VERSION=v1.0.0-alpha.5
curl -L "https://github.com/sven1103-agent/opencode-config-cli/releases/download/${VERSION}/oc_${VERSION#v}_darwin_arm64.tar.gz" | tar xz
mv oc ~/.local/bin/
```

Set `VERSION` on its own line before running `curl`. Do not collapse this into `VERSION=... curl ...`, because `${VERSION}` is expanded by the shell before that temporary assignment takes effect.

**Available platforms:**
- `darwin_amd64` — macOS Intel
- `darwin_arm64` — macOS Apple Silicon
- `linux_amd64` — Linux x86_64
- `linux_arm64` — Linux ARM64

**Add to PATH:**

```sh
export PATH="$HOME/.local/bin:$PATH"
```

To make this permanent, add the line above to your shell profile (`~/.zshrc` or `~/.bashrc`).

## Verify Installation

```sh
oc version
```

You should see output like `oc version 1.0.0-alpha.5`.

## PATH Issues

If you see "command not found: oc" after installation:

1. **For go install**: Ensure `$HOME/go/bin` is in your PATH
2. **For manual download**: Ensure `~/.local/bin` is in your PATH

```sh
# Add to ~/.zshrc or ~/.bashrc
export PATH="$HOME/go/bin:$HOME/.local/bin:$PATH"
```

## Upgrading

Re-run `go install` with the desired version:

```sh
go install github.com/sven1103-agent/opencode-config-cli@latest
```

Or download a new release manually from [GitHub Releases](https://github.com/sven1103-agent/opencode-config-cli/releases).

## Next Steps

- [Configure a config bundle](config-bundles.md)
- [CLI Reference](cli-reference.md)
