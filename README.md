# OpenCode Config CLI

[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sven1103-agent/opencode-config-cli?logo=go)](https://github.com/sven1103-agent/opencode-config-cli/blob/main/go.mod)
[![Version](https://img.shields.io/github/v/release/sven1103-agent/opencode-config-cli?include_prereleases&label=version)](https://github.com/sven1103-agent/opencode-config-cli/releases)
[![CI](https://github.com/sven1103-agent/opencode-config-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/sven1103-agent/opencode-config-cli/actions/workflows/ci.yml)
[![E2E CLI](https://github.com/sven1103-agent/opencode-config-cli/actions/workflows/e2e-cli.yml/badge.svg)](https://github.com/sven1103-agent/opencode-config-cli/actions/workflows/e2e-cli.yml)

**Manage OpenCode configuration bundles and schema-validated multi-agent workflows**

## What is this?

A CLI tool (`oc`) that manages OpenCode [configuration bundles](docs/config-bundles.md) from external sources, enabling versioned, validated configs for AI agents.

## Quick Start (30 seconds)

```sh
# Install via Go (macOS/Linux)
go install github.com/sven1103-agent/opencode-config-cli@latest

# Create an alias to the default installation location
alias oc='$HOME/go/bin/opencode-config-cli'

# Register a config bundle
oc source add qbicsoftware/opencode-config-bundle --name qbic

# Apply a preset
oc bundle apply qbic --preset mixed --project-root .
```

## Installation

Detailed installation guide: [docs/installation.md](docs/installation.md)

Install methods:
- `go install` (recommended)
- Manual download from GitHub Releases

## Shell Completion

Enable tab completion for your shell:

```sh
# Bash: source on the fly, or install permanently
source <(oc completion bash)
oc completion bash | sudo tee /etc/bash_completion.d/oc > /dev/null

# Zsh: source on the fly (recommended)
source <(oc completion zsh)
# If you use alias (alias oc="opencode-config-cli"), also add:
compdef _oc opencode-config-cli

# Or save to completions dir:
# oc completion zsh > ~/.zsh/completions/_oc
# Add to ~/.zshrc: fpath=(~/.zsh/completions $fpath)
# Add: compdef _oc opencode-config-cli
# Clear cache: rm -f ~/.zcompdump && exec zsh

# Fish
oc completion fish > ~/.config/fish/completions/oc.fish
```

Run `oc completion --help` for full instructions.

## Key Concepts

| Concept | Description |
|---------|-------------|
| **Sources** | Registered bundle repositories (GitHub releases) |
| **Bundles** | Versioned, schema-validated OpenCode configs |
| **Presets** | Named configurations (e.g., `mixed`, `openai`) |

## Available Bundles

- [qbicsoftware/opencode-config-bundle](https://github.com/qbicsoftware/opencode-config-bundle) — Official bundle with multiple presets

## Documentation

| Guide | Description |
|-------|-------------|
| [docs/README.md](docs/README.md) | User documentation hub |
| [docs/installation.md](docs/installation.md) | Install on macOS/Linux |
| [docs/config-bundles.md](docs/config-bundles.md) | Understand bundles + create your own |
| [docs/cli-reference.md](docs/cli-reference.md) | Full command reference |
| [docs/troubleshooting.md](docs/troubleshooting.md) | FAQ and common issues |

## Legacy (Bash Version)

The original Bash-based helper is deprecated. Use the Go CLI (`oc`) instead.

Archive documentation: [docs/legacy/bash-helper.md](docs/legacy/bash-helper.md)

---

## License

AGPL-3.0 — see the [LICENSE](LICENSE) file for details.
