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

> **Last updated**: 2026-03-28

| Status | Count | Legend |
|--------|-------|--------|
| 🔄 In Progress | 1 | Currently being implemented |
| ⏳ Open | 10 | Not yet started |

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
- `US-040`: ✅ Done (merged PR #83)
- `US-041`: ✅ Done (merged PR #84)

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

**Status**: 🔄 In Progress

Goal:
- Implement the core commands users need day-to-day

Primary stories:
- `US-042` - Implement init command
- `US-043` - Implement preset list and preset use commands
- `US-044` - Implement schema install and validate commands

Implementation:
- `US-042`: 🔄 In Progress (PR #85)
- `US-043`: ⏳ Open
- `US-044`: ⏳ Open

Why next:
- These are the commands users need most
- Get to minimum viable product quickly

Exit criteria:
- `opencode-helper init` copies preset and installs schemas
- `opencode-helper preset list` shows bundled presets with descriptions
- `opencode-helper preset use <name>` applies a preset
- `opencode-helper schema install` installs schemas to .opencode/
- `opencode-helper schema validate` checks installed schemas
- All commands support --dry-run, --force flags

---

### M3 - Extended Commands

**Status**: ⏳ Open

Goal:
- Implement config source management and bundle operations

Primary stories:
- `US-045` - Implement source commands (add, list, remove)
- `US-046` - Implement bundle commands (apply, status, update)
- `US-047` - Implement update command

Implementation:
- `US-045`: ⏳ Open
- `US-046`: ⏳ Open
- `US-047`: ⏳ Open

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

### M4 - Polish & Distribution

**Status**: ⏳ Open

Goal:
- Add polish features and prepare for distribution

Primary stories:
- `US-048` - Add shell completions
- `US-049` - Add --interactive flag for TTY mode
- `US-050` - Set up goreleaser
- `US-051` - Create Homebrew tap

Implementation:
- `US-048`: ⏳ Open
- `US-049`: ⏳ Open
- `US-050`: ⏳ Open
- `US-051`: ⏳ Open

Why last:
- These enhance the experience but are not required for core functionality
- Distribution (goreleaser, Homebrew) should come after commands are stable

Exit criteria:
- Shell completions for bash, zsh, fish
- Interactive mode with --interactive flag for relevant commands
- goreleaser configured for automated releases
- Homebrew tap created/updated

---

## Story Dependency View

Foundation first:
- `US-040` -> `US-041`

MVP commands depend on foundation:
- `US-041` -> `US-042`, `US-043`, `US-044`

Extended commands depend on MVP:
- `US-042`, `US-043`, `US-044` -> `US-045` -> `US-046` -> `US-047`

Polish depends on extended:
- `US-047` -> `US-048`, `US-049` -> `US-050` -> `US-051`

---

## Recommended Primary Story Order

| # | Story | Status | Description |
|---|-------|--------|-------------|
 | 1 | `US-040` | ✅ Done | Set up Go project structure (PR #83) |
 | 2 | `US-041` | ✅ Done | Add GitHub Actions CI (PR #84) |
 | 3 | `US-042` | 🔄 In Progress | Implement init command (PR #85) |
 | 4 | `US-043` | ⏳ Open | Implement preset list/use |
 | 5 | `US-044` | ⏳ Open | Implement schema install/validate |
 | 6 | `US-045` | ⏳ Open | Implement source commands |
 | 7 | `US-046` | ⏳ Open | Implement bundle commands |
 | 8 | `US-047` | ⏳ Open | Implement update command |
 | 9 | `US-048` | ⏳ Open | Add shell completions |
 | 10 | `US-049` | ⏳ Open | Add --interactive flag |
 | 11 | `US-050` | ⏳ Open | Set up goreleaser |
 | 12 | `US-051` | ⏳ Open | Create Homebrew tap |
 | 13 | `US-052` | ⏳ Open | Add example config bundle as release assets |

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
