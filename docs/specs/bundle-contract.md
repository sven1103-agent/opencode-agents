# Bundle Contract

This document defines the contract that configuration bundles must comply with to work with the `oc` CLI.

## Ownership

The bundle contract is **owned by the opencode-config-cli repository** where the CLI lives. All schema definitions, validation rules, and contract changes are managed here.

Configuration bundles (like [qbicsoftware/opencode-config-bundle](https://github.com/qbicsoftware/opencode-config-bundle)) are **consumers** of this contract.

---

## Bundle Manifest Schema

### File Location

```
<bundle-root>/
  opencode-bundle.manifest.json  <- required at bundle root
  <preset-entrypoint>.json       <- preset files referenced by manifest
  .opencode/schemas/
    handoff.schema.json          <- canonical handoff contract
    result.schema.json           <- canonical result contract
```

### Schema (JSON Schema Draft 2020-12)

The schema is **embedded in the CLI binary** for offline validation.

The single source of truth is `internal/bundle/1.0.0.schema.json` in the repository.
This file is embedded into the `oc` CLI binary at build time.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["manifest_version", "bundle_name", "bundle_version", "presets"],
  "properties": {
    "manifest_version": {
      "type": "string",
      "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$",
      "description": "Semantic version of the bundle manifest schema (e.g., 1.0.0). CLI supports this version and N-1."
    },
    "bundle_name": {
      "type": "string",
      "description": "Stable identifier for the bundle"
    },
    "bundle_version": {
      "type": "string",
      "description": "Release tag or version string (e.g., v1.0.0)"
    },
    "source_repo": {
      "type": "string",
      "description": "URL of source repository"
    },
    "source_commit": {
      "type": "string",
      "description": "Git commit SHA"
    },
    "release_tag": {
      "type": "string",
      "description": "Release tag"
    },
    "presets": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["name", "description", "entrypoint"],
        "properties": {
          "name": {
            "type": "string",
            "description": "Stable preset identifier (lowercase, hyphens allowed)"
          },
          "description": {
            "type": "string",
            "description": "Short human-readable description"
          },
          "entrypoint": {
            "type": "string",
            "description": "Path to preset JSON file relative to bundle root"
          },
          "prompt_files": {
            "type": "array",
            "items": { "type": "string" },
            "description": "Array of prompt file paths to materialize (can be empty)"
          }
        }
      }
    }
  }
}
```

### Required Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `manifest_version` | string (semver) | Yes | Must be `1.0.0` or compatible. CLI supports this version and N-1. |
| `bundle_name` | string | Yes | Stable identifier for the bundle |
| `bundle_version` | string | Yes | Release tag or version string |
| `presets` | array | Yes | Array of preset descriptor objects |

### Preset Descriptor Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Stable identifier (lowercase, hyphens allowed) |
| `description` | string | Yes | Short human-readable description |
| `entrypoint` | string | Yes | Path to preset JSON file relative to bundle root |
| `prompt_files` | array | Yes | Array of prompt file paths (can be empty) |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `source_repo` | string | URL of source repository |
| `source_commit` | string | Git commit SHA |
| `release_tag` | string | Release tag (often matches `bundle_version`) |

---

## CLI Validation Rules

The CLI rejects a bundle when:

- `opencode-bundle.manifest.json` is missing from bundle root
- Manifest is not valid JSON
- `manifest_version` is unsupported (not semver-compatible with CLI's supported versions)
- Any required top-level field is missing
- Any preset is missing required fields
- An `entrypoint` path does not exist in the bundle
- A `prompt_files` entry does not exist in the bundle

---

## Future Extensions

### US-EXT-001 - Support Custom Tools in Bundles

**Status**: Backlog  
**Milestone**: v1.1.0

Add optional `tools` field to support custom OpenCode tools in bundles:

```json
{
  "manifest_version": 1,
  "bundle_version": "v1.1.0",
  "tools": [
    {
      "name": "database",
      "entrypoint": ".opencode/tools/database.ts",
      "description": "Query the project database"
    }
  ]
}
```

See [Custom Tools Docs](https://opencode.ai/docs/custom-tools) for background on how OpenCode tools work.

---

*Contract version: 1.0.0*
*Last updated: 2026-04-01*
