# OpenCode Helper CLI

## Contents

- [Overview](#overview)
- [Decision Log](#decision-log)
- [PRD](#prd)
- [Scope](#scope)
- [Requirements](#requirements)
  - [Functional Requirements](#functional-requirements)
  - [Non-Functional Requirements](#non-functional-requirements)
- [First Features](#first-features)
- [Traceability Matrix](#traceability-matrix)
- [User Story Backlog Placeholder](#user-story-backlog-placeholder)
- [Open Questions for Post-V1](#open-questions-for-post-v1)

## Overview

This document is the single traceable product source for the OpenCode helper CLI. It captures the initial product direction, requirements, first feature set, and the traceability structure that later user stories must follow.

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

---

## PRD

### <a id="prd-001"></a>PRD-001 - Product Goal

Provide a small helper CLI that bootstraps and maintains a local OpenCode project setup using official bundled config presets and official bundled inter-agent schemas, with safe validation and self-update behavior tied to CLI releases.

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

### Out of Scope

- Direct end-user fetching of raw assets from GitHub releases
- Live remote schema syncing during normal setup
- Automatic migration of arbitrary user-customized configs
- Windows support unless explicitly added later
- Full remote orchestration or hosted service behavior

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

Depends on:
- [REQ-F-015](#req-f-015)
- [REQ-F-017](#req-f-017)

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

---

## Traceability Matrix

| ID | Type | Links To |
|---|---|---|
| [DEC-001](#dec-001) | Decision | [PRD-001](#prd-001), [REQ-F-001](#req-f-001), [REQ-F-002](#req-f-002), [REQ-F-008](#req-f-008), [REQ-NF-002](#req-nf-002), [REQ-NF-003](#req-nf-003) |
| [PRD-001](#prd-001) | PRD | [REQ-F-001](#req-f-001) to [REQ-F-014](#req-f-014), [REQ-NF-001](#req-nf-001) to [REQ-NF-007](#req-nf-007) |
| [REQ-F-001](#req-f-001) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-002](#feat-002), [FEAT-003](#feat-003) |
| [REQ-F-001a](#req-f-001a) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-002](#feat-002) |
| [REQ-F-002](#req-f-002) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-004](#feat-004) |
| [REQ-F-003](#req-f-003) | Functional Requirement | [FEAT-002](#feat-002) |
| [REQ-F-004](#req-f-004) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-003](#feat-003) |
| [REQ-F-005](#req-f-005) | Functional Requirement | [FEAT-001](#feat-001), [FEAT-004](#feat-004) |
| [REQ-F-006](#req-f-006) | Functional Requirement | [FEAT-005](#feat-005) |
| [REQ-F-007](#req-f-007) | Functional Requirement | [FEAT-006](#feat-006) |
| [REQ-F-008](#req-f-008) | Functional Requirement | [FEAT-008](#feat-008) |
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
| [REQ-F-019](#req-f-019) | Functional Requirement | [FEAT-009](#feat-009), [US-020](#us-020) |
| [REQ-NF-001](#req-nf-001) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-008](#feat-008) |
| [REQ-NF-002](#req-nf-002) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-007](#feat-007) |
| [REQ-NF-003](#req-nf-003) | Non-Functional Requirement | [FEAT-008](#feat-008), [FEAT-010](#feat-010) |
| [REQ-NF-004](#req-nf-004) | Non-Functional Requirement | [FEAT-005](#feat-005), [FEAT-006](#feat-006) |
| [REQ-NF-005](#req-nf-005) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-009](#feat-009) |
| [REQ-NF-006](#req-nf-006) | Non-Functional Requirement | [FEAT-005](#feat-005), [FEAT-007](#feat-007) |
| [REQ-NF-007](#req-nf-007) | Non-Functional Requirement | [FEAT-009](#feat-009) |
| [REQ-NF-008](#req-nf-008) | Non-Functional Requirement | [FEAT-010](#feat-010), [US-016](#us-016) |
| [REQ-NF-009](#req-nf-009) | Non-Functional Requirement | [FEAT-010](#feat-010), [US-019](#us-019) |
| [FEAT-009](#feat-009) | Feature | [REQ-F-012](#req-f-012), [REQ-F-013](#req-f-013), [REQ-F-014](#req-f-014), [REQ-F-015](#req-f-015), [REQ-F-016](#req-f-016), [REQ-F-017](#req-f-017), [REQ-F-018](#req-f-018), [REQ-F-019](#req-f-019), [US-012](#us-012), [US-013](#us-013), [US-017](#us-017), [US-018](#us-018), [US-019](#us-019), [US-020](#us-020) |
| [FEAT-010](#feat-010) | Feature | [REQ-F-015](#req-f-015), [REQ-F-016](#req-f-016), [REQ-F-018](#req-f-018), [REQ-NF-008](#req-nf-008), [REQ-NF-009](#req-nf-009), [US-016](#us-016), [US-019](#us-019) |
| [US-012](#us-012) | User Story | [FEAT-009](#feat-009), [REQ-F-012](#req-f-012), [REQ-F-014](#req-f-014) |
| [US-013](#us-013) | User Story | [FEAT-009](#feat-009), [REQ-F-012](#req-f-012), [REQ-F-013](#req-f-013) |
| [US-014](#us-014) | User Story | [FEAT-009](#feat-009), [REQ-NF-007](#req-nf-007) |
| [US-015](#us-015) | User Story | [FEAT-009](#feat-009), [REQ-NF-007](#req-nf-007) |
| [US-016](#us-016) | User Story | [FEAT-010](#feat-010), [REQ-F-015](#req-f-015), [REQ-F-016](#req-f-016), [REQ-NF-008](#req-nf-008) |
| [US-017](#us-017) | User Story | [FEAT-009](#feat-009), [REQ-F-015](#req-f-015), [REQ-F-016](#req-f-016) |
| [US-018](#us-018) | User Story | [FEAT-009](#feat-009), [REQ-F-017](#req-f-017) |
| [US-019](#us-019) | User Story | [FEAT-009](#feat-009), [FEAT-010](#feat-010), [REQ-F-018](#req-f-018), [REQ-NF-009](#req-nf-009) |
| [US-020](#us-020) | User Story | [FEAT-009](#feat-009), [REQ-F-017](#req-f-017), [REQ-F-019](#req-f-019) |

---

## User Story Backlog Placeholder

User stories will be added later and must:
- have stable IDs using the format `US-###`
- reference at least one feature ID
- reference at least one requirement ID
- include acceptance criteria
- identify whether the story is user-facing, maintainer-facing, or release-engineering-facing

Traceability workflow:
- PRs that implement a user story must reference the story ID in the PR title and body.
- PR title format: `US-###: <short title>`
- PR body must contain a line: `Implements: US-###`
- The implementation PR must update the corresponding story entry in this document:
  - set `Status: Done`
  - set `PR: [#<number>](<url>)` (preferred) or `PR: #<number>`
- Practical: if the PR number/URL is not known yet, open the PR with `PR: TBD`, then update the same PR branch before merge to replace `TBD` with the actual PR reference.

### User Stories (Iteration 1)

### US-001 - List available bundled presets with descriptions

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

### US-002 - Initialize a project from bundled assets (preset + schemas + manifest)

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

### US-003 - Apply a selected preset to a chosen output path

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

### US-004 - Install bundled schemas into project-local scope

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

### US-005 - Validate setup health (missing vs drift) with stable exit codes

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

### US-006 - Provide diagnostics-oriented guidance for invalid or drifted states

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

### US-007 - Report CLI version and bundled asset identity/provenance

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

### US-008 - Default-safe file handling with explicit overwrite opt-in

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

### US-009 - Support `--project-root` and `--output` for non-default layouts

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

### US-010 - Offline-friendly operation after installation

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

### US-011 - Deterministic bundled asset set per release

Priority:
- P1

Status:
- Done

PR:
- https://github.com/sven1103-agent/opencode-agents/pull/32

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
- https://github.com/sven1103-agent/opencode-agents/pull/15

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
- https://github.com/sven1103-agent/opencode-agents/pull/19

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
- https://github.com/sven1103-agent/opencode-agents/pull/22

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
- #24

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
- Planned

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
- Given an older release is already installed locally, when I choose to reactivate it, then rollback completes without re-downloading the older release bundle.

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

## Open Questions for Post-V1

- Should the CLI support exporting presets under names other than `opencode.json`?
- Should schema install support both project-local and user-global scopes in V1, or should V1 stay project-local only?
- Should `self-update` support pinned release channels or version constraints?
- Should the CLI embed asset checksums and manifest metadata directly for stronger provenance reporting?
- Should custom user overlays on top of official presets be supported in V2?
