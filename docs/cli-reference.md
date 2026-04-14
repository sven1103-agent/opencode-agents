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
oc source list [--with-presets]
```

**Options:**
- `--with-presets` — Show all presets from all registered sources in a flat table

**Examples:**

```sh
# List all registered sources
oc source list

# List all sources with their available presets
oc source list --with-presets
```

### oc source remove

Remove a registered source.

```sh
oc source remove <name>
```

## Bundle Commands

Install and manage config bundles.

### oc bundle install

Install a preset from a registered source to your project.

```sh
oc bundle install <source-ref> [--preset <preset>] [--auto] --project-root <path>
```

**Options:**
- `--preset` — Preset name to install (required in non-interactive mode)
- `--auto` — Run in non-interactive mode (disables prompts)
- `--project-root` — Target directory (default: `.`)
- `--force` — Overwrite existing files
- `--dry-run` — Show what would be done without doing it

**Examples:**

```sh
# Install a preset from a source
oc bundle install qbic --preset mixed --project-root ./myproject

# Install with interactive preset selection
oc bundle install qbic --project-root ./myproject

# Force overwrite existing config
oc bundle install qbic --preset mixed --force
```

### oc bundle status

Show provenance of installed bundles.

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

**Example:**

```sh
oc bundle update qbic
```

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
