# CLI Reference

Complete reference for the `oc` CLI commands.

## Global Options

| Option | Description |
|--------|-------------|
| `--help` | Show help |
| `--version` | Show version |

## Source Commands

Manage registered config sources (GitHub repositories).

### oc source add

Register a new config source.

```sh
oc source add <source> --name <name>
```

**Arguments:**
- `source` — GitHub repo (`owner/repo`), full URL, or release URL

**Options:**
- `--name` — Local name for the source (required)

**Examples:**

```sh
# Register a repo (uses latest release)
oc source add qbicsoftware/opencode-config-bundle --name qbic

# Register a specific release
oc source add https://github.com/qbicsoftware/opencode-config-bundle/releases/tag/v1.2.3 --name qbic-v123
```

### oc source list

List all registered sources.

```sh
oc source list
```

### oc source remove

Remove a registered source.

```sh
oc source remove <name>
```

## Bundle Commands

Apply and manage config bundles.

### oc bundle apply

Apply a preset from a registered source.

```sh
oc bundle apply <source-ref> [--preset <preset>] [--auto] --project-root <path>
```

**Options:**
- `--preset` — Preset name to apply
- `--auto` — Disable interactive preset selection in TTY mode
- `--project-root` — Target directory (default: `.`)

**Example:**

```sh
oc bundle apply qbic --preset mixed --project-root ./myproject
oc bundle apply qbic --project-root ./myproject
```

### oc bundle status

Show provenance of applied bundles.

```sh
oc bundle status --project-root <path>
```

**Example:**

```sh
oc bundle status --project-root ./myproject
```

### oc bundle update

Check for and apply updates from update-capable sources.

```sh
oc bundle update <source-id>
```

## Preset Commands

### oc preset list

List available presets.

```sh
oc preset list --sources
```

Shows presets from all registered sources with source and bundle context.

### oc completion

Generate shell completions.

```sh
oc completion bash
oc completion zsh
oc completion fish
```

## Migration Commands

### oc migrate legacy-config

Migrate a V1 legacy project to V2.

```sh
oc migrate legacy-config --project-root <path>
```

**Example:**

```sh
oc migrate legacy-config --project-root ./myproject
```

This migrates projects using the old `.opencode/opencode-helper-manifest.tsv` format to the new V2 format.

## Other Commands

### oc version

Show CLI version.

```sh
oc version
```

### oc update

Update the CLI to the latest version.

```sh
oc update
```

Note: This only works if you installed via the installer script.
