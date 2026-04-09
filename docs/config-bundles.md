# Config Bundles

Config bundles are versioned, schema-validated OpenCode configuration packages distributed via GitHub releases.

## What is a Bundle?

A bundle is a tar archive containing:
- `opencode-bundle.manifest.json` — Bundle metadata
- Preset JSON files (e.g., `opencode.mixed.json`, `opencode.openai.json`)
- Optional: `.opencode/schemas/` for session artifacts

## Available Bundles

| Bundle | Description |
|--------|-------------|
| [qbicsoftware/opencode-config-bundle](https://github.com/qbicsoftware/opencode-config-bundle) | Official bundle with multiple presets |

## Using a Bundle

### 1. Register a Source

```sh
oc source add qbicsoftware/opencode-config-bundle --name qbic
```

Source formats:
- `owner/repo` — e.g., `qbicsoftware/opencode-config-bundle`
- `github.com/owner/repo`
- `https://github.com/owner/repo/releases/tag/v1.0.0`

### 2. List Available Presets

```sh
oc preset list --sources
```

### 3. Apply a Preset

```sh
oc bundle apply qbic --preset mixed --project-root .
oc bundle apply qbic --project-root .
```

## Creating Your Own Bundle

### 1. Create the Manifest

Create `opencode-bundle.manifest.json` at your bundle root:

```json
{
  "manifest_version": "1.0.0",
  "bundle_name": "my-bundle",
  "bundle_version": "v1.0.0",
  "presets": [
    {
      "name": "openai",
      "description": "OpenAI-based configuration",
      "entrypoint": "opencode.openai.json",
      "prompt_files": []
    }
  ]
}
```

### 2. Required Fields

| Field | Description |
|-------|-------------|
| `manifest_version` | Semantic version (e.g., `1.0.0`) |
| `bundle_name` | Stable identifier (lowercase, hyphens allowed) |
| `bundle_version` | Release tag |
| `presets` | Array of preset objects |

### 3. Preset Descriptor Fields

| Field | Description |
|-------|-------------|
| `name` | Stable preset ID |
| `description` | Short description |
| `entrypoint` | Path to preset JSON file |
| `prompt_files` | Array of prompt file paths (can be empty) |

### 4. Publish as GitHub Release

1. Create a GitHub release
2. Attach a `.tar.gz` archive containing:
   - `opencode-bundle.manifest.json` at root
   - All preset files
3. (Optional) Add a `-checksums.txt` for integrity verification

### 5. Register Your Bundle

```sh
oc source add your-username/your-bundle --name mybundle
oc bundle apply mybundle --preset <preset-name> --project-root .
