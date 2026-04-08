# OpenCode Config CLI

**Manage OpenCode configuration bundles and schema-validated multi-agent workflows**

## What is this?

A CLI tool (`oc`) that manages OpenCode configuration bundles from external sources, enabling versioned, validated configs for AI agents.

## Quick Start (30 seconds)

```sh
# Install (macOS/Linux)
curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh

# Register a config bundle
oc source add qbicsoftware/opencode-config-bundle --name qbic

# Apply a preset
oc bundle apply qbic --preset mixed --project-root .
```

## Installation

Detailed installation guide: [docs/installation.md](docs/installation.md)

Install methods:
- Installer script (recommended)
- `go install`
- Manual download

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
