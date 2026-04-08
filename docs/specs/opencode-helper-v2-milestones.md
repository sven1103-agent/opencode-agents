# OpenCode Helper CLI - V2 Milestone Plan

## Purpose

This document sequences the V2+ config-source architecture defined in `docs/opencode-helper-cli.md` into implementation-sized milestones.

The goal is to keep the main PRD focused on product contract and traceability while this document captures rollout order, dependency shaping, and migration strategy.

## Planning Principles

- Lock the bundle contract before implementing multiple source types.
- Prefer local-first source support before remote/update-capable sources.
- Keep V1 setups working throughout the V2 rollout.
- Add update-capable sources only after provenance and normalized resolution are in place.
- Keep user stories small enough that each PR can implement one primary story end-to-end.

## Implementation Status

> **Last updated**: 2026-03-27

| Status | Count | Legend |
|--------|-------|--------|
| ✅ Merged | 14 | Completed and merged to main |
| 🔄 In Progress | 0 | Currently being implemented |
| ⏳ Open | 0 | Not yet started |

---

## V2 Baseline Release - 2026.03.27

**Release Date**: 2026-03-27

This is the initial V2 baseline release with full config-source management support.

### Included Stories

| # | Story | PR |
|---|-------|-----|
| 1 | US-032 - Publish a deterministic bundle manifest | #49, #50 |
| 2 | US-036 - Self-update preserves legacy V1 setups | #52 |
| 3 | US-026 - Discover presets from registered sources | #54 |
| 4 | US-033 - Register a local directory as config source | #55 |
| 5 | US-027 - List available presets from registered sources | #56 |
| 6 | US-028 - Apply preset from registered source to project | #57 |
| 7 | US-029 - Materialize referenced prompt files | (via US-028) |
| 8 | US-034 - Register a local tarball as config source | #62 |
| 9 | US-038 - Surface source capabilities consistently | #63 |
| 10 | US-035 - GitHub release bundle source support | #64 |
| 11 | US-039 - Verify remote bundle integrity | #65 |
| 12 | US-030 - Show provenance for applied bundles | #66 |
| 13 | US-031 - Prompt before bundle updates | #67 |
| 14 | US-037 - Legacy config migration | #68 |

### New Features

- **Config Source Management**: Register local directories, tarballs, or GitHub releases as config sources
- **Bundle Commands**: `bundle apply`, `bundle status`, `bundle update`
- **GitHub Integration**: Register and apply bundles from GitHub releases with version selection
- **Integrity Verification**: SHA256 checksums verify remote bundles before apply
- **Provenance Tracking**: Track which source/version was applied to each project
- **Update Prompts**: Confirmation required before updating bundles
- **Legacy Migration**: Explicit migration path from V1 to V2

### Upgrade from V1

```sh
# Update CLI to V2
opencode-helper self-update

# Or fresh install
curl -fsSL https://github.com/sven1103-agent/opencode-config-cli/releases/latest/download/install.sh | sh

# For existing V1 projects, migration is optional
opencode-helper migrate legacy-config --project-root ./myproject
```

---

## Recommended Milestones

### M0 - Contract and Migration Foundation

**Status**: ✅ Completed

Goal:
- Freeze the minimum V2 contract so later CLI work does not guess at manifest structure or migration behavior.

Primary stories:
- `US-032` - Publish a deterministic bundle manifest
- `US-036` - Self-update preserves legacy V1 setups

