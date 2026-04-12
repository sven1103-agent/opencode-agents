# OpenCode Config CLI - Documentation

Welcome to the OpenCode Config CLI documentation. This CLI manages OpenCode configuration bundles from external sources, enabling versioned, validated configs for AI agents.

## Quick Links

| Topic | Description |
|-------|-------------|
| [Installation](installation.md) | Install the CLI on macOS or Linux |
| [Config Bundles](config-bundles.md) | Understand bundles and create your own |
| [CLI Demo Playbook](demo-playbook.md) | Repeatable script for recording README terminal demos |
| [CLI Reference](cli-reference.md) | Full command reference |
| [Troubleshooting](troubleshooting.md) | Common issues and solutions |
| [Legacy Docs](legacy/bash-helper.md) | Deprecated Bash version |

## Key Concepts

- **Sources** — Registered bundle repositories (GitHub releases)
- **Bundles** — Versioned, schema-validated OpenCode configurations
- **Presets** — Named configurations (e.g., `mixed`, `openai`)

## Getting Started

```sh
# Install via Go
go install github.com/sven1103-agent/opencode-config-cli@latest

# Register a config bundle
oc source add qbicsoftware/opencode-config-bundle --name qbic

# Apply a preset
oc bundle apply qbic --preset mixed --project-root .

# Or let the CLI prompt for a preset in a TTY
oc bundle apply qbic --project-root .
```

## Available Bundles

- [qbicsoftware/opencode-config-bundle](https://github.com/qbicsoftware/opencode-config-bundle) — Official configuration bundle with multiple presets

## Need Help?

- Check the [troubleshooting guide](troubleshooting.md) for common issues
- Open an issue at [github.com/sven1103-agent/opencode-config-cli/issues](https://github.com/sven1103-agent/opencode-config-cli/issues)
