# Installation

This guide covers how to install the OpenCode Config CLI (`oc`) on macOS and Linux.

## Prerequisites

- A Unix-like operating system (macOS or Linux)
- `curl` or `wget` for downloading
- `tar` for extracting archives
- Go 1.21+ (for `go install` method)

## Install Methods

### Method 1: Installer Script (Recommended)

```sh
curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh
```

This downloads and runs the installer, which places the `oc` binary in `~/.local/bin/`.

**Add to PATH:**

```sh
export PATH="$HOME/.local/bin:$PATH"
```

To make this permanent, add the line above to your shell profile (`~/.zshrc` or `~/.bashrc`).

**Custom install directory:**

```sh
curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh -s -- --bin-dir "$HOME/bin"
```

### Method 2: Go Install

If you have Go installed:

```sh
go install github.com/sven1103-agent/opencode-config-cli@latest
```

This installs the `oc` binary to `$GOPATH/bin` (or `$HOME/go/bin`).

**Note:** `@latest` follows Go module version selection. It does not mean "latest GitHub prerelease".

**Install a specific version:**

```sh
go install github.com/sven1103-agent/opencode-config-cli@v1.0.0-alpha.4
```

### Method 3: Manual Download

Download the correct tarball for your platform from [GitHub Releases](https://github.com/sven1103-agent/opencode-config-cli/releases):

```sh
# Example: macOS ARM64
VERSION=v1.0.0-alpha.4
curl -L "https://github.com/sven1103-agent/opencode-config-cli/releases/download/${VERSION}/oc_${VERSION#v}_darwin_arm64.tar.gz" | tar xz
mv oc ~/.local/bin/
```

**Available platforms:**
- `darwin_amd64` — macOS Intel
- `darwin_arm64` — macOS Apple Silicon
- `linux_amd64` — Linux x86_64
- `linux_arm64` — Linux ARM64

## Verify Installation

```sh
oc version
```

You should see output like `oc version 1.0.0-alpha.4`.

## Version Pinning

### With installer script:

```sh
curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh -s -- --version v1.0.0-alpha.4
```

Or via environment variable:

```sh
OPENCODE_HELPER_VERSION=v1.0.0-alpha.4 sh -c 'curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh'
```

## Upgrading

Re-run the installer to upgrade to the latest version:

```sh
curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh
```

Or use `go install` with the desired version:

```sh
go install github.com/sven1103-agent/opencode-config-cli@latest
```

## Next Steps

- [Configure a config bundle](config-bundles.md)
- [CLI Reference](cli-reference.md)
