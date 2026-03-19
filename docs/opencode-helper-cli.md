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
- install one official tool
- inspect all available presets with a short description of each preset's purpose
- apply a supported local configuration
- install required schemas
- validate setup health
- update to the latest supported release

### <a id="prd-005"></a>PRD-005 - Success Criteria

V1 is successful when:
- A user can install the CLI without cloning this repository
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

---

## Traceability Matrix

| ID | Type | Links To |
|---|---|---|
| [DEC-001](#dec-001) | Decision | [PRD-001](#prd-001), [REQ-F-001](#req-f-001), [REQ-F-002](#req-f-002), [REQ-F-008](#req-f-008), [REQ-NF-002](#req-nf-002), [REQ-NF-003](#req-nf-003) |
| [PRD-001](#prd-001) | PRD | [REQ-F-001](#req-f-001) to [REQ-F-011](#req-f-011), [REQ-NF-001](#req-nf-001) to [REQ-NF-006](#req-nf-006) |
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
| [REQ-NF-001](#req-nf-001) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-008](#feat-008) |
| [REQ-NF-002](#req-nf-002) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-007](#feat-007) |
| [REQ-NF-003](#req-nf-003) | Non-Functional Requirement | [FEAT-008](#feat-008) |
| [REQ-NF-004](#req-nf-004) | Non-Functional Requirement | [FEAT-005](#feat-005), [FEAT-006](#feat-006) |
| [REQ-NF-005](#req-nf-005) | Non-Functional Requirement | [FEAT-001](#feat-001) to [FEAT-008](#feat-008) |
| [REQ-NF-006](#req-nf-006) | Non-Functional Requirement | [FEAT-005](#feat-005), [FEAT-007](#feat-007) |

---

## User Story Backlog Placeholder

User stories will be added later and must:
- have stable IDs using the format `US-###`
- reference at least one feature ID
- reference at least one requirement ID
- include acceptance criteria
- identify whether the story is user-facing, maintainer-facing, or release-engineering-facing

### User Stories (Iteration 1)

### US-001 - List available bundled presets with descriptions

Priority:
- P0

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

Template:

### US-001 - Title

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
