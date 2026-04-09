# OpenCode Helper CLI - Go Migration Milestone Plan

## Purpose

This document sequences the Go migration of the OpenCode helper CLI into implementation-sized milestones.

The goal is to migrate from bash to Go for better JSON handling, testability, and CLI experience, while maintaining feature parity and enabling future improvements like Homebrew distribution.

## Planning Principles

- Feature-by-feature migration: implement one command at a time
- Maintain backward compatibility in behavior (even if not exact exit codes)
- Use modern Go practices: cobra, slog, standard library where possible
- Prioritize user-facing commands first
- Add polish (completions, goreleaser, Homebrew) after core functionality

## Implementation Status

> **Last updated**: 2026-04-01

| Status | Count | Legend |
|--------|-------|--------|
| ✅ Done | 9 | Completed and merged |
| ❌ Out of Scope | 1 | Not applicable to CLI |
| ⏳ Open | 3 | Not yet started |

---

## Recommended Milestones

### M1 - Foundation

**Status**: ✅ Done

Goal:
- Set up Go project structure with cobra scaffolding
- Establish CI pipeline

Primary stories:
- `US-040` - Set up Go project structure
- `US-041` - Add GitHub Actions CI

Implementation:
- `US-040`: ✅ Done (PR #83)
- `US-041`: ✅ Done (PR #84)

Why first:
- Foundation must be solid before implementing commands
- CI catches issues early in the migration

Exit criteria:
- Go module created with proper module path
- Cobra CLI framework with root command
- help and version subcommands work
- Project structure follows feature-based organization
- Dependencies: spf13/cobra, stretchr/testify, gojsonschema
- CI runs go test, go build, golangci-lint on push/PR

---

### M2 - MVP Commands

**Status**: ✅ Done

Goal:
- Implement the core commands users need day-to-day

Primary stories:
- `US-042` - Implement init command
- `US-043` - Implement preset list and preset use commands

Implementation:
- `US-042`: ✅ Done (PR #87)
- `US-043`: ✅ Done (PR #88)
- `US-044`: ❌ Out of Scope (bundle maintainer concern, not CLI)

Why next:
- These are the commands users need most
- Get to minimum viable product quickly

Exit criteria:
- `opencode-helper init` copies preset and installs schemas
- `opencode-helper preset list` shows bundled presets with descriptions
- `opencode-helper preset use <name>` applies a preset
- `opencode-helper schema install` installs schemas to .opencode/
- All commands support --dry-run, --force flags

Note: Schema validation is out of scope - it's a bundle maintainer responsibility.

---

### M3 - Extended Commands

**Status**: ✅ Done

Goal:
- Implement config source management and bundle operations

Primary stories:
- `US-045` - Implement source commands (add, list, remove)
- `US-046` - Implement bundle commands (apply, status, update)
- `US-047` - Implement update command

Implementation:
- `US-045`: ✅ Done (PR #94)
- `US-046`: ✅ Done (PR #95)
- `US-047`: ✅ Done (PR #93)

Why here:
- These depend on MVP commands being stable
- Source registry provides foundation for bundle operations

Exit criteria:
- `opencode-helper source add <location>` registers a config source
- `opencode-helper source list` shows registered sources
- `opencode-helper source remove <id>` unregisters a source
- Registry stored in XDG compliance location (~/.config/opencode-helper/)
- `opencode-helper bundle apply` applies preset from registered source
- `opencode-helper bundle status` shows provenance
- `opencode-helper bundle update` checks for and applies updates
- `opencode-helper update` checks for new CLI version

---

### M4 - Polish, Distribution, and Coverage

**Status**: 🔄 In Progress

Goal:
- Add polish features, release distribution, and shipped-binary coverage

Primary stories:
- `US-048` - Add shell completions
- `US-049` - Add human-first interactive defaults for TTY source-aware flows
- `US-050` - Set up goreleaser
- `US-051` - Create Homebrew tap
- `US-052` - Add CI-backed end-to-end coverage for the Go CLI

Implementation:
- `US-048`: ⏳ Open
- `US-049`: 🔄 In progress
- `US-050`: ✅ Done (PR #90)
- `US-051`: ⏳ Open
- `US-052`: ✅ Done (PR #120)

Why last:
- These enhance the experience but are not required for core functionality
- Distribution (goreleaser, Homebrew) should come after commands are stable
- Black-box coverage should validate the shipped `oc` binary after command behavior is stable

Exit criteria:
- Shell completions for bash, zsh, fish
- Human-first interactive defaults for relevant source-aware TTY commands, with explicit non-interactive opt-out for automation
- goreleaser configured for automated releases
- Homebrew tap created/updated
- Dedicated GitHub Actions workflow builds `oc` and runs black-box E2E tests
- E2E workflow runs on Linux and macOS
- E2E tests cover local directory and local archive source flows
- E2E coverage inventory is documented and kept current

---

## Story Dependency View

Foundation first:
- `US-040` -> `US-041`

MVP commands depend on foundation:
- `US-041` -> `US-042`, `US-043`

Extended commands depend on MVP:
- `US-042`, `US-043` -> `US-045` -> `US-046` -> `US-047`

Polish depends on extended:
- `US-047` -> `US-048`, `US-049` -> `US-050` -> `US-051`

Shipped-binary coverage depends on core command stability:
- `US-045`, `US-046` -> `US-052`

---

## Recommended Primary Story Order

| # | Story | Status | Description |
|---|-------|--------|-------------|
| # | Story | Status | Description |
|---|-------|--------|-------------|
| 1 | `US-040` | ✅ Done | Set up Go project structure (PR #83) |
| 2 | `US-041` | ✅ Done | Add GitHub Actions CI (PR #84) |
| 3 | `US-042` | ✅ Done | Implement init command (PR #87) |
| 4 | `US-043` | ✅ Done | Implement preset list/use (PR #88) |
| 5 | `US-044` | ❌ Out of Scope | Schema validation (bundle maintainer) |
| 6 | `US-045` | ✅ Done | Implement source commands (PR #94) |
| 7 | `US-046` | ✅ Done | Implement bundle commands (PR #95) |
| 8 | `US-047` | ✅ Done | Implement update command (PR #93) |
| 9 | `US-048` | ⏳ Open | Add shell completions |
| 10 | `US-049` | 🔄 In progress | Add human-first interactive defaults for TTY source-aware flows |
| 11 | `US-050` | ✅ Done | Set up goreleaser (PR #90) |
| 12 | `US-051` | ⏳ Open | Create Homebrew tap |
| 13 | `US-052` | ✅ Done | Add CI-backed end-to-end coverage for the Go CLI (PR #120) |

---

## Migration Notes

### What Stays in Bash

The following remain in bash (maintainer-only tools):
- Release build script (`scripts/release/build-opencode-helper-bundle.sh`)
- Any future CI/CD automation scripts

### What Changes

- **Distribution**: From shell installer to Go binary + goreleaser + Homebrew
- **Error handling**: Standard Go exit codes (0, 1, 2) instead of bash-specific codes
- **Config storage**: XDG compliance maintained
- **Commands**: Same commands, possibly simplified structure

### What Improves

- **JSON handling**: Native Go JSON libraries
- **Testability**: Go unit tests with testify
- **CLI experience**: Cobra's built-in help, completions, validation
- **Cross-compilation**: Easy builds for multiple platforms via goreleaser
- **Distribution**: Homebrew tap for easy installation

---

## Future Considerations (Post-Migration)

After the initial Go migration, consider:
- Adding JSON schema validation for config files (currently not in scope)
- Supporting Windows (currently out of scope)
- Adding auth flows for private GitHub repos
- Implementing config file (e.g., ~/.config/opencode-helper/config.yaml) for preferences
