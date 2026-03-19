# AGENTS.md

## Scope

This repository contains OpenCode agent configuration assets and supporting setup materials.

Primary repo contents:
- `opencode.openai.json` and `opencode.mixed.json` are the main OpenCode config presets
- `.opencode/schemas/` contains the vendored inter-agent handoff and result schemas used by this repo's workflow
- `scripts/opencode-helper` is the iteration-1 helper CLI for bootstrapping `opencode.json` + installing/validating local schemas
- `docs/opencode-helper-cli.md` is the traceable product document for the planned helper CLI
- `README.md` explains the multi-agent configuration approach and how the repo is intended to be used

## Working Expectations

- Keep changes aligned with the repo's planning-first, schema-validated workflow
- Preserve the role of `.opencode/schemas/handoff.schema.json` and `.opencode/schemas/result.schema.json` as the local canonical artifact contracts for this repo
- Prefer additive, traceable documentation changes over informal notes scattered across files
- Treat helper CLI planning in `docs/opencode-helper-cli.md` as the source of truth for requirements, features, and later user stories
- Avoid changing unrelated session artifacts under `.opencode/sessions/`

## Quick Orientation

- If the task is about agent behavior or configuration presets, inspect `README.md`, `opencode.openai.json`, and `opencode.mixed.json`
- If the task is about local schema installation or validation, inspect `scripts/opencode-helper`
- If the task is about future helper CLI behavior, update `docs/opencode-helper-cli.md` first so product decisions remain traceable