Implementation:
- `US-032`: Merged via [PR #49](https://github.com/sven1103/opencode-agents/pull/49), [PR #50](https://github.com/sven1103/opencode-agents/pull/50)
- `US-036`: Merged via [PR #52](https://github.com/sven1103/opencode-agents/pull/52)

Why first:
- `US-032` establishes the manifest contract that all sources rely on.
- `US-036` protects existing users before the CLI starts exposing V2 behavior.

Exit criteria:
- Manifest path, required fields, and validation behavior are stable.
- V1 users can upgrade through `self-update` without breaking existing projects.
- Legacy setups remain usable and receive explicit migration guidance.

### M1 - Source Registry and Local Directory Workflow

**Status**: ✅ Completed

Goal:
- Deliver the smallest useful V2 flow using local directories as config sources.

Primary stories:
- `US-026` - Register a config source
- `US-033` - Register and use a local directory source

Implementation:
- `US-026`: Merged via [PR #54](https://github.com/sven1103/opencode-agents/pull/54)
- `US-033`: Merged via [PR #55](https://github.com/sven1103/opencode-agents/pull/55)

Why next:
- Local directories are the cheapest development loop for maintainers and users.
- Source registration and local-directory support provide the smallest useful V2 development loop.

Exit criteria:
- Users can add/remove/list sources.
- Users can register a valid local bundle directory.

### M2 - Discovery, Apply, and Provenance Persistence

**Status**: 🔄 Partially Complete

Goal:
- Turn source registration into real project setup while ensuring provenance is captured during apply.

Primary stories:
- `US-027` - List presets across registered config sources
- `US-028` - Apply a preset from a specific bundle source/version
- `US-029` - Materialize referenced prompt files when applying a bundle

Implementation:
- `US-027`: ✅ Merged via [PR #56](https://github.com/sven1103/opencode-agents/pull/56)
- `US-028`: ✅ Merged via [PR #57](https://github.com/sven1103/opencode-agents/pull/57)
- `US-029`: ✅ Done - Verified working (implemented in US-028)

Why here:
- Discovery validates that manifest parsing, source identity, and registry state all work together.
- Apply and prompt materialization make the V2 flow usable end-to-end.
- Provenance must be persisted during apply before later status and update flows can rely on it.

Exit criteria:
- Presets from registered local directory sources are discoverable with source/version context.
- Applying a preset from a registered source writes the preset via the manifest `entrypoint`.
- Referenced prompt files are copied into the target project.
- Apply persists provenance needed for later `bundle status` and update behavior.
- `--dry-run` and overwrite-safety rules work for the V2 apply flow.

### M3 - Local `.tar.gz` Bundle Distribution

**Status**: ⏳ Open

Goal:
- Support portable offline bundle sharing without requiring a remote service.

Primary stories:
- `US-034` - Install and apply a bundle from a local archive

Implementation:
- `US-034`: ✅ Done - PR #62

Why here:
- Local archives are a natural next step after local directories.
- This milestone also exercises normalized bundle resolution independent of source transport.

Exit criteria:
- Users can register a local `.tar.gz` bundle as a source.
- The CLI can resolve archives whose bundle sits at archive root or under one top-level extracted directory.
- Apply behavior matches the local directory flow after normalization.

### M4 - Remote Source Foundation

**Status**: ⏳ Open

Goal:
- Add the first official remote source type with explicit capability handling and integrity verification.

Primary stories:
- `US-038` - Surface source capabilities consistently
- `US-035` - Install and apply a bundle from a GitHub release source
- `US-039` - Verify remote bundle integrity before apply or update

Implementation:
- `US-035`: ✅ Merged (PR #64)
- `US-038`: ✅ Merged (PR #63)
- `US-039`: ✅ Merged (PR #65)

Why here:
- Remote release support depends on the manifest contract, normalized resolution, and local apply flow already being stable.
- Capability handling must be explicit before update behavior is exposed.
- Integrity verification should be in place before remote bundles are trusted.
- `US-035` and `US-039` should be treated as one releasability gate for remote source support.

Exit criteria:
- Users can register a GitHub release bundle source and apply a selected release version.
- The CLI exposes which sources support discovery, version selection, and updates.
- Remote bundle apply fails closed when integrity verification fails.

### M5 - Status, Updates, and Explicit Migration

**Status**: ⏳ Open

Goal:
- Add the final user-facing control surfaces around already persisted provenance, update behavior, and explicit migration from legacy setups.

Primary stories:
- `US-030` - Show provenance for installed bundles and applied presets
- `US-031` - Prompt before updating to a newer bundle release
- `US-037` - Explicitly migrate a legacy bundled-preset project to config-source management

Implementation:
- `US-030`: ✅ Merged (PR #66)
- `US-031`: ✅ Merged (PR #67)
- `US-037`: ✅ Merged (PR #68)

Why here:
- `bundle status` can now read provenance written by the apply flow.
- Update prompting should come only after capability handling and remote integrity verification are in place.
- Explicit migration is safer after the V2 config-source model is already working end-to-end.

Exit criteria:
- `bundle status` reports source identifier, source location, bundle version, and commit metadata when available.
- Missing provenance is reported with actionable guidance.
- Update prompts appear only for update-capable sources.
- Legacy users can explicitly migrate into config-source-managed state without risking implicit rewrite.

## Story Dependency View

Core contract and safety:
- `US-032` -> `US-026`, `US-027`, `US-028`, `US-029`, `US-033`, `US-034`, `US-035`, `US-039`
- `US-036` -> all rollout milestones involving shipped V2 CLI behavior

Local-first path:
- `US-026` -> `US-033` -> `US-027` -> `US-028` -> `US-029`

Archive path:
- `US-028` + `US-032` -> `US-034`

Remote/update path:
- `US-038` -> `US-035` -> `US-039` -> `US-031`

Legacy migration path:
- `US-036` -> `US-037`

## Recommended Primary Story Order

| # | Story | Status | PR Reference |
|---|-------|--------|---------------|
| 1 | `US-032` | ✅ Merged | PR #49, #50 |
| 2 | `US-036` | ✅ Merged | PR #52 |
| 3 | `US-026` | ✅ Merged | PR #54 |
| 4 | `US-033` | ✅ Merged | PR #55 |
| 5 | `US-027` | ✅ Merged | PR #56 |
| 6 | `US-028` | ✅ Merged | PR #57 |
| 7 | `US-029` | ✅ Done | (implemented in US-028) |
| 8 | `US-034` | ✅ Merged | PR #62 |
| 9 | `US-038` | ✅ Merged | PR #63 |
| 10 | `US-035` | ✅ Merged | PR #64 |
| 11 | `US-039` | ✅ Merged | PR #65 |
| 12 | `US-030` | ✅ Merged | PR #66 |
| 13 | `US-031` | ✅ Merged | PR #67 |
| 14 | `US-037` | ✅ Merged | PR #68 |

## Notes

- This milestone plan is additive and should be updated as stories are refined or split.
- If remote source auth, generic HTTPS bundles, or additional archive formats are introduced later, they should be added as new milestones rather than squeezed into the first V2 baseline.
