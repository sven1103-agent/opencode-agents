# OpenCode Helper CLI

## Contents

- [Overview](#overview)
- [Decision Log](#decision-log)
- [Release Intake](#release-intake)
- [PRD](#prd)
- [Scope](#scope)
- [Requirements](#requirements)
  - [Functional Requirements](#functional-requirements)
  - [Non-Functional Requirements](#non-functional-requirements)
- [First Features](#first-features)
- [V2 Milestone Plan](opencode-helper-v2-milestones.md)
- [Go Migration Milestone Plan](opencode-helper-go-migration-milestones.md)
- [Traceability Matrix](#traceability-matrix)
- [User Story Backlog Placeholder](#user-story-backlog-placeholder)
- [Open Questions for Post-V1](#open-questions-for-post-v1)

## Overview

This document is the single traceable product source for the OpenCode helper CLI. It captures the initial product direction, requirements, first feature set, and the traceability structure that later user stories must follow.

Related planning docs:
- `docs/opencode-helper-v2-milestones.md` - milestone sequencing and rollout plan for the V2+ config-source architecture

Status: Draft for V1  
Owner: TBD  
Scope baseline: V1

---

## Decision Log

### <a id="dec-001"></a>DEC-001 - CLI is the sole official distribution channel for V1

The helper CLI is the only supported distribution mechanism for official OpenCode setup assets in V1.

Included in the CLI release bundle:
- Official config presets
- Official inter-agent schema contracts
- Version metadata needed for validation and provenance

Implications:
- Users do not need to clone this repository to use the official setup assets
- The CLI release is the canonical end-user artifact
- `self-update` updates the CLI and bundled assets together
- Live remote schema fetching is out of scope for V1

### <a id="dec-002"></a>DEC-002 - Config bundles are distributed via versioned config sources (V2+)

Starting in V2+, official OpenCode configuration assets (presets + prompts) are distributed via one or more separate versioned config sources with independent release flow.

Supported source types may include:
- local directories
- local `.tar.gz` archive files
- versioned remote release assets (for example GitHub release bundles)

The helper CLI remains the supported installer/applier:
- The CLI can register one or more config sources.
- The CLI resolves and caches config bundles from supported source types.
- The CLI applies a chosen bundle (and preset within it) into a target project.
- The CLI reports provenance for installed bundles/presets/prompts.

Interaction with DEC-001:
- DEC-001 remains the V1 baseline.
- DEC-002 supersedes DEC-001 for V2+ scope.

---

## Release Intake

### <a id="release-intake-001"></a>Release Intake For 1.0.0

The `1.0.0` GitHub milestone is the sole curated release todo list for general availability.

Lightweight maintainer process:

- prerelease testers are encouraged to open normal GitHub issues for any bug or rough edge they find
- the repository provides a small prerelease bug report issue template to make reports easier to reproduce
- only the maintainer assigns issues to milestone `1.0.0`
- only issues that must land before GA are added to milestone `1.0.0`
- only clear GA blockers receive the `release:blocker` label
- valid but non-GA issues may remain outside the milestone or receive `release:post-1.0.0`

Prerelease release notes should include a short call for smoke testing and issue filing, asking reporters to provide:

- prerelease version tested
- OS and shell
- install method
- command run
- expected behavior
- actual behavior
- reproduction steps

---

## PRD

### <a id="prd-001"></a>PRD-001 - Product Goal

Provide a small helper CLI that bootstraps and maintains a local OpenCode project setup using official bundled config presets and official bundled inter-agent schemas, with safe validation and self-update behavior tied to CLI releases.

V2+ direction (see [DEC-002](#dec-002)):
- The CLI bootstraps and maintains a local OpenCode project setup by installing and applying config bundles from registered config sources.
- The CLI may still ship with minimal built-in defaults, but official presets/prompts are expected to come from external config sources.

### <a id="prd-002"></a>PRD-002 - Problem Statement

Current OpenCode project setup depends on repository-local assets and setup knowledge that are not yet exposed through a single user-facing distribution and setup workflow. Users may not want to clone the repository locally just to obtain supported presets and schema files.

### <a id="prd-003"></a>PRD-003 - Target Users

- Maintainers who publish and evolve the official OpenCode setup
- Developers who want to initialize a project with a supported OpenCode setup
- Contributors who need a repeatable, validated local setup flow

### <a id="prd-004"></a>PRD-004 - User Value

The helper CLI should let a user:
- run a low-friction installer wizard (`opencode-helper-install`) to install `opencode-helper` globally and ensure it is available on `PATH`
- inspect all available presets with a short description of each preset's purpose
- apply a supported local configuration
- install required schemas
- validate setup health
- update to the latest supported release

V2+ user value additions:
- register one or more config sources
- browse presets across sources (with source + bundle release context)
- install and apply a chosen config bundle release (including referenced prompt files)
- see provenance of installed bundle/preset/prompts
- check for newer compatible bundle releases from update-capable sources and get prompted before updating

### <a id="prd-005"></a>PRD-005 - Success Criteria

V1 is successful when:
- A user can install `opencode-helper` without cloning this repository (via the `opencode-helper-install` wizard)
- A user can complete a global install such that `opencode-helper` works in a new terminal session (with `PATH` updated by the installer)
- A user can initialize a project from bundled presets and schemas
- A user can list all bundled presets and understand each preset's purpose before selecting one
- A user can validate whether the project setup is healthy or drifted
- A user can update the CLI and bundled assets via GitHub release upgrade
- Installed project assets can be traced back to the CLI release version

---

## Scope

### In Scope

- Bundling official config presets with the CLI release
- Bundling all available official OpenCode config variants with the CLI release
- Bundling official inter-agent schemas with the CLI release
- Project-local setup from bundled assets
- Preset discovery with short human-readable descriptions or purpose text
- Safe install/apply behavior for project files
- Validation and diagnostics
- Self-update from the latest GitHub release
- Version/provenance reporting for bundled assets

### In Scope (V2+ config bundles)

- Registering and managing one or more config sources
- Resolving/caching config bundles from supported source types
- Preset discovery across registered sources (including descriptions and provenance context)
- Applying a selected bundle release (and preset) to a project
- Installing prompt files referenced by presets (e.g. `{file:./prompts/...}`) as part of applying a bundle
- Reporting provenance for installed bundle/preset/prompt files
- Checking for newer compatible bundle releases from update-capable sources and prompting the user before updating installed bundles

### Out of Scope

- Direct end-user fetching of raw assets from GitHub releases
- Live remote schema syncing during normal setup
- Automatic migration of arbitrary user-customized configs
- Windows support unless explicitly added later
- Full remote orchestration or hosted service behavior

V2+ out of scope (initially):
- Automatically merging arbitrary user overlays with bundle upgrades
- Mandatory online operation after bundles are installed/cached
- Config repo authentication flows beyond what the underlying git/http tooling supports by default

Supported source types for the first V2+ scope baseline:
- local directory source
- local `.tar.gz` archive source
- GitHub release bundle source

---

## Requirements

### Functional Requirements

#### <a id="req-f-001"></a>REQ-F-001 - Bundled Presets

The CLI shall include all officially supported OpenCode config presets as bundled release assets.

Depends on:
- [DEC-001](#dec-001)

#### <a id="req-f-001a"></a>REQ-F-001a - Preset Metadata

The CLI shall bundle short human-readable metadata for each preset, including at least a name and a brief description of its purpose.

Depends on:
- [REQ-F-001](#req-f-001)

#### <a id="req-f-002"></a>REQ-F-002 - Bundled Schemas

The CLI shall include the official inter-agent handoff and result schemas as bundled release assets.

Depends on:
- [DEC-001](#dec-001)

#### <a id="req-f-003"></a>REQ-F-003 - Preset Discovery

The CLI shall let users list the bundled config presets available in the installed CLI release.

The preset list shall show each preset together with its short description or purpose.

Depends on:
- [REQ-F-001](#req-f-001)
- [REQ-F-001a](#req-f-001a)

#### <a id="req-f-004"></a>REQ-F-004 - Preset Application

The CLI shall let users materialize a selected bundled preset into a target project as `opencode.json` or equivalent configured output.

Depends on:
- [REQ-F-001](#req-f-001)
- [REQ-F-003](#req-f-003)

#### <a id="req-f-005"></a>REQ-F-005 - Schema Installation

The CLI shall let users install the bundled schemas into the target project or supported scope.

Depends on:
- [REQ-F-002](#req-f-002)

#### <a id="req-f-006"></a>REQ-F-006 - Setup Validation

The CLI shall validate that the target project contains the required OpenCode setup files and that installed assets are compatible with the bundled release contents.

Depends on:
- [REQ-F-001](#req-f-001)
- [REQ-F-002](#req-f-002)

#### <a id="req-f-007"></a>REQ-F-007 - Diagnostics

The CLI shall provide diagnostics for missing files, drift, invalid setup state, and likely operator mistakes.

Depends on:
- [REQ-F-006](#req-f-006)

#### <a id="req-f-008"></a>REQ-F-008 - Self Update

The CLI shall update itself from the latest GitHub release of the project.

Depends on:
- [DEC-001](#dec-001)

#### <a id="req-f-009"></a>REQ-F-009 - Version Reporting

The CLI shall report its own version and the version or identity of the bundled asset set.

Depends on:
- [REQ-F-001](#req-f-001)
- [REQ-F-002](#req-f-002)
- [REQ-F-008](#req-f-008)

#### <a id="req-f-010"></a>REQ-F-010 - Safe File Handling

The CLI shall avoid overwriting existing user files by default and shall require explicit opt-in for destructive replacement behavior.

Depends on:
- [REQ-F-004](#req-f-004)
- [REQ-F-005](#req-f-005)

#### <a id="req-f-011"></a>REQ-F-011 - Guided Preset Selection

The CLI shall support an explicit preset-selection flow based on the bundled preset list so that users can choose from the available presets before applying one.

Depends on:
- [REQ-F-003](#req-f-003)
- [REQ-F-004](#req-f-004)

#### <a id="req-f-012"></a>REQ-F-012 - Install Flow Behavior Matrix

The installer wizard (`opencode-helper-install`) shall install `opencode-helper` globally on macOS and Linux with zero prompts in non-interactive contexts (including piped `curl|sh` and `wget|sh`).

Prompt behavior:
- non-interactive mode (`--yes` or non-TTY stdin): no prompts
- interactive TTY mode: prompt only for install location when default `~/.local/bin` is inconvenient and `--bin-dir` is not provided

The installer shall detect the active user shell and update `PATH` by editing the correct shell config file when shell detection succeeds.

Shell config targets and selection rules:
- zsh: write `PATH` updates to `~/.zshrc`
- bash (macOS): write to the first match in order: `~/.bash_profile` (if exists), else `~/.bashrc` (if exists), else create and write `~/.bash_profile`
- bash (Linux): write to the first match in order: `~/.bashrc` (if exists), else `~/.bash_profile` (if exists), else create and write `~/.bashrc`
- fish: write to `~/.config/fish/config.fish` (create parent dir and file if missing)

Depends on:
- [REQ-NF-005](#req-nf-005)

#### <a id="req-f-013"></a>REQ-F-013 - Non-Interactive Install Flags

The installer wizard (`opencode-helper-install`) shall support non-interactive install flags suitable for CI/power users, including at least:
- `--yes`
- `--bin-dir <path>`

In non-interactive mode (`--yes` or non-TTY stdin), the CLI shall not prompt. Shell config choice shall be derived from `$SHELL` using the selection rules in [REQ-F-012](#req-f-012). If `$SHELL` is unset or unsupported, installation shall continue but PATH file edits shall be skipped with a warning and manual PATH guidance.

Optionally:
- `--dry-run` (if implemented): print planned actions (install path, sudo/no-sudo decision, shell config file path, and the exact `PATH` line/block to be added/updated) and make no changes; must not request `sudo`.

Depends on:
- [REQ-F-012](#req-f-012)

#### <a id="req-f-014"></a>REQ-F-014 - Idempotent Shell Config Edits and Privilege Boundaries

The installer wizard (`opencode-helper-install`) shall make idempotent shell config edits using stable markers and shall avoid duplicate `PATH` entries.

Shell config edits shall be managed via a single marker block with these exact sentinel lines:
- `# opencode-helper: BEGIN managed PATH`
- `# opencode-helper: END managed PATH`

On re-run, if the marker block exists, the CLI shall update the contents within the block in place (and shall not add a second block). The resulting config must include exactly one occurrence of the chosen `--bin-dir` in the `PATH` update it writes.

The CLI shall use `sudo` only when the chosen install directory is privileged. For this requirement, a privileged install directory is any `--bin-dir` that:
- is under one of: `/usr/local`, `/opt/homebrew`, `/usr/bin`, `/bin`, `/sbin`, `/opt`, or
- is not writable by the invoking user

When `sudo` is required, it shall apply only to writing the installed binary into the privileged `--bin-dir`; shell config edits shall be performed as the invoking user.

Depends on:
- [REQ-F-012](#req-f-012)

#### <a id="req-f-015"></a>REQ-F-015 - Release Bundle Installation Source

The installer wizard (`opencode-helper-install`) shall install official helper releases by fetching versioned release assets from GitHub Releases rather than copying repo-local files.

The default install target shall be the latest supported GitHub release.

The installed release bundle shall contain:
- the `opencode-helper` executable entrypoint
- all officially supported bundled config presets for that release
- the official bundled inter-agent schemas for that release
- release metadata sufficient to report installed release provenance

Depends on:
- [DEC-001](#dec-001)
- [REQ-F-001](#req-f-001)
- [REQ-F-002](#req-f-002)

#### <a id="req-f-016"></a>REQ-F-016 - Release Bundle Checksum Verification

Before installing a downloaded release bundle, the installer wizard (`opencode-helper-install`) shall verify the bundle checksum against a checksum manifest published with the same GitHub release.

Verification requirements:
- the checksum algorithm shall be SHA-256
- the installer shall fail closed if the checksum manifest is missing, unreadable, lacks a matching bundle entry, or the computed checksum does not match the published checksum
- no bundle contents shall be activated for use if checksum verification fails

Depends on:
- [REQ-F-015](#req-f-015)

#### <a id="req-f-017"></a>REQ-F-017 - Release Selection

The installer wizard (`opencode-helper-install`) shall install the latest supported release by default and shall also support explicit installation of a user-selected older release tag.

At minimum, the installer shall support:
- default latest-release installation behavior
- a release selection flag `--version <tag>`

Depends on:
- [REQ-F-015](#req-f-015)

#### <a id="req-f-018"></a>REQ-F-018 - Installed Release Visibility and Update Awareness

The installer wizard (`opencode-helper-install`) and the helper CLI shall make the installed release visible to the user.

At minimum:
- the installer shall print the target release before installation and the active installed release after installation
- the helper CLI shall report the installed release tag together with bundled asset provenance metadata
- the helper CLI should indicate whether the active release is the latest supported release when network access is available for that check

Depends on:
- [REQ-F-015](#req-f-015)
- [REQ-F-009](#req-f-009)

#### <a id="req-f-019"></a>REQ-F-019 - Side-by-Side Versioned Installs and Rollback

The installer wizard (`opencode-helper-install`) shall install each helper release into a versioned side-by-side location and activate one release via a stable `current` symlink or equivalent indirection.

Rollback requirements:
- installing a newer release shall not require deleting previously installed releases by default
- activation of the selected release shall update the stable launcher/symlink to the chosen versioned install root
- users shall be able to revert to an already installed older release without re-downloading assets if that older release remains present locally
- helper CLI shall expose `opencode-helper release use <tag>` to reactivate an already installed local release without network access

Depends on:
- [REQ-F-015](#req-f-015)
- [REQ-F-017](#req-f-017)

#### <a id="req-f-020"></a>REQ-F-020 - Bootstrap Install Script

The helper CLI shall publish a standalone bootstrap install script (`install.sh`) as a separate release asset alongside the release bundle, enabling the standard `curl|sh` one-liner distribution:

```sh
curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh
```

The bootstrap script shall:

- Be self-contained and small (under 100 lines of shell)
- Use only `curl` or `wget` (whichever is available) to download assets
- Accept the same flags as `opencode-helper-install` (`--bin-dir`, `--version`)
- Accept `OPENCODE_HELPER_VERSION` as an environment variable to pin the release tag
- Fetch `opencode-helper-<tag>-checksums.txt` from the same release and verify the downloaded `opencode-helper-install` before executing it
- Pass through all arguments to `opencode-helper-install --yes`

The bootstrap script is itself a release asset. It is:
- Included in the release tarball bundle at the bundle root
- Also uploaded as a separate release asset file so it can be downloaded independently from the GitHub Releases page

The bootstrap script shall NOT download the tarball. Its single responsibility is to download, verify, and delegate to `opencode-helper-install`.

Depends on:
- [DEC-001](#dec-001)
- [REQ-F-015](#req-f-015)
- [REQ-F-016](#req-f-016)

#### <a id="req-f-021"></a>REQ-F-021 - Config Source Registry

The CLI shall support registering one or more config sources in a global (user-level) registry so that config bundles can be discovered and installed.

Registry items must include at least:
- a source location
- a source type, or enough information for the CLI to determine the source type during validation
- a stable local identifier derived from the source manifest or assigned by the CLI

Depends on:
- [DEC-002](#dec-002)

#### <a id="req-f-022"></a>REQ-F-022 - Config Bundle Manifest Contract

Each config bundle shall include a machine-readable manifest file at a fixed path:
- `opencode-bundle.manifest.json`

The manifest shall use JSON and a versioned top-level structure that the CLI can validate.

Required top-level fields:
- `manifest_version` - manifest format version understood by the CLI
- `bundle_name` - stable bundle identifier
- `bundle_version` - bundle release/tag or other source-level version identifier
- `presets` - array of preset descriptors

Each preset descriptor shall include:
- `name` - stable preset identifier
- `description` - short human-readable description
- `entrypoint` - path to the preset config file within the bundle
- `prompt_files` - array of prompt file paths that must be materialized when the preset is applied

Optional provenance fields:
- `source_repo`
- `source_commit`
- `release_tag`

The CLI shall reject a bundle when:
- the manifest file is missing
- the manifest is not valid JSON
- `manifest_version` is unsupported
- any required field is missing
- a declared `entrypoint` or `prompt_files` path does not exist in the bundle

Note: a published JSON Schema for this manifest may be added later, but the file location, field names, and validation rules above are normative for V2+ behavior.

Example manifest:

```json
{
  "manifest_version": 1,
  "bundle_name": "official-openai",
  "bundle_version": "v2.0.0",
  "source_repo": "https://github.com/example/opencode-configs",
  "source_commit": "abc1234",
  "release_tag": "v2.0.0",
  "presets": [
    {
      "name": "openai",
      "description": "OpenAI-based planning-first preset",
      "entrypoint": "presets/opencode.openai.json",
      "prompt_files": [
        "prompts/coding-boss.txt",
        "prompts/planner.txt",
        "prompts/code-reviewer.txt"
      ]
    }
  ]
}
```

Depends on:
- [DEC-002](#dec-002)

---

### <a id="bundle-manifest-reference"></a>Bundle Manifest Reference

The `opencode-bundle.manifest.json` file is the core contract between the config bundle and the V2 CLI. It must be present at the root of every V2-compatible bundle.

#### File Location

```
<bundle-root>/
  opencode-bundle.manifest.json  <- required at bundle root
  opencode.openai.json           <- preset files referenced by manifest
  opencode.mixed.json
  ...
```

#### Required Structure

```json
{
  "manifest_version": 1,
  "bundle_name": "opencode-helper",
  "bundle_version": "v1.0.0",
  "presets": [
    {
      "name": "openai",
      "description": "OpenAI GPT-5 based coding and docs agents",
      "entrypoint": "opencode.openai.json",
      "prompt_files": []
    }
  ]
}
```

#### Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `manifest_version` | integer | Yes | Currently must be `1`. CLI rejects unsupported versions. |
| `bundle_name` | string | Yes | Stable identifier for the bundle (e.g., "opencode-helper") |
| `bundle_version` | string | Yes | Release tag or version string (e.g., "v1.0.0") |
| `presets` | array | Yes | Array of preset descriptor objects |

#### Preset Descriptor

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Stable identifier (lowercase, hyphens allowed) |
| `description` | string | Yes | Short human-readable description |
| `entrypoint` | string | Yes | Path to preset JSON file relative to bundle root |
| `prompt_files` | array | Yes | Array of prompt file paths to materialize (can be empty) |

#### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `source_repo` | string | URL of source repository |
| `source_commit` | string | Git commit SHA |
| `release_tag` | string | Release tag (often matches `bundle_version`) |

#### Current Bundle Example

The opencode-helper release bundle currently includes 5 presets:

```json
{
  "manifest_version": 1,
  "bundle_name": "opencode-helper",
  "bundle_version": "v1.0.0",
  "presets": [
    {"name": "openai", "description": "OpenAI GPT-5 based coding and docs agents", "entrypoint": "opencode.openai.json", "prompt_files": []},
    {"name": "mixed", "description": "Claude + GPT-5 hybrid coding and docs agents", "entrypoint": "opencode.mixed.json", "prompt_files": []},
    {"name": "big-pickle", "description": "Big-Pickle model based coding and docs agents", "entrypoint": "opencode.big-pickle.json", "prompt_files": []},
    {"name": "minimax", "description": "MiniMax M2.5 based coding and docs agents", "entrypoint": "opencode.minimax.json", "prompt_files": []},
    {"name": "kimi", "description": "Kimi K2.5 based coding and docs agents", "entrypoint": "opencode.kimi.json", "prompt_files": []}
  ]
}
```

#### CLI Validation Rules

The CLI rejects a bundle when:
- `opencode-bundle.manifest.json` is missing from bundle root
- Manifest is not valid JSON
- `manifest_version` is unsupported (not `1`)
- Any required top-level field is missing (`manifest_version`, `bundle_name`, `bundle_version`, `presets`)
- Any preset is missing required fields (`name`, `description`, `entrypoint`, `prompt_files`)
- An `entrypoint` path does not exist in the bundle
- A `prompt_files` entry does not exist in the bundle

#### Generating the Manifest

For the opencode-helper official bundle, the manifest is generated by `scripts/release/build-opencode-helper-bundle.sh`. Third-party bundle creators should include a manually-authored `opencode-bundle.manifest.json` at their bundle root following this contract.

---

#### <a id="req-f-023"></a>REQ-F-023 - Preset Discovery Across Config Sources

The CLI shall let users list presets across all registered config sources and shall display each preset together with:
- short description
- source reference suitable for later apply commands (at minimum the source identifier, and when available a human-friendly source name)
- bundle release/version

Depends on:
- [REQ-F-021](#req-f-021)
- [REQ-F-022](#req-f-022)

#### <a id="req-f-024"></a>REQ-F-024 - Resolve, Install, and Apply a Bundle

The CLI shall let users choose a config source and bundle or version selector, resolve/cache it locally, apply a selected preset from that bundle into a target project, and persist bundle provenance for later status/update behavior.

Human-centered interaction requirements:
- source selection for user-facing commands shall accept a source reference that may be either the registered source identifier or a unique human-friendly source name
- when a human-friendly source name matches multiple registered sources, the CLI shall fail closed with an ambiguity error that lists the matching sources instead of choosing implicitly
- in TTY mode, when a command has enough context to guide a human through preset selection safely, the CLI should default to prompting instead of requiring the user to memorize bundle identifiers and preset names
- commands that support this guided TTY behavior shall provide an explicit non-interactive opt-out such as `--auto` so automation remains deterministic
- outside TTY mode, commands shall not prompt implicitly and shall instead fail with an actionable error when required selection input is missing

Depends on:
- [REQ-F-021](#req-f-021)
- [REQ-F-022](#req-f-022)
- [REQ-F-010](#req-f-010)

#### <a id="req-f-025"></a>REQ-F-025 - Install Referenced Prompt Files

When applying a preset from a config bundle, the CLI shall also materialize any prompt files referenced by that preset into the target project so that the resulting `opencode.json` can resolve `{file:...}` references locally.

Depends on:
- [REQ-F-024](#req-f-024)

#### <a id="req-f-026"></a>REQ-F-026 - Bundle Provenance Reporting

The CLI shall report previously persisted provenance for installed config bundles and for applied presets/prompts, including at least:
- source identifier and source location
- installed bundle release/version
- a source commit identifier when available

Depends on:
- [REQ-F-021](#req-f-021)
- [REQ-F-024](#req-f-024)

#### <a id="req-f-027"></a>REQ-F-027 - Bundle Update Check with User Prompt

For update-capable config sources, the CLI shall detect when a newer compatible bundle release is available for an installed bundle and shall prompt the user before performing a bundle update.

Depends on:
- [REQ-F-021](#req-f-021)
- [REQ-F-024](#req-f-024)
- [REQ-F-026](#req-f-026)

#### <a id="req-f-028"></a>REQ-F-028 - Supported Bundle Source Types

The first V2+ scope baseline shall support at least these config source types:
- local directory source
- local `.tar.gz` archive source
- GitHub release bundle source

Depends on:
- [DEC-002](#dec-002)

#### <a id="req-f-029"></a>REQ-F-029 - Bundle Resolution Normalization

The CLI shall resolve every supported config source into a normalized local bundle root before validation and apply.

Normalization requirements:
- exactly one `opencode-bundle.manifest.json` must be discoverable for the resolved bundle
- the CLI shall tolerate bundle content at the source root or nested under a single extracted top-level directory
- after resolution, manifest validation and apply behavior shall be source-type independent

Depends on:
- [REQ-F-022](#req-f-022)
- [REQ-F-024](#req-f-024)
- [REQ-F-028](#req-f-028)

#### <a id="req-f-030"></a>REQ-F-030 - V1 to V2 Upgrade Compatibility and Migration

The CLI shall allow an existing V1 user to upgrade to a V2+-capable CLI via the existing `self-update` flow without breaking an already working project setup.

Upgrade compatibility requirements:
- `self-update` shall continue to upgrade the CLI in place using the existing supported workflow
- existing V1 project setups shall remain usable after upgrading the CLI
- the V2+ CLI shall detect legacy bundled-preset project state when relevant and provide clear guidance for adopting config sources
- if the CLI supports migration from legacy bundled-preset state into config-source-managed state, that migration shall require explicit user action and shall not be forced implicitly during CLI upgrade

Depends on:
- [REQ-F-008](#req-f-008)
- [REQ-F-024](#req-f-024)
- [DEC-002](#dec-002)

#### <a id="req-f-031"></a>REQ-F-031 - Config Source Capability Model

The CLI shall classify each registered config source by capability so source-type-specific behavior is explicit and predictable.

At minimum, the capability model shall distinguish whether a source supports:
- bundle discovery
- version selection
- update checks

Commands that depend on a capability shall surface a clear error when invoked against a source that does not support it.

Depends on:
- [REQ-F-021](#req-f-021)
- [REQ-F-028](#req-f-028)

#### <a id="req-f-032"></a>REQ-F-032 - Remote Bundle Integrity Verification

For remote config sources, the CLI shall verify remote bundle integrity before applying or updating a bundle.

Verification requirements:
- the remote source shall expose integrity metadata sufficient for the CLI to verify the downloaded bundle
- bundle apply or update shall fail closed if integrity metadata is missing, unreadable, or does not match the downloaded bundle
- no downloaded remote bundle contents shall be activated for use if integrity verification fails

Depends on:
- [REQ-F-024](#req-f-024)
- [REQ-F-027](#req-f-027)
- [REQ-F-028](#req-f-028)

### Non-Functional Requirements

#### <a id="req-nf-001"></a>REQ-NF-001 - Minimal Footprint

The CLI should remain small and operationally simple, with minimal runtime dependencies.

#### <a id="req-nf-002"></a>REQ-NF-002 - Offline-Friendly Operation

Normal project setup after CLI installation should not require network access.

Depends on:
- [DEC-001](#dec-001)

#### <a id="req-nf-003"></a>REQ-NF-003 - Deterministic Release Bundle

Each CLI release shall pin a specific compatible set of bundled presets and schemas.

Depends on:
- [DEC-001](#dec-001)

#### <a id="req-nf-004"></a>REQ-NF-004 - Clear Automation Semantics

The CLI should provide stable exit codes and machine-friendly command behavior for validation and diagnostics flows.

#### <a id="req-nf-005"></a>REQ-NF-005 - Cross-Platform Baseline

V1 should support macOS and Linux.

#### <a id="req-nf-006"></a>REQ-NF-006 - Traceable Provenance

Installed assets should be attributable to the CLI release that provided them.

Depends on:
- [REQ-F-009](#req-f-009)

#### <a id="req-nf-007"></a>REQ-NF-007 - Installer Output UX (Theme, Icons, Accessibility)

The installer wizard (`opencode-helper-install`) should provide a pleasant, themed terminal experience that is still accessible and robust across terminals and automation.

Requirements:
- Use a small, consistent set of plain Unicode status icons as line prefixes (e.g. `✓` ok, `!` warning, `x` error, `?` question, `i` info).
- Color is optional, but output shall never rely on color alone; each status line must include a textual label (e.g. `ok:`, `warn:`, `error:`).
- Icons shall be used sparingly (at most one status icon per line, as a prefix).
- Support explicit fallbacks:
  - `--no-color` disables ANSI colors.
  - `--ascii` disables Unicode icons (and uses ASCII-only prefixes such as `[ok]`, `[warn]`, `[err]`, `[?]`, `[i]`).
- When stdout is not a TTY (e.g. piped output / CI logs), color and Unicode icons should be disabled by default.

#### <a id="req-nf-008"></a>REQ-NF-008 - Deterministic Release Packaging

Each published helper release shall produce a deterministic release bundle, release metadata manifest, and checksum manifest from the tagged source tree used for that release.

The release packaging process shall ensure that the bundled helper executable, presets, schemas, and provenance metadata all correspond to the same tagged release.

Depends on:
- [REQ-NF-003](#req-nf-003)
- [REQ-F-015](#req-f-015)
- [REQ-F-016](#req-f-016)

#### <a id="req-nf-009"></a>REQ-NF-009 - Auditable Installed Provenance

The installed helper should expose enough local provenance metadata to identify the active release tag, source commit, and bundled asset set without requiring access to the source repository.

Depends on:
- [REQ-F-018](#req-f-018)
- [REQ-NF-006](#req-nf-006)

---

## First Features

### <a id="feat-001"></a>FEAT-001 - Project Initialization

Description:
- Bootstrap a target project using bundled official assets and, when needed, guide the user through choosing from the available presets

Likely command shape:
- `opencode-helper init`

Satisfies:
- [REQ-F-001](#req-f-001)
- [REQ-F-001a](#req-f-001a)
- [REQ-F-002](#req-f-002)
- [REQ-F-003](#req-f-003)
- [REQ-F-004](#req-f-004)
- [REQ-F-005](#req-f-005)
- [REQ-F-010](#req-f-010)
- [REQ-F-011](#req-f-011)

### <a id="feat-002"></a>FEAT-002 - Preset Listing

Description:
- List bundled config presets included in the current CLI release together with each preset's short purpose

Likely command shape:
- `opencode-helper preset list`

Satisfies:
- [REQ-F-003](#req-f-003)
- [REQ-F-001a](#req-f-001a)
- [REQ-F-009](#req-f-009)

### <a id="feat-003"></a>FEAT-003 - Preset Selection and Apply

Description:
- Let users choose from the available presets and apply the selected bundled preset to a target project

Likely command shape:
- `opencode-helper preset use <name>`

Satisfies:
- [REQ-F-004](#req-f-004)
- [REQ-F-010](#req-f-010)
- [REQ-F-011](#req-f-011)

### <a id="feat-004"></a>FEAT-004 - Schema Install

Description:
- Install the bundled handoff/result schemas into the project or supported scope

Likely command shape:
- `opencode-helper schema install`

Satisfies:
- [REQ-F-005](#req-f-005)
- [REQ-F-010](#req-f-010)

### <a id="feat-005"></a>FEAT-005 - Validation

Description:
- Validate that the local project setup matches supported expectations

Likely command shape:
- `opencode-helper validate`

Satisfies:
- [REQ-F-006](#req-f-006)
- [REQ-F-009](#req-f-009)
- [REQ-NF-004](#req-nf-004)
- [REQ-NF-006](#req-nf-006)

### <a id="feat-006"></a>FEAT-006 - Doctor

Description:
- Diagnose drift, missing files, incompatible state, and likely remediation paths

Likely command shape:
- `opencode-helper doctor`

Satisfies:
- [REQ-F-007](#req-f-007)
- [REQ-NF-004](#req-nf-004)

### <a id="feat-007"></a>FEAT-007 - Version

Description:
- Report CLI version and bundled asset identity

Likely command shape:
- `opencode-helper version`

Satisfies:
- [REQ-F-009](#req-f-009)
- [REQ-NF-006](#req-nf-006)

### <a id="feat-008"></a>FEAT-008 - Self Update

Description:
- Update the CLI from the latest GitHub release, including its bundled asset set

Likely command shape:
- `opencode-helper self-update`

Satisfies:
- [REQ-F-008](#req-f-008)
- [REQ-F-009](#req-f-009)
- [REQ-NF-003](#req-nf-003)

### <a id="feat-009"></a>FEAT-009 - Install Wizard

Description:
- Install `opencode-helper` globally with zero-prompt piped install support and TTY-only install-dir prompting when the default path is inconvenient

Likely command shape:
- `opencode-helper-install`
- `opencode-helper release use <tag>`

Satisfies:
- [REQ-F-012](#req-f-012)
- [REQ-F-013](#req-f-013)
- [REQ-F-014](#req-f-014)
- [REQ-F-015](#req-f-015)
- [REQ-F-016](#req-f-016)
- [REQ-F-017](#req-f-017)
- [REQ-F-018](#req-f-018)
- [REQ-F-019](#req-f-019)
- [REQ-NF-005](#req-nf-005)
- [REQ-NF-007](#req-nf-007)

### <a id="feat-010"></a>FEAT-010 - Release Packaging and Provenance

Description:
- Build and publish deterministic helper release bundles, metadata manifests, and checksum manifests so installer and runtime provenance remain traceable and verifiable

Likely command shape:
- release pipeline / GitHub Actions workflow

Satisfies:
- [REQ-F-015](#req-f-015)
- [REQ-F-016](#req-f-016)
- [REQ-F-018](#req-f-018)
- [REQ-NF-003](#req-nf-003)
- [REQ-NF-008](#req-nf-008)
- [REQ-NF-009](#req-nf-009)

### <a id="feat-011"></a>FEAT-011 - Config Source Management

Description:
- Manage a user-level registry of config sources (add/remove/list)

Likely command shape:
- `opencode-helper source add <location> [--type <type>]`
- `opencode-helper source list`
- `opencode-helper source remove <source-id>`

Satisfies:
- [REQ-F-021](#req-f-021)
- [REQ-F-028](#req-f-028)

### <a id="feat-012"></a>FEAT-012 - Bundle Install and Apply

Description:
- Resolve/cache a chosen config bundle from a registered config source and apply a preset from that bundle into a target project

Likely command shape:
- `opencode-helper bundle install <source-id> [--version <tag>]`
- `opencode-helper bundle apply <source-ref> [--version <tag>] [--preset <name>] [--auto]`

Satisfies:
- [REQ-F-022](#req-f-022)
- [REQ-F-023](#req-f-023)
- [REQ-F-024](#req-f-024)
- [REQ-F-025](#req-f-025)
- [REQ-F-029](#req-f-029)
- [REQ-F-010](#req-f-010)

### <a id="feat-013"></a>FEAT-013 - Bundle Updates

Description:
- Check for newer compatible bundle releases and prompt the user before updating installed bundles

Likely command shape:
- `opencode-helper bundle update <source-id> [--yes]`

Satisfies:
- [REQ-F-027](#req-f-027)
- [REQ-F-031](#req-f-031)

### <a id="feat-014"></a>FEAT-014 - Bundle Provenance

Description:
- Persist and report provenance for installed bundles and applied presets/prompts

Likely command shape:
- `opencode-helper bundle status`

Satisfies:
- [REQ-F-024](#req-f-024)
- [REQ-F-026](#req-f-026)

### <a id="feat-015"></a>FEAT-015 - Source-Type Resolution

Description:
- Resolve supported source types (local directory, local `.tar.gz` archive, GitHub release bundle) into a normalized local bundle root before validation and apply

Likely command shape:
- internal resolution layer used by `source add`, `bundle install`, and `bundle apply`

Satisfies:
- [REQ-F-028](#req-f-028)
- [REQ-F-029](#req-f-029)

### <a id="feat-016"></a>FEAT-016 - Legacy Upgrade and Migration

Description:
- Preserve working V1 setups through CLI self-update and provide legacy-state detection plus migration guidance

Likely command shape:
- `opencode-helper self-update`
- `opencode-helper migrate legacy-config`

Satisfies:
- [REQ-F-030](#req-f-030)

### <a id="feat-017"></a>FEAT-017 - Source Capability Handling

Description:
- Determine and surface source capabilities so version selection and update behavior only appear where supported

Likely command shape:
- capability detection during `source add`
- capability-aware behavior in `source list`, `bundle apply`, and `bundle update`

Satisfies:
- [REQ-F-031](#req-f-031)

### <a id="feat-018"></a>FEAT-018 - Remote Bundle Verification

Description:
- Verify remote bundle integrity before applying or updating remote bundles

Likely command shape:
- verification step inside remote `bundle apply` and `bundle update` flows

Satisfies:
- [REQ-F-032](#req-f-032)

---

## Traceability Matrix

| ID | Type | Links To |
|---|---|---|
| [DEC-001](#dec-001) | Decision | [PRD-001](#prd-001), [REQ-F-001](#req-f-001), [REQ-F-002](#req-f-002), [REQ-F-008](#req-f-008), [REQ-NF-002](#req-nf-002), [REQ-NF-003](#req-nf-003) |
| [DEC-002](#dec-002) | Decision | [PRD-001](#prd-001), [REQ-F-021](#req-f-021) to [REQ-F-032](#req-f-032) |
| [PRD-001](#prd-001) | PRD | [REQ-F-001](#req-f-001) to [REQ-F-032](#req-f-032), [REQ-NF-001](#req-nf-001) to [REQ-NF-009](#req-nf-009) |
| [REQ-F-001](#req-f-001) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-002](#feat-002), [FEAT-003](#feat-003) |
| [REQ-F-001a](#req-f-001a) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-002](#feat-002) |
| [REQ-F-002](#req-f-002) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-004](#feat-004) |
| [REQ-F-003](#req-f-003) | Functional Requirement | [FEAT-002](#feat-002) |
| [REQ-F-004](#req-f-004) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-003](#feat-003) |
| [REQ-F-005](#req-f-005) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-004](#feat-004) |
| [REQ-F-006](#req-f-006) | Functional Requirement | [FEAT-005](#feat-005) |
| [REQ-F-007](#req-f-007) | Functional Requirement | [FEAT-006](#feat-006) |
| [REQ-F-008](#req-f-008) | Functional Requirement | [FEAT-008](#feat-008), [US-023](#us-023) |
| [REQ-F-009](#req-f-009) | Functional Requirement | [FEAT-002](#feat-002), [FEAT-005](#feat-005), [FEAT-007](#feat-007), [FEAT-008](#feat-008) |
| [REQ-F-010](#req-f-010) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-003](#feat-003), [FEAT-004](#feat-004) |
| [REQ-F-011](#req-f-011) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-003](#feat-003) |
| [REQ-F-012](#req-f-012) | Functional Requirement | [FEAT-009](#feat-009) |
| [REQ-F-013](#req-f-013) | Functional Requirement | [FEAT-009](#feat-009) |
| [REQ-F-014](#req-f-014) | Functional Requirement | [FEAT-009](#feat-009) |
| [REQ-F-015](#req-f-015) | Functional Requirement | [FEAT-009](#feat-009), [FEAT-010](#feat-010), [US-016](#us-016), [US-017](#us-017), [US-018](#us-018), [US-019](#us-019), [US-020](#us-020) |
| [REQ-F-016](#req-f-016) | Functional Requirement | [FEAT-009](#feat-009), [FEAT-010](#feat-010), [US-016](#us-016), [US-017](#us-017) |
| [REQ-F-017](#req-f-017) | Functional Requirement | [FEAT-009](#feat-009), [US-018](#us-018), [US-020](#us-020) |
| [REQ-F-018](#req-f-018) | Functional Requirement | [FEAT-009](#feat-009), [FEAT-010](#feat-010), [US-019](#us-019) |
| [REQ-F-019](#req-f-019) | Functional Requirement | [FEAT-009](#feat-009), [US-020](#us-020), [US-023](#us-023) |
| [REQ-F-020](#req-f-020) | Functional Requirement | [FEAT-009](#feat-009), [US-021](#us-021) |
| [REQ-F-021](#req-f-021) | Functional Requirement | [FEAT-011](#feat-011), [US-026](#us-026) |
| [REQ-F-022](#req-f-022) | Functional Requirement | [FEAT-012](#feat-012), [US-027](#us-027), [US-032](#us-032) |
| [REQ-F-023](#req-f-023) | Functional Requirement | [FEAT-012](#feat-012), [US-027](#us-027) |
| [REQ-F-024](#req-f-024) | Functional Requirement | [FEAT-012](#feat-012), [FEAT-014](#feat-014), [US-028](#us-028), [US-037](#us-037) |
| [REQ-F-025](#req-f-025) | Functional Requirement | [FEAT-012](#feat-012), [US-028](#us-028), [US-029](#us-029) |
| [REQ-F-026](#req-f-026) | Functional Requirement | [FEAT-014](#feat-014), [US-030](#us-030) |
| [REQ-F-027](#req-f-027) | Functional Requirement | [FEAT-013](#feat-013), [US-031](#us-031) |
| [REQ-F-028](#req-f-028) | Functional Requirement | [FEAT-011](#feat-011), [FEAT-015](#feat-015), [US-033](#us-033), [US-034](#us-034), [US-035](#us-035) |
| [REQ-F-029](#req-f-029) | Functional Requirement | [FEAT-012](#feat-012), [FEAT-015](#feat-015), [US-028](#us-028), [US-034](#us-034), [US-035](#us-035) |
| [REQ-F-030](#req-f-030) | Functional Requirement | [FEAT-016](#feat-016), [US-036](#us-036), [US-037](#us-037) |
| [REQ-F-031](#req-f-031) | Functional Requirement | [FEAT-013](#feat-013), [FEAT-017](#feat-017), [US-031](#us-031), [US-038](#us-038) |
| [REQ-F-032](#req-f-032) | Functional Requirement | [FEAT-018](#feat-018), [US-031](#us-031), [US-039](#us-039) |
| [REQ-NF-001](#req-nf-001) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-008](#feat-008) |
| [REQ-NF-002](#req-nf-002) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-007](#feat-007) |
| [REQ-NF-003](#req-nf-003) | Non-Functional Requirement | [FEAT-008](#feat-008), [FEAT-010](#feat-010) |
| [REQ-NF-004](#req-nf-004) | Non-Functional Requirement | [FEAT-005](#feat-005), [FEAT-006](#feat-006) |
| [REQ-NF-005](#req-nf-005) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-009](#feat-009) |
| [REQ-NF-006](#req-nf-006) | Non-Functional Requirement | [FEAT-005](#feat-005), [FEAT-007](#feat-007) |
| [REQ-NF-007](#req-nf-007) | Non-Functional Requirement | [FEAT-009](#feat-009) |
| [REQ-NF-008](#req-nf-008) | Non-Functional Requirement | [FEAT-010](#feat-010), [US-016](#us-016) |
| [REQ-NF-009](#req-nf-009) | Non-Functional Requirement | [FEAT-010](#feat-010), [US-019](#us-019) |
| [FEAT-001](#feat-001) | Feature | [REQ-F-001](#req-f-001), [REQ-F-001a](#req-f-001a), [REQ-F-002](#req-f-002), [REQ-F-003](#req-f-003), [REQ-F-004](#req-f-004), [REQ-F-005](#req-f-005), [REQ-F-010](#req-f-010), [REQ-F-011](#req-f-011), [US-002](#us-002), [US-008](#us-008), [US-009](#us-009), [US-010](#us-010) |
| [FEAT-002](#feat-002) | Feature | [REQ-F-001](#req-f-001), [REQ-F-001a](#req-f-001a), [REQ-F-003](#req-f-003), [REQ-F-009](#req-f-009), [US-001](#us-001), [US-010](#us-010) |
| [FEAT-003](#feat-003) | Feature | [REQ-F-004](#req-f-004), [REQ-F-010](#req-f-010), [REQ-F-011](#req-f-011), [US-003](#us-003), [US-008](#us-008), [US-009](#us-009), [US-010](#us-010) |
| [FEAT-004](#feat-004) | Feature | [REQ-F-002](#req-f-002), [REQ-F-005](#req-f-005), [REQ-F-010](#req-f-010), [US-004](#us-004), [US-008](#us-008), [US-010](#us-010) |
| [FEAT-005](#feat-005) | Feature | [REQ-F-006](#req-f-006), [REQ-F-009](#req-f-009), [REQ-NF-004](#req-nf-004), [REQ-NF-006](#req-nf-006), [US-005](#us-005), [US-010](#us-010) |
| [FEAT-006](#feat-006) | Feature | [REQ-F-007](#req-f-007), [REQ-NF-004](#req-nf-004), [US-006](#us-006) |
| [FEAT-007](#feat-007) | Feature | [REQ-F-009](#req-f-009), [REQ-NF-006](#req-nf-006), [US-007](#us-007), [US-011](#us-011) |
| [FEAT-008](#feat-008) | Feature | [REQ-F-008](#req-f-008), [REQ-F-009](#req-f-009), [REQ-NF-003](#req-nf-003), [US-023](#us-023) |
| [FEAT-009](#feat-009) | Feature | [REQ-F-012](#req-f-012), [REQ-F-013](#req-f-013), [REQ-F-014](#req-f-014), [REQ-F-015](#req-f-015), [REQ-F-016](#req-f-016), [REQ-F-017](#req-f-017), [REQ-F-018](#req-f-018), [REQ-F-019](#req-f-019), [REQ-F-020](#req-f-020), [US-012](#us-012), [US-013](#us-013), [US-014](#us-014), [US-015](#us-015), [US-017](#us-017), [US-018](#us-018), [US-019](#us-019), [US-020](#us-020), [US-021](#us-021) |
| [FEAT-010](#feat-010) | Feature | [REQ-F-015](#req-f-015), [REQ-F-016](#req-f-016), [REQ-F-018](#req-f-018), [REQ-NF-008](#req-nf-008), [REQ-NF-009](#req-nf-009), [US-016](#us-016), [US-019](#us-019) |
| [FEAT-011](#feat-011) | Feature | [REQ-F-021](#req-f-021), [REQ-F-028](#req-f-028), [US-026](#us-026), [US-033](#us-033) |
| [FEAT-012](#feat-012) | Feature | [REQ-F-022](#req-f-022), [REQ-F-023](#req-f-023), [REQ-F-024](#req-f-024), [REQ-F-025](#req-f-025), [REQ-F-029](#req-f-029), [REQ-F-010](#req-f-010), [US-027](#us-027), [US-028](#us-028), [US-029](#us-029), [US-032](#us-032), [US-034](#us-034), [US-035](#us-035) |
| [FEAT-013](#feat-013) | Feature | [REQ-F-027](#req-f-027), [REQ-F-031](#req-f-031), [US-031](#us-031) |
| [FEAT-014](#feat-014) | Feature | [REQ-F-024](#req-f-024), [REQ-F-026](#req-f-026), [US-028](#us-028), [US-030](#us-030) |
| [FEAT-015](#feat-015) | Feature | [REQ-F-028](#req-f-028), [REQ-F-029](#req-f-029), [US-033](#us-033), [US-034](#us-034), [US-035](#us-035) |
| [FEAT-016](#feat-016) | Feature | [REQ-F-030](#req-f-030), [US-036](#us-036), [US-037](#us-037) |
| [FEAT-017](#feat-017) | Feature | [REQ-F-031](#req-f-031), [US-035](#us-035), [US-038](#us-038) |
| [FEAT-018](#feat-018) | Feature | [REQ-F-032](#req-f-032), [US-031](#us-031), [US-039](#us-039) |
| [US-001](#us-001) | User Story | [FEAT-002](#feat-002), [REQ-F-003](#req-f-003), [REQ-F-001a](#req-f-001a), [REQ-F-009](#req-f-009) |
| [US-002](#us-002) | User Story | [FEAT-001](#feat-001), [REQ-F-004](#req-f-004), [REQ-F-005](#req-f-005), [REQ-F-010](#req-f-010), [REQ-F-011](#req-f-011), [REQ-NF-002](#req-nf-002) |
| [US-003](#us-003) | User Story | [FEAT-003](#feat-003), [REQ-F-004](#req-f-004), [REQ-F-010](#req-f-010), [REQ-F-011](#req-f-011) |
| [US-004](#us-004) | User Story | [FEAT-004](#feat-004), [REQ-F-005](#req-f-005), [REQ-F-010](#req-f-010), [REQ-NF-002](#req-nf-002) |
| [US-005](#us-005) | User Story | [FEAT-005](#feat-005), [REQ-F-006](#req-f-006), [REQ-F-007](#req-f-007), [REQ-NF-004](#req-nf-004), [REQ-NF-006](#req-nf-006) |
| [US-006](#us-006) | User Story | [FEAT-006](#feat-006), [REQ-F-007](#req-f-007), [REQ-NF-004](#req-nf-004) |
| [US-007](#us-007) | User Story | [FEAT-007](#feat-007), [REQ-F-009](#req-f-009), [REQ-NF-006](#req-nf-006), [REQ-NF-003](#req-nf-003) |
| [US-008](#us-008) | User Story | [FEAT-001](#feat-001), [FEAT-003](#feat-003), [FEAT-004](#feat-004), [REQ-F-010](#req-f-010) |
| [US-009](#us-009) | User Story | [FEAT-001](#feat-001), [FEAT-003](#feat-003), [REQ-F-004](#req-f-004) |
| [US-010](#us-010) | User Story | [FEAT-001](#feat-001), [FEAT-002](#feat-002), [FEAT-003](#feat-003), [FEAT-004](#feat-004), [FEAT-005](#feat-005), [FEAT-007](#feat-007), [REQ-NF-002](#req-nf-002), [DEC-001](#dec-001) |
| [US-011](#us-011) | User Story | [FEAT-007](#feat-007), [REQ-NF-003](#req-nf-003), [REQ-F-001](#req-f-001), [REQ-F-002](#req-f-002) |
| [US-012](#us-012) | User Story | [FEAT-009](#feat-009), [REQ-F-012](#req-f-012), [REQ-F-014](#req-f-014) |
| [US-013](#us-013) | User Story | [FEAT-009](#feat-009), [REQ-F-012](#req-f-012), [REQ-F-013](#req-f-013) |
| [US-014](#us-014) | User Story | [FEAT-009](#feat-009), [REQ-NF-007](#req-nf-007) |
| [US-015](#us-015) | User Story | [FEAT-009](#feat-009), [REQ-NF-007](#req-nf-007) |
| [US-016](#us-016) | User Story | [FEAT-010](#feat-010), [REQ-F-015](#req-f-015), [REQ-F-016](#req-f-016), [REQ-NF-008](#req-nf-008) |
| [US-017](#us-017) | User Story | [FEAT-009](#feat-009), [REQ-F-015](#req-f-015), [REQ-F-016](#req-f-016) |
| [US-018](#us-018) | User Story | [FEAT-009](#feat-009), [REQ-F-017](#req-f-017) |
| [US-019](#us-019) | User Story | [FEAT-009](#feat-009), [FEAT-010](#feat-010), [REQ-F-018](#req-f-018), [REQ-NF-009](#req-nf-009) |
| [US-020](#us-020) | User Story | [FEAT-009](#feat-009), [REQ-F-017](#req-f-017), [REQ-F-019](#req-f-019) |
| [US-021](#us-021) | User Story | [FEAT-009](#feat-009), [REQ-F-020](#req-f-020) |
| [US-023](#us-023) | User Story | [FEAT-008](#feat-008), [REQ-F-008](#req-f-008), [REQ-F-019](#req-f-019) |
| [US-024](#us-024) | User Story | [FEAT-003](#feat-003), [REQ-F-003](#req-f-003), [REQ-F-004](#req-f-004), [REQ-F-010](#req-f-010) |
| [US-025](#us-025) | User Story | [FEAT-003](#feat-003), [REQ-F-004](#req-f-004), [REQ-NF-004](#req-nf-004) |
| [US-026](#us-026) | User Story | [FEAT-011](#feat-011), [REQ-F-021](#req-f-021) |
| [US-027](#us-027) | User Story | [FEAT-012](#feat-012), [REQ-F-022](#req-f-022), [REQ-F-023](#req-f-023) |
| [US-028](#us-028) | User Story | [FEAT-012](#feat-012), [REQ-F-024](#req-f-024), [REQ-F-025](#req-f-025) |
| [US-029](#us-029) | User Story | [FEAT-012](#feat-012), [REQ-F-025](#req-f-025) |
| [US-030](#us-030) | User Story | [FEAT-014](#feat-014), [REQ-F-026](#req-f-026) |
| [US-031](#us-031) | User Story | [FEAT-013](#feat-013), [FEAT-018](#feat-018), [REQ-F-027](#req-f-027), [REQ-F-031](#req-f-031), [REQ-F-032](#req-f-032) |
| [US-032](#us-032) | User Story | [FEAT-012](#feat-012), [REQ-F-022](#req-f-022) |
| [US-033](#us-033) | User Story | [FEAT-011](#feat-011), [FEAT-015](#feat-015), [REQ-F-028](#req-f-028) |
| [US-034](#us-034) | User Story | [FEAT-012](#feat-012), [FEAT-015](#feat-015), [REQ-F-028](#req-f-028), [REQ-F-029](#req-f-029) |
| [US-035](#us-035) | User Story | [FEAT-012](#feat-012), [FEAT-015](#feat-015), [FEAT-017](#feat-017), [REQ-F-028](#req-f-028), [REQ-F-029](#req-f-029), [REQ-F-031](#req-f-031) |
| [US-036](#us-036) | User Story | [FEAT-016](#feat-016), [REQ-F-030](#req-f-030), [REQ-F-008](#req-f-008) |
| [US-037](#us-037) | User Story | [FEAT-016](#feat-016), [REQ-F-030](#req-f-030), [REQ-F-024](#req-f-024) |
| [US-038](#us-038) | User Story | [FEAT-017](#feat-017), [REQ-F-031](#req-f-031) |
| [US-039](#us-039) | User Story | [FEAT-018](#feat-018), [REQ-F-032](#req-f-032) |
| [US-040](#us-040) | User Story | (Go Migration) |
| [US-041](#us-041) | User Story | (Go Migration) |
| [US-042](#us-042) | User Story | (Go Migration) |
| [US-043](#us-043) | User Story | (Go Migration) |
| [US-044](#us-044) | User Story | (Go Migration) |
| [US-045](#us-045) | User Story | (Go Migration) |
| [US-046](#us-046) | User Story | (Go Migration) |
| [US-047](#us-047) | User Story | (Go Migration) |
| [US-048](#us-048) | User Story | (Go Migration) |
| [US-049](#us-049) | User Story | (Go Migration) |
| [US-050](#us-050) | User Story | (Go Migration) |
| [US-051](#us-051) | User Story | (Go Migration) |
| [US-052](#us-052) | User Story | (Go Migration) |

---

### <a id="us-040"></a>US-040 - Set up Go project structure

Priority:
- P0

Status:
- Open

Type:
- Maintainer-facing

Story:
- As a maintainer, I want the CLI implemented in Go so that I can benefit from better JSON handling, testability, and a more robust CLI experience.

Acceptance criteria:
- Go module created with proper module path
- Cobra CLI framework set up with root command
- help and version subcommands work
- Project structure follows feature-based organization
- Dependencies: spf13/cobra, stretchr/testify, gojsonschema

---

### <a id="us-041"></a>US-041 - Add GitHub Actions CI

Priority:
- P1

Status:
- Open

Type:
- Maintainer-facing

Related stories:
- [US-040](#us-040)

Story:
- As a maintainer, I want automated CI so that I can catch issues early.

Acceptance criteria:
- GitHub Actions workflow runs on push/PR
- Runs go test ./...
- Runs go build ./...
- Runs golangci-lint

---

### <a id="us-042"></a>US-042 - Implement init command

Priority:
- P0

Status:
- Open

Type:
- User-facing

Related stories:
- [US-040](#us-040)

Story:
- As a user, I want to initialize a project with a preset and schemas in one step.

Acceptance criteria:
- `opencode-helper init [--project-root PATH] [--preset NAME] [--output PATH] [--force] [--dry-run]` works
- Copies preset file to project
- Installs schemas to .opencode/
- Writes helper manifest

---

### <a id="us-043"></a>US-043 - Implement preset list and preset use commands

Priority:
- P0

Status:
- Open

Type:
- User-facing

Related stories:
- [US-040](#us-040)

Story:
- As a user, I want to list and apply bundled presets.

Acceptance criteria:
- `opencode-helper preset list` lists all bundled presets with descriptions
- `opencode-helper preset use <name> [--project-root PATH] [--output PATH] [--force] [--dry-run]` applies a preset

---

### <a id="us-044"></a>US-044 - Implement schema install and validate commands

Priority:
- P0

Status:
- Closed - Out of Scope (2026-03-30)

Type:
- User-facing

Related stories:
- [US-040](#us-040)

Story:
- ~~As a user, I want to install and validate schemas in my project.~~

Note:
- Schema validation is a bundle maintainer concern, not a CLI concern. The CLI should remain agnostic to bundle implementation details. Schema validation was removed from scope.

---

### <a id="us-045"></a>US-045 - Implement source commands (add, list, remove)

Priority:
- P2

Status:
- Open

Type:
- User-facing

Related stories:
- [US-040](#us-040)

Story:
- As a user, I want to manage config sources.

Acceptance criteria:
- `opencode-helper source add <location> [--name NAME]` registers a source
- `opencode-helper source list` lists registered sources
- `opencode-helper source remove <id>` removes a source
- Stores registry in XDG compliance location (~/.config/opencode-helper/)

---

### <a id="us-046"></a>US-046 - Implement bundle commands (apply, status, update)

Priority:
- P2

Status:
- Open

Type:
- User-facing

Related stories:
- [US-040](#us-040), [US-045](#us-045)

Story:
- As a user, I want to apply bundles from registered sources and track provenance.

Acceptance criteria:
- `opencode-helper bundle apply <source-id> --preset <name> [--project-root PATH] [--force] [--dry-run]` applies a preset from a bundle
- `opencode-helper bundle status [--project-root PATH]` shows provenance
- `opencode-helper bundle update <source-id> [--yes]` checks for and applies updates

---

### <a id="us-047"></a>US-047 - Implement update command

Priority:
- P2

Status:
- Open

Type:
- User-facing

Related stories:
- [US-040](#us-040)

Story:
- As a user, I want to check if a new version is available.

Acceptance criteria:
- `opencode-helper update` checks GitHub for newer version
- Reports current version vs latest available
- Tells user to upgrade via package manager

---

### <a id="us-048"></a>US-048 - Add shell completions

Priority:
- P3

Status:
- Open

Type:
- User-facing

Related stories:
- [US-040](#us-040)

Story:
- As a user, I want shell completions for faster CLI usage.

Acceptance criteria:
- `opencode-helper completion bash` generates bash completions
- `opencode-helper completion zsh` generates zsh completions
- `opencode-helper completion fish` generates fish completions

---

### <a id="us-049"></a>US-049 - Add human-first interactive defaults for TTY source-aware flows

Priority:
- P3

Status:
- In review

PR:
- [#141](https://github.com/sven1103-agent/opencode-config-cli/pull/141)
- TBD

Type:
- User-facing

Related features:
- [FEAT-012](#feat-012)

Related requirements:
- [REQ-F-023](#req-f-023)
- [REQ-F-024](#req-f-024)

Story:
- As a user working in a terminal, I want source-aware commands to guide me by default so that I can apply a preset without memorizing opaque source IDs and preset names.

Acceptance criteria:
- Commands that support guided source-aware selection use human-first behavior in TTY mode when required selection input is omitted.
- Commands that support guided source-aware selection provide an explicit non-interactive opt-out such as `--auto`.
- In non-TTY mode, commands remain non-interactive and require explicit arguments.
- When prompting is used in TTY mode, invalid selections produce a clear error and re-prompt.

---

### <a id="us-050"></a>US-050 - Set up goreleaser

Priority:
- P2

Status:
- Open

Type:
- Maintainer-facing

Related stories:
- [US-040](#us-040)

Story:
- As a maintainer, I want automated releases with goreleaser.

Acceptance criteria:
- .goreleaser.yaml configured
- Creates tar.gz archives with bundled presets
- Generates checksums
- Integrates with GitHub releases

---

### <a id="us-051"></a>US-051 - Create Homebrew tap

Priority:
- P2

Status:
- Open

Type:
- User-facing

Related stories:
- [US-050](#us-050)

Story:
- As a user, I want to install via Homebrew.

Acceptance criteria:
- Homebrew tap created/updated
- Installation via `brew install sven1103-agent/opencode-config-cli/oc` works
- Formula stays up to date with releases

---

### <a id="us-052"></a>US-052 - Add CI-backed end-to-end coverage for the Go CLI

Priority:
- P1

Status:
- Done

PR:
- [#120](https://github.com/sven1103-agent/opencode-config-cli/pull/120)

Type:
- Release-engineering-facing

Related stories:
- [US-041](#us-041)
- [US-045](#us-045)
- [US-046](#us-046)

Story:
- As a maintainer, I want CI to execute black-box tests against the built `oc` binary so that shipped command wiring, filesystem behavior, and user-facing flows are validated before merge.

Acceptance criteria:
- A dedicated GitHub Actions workflow builds `oc` and runs CLI E2E tests on push and pull request.
- E2E tests execute the compiled binary via `os/exec` instead of calling command internals directly.
- The workflow runs on `ubuntu-latest` and `macos-latest`.
- Initial E2E coverage includes `version`, local source registration, local bundle apply, local archive apply, provenance status, and key offline failure paths.
- Covered and deferred scenarios are tracked in a dedicated E2E coverage reference document.

---

## User Story Backlog Placeholder

User stories will be added later and must:
- have stable IDs using the format `US-###`
- reference at least one feature ID
- reference at least one requirement ID
- include acceptance criteria
- identify whether the story is user-facing, maintainer-facing, or release-engineering-facing

Traceability workflow:
- Every implementation PR must declare exactly one primary user story.
- Primary story selection rules (in order):
  - If the user request names a `US-###`, that is the primary story.
  - Otherwise, choose the one story whose acceptance criteria the PR is intended to complete end-to-end.
  - If the change spans multiple stories and none is clearly primary, split into multiple PRs (recommended) or explicitly pick a primary and list the rest as secondary (see below).
- PR title format (primary story only): `US-###: <short title>`
- PR body must contain a primary line: `Implements: US-###`
- If the PR also advances other stories, include additional lines:
  - `Also advances: US-###, US-###` (optional)
  - Do not mark secondary stories `Done` unless their acceptance criteria are fully satisfied by this PR.
- Updating story entries in this document:
  - Primary story: update the story entry with the PR reference.
    - If all acceptance criteria are satisfied, set `Status: Done`.
    - If not all acceptance criteria are satisfied, set `Status: In progress` or `Status: In review` (choose the closest), and keep it out of `Done`.
  - Secondary stories (if listed): only add the PR reference if the PR is intentionally tracked there; never set `Status: Done` unless fully complete.
- PR reference format:
  - Preferred: `PR: [#<number>](<url>)`
  - Acceptable: `PR: #<number>`
- Practical: if the PR number/URL is not known yet, open the PR with `PR: TBD`, then update the same PR branch before merge to replace `TBD` with the actual PR reference.

### User Stories (Iteration 1)

Note:
- `US-001` through `US-011` were completed before the current one-primary-story-per-PR traceability rule was adopted, so several of those stories share legacy implementation PR references.

### <a id="us-001"></a>US-001 - List available bundled presets with descriptions

Priority:
- P0

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-002](#feat-002)

Related requirements:
- [REQ-F-003](#req-f-003)
- [REQ-F-001a](#req-f-001a)
- [REQ-F-009](#req-f-009)

Story:
- As a developer, I want to list all available official presets with a short purpose description so that I can pick an appropriate starting configuration.

Acceptance criteria:
- `opencode-helper preset list` prints each preset + description.
- Unknown flags exit nonzero.

---

### <a id="us-002"></a>US-002 - Initialize a project from bundled assets (preset + schemas + manifest)

Priority:
- P0

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-001](#feat-001)

Related requirements:
- [REQ-F-004](#req-f-004)
- [REQ-F-005](#req-f-005)
- [REQ-F-010](#req-f-010)
- [REQ-F-011](#req-f-011)
- [REQ-NF-002](#req-nf-002)

Story:
- As a developer, I want to initialize a project using an official preset and bundled schemas so that my repo has a working, validated OpenCode setup without cloning a source repo.

Acceptance criteria:
- `init --project-root <dir> --preset <name>` writes output and installs schemas.
- `--dry-run` makes no changes.
- Overwrite is blocked without `--force`.

---

### <a id="us-003"></a>US-003 - Apply a selected preset to a chosen output path

Priority:
- P0

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-003](#feat-003)

Related requirements:
- [REQ-F-004](#req-f-004)
- [REQ-F-010](#req-f-010)
- [REQ-F-011](#req-f-011)

Story:
- As a developer, I want to apply a specific preset to a chosen output path so that I can adopt official configuration in an existing repo.

Acceptance criteria:
- `preset use <name> --output <path>` writes the preset.
- Overwrite is blocked without `--force`.
- `--dry-run` makes no changes.

---

### <a id="us-004"></a>US-004 - Install bundled schemas into project-local scope

Priority:
- P0

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-004](#feat-004)

Related requirements:
- [REQ-F-005](#req-f-005)
- [REQ-F-010](#req-f-010)
- [REQ-NF-002](#req-nf-002)

Story:
- As a developer, I want to install the official handoff/result schemas into my project so that orchestration artifacts can be validated locally.

Acceptance criteria:
- `schema install --project-root <dir>` installs schemas locally.
- `--dry-run` makes no changes.

---

### <a id="us-005"></a>US-005 - Validate setup health (missing vs drift) with stable exit codes

Priority:
- P0

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-005](#feat-005)

Related requirements:
- [REQ-F-006](#req-f-006)
- [REQ-F-007](#req-f-007)
- [REQ-NF-004](#req-nf-004)
- [REQ-NF-006](#req-nf-006)

Story:
- As a developer (or CI), I want to validate that required setup files exist and match the bundled release expectations so that I can detect missing installs and drift reliably.

Acceptance criteria:
- Healthy setup exits 0.
- Missing setup exits a stable missing exit code and lists missing items.
- Drifted setup exits a stable drift exit code and prints diagnostics.

---

### <a id="us-006"></a>US-006 - Provide diagnostics-oriented guidance for invalid or drifted states

Priority:
- P1

Status:
- Done

PR:
- [#28](https://github.com/sven1103-agent/opencode-agents/pull/28)

Type:
- User-facing

Related features:
- [FEAT-006](#feat-006)

Related requirements:
- [REQ-F-007](#req-f-007)
- [REQ-NF-004](#req-nf-004)

Story:
- As a developer, I want actionable diagnostics for missing files, drift, or invalid setup so that I can quickly remediate mistakes.

Acceptance criteria:
- Missing output suggests a fix.
- Drift output suggests remediation.
- Exit codes remain correct.

---

### <a id="us-007"></a>US-007 - Report CLI version and bundled asset identity/provenance

Priority:
- P0

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-007](#feat-007)

Related requirements:
- [REQ-F-009](#req-f-009)
- [REQ-NF-006](#req-nf-006)
- [REQ-NF-003](#req-nf-003)

Story:
- As a developer, I want the CLI to report its version and the identity of its bundled presets/schemas so that I can trace installed assets back to a specific release.

Acceptance criteria:
- `version` prints the CLI version and identifiers for bundled presets and schemas.
- Missing bundled asset yields a missing exit code and a message.

---

### <a id="us-008"></a>US-008 - Default-safe file handling with explicit overwrite opt-in

Priority:
- P0

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-001](#feat-001)
- [FEAT-003](#feat-003)
- [FEAT-004](#feat-004)

Related requirements:
- [REQ-F-010](#req-f-010)

Story:
- As a developer, I want the CLI to avoid overwriting existing files by default so that I can run it safely in repos that already have configuration.

Acceptance criteria:
- Without `--force`, no overwrite occurs and the command returns an overwrite-blocked exit code.
- With `--force`, overwrite occurs and the command exits 0.

---

### <a id="us-009"></a>US-009 - Support `--project-root` and `--output` for non-default layouts

Priority:
- P0

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-001](#feat-001)
- [FEAT-003](#feat-003)

Related requirements:
- [REQ-F-004](#req-f-004)

Story:
- As a developer, I want to target a specific project root and config output path so that the helper works in monorepos and non-standard directory structures.

Acceptance criteria:
- Writes occur under the provided `--project-root`.
- Output respects `--output`.
- Invalid output errors with a usage/runtime exit code.

---

### <a id="us-010"></a>US-010 - Offline-friendly operation after installation

Priority:
- P1

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- User-facing

Related features:
- [FEAT-001](#feat-001)
- [FEAT-002](#feat-002)
- [FEAT-003](#feat-003)
- [FEAT-004](#feat-004)
- [FEAT-005](#feat-005)
- [FEAT-007](#feat-007)

Related requirements:
- [REQ-NF-002](#req-nf-002)
- [DEC-001](#dec-001)

Story:
- As a developer, I want normal setup operations to work without network access so that I can bootstrap and validate projects in restricted environments.

Acceptance criteria:
- `preset list`, `init`, `preset use`, `schema install`, `validate`, and `version` complete offline.

---

### <a id="us-011"></a>US-011 - Deterministic bundled asset set per release

Priority:
- P1

Status:
- Done

PR:
- [#12](https://github.com/sven1103-agent/opencode-agents/pull/12)

Type:
- Maintainer-facing

Related features:
- [FEAT-007](#feat-007)

Related requirements:
- [REQ-NF-003](#req-nf-003)
- [REQ-F-001](#req-f-001)
- [REQ-F-002](#req-f-002)

Story:
- As a maintainer, I want each CLI release to pin a specific set of presets and schemas so that setup and validation results are reproducible.

Acceptance criteria:
- `version` provenance is stable per release.
- The same release yields identical `preset list` output.

---

### <a id="us-012"></a>US-012 - Install wizard performs global install and updates PATH

Priority:
- P0

Status:
- Done

PR:
- [#15](https://github.com/sven1103-agent/opencode-agents/pull/15)

Type:
- User-facing

Related features:
- [FEAT-009](#feat-009)

Related requirements:
- [REQ-F-012](#req-f-012)
- [REQ-F-014](#req-f-014)

Story:
- As a developer, I want an interactive install wizard (`opencode-helper-install`) that installs `opencode-helper` globally and updates my shell config so that I can run `opencode-helper` from a new terminal without manual PATH steps.

Acceptance criteria:
- Given `opencode-helper-install` installs into an unprivileged `--bin-dir` (e.g. `~/.local/bin`), when the wizard completes, then opening a new terminal session allows running `opencode-helper` successfully without manual `PATH` steps.
- Given the user selects a privileged `--bin-dir` (e.g. `/usr/local/bin`), when the install step writes the binary, then the flow requests `sudo`; and when the flow edits shell config, then it does not use `sudo`.
- Given the wizard has been run once, when it is re-run with the same `--bin-dir`, then the shell config contains exactly one managed marker block and the resulting `PATH` update includes the `--bin-dir` exactly once.
- The shell config file modified contains the exact sentinel lines `# opencode-helper: BEGIN managed PATH` and `# opencode-helper: END managed PATH`.

---

### <a id="us-013"></a>US-013 - Install flow handles shell detection and PATH updates safely

Priority:
- P0

Status:
- Done

PR:
- [#19](https://github.com/sven1103-agent/opencode-agents/pull/19)

Type:
- User-facing

Related features:
- [FEAT-009](#feat-009)

Related requirements:
- [REQ-F-012](#req-f-012)
- [REQ-F-013](#req-f-013)

Story:
- As a developer, I want `opencode-helper-install` to detect my shell and update the correct shell config when possible so that PATH setup is reliable while non-interactive installs still complete.

Acceptance criteria:
- Given the active shell is zsh, when the wizard updates `PATH`, then it writes to `~/.zshrc`.
- Given the active shell is bash on macOS, when the wizard updates `PATH`, then it writes to the first match in order: `~/.bash_profile` (if exists), else `~/.bashrc` (if exists), else creates and writes `~/.bash_profile`.
- Given the active shell is bash on Linux, when the wizard updates `PATH`, then it writes to the first match in order: `~/.bashrc` (if exists), else `~/.bash_profile` (if exists), else creates and writes `~/.bashrc`.
- Given the active shell is fish, when the wizard updates `PATH`, then it writes to `~/.config/fish/config.fish`.
- Given `SHELL` is set to a supported shell, when running `opencode-helper-install --yes --bin-dir <path>`, then the install completes without prompts and updates the shell config derived from `SHELL`.
- Given `SHELL` is unset (or unsupported), when running `opencode-helper-install --yes --bin-dir <path>`, then the install succeeds, skips shell config edits, and prints manual PATH instructions.

---

### <a id="us-014"></a>US-014 - Installer output is themed and readable

Priority:
- P1

Status:
- Done

PR:
- [#22](https://github.com/sven1103-agent/opencode-agents/pull/22)

Type:
- User-facing

Related features:
- [FEAT-009](#feat-009)

Related requirements:
- [REQ-NF-007](#req-nf-007)

Story:
- As a developer, I want `opencode-helper-install` to present a clear, themed terminal wizard with status indicators so that I feel confident and engaged during installation.

Acceptance criteria:
- Status lines use a consistent textual label set (at least `ok:`, `warn:`, `error:`) and may include a single plain Unicode icon prefix.
- No output relies on color alone to convey meaning.

---

### <a id="us-015"></a>US-015 - Installer output has accessible fallbacks for logs and CI

Priority:
- P1

Status:
- Done

PR:
- [#33](https://github.com/sven1103-agent/opencode-agents/pull/33)

Type:
- User-facing

Related features:
- [FEAT-009](#feat-009)

Related requirements:
- [REQ-NF-007](#req-nf-007)

Story:
- As a developer, I want to disable colors and Unicode icons (or have them auto-disabled in non-interactive logs) so that the installer output stays readable in any terminal, CI system, or pasted log.

Acceptance criteria:
- Given `opencode-helper-install --no-color`, then output contains no ANSI color escape sequences.
- Given `opencode-helper-install --ascii`, then output contains no non-ASCII characters.
- Given stdout is not a TTY, then the installer disables colors and Unicode icons by default.

---

### <a id="us-016"></a>US-016 - Publish deterministic release bundles with metadata and checksums

Priority:
- P0

Status:
- Done

PR:
- [#24](https://github.com/sven1103-agent/opencode-agents/pull/24)

Type:
- Release-engineering-facing

Related features:
- [FEAT-010](#feat-010)

Related requirements:
- [REQ-F-015](#req-f-015)
- [REQ-F-016](#req-f-016)
- [REQ-NF-008](#req-nf-008)

Story:
- As a maintainer, I want each tagged helper release to publish a deterministic bundle, release metadata manifest, and SHA-256 checksum manifest so that installer downloads are version-locked and verifiable.

Acceptance criteria:
- Given a tagged helper release, when the release pipeline runs, then it publishes a versioned helper bundle as a GitHub release asset.
- Given that same tagged release, when the pipeline completes, then the published assets include release metadata identifying at least the release tag and source commit.
- Given that same tagged release, when the pipeline completes, then the published assets include a SHA-256 checksum manifest that contains the helper bundle checksum.

---

### <a id="us-017"></a>US-017 - Installer fetches and verifies a release bundle before activation

Priority:
- P0

Status:
- Done

PR:
- [#28](https://github.com/sven1103-agent/opencode-agents/pull/28)

Type:
- User-facing

Related features:
- [FEAT-009](#feat-009)

Related requirements:
- [REQ-F-015](#req-f-015)
- [REQ-F-016](#req-f-016)

Story:
- As a developer, I want `opencode-helper-install` to fetch an official release bundle and verify its SHA-256 checksum before installation so that I can trust the installed helper assets.

Acceptance criteria:
- Given no explicit version selection, when `opencode-helper-install` runs successfully, then it downloads the latest supported helper release bundle from GitHub Releases.
- Given a downloaded bundle and checksum manifest, when the installer verifies the bundle, then installation proceeds only if the SHA-256 checksum matches the published checksum entry.
- Given a missing checksum manifest, missing bundle entry, or checksum mismatch, when verification runs, then the installer exits non-zero and does not activate the downloaded bundle.

---

### <a id="us-018"></a>US-018 - Installer supports explicit release selection

Priority:
- P0

Status:
- Done

PR:
- [#32](https://github.com/sven1103-agent/opencode-agents/pull/32)

Type:
- User-facing

Related features:
- [FEAT-009](#feat-009)

Related requirements:
- [REQ-F-017](#req-f-017)

Story:
- As a developer, I want `opencode-helper-install` to let me choose which helper release to install (defaulting to latest) so that I can pin my environment to a known version when needed.

Acceptance criteria:
- Given no explicit version selection, when `opencode-helper-install` runs in either interactive or non-interactive mode, then it installs the latest supported helper release by default.
- Given `opencode-helper-install --version <tag>`, when `<tag>` exists as a supported GitHub release, then the installer downloads and installs that exact release.
- Given `opencode-helper-install --version <tag>`, when `<tag>` does not exist or is not installable, then the installer exits non-zero with a clear error.

---

### <a id="us-019"></a>US-019 - Installer and helper report active release provenance

Priority:
- P1

Status:
- Done

PR:
- [#33](https://github.com/sven1103-agent/opencode-agents/pull/33)

Type:
- User-facing

Related features:
- [FEAT-009](#feat-009)
- [FEAT-010](#feat-010)

Related requirements:
- [REQ-F-018](#req-f-018)
- [REQ-NF-009](#req-nf-009)

Story:
- As a developer, I want the installer and installed helper to show the active release tag and bundled asset provenance so that I can confirm what version is installed and trace it back to a release.

Acceptance criteria:
- Given a successful install, when the installer completes, then it prints the target release before install and the active installed release after activation.
- Given an installed helper release, when I run the helper version-reporting command, then it prints the active release tag and bundled asset provenance metadata.
- When network access is available for release lookup, the helper may also report whether the active release is the latest supported release.

---

### <a id="us-020"></a>US-020 - Side-by-side installs support local rollback via stable symlink

Priority:
- P0

Status:
- In progress

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-009](#feat-009)

Related requirements:
- [REQ-F-017](#req-f-017)
- [REQ-F-019](#req-f-019)

Story:
- As a developer, I want helper releases to be installed side-by-side and activated through a stable symlink so that I can roll back to an already installed older release quickly.

Acceptance criteria:
- Given an existing installed helper release, when a newer release is installed, then the older release remains on disk by default.
- Given multiple installed releases, when the active release changes, then the stable launcher or `current` symlink points to the selected versioned install root.
- Given an older release is already installed locally, when I run `opencode-helper release use <tag>`, then rollback completes without re-downloading the older release bundle.

Template:

### US-001 - Title

Status:
- Planned | In progress | In review | Done

PR:
- TBD | #123 | https://...

Type:
- User-facing | Maintainer-facing | Release-engineering-facing

Related features:
- FEAT-xxx

Related requirements:
- REQ-F-xxx
- REQ-NF-xxx

Story:
- As a ...
- I want ...
- So that ...

Acceptance criteria:
- Given ...
- When ...
- Then ...

---

### <a id="us-021"></a>US-021 - Bootstrap install.sh enables curl|sh distribution

Priority:
- P0

Status:
- Done

PR:
- [#39](https://github.com/sven1103-agent/opencode-agents/pull/39)

Type:
- User-facing | Release-engineering-facing

Related features:
- [FEAT-009](#feat-009)

Related requirements:
- [REQ-F-020](#req-f-020)

Story:
- As a developer, I want to install `opencode-helper` with a single `curl|sh` command without cloning the repository, so that getting started with the helper CLI is frictionless.

Acceptance criteria:
- Given a macOS or Linux system with `curl` or `wget` available, when I run `curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh`, then `opencode-helper` is installed to `~/.local/bin/opencode-helper` and the helper is on `PATH`.
- Given the bootstrap install script, when I run `curl .../install.sh | sh --version v1.0.0`, then the installed helper is from that exact release tag.
- Given the bootstrap install script, when the downloaded `opencode-helper-install` fails its checksum verification, then the bootstrap exits non-zero without executing the installer.
- Given `OPENCODE_HELPER_VERSION` is set in the environment, when the bootstrap runs, then it uses that tag instead of `latest`.

---

### <a id="us-022"></a>US-022 - Self-documenting CLI commands enable confident tool discovery

Priority:
- P1

Status:
- Done

PR:
- [#45](https://github.com/sven1103-agent/opencode-agents/pull/45)

Type:
- User-facing

Related features:
- [FEAT-001](#feat-001) to [FEAT-007](#feat-007)

Related requirements:
- [REQ-NF-004](#req-nf-004)

Story:
- As a developer, I want the CLI to explain what each command and subcommand does directly in the help output so that I can understand the tool's capabilities without guessing or consulting external documentation.

Acceptance criteria:
- `opencode-helper help` shows each command with a brief description of its purpose.
- `opencode-helper <command> --help` shows a description followed by usage syntax.
- Commands without descriptions are never shipped.

---

### <a id="us-023"></a>US-023 - Self-update preserves existing installation setup and location

Priority:
- P0

Status:
- Planned

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-008](#feat-008)

Related requirements:
- [REQ-F-008](#req-f-008)
- [REQ-F-019](#req-f-019)

Story:
- As a developer, I want to run `opencode-helper self-update` from an existing installation to upgrade to a newer release while preserving my established installation setup and location, so that my environment remains stable and the tool stays where I expect it.

Acceptance criteria:
- Given an existing helper installation at a known `--bin-dir` location, when I run `opencode-helper self-update`, then the updated helper is installed to the same `--bin-dir` location without requiring re-configuration.
- Given an existing helper installation with a stable symlink or launcher, when `self-update` installs a newer release, then the stable launcher/symlink is updated to point to the newly installed release.
- Given an existing helper installation, when `self-update` completes successfully, then previously installed releases remain on disk and are available for rollback via `opencode-helper release use <tag>`.
- Given `opencode-helper self-update --version <tag>`, when `<tag>` exists as an installable GitHub release, then self-update installs that specific version to the existing `--bin-dir` location.
- Given `opencode-helper self-update`, when a newer release is available, then the CLI downloads, verifies, and activates the newer release; when no newer release is available, then the CLI reports that the current release is up-to-date and exits 0.
- Given `opencode-helper self-update`, when network access is unavailable or the release fetch fails, then the CLI exits non-zero with a clear error message and the previous installation remains intact.

---

### <a id="us-024"></a>US-024 - Interactive preset selection without memorization

Priority:
- P1

Status:
- In progress

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-003](#feat-003)

Related requirements:
- [REQ-F-003](#req-f-003)
- [REQ-F-004](#req-f-004)
- [REQ-F-010](#req-f-010)

Story:
- As a developer, I want to select and apply a preset interactively from a numbered list so that I can switch configurations without remembering preset names.

Rationale:
Currently, switching presets requires two steps: `preset list` (to see names) then `preset use <name>` (to apply). This adds friction.

Acceptance criteria:
- Given `opencode-helper preset switch` is invoked in a TTY, when the command runs, then it displays a numbered list of all available presets with their descriptions and indicates the currently selected preset (if a valid manifest exists).
- Given `opencode-helper preset switch` is invoked with piped stdin, when the input is a number (1-5), then the preset at that position is applied.
- Given `opencode-helper preset switch` is invoked with piped stdin, when the input is a preset name or partial match, then the matching preset is applied.
- Given an invalid selection, when in TTY mode, then the command shows an error and re-prompts for input.
- Given an invalid selection, when in non-TTY mode, then the command exits non-zero with a clear error.
- Given `opencode-helper preset switch`, when the selected preset is already applied and no overwrite is desired, then overwrite is blocked without `--force`.
- Given `opencode-helper preset switch --dry-run`, then the command shows what would happen without writing any files.
- Given `opencode-helper preset switch --force`, then the command overwrites the existing preset.
- Given `opencode-helper preset switch --project-root <path>`, when a valid manifest exists, then the current preset is indicated in the interactive list.
- Given `opencode-helper preset switch --project-root <path>`, when no manifest is found, then the command warns the user but proceeds to display the preset list without a current indicator.

---

### <a id="us-025"></a>US-025 - Show current active preset for a project

Priority:
- P1

Status:
- Done

PR:
- [#44](https://github.com/sven1103-agent/opencode-agents/pull/44)

Type:
- User-facing

Related features:
- [FEAT-003](#feat-003)

Related requirements:
- [REQ-F-004](#req-f-004)
- [REQ-NF-004](#req-nf-004)

Story:
- As a developer, I want to query the current preset for a project so that I always know which configuration is active.

Rationale:
The manifest already stores `preset_name`, but there is no command to read it. Users lack context when running validation or other commands.

Acceptance criteria:
- Given a project with a valid manifest, when `opencode-helper preset current` runs, then it outputs the preset name, description, source file, and the helper version that applied it.
- Given a project with a missing manifest, when `opencode-helper preset current` runs, then it exits with a stable missing exit code and prints an actionable message.
- Given a project with a drifted or invalid manifest, when `opencode-helper preset current` runs, then it exits with a stable drift exit code and prints diagnostics.
- Given `opencode-helper preset current --project-root <path>`, then it reads the manifest from the specified project root.

---

### <a id="us-026"></a>US-026 - Register a config source

Priority:
- P0

Status:
- Planned

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-011](#feat-011)

Related requirements:
- [REQ-F-021](#req-f-021)

Story:
- As a developer, I want to register a config source by location so that the helper can discover and install config bundles from it.

Acceptance criteria:
- Given a valid config source location, when I run `opencode-helper source add <location>`, then the source is added to the user-level registry and `opencode-helper source list` includes its resolved source identifier, source type, and location.
- Given a source cannot be validated or its bundle manifest is missing/invalid, when I run `opencode-helper source add <location>`, then the command exits non-zero with a clear validation error and does not add the source.
- Given one or more sources exist, when I run `opencode-helper source remove <source-id>`, then that source is removed from the registry and no longer appears in `source list`.

---

### <a id="us-027"></a>US-027 - List presets across registered config sources

Priority:
- P0

Status:
- Planned

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-012](#feat-012)

Related requirements:
- [REQ-F-022](#req-f-022)
- [REQ-F-023](#req-f-023)

Story:
- As a developer, I want to list presets available across my registered config sources so that I can choose a preset with clear source and version context.

Acceptance criteria:
- Given at least one registered config source with an installable bundle containing `opencode-bundle.manifest.json`, when I run `opencode-helper preset list --sources`, then the output includes each preset `name` and `description` together with its source reference and `bundle_version`.
- Given a source resolves to a bundle with a missing or invalid manifest, when the CLI inspects that source, then it excludes that bundle from preset discovery and reports an actionable validation error.
- Given no config sources are registered, when I run `opencode-helper preset list --sources`, then the command exits with the documented status and prints actionable guidance to add a source.

---

### <a id="us-028"></a>US-028 - Apply a preset from a specific bundle source/version

Priority:
- P0

Status:
- Planned

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-012](#feat-012)

Related requirements:
- [REQ-F-024](#req-f-024)
- [REQ-F-025](#req-f-025)

Story:
- As a developer, I want to apply a preset from a chosen config bundle release so that my project uses a specific versioned configuration.

Acceptance criteria:
- Given a registered source with a unique human-friendly name, when I run `opencode-helper bundle apply <source-name> --preset <name> --project-root <path>`, then the CLI resolves that source by name and applies the preset without requiring the opaque source identifier.
- Given a registered source identifier, when I run `opencode-helper bundle apply <source-id> --preset <name> --project-root <path>`, then the existing identifier-based workflow continues to work.
- Given an interactive TTY and a valid source reference, when I run `opencode-helper bundle apply <source-ref> --project-root <path>` without `--preset`, then the CLI displays a numbered preset selection list for that source and applies the selected preset.
- Given an interactive TTY and a valid source reference, when I provide an invalid preset selection, then the CLI shows a clear error and re-prompts.
- Given non-interactive execution, when I run `opencode-helper bundle apply <source-ref> --project-root <path>` without `--preset`, then the command exits non-zero with an actionable message explaining that `--preset` is required outside interactive mode.
- Given an interactive TTY, when I run `opencode-helper bundle apply <source-ref> --auto --preset <name> --project-root <path>`, then the command does not prompt and applies the preset directly.
- Given two registered sources share the same human-friendly name, when I run `opencode-helper bundle apply <duplicate-name> --preset <name>`, then the command exits non-zero and lists the ambiguous matches instead of choosing one implicitly.
- Given the target project contains existing files that would be overwritten, when I run the command without `--force`, then the command exits non-zero and makes no destructive changes.
- Given the command is run with `--dry-run`, when it completes, then it prints the planned changes and writes no files.

---

### <a id="us-029"></a>US-029 - Materialize referenced prompt files when applying a bundle

Priority:
- P0

Status:
- Planned

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-012](#feat-012)

Related requirements:
- [REQ-F-025](#req-f-025)

Story:
- As a developer, I want referenced prompt files to be installed alongside the applied preset so that `{file:...}` prompt references resolve locally without extra manual steps.

Acceptance criteria:
- Given a preset contains `{file:./prompts/...}` references and the bundle manifest declares the required `prompt_files`, when I apply the preset, then those prompt files are written into the project at the expected relative paths.
- Given a referenced prompt file already exists in the project, when applying without `--force`, then the command exits non-zero and does not overwrite the existing prompt file.

---

### <a id="us-030"></a>US-030 - Show provenance for installed bundles and applied presets

Priority:
- P1

Status:
- Done

PR:
- #66

Type:
- User-facing

Related features:
- [FEAT-014](#feat-014)

Related requirements:
- [REQ-F-026](#req-f-026)

Story:
- As a developer, I want to see which config source and bundle version my project configuration came from so that I can audit and reproduce the setup.

Acceptance criteria:
- Given a project initialized/applied from a bundle with persisted provenance, when I run `opencode-helper bundle status --project-root <path>`, then it prints the source identifier, source location, bundle release/version, and a commit identifier when available.
- Given no bundle provenance exists for the project, when I run the status command, then it exits with a stable missing exit code and prints guidance to apply a bundle.

---

### <a id="us-031"></a>US-031 - Prompt before updating to a newer bundle release

Priority:
- P1

Status:
- Done

PR:
- #67

Type:
- User-facing

Related features:
- [FEAT-013](#feat-013)
- [FEAT-018](#feat-018)

Related requirements:
- [REQ-F-027](#req-f-027)
- [REQ-F-031](#req-f-031)
- [REQ-F-032](#req-f-032)

Story:
- As a developer, I want the helper to detect newer compatible bundle releases and ask me before updating so that I stay in control of configuration changes.

Acceptance criteria:
- Given an installed bundle from an update-capable source and network access, when a newer compatible bundle release exists, then `opencode-helper bundle update <source-id>` prompts for confirmation before downloading/applying updates.
- Given I decline the prompt, when the command exits, then no changes are made and the currently installed bundle remains active.
- Given the selected source does not support update checks, when I run the update command, then the CLI exits non-zero with a clear capability error and leaves existing installed bundles unchanged.
- Given integrity verification for the downloaded remote bundle fails, when I run the update command, then the CLI exits non-zero and does not activate the downloaded bundle.
- Given network access is unavailable (or release lookup fails), when I run the update command, then it exits non-zero with a clear error and leaves existing installed bundles unchanged.

---

### <a id="us-032"></a>US-032 - Publish a deterministic bundle manifest

Priority:
- P1

Status:
- Done

PR:
- [#50](https://github.com/sven1103-agent/opencode-agents/pull/50)

Type:
- Maintainer-facing

Related features:
- [FEAT-012](#feat-012)

Related requirements:
- [REQ-F-022](#req-f-022)

Story:
- As a config bundle maintainer, I want each bundle to include a deterministic manifest that declares available presets and required prompt files so that the helper can discover and install the bundle reliably.

Acceptance criteria:
- Given a config bundle is inspected, when its contents are read, then `opencode-bundle.manifest.json` is present at the bundle root.
- Given the manifest is parsed as JSON, when inspected, then it includes `manifest_version`, `bundle_name`, `bundle_version`, and `presets`.
- Given each preset entry is inspected, then it includes `name`, `description`, `entrypoint`, and `prompt_files`.
- Given a preset declares `entrypoint` or `prompt_files`, when the bundle contents are inspected, then those paths exist in the bundle.

---

### <a id="us-033"></a>US-033 - Register and use a local directory source

Priority:
- P0

Status:
- Planned

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-011](#feat-011)
- [FEAT-015](#feat-015)

Related requirements:
- [REQ-F-028](#req-f-028)

Story:
- As a developer, I want to use a local directory as a config source so that I can iterate on bundles without publishing an archive or remote release.

Acceptance criteria:
- Given a local directory containing a valid `opencode-bundle.manifest.json`, when I run `opencode-helper source add ./path/to/bundle`, then the CLI registers it as a local-directory source.
- Given that registered local-directory source, when I run `opencode-helper preset list --sources`, then presets from that directory appear with the correct source identifier and bundle version.
- Given the local directory does not contain a valid manifest, when I add it as a source, then the CLI exits non-zero with a clear validation error.

---

### <a id="us-034"></a>US-034 - Install and apply a bundle from a local archive

Priority:
- P0

Status:
- Done

PR:
- #62

Type:
- User-facing

Related features:
- [FEAT-012](#feat-012)
- [FEAT-015](#feat-015)
- [FEAT-017](#feat-017)

Related requirements:
- [REQ-F-028](#req-f-028)
- [REQ-F-029](#req-f-029)

Story:
- As a developer, I want to use a local `.tar.gz` archive file as a config source so that I can install and share versioned bundles without requiring a live remote repo.

Acceptance criteria:
- Given a local archive containing a valid bundle, when I run `opencode-helper source add ./bundle.tar.gz`, then the CLI registers it as a local-archive source.
- Given that registered local-archive source, when I run `opencode-helper bundle apply <source-id> --preset <name> --project-root <path>`, then the CLI extracts or resolves the archive to a normalized bundle root, validates the manifest, and applies the selected preset.
- Given the archive contains the bundle at its root or under a single top-level extracted directory, when the CLI resolves it, then apply behavior is identical.

---

### <a id="us-035"></a>US-035 - Install and apply a bundle from a GitHub release source

Priority:
- P1

Status:
- In Review

PR:
- #64

Type:
- User-facing

Related features:
- [FEAT-012](#feat-012)
- [FEAT-015](#feat-015)

Related requirements:
- [REQ-F-028](#req-f-028)
- [REQ-F-029](#req-f-029)
- [REQ-F-031](#req-f-031)

Story:
- As a developer, I want to use a GitHub release bundle as a config source so that I can install officially published bundles from a remote release source.

Acceptance criteria:
- Given a GitHub release source pointing at a valid bundle asset, when I register it, then the CLI records it as a GitHub-release source and can discover presets from its manifest.
- Given that registered GitHub-release source, when I run `opencode-helper bundle apply <source-id> --version <tag> --preset <name> --project-root <path>`, then the CLI downloads or reuses the selected release asset, resolves it to a normalized bundle root, validates the manifest, and applies the preset.

Notes:
- The CLI accepts GitHub release sources as bare `owner/repo`, `github.com/owner/repo`, or a `https://github.com/<owner>/<repo>/releases/tag/<tag>` URL.
- A valid GitHub release source must publish at least one bundle `.tar.gz` asset whose extracted contents contain `opencode-bundle.manifest.json` either at archive root or under a single top-level directory.
- If the release also publishes a matching `-checksums.txt` asset, the CLI verifies the downloaded bundle before apply and fails closed on checksum mismatch.

---

### <a id="us-036"></a>US-036 - Self-update preserves legacy V1 setups

Priority:
- P0

Status:
- Planned

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-016](#feat-016)

Related requirements:
- [REQ-F-030](#req-f-030)
- [REQ-F-008](#req-f-008)

Story:
- As a developer using an existing V1 helper installation, I want to upgrade via `self-update` without breaking my current project setup, so that I can adopt the V2+ CLI safely.

Acceptance criteria:
- Given an existing V1 helper installation, when I run `opencode-helper self-update` and the latest available version is V2+-capable, then the CLI upgrades successfully using the existing self-update workflow.
- Given a project that uses a legacy bundled-preset setup, when I run the upgraded V2+ CLI, then the project remains usable without immediate migration.
- Given the upgraded V2+ CLI detects a legacy bundled-preset setup, when I run a relevant setup or status command, then the CLI reports that the project is using legacy setup state and provides actionable migration guidance.

---

### <a id="us-037"></a>US-037 - Explicitly migrate a legacy bundled-preset project to config-source management

Priority:
- P1

Status:
- Done

PR:
- #68

Type:
- User-facing

Related features:
- [FEAT-016](#feat-016)

Related requirements:
- [REQ-F-030](#req-f-030)
- [REQ-F-024](#req-f-024)

Story:
- As a developer with a legacy bundled-preset project, I want migration into config-source-managed setup to require an explicit action, so that my working project is never silently rewritten.

Acceptance criteria:
- Given a legacy bundled-preset project, when I do not explicitly invoke a migration command or confirm migration, then the CLI does not rewrite the project into a config-source-managed state implicitly.
- Given a legacy bundled-preset project, when I run `opencode-helper migrate legacy-config` and confirm the operation, then the CLI writes the config-source-managed state needed for the selected source and preserves traceable provenance for the migrated setup.
- Given migration cannot complete safely, when the migration command runs, then the CLI exits non-zero with actionable diagnostics and leaves the legacy project setup intact.

---

### <a id="us-038"></a>US-038 - Surface source capabilities consistently

Priority:
- P1

Status:
- Planned

PR:
- TBD

Type:
- User-facing

Related features:
- [FEAT-017](#feat-017)

Related requirements:
- [REQ-F-031](#req-f-031)

Story:
- As a developer, I want the CLI to distinguish which sources support discovery, version selection, and updates, so that commands behave predictably across local and remote source types.

Acceptance criteria:
- Given a registered source, when I run `opencode-helper source list`, then the CLI shows which relevant capabilities that source supports.
- Given I invoke a capability-dependent command against a source that does not support it, when the command runs, then the CLI exits non-zero with a clear capability error.
- Given a local directory or local `.tar.gz` source, when the CLI inspects it, then update-check capability is not reported unless explicitly supported by that source type in a future iteration.

---

### <a id="us-039"></a>US-039 - Verify remote bundle integrity before apply or update

Priority:
- P0

Status:
- Done

PR:
- #65

Type:
- User-facing

Related features:
- [FEAT-018](#feat-018)

Related requirements:
- [REQ-F-032](#req-f-032)

Story:
- As a developer, I want remote config bundles to be verified before they are applied or updated, so that I can trust downloaded bundle contents.

Acceptance criteria:
- Given a remote bundle source with integrity metadata, when I run a remote apply or update flow, then the CLI verifies the downloaded bundle before activation.
- Given integrity metadata is missing, unreadable, or does not match the downloaded bundle, when verification runs, then the CLI exits non-zero and does not activate the downloaded bundle.
- Given verification succeeds, when the remote bundle is applied or updated, then the verified bundle is the one resolved into the normalized local bundle root.

---

## Open Questions for Post-V1

- Should the CLI support exporting presets under names other than `opencode.json`?
- Should schema install support both project-local and user-global scopes in V1, or should V1 stay project-local only?
- Should `self-update` support pinned release channels or version constraints?
- Should the CLI embed asset checksums and manifest metadata directly for stronger provenance reporting?
- Should custom user overlays on top of official presets be supported in V2?
