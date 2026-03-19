# OpenCode AI Agents — Planning-First Multi-Tier Configuration

Drop-in OpenCode agent configs that route work through schema-validated JSON handoff artifacts before planning, execution, and review. The goal is fewer ambiguous changes, less rework, and tighter safety boundaries with a small, explicit contract between agents.

This repository contains a planning-first, multi-tier agent configuration for [OpenCode AI](https://opencode.ai), plus an OpenAI-model variant in `opencode.openai.json`. It defines specialized agents across four functional tiers — routing/orchestration, planning, execution, and validation — designed to minimize cost while preserving quality at every decision point.

Repository configuration files:

- `opencode.mixed.json` — mixed model stack (routing/planning/review vs. code execution)
- `opencode.openai.json` — OpenAI-based variant (including docs routing examples in this README)

Product planning:

- `docs/opencode-helper-cli.md` — traceable PRD, requirements, features, and later user-story source for the helper CLI

---

## Design Philosophy

### Planning-First Execution

The core insight behind this configuration is that **unplanned implementation is expensive to undo**. Before any file is touched, the system asks: is this task concrete and scoped enough to implement directly? If not, a dedicated planning agent runs first.

This separation:

- **Prevents scope creep** — the plan defines the boundary before a single line of code is written
- **Enables cheap execution** — a validated plan can be handed to a lower-cost model with confidence
- **Allows early escalation** — if the plan reveals unexpected complexity, routing can adjust before implementation begins
- **Reduces rework** — reviewers receive a structured summary of what changed and why, not just a diff

### Model Tier Strategy

Four model tiers are used, selected on the principle: **use the cheapest model that can do the job correctly**.

| Tier | Model | Used For |
|------|-------|----------|
| **Standard** | `claude-sonnet-4-6` | Planning, routing decisions, review, senior implementation |
| **Fast** | `claude-haiku-4-5` | Cheap routing (docs), narrow doc edits |
| **Mini** | `gpt-5.1-codex-mini` | Trivial and localized code edits only |
| **Codex** | `gpt-5.3-codex` | Primary implementation execution |

**Rationale by role:**

- **Routing agents** (`coding-boss`, `docs`) need just enough judgment to classify a task — haiku is sufficient for docs; sonnet is used for coding-boss because misrouting a non-trivial coding task is costlier than misrouting a docs task.
- **Planning agents** (`planner`, `docs-planner`) require genuine reasoning over the codebase — sonnet is warranted.
- **Execution agents** use code-optimized models. `gpt-5.1-codex-mini` handles trivial edits cheaply. `gpt-5.3-codex` handles the primary implementation workload.
- **Validation agents** (`code-reviewer`, `docs-reviewer`, `agent-architect`) need judgment and nuance — sonnet throughout.

---

## Agent Reference

### Tier 1 — Routing Agents (Entry Points)

---

#### `coding-boss`

> **Routes coding work by phase: planning, trivial implementation, normal implementation, and review**

- **Model:** `claude-sonnet-4-6`
- **Mode:** Primary (default entry point)
- **Color:** `#BEEE62`

`coding-boss` is the default agent for all coding tasks. It classifies the incoming request and delegates; it never writes code, but it does write/update session artifacts under `.opencode/sessions/<session_id>/handoffs/` and `.opencode/sessions/<session_id>/results/`.

**Routing decision tree:**

```
Incoming coding task
│
├─ NOT implementation-ready?
│    └─→ delegate to planner
│
├─ Implementation-ready AND trivial/localized?
│    └─→ delegate to implementer-small
│
└─ Implementation-ready but non-trivial?
     ├─→ delegate to implementer
     └─→ after completion, delegate to code-reviewer
```

**Implementation-ready** means ALL of the following are true:
- The change is concrete and narrowly scoped
- The affected subsystem or files are obvious
- No architectural decision is required
- No public API, schema, migration, security boundary, or cross-service contract is affected

**Trivial/localized** (eligible for `implementer-small`) means:
- Limited to one or two closely related files
- Small mechanical change
- No architecture or debugging involved

**Permission boundaries:**

```json
"permission": {
  "task": {
    "*": "deny",
    "planner": "allow",
    "implementer-small": "allow",
    "implementer": "allow",
    "code-reviewer": "allow"
  },
  "write": "allow",
  "edit": "allow"
}
```

`coding-boss` operates on an explicit task allow-list. It cannot delegate to any agent not listed. Its file access is intended for maintaining `.opencode/sessions/<session_id>/{handoffs,results}/*` artifacts only.

---

#### `docs`

> **Routes documentation tasks with planner-first policy**

- **Model:** `claude-haiku-4-5` (cheapest tier — routing logic is simple)
- **Mode:** Primary
- **Color:** `#D74E09`

`docs` is the entry point for all documentation work. Like `coding-boss`, it never writes documentation itself — it classifies and delegates, and maintains `.opencode/sessions/<session_id>/{handoffs,results}/*` artifacts.

**Planner-first policy:** Unless the request is a trivial wording fix, typo correction, or formatting cleanup with an obvious target file, `docs` routes to `docs-planner` first. File reading alone is not sufficient justification to skip planning — if understanding requires reading multiple files or synthesizing behavior, `docs-planner` must run first.

**Routing decision tree:**

```
Incoming documentation task
│
├─ AGENTS.md or multi-agent workflow design?
│    └─→ agent-architect
│
├─ Narrow, explicit, no synthesis needed?
│    └─→ docs-writer-fast
│
└─ Everything else (architecture, onboarding, migration,
   feature docs, multi-file synthesis)
      └─→ docs-planner → docs-writer-fast
           └─→ (optional) docs-reviewer
```

**Permission boundaries:**

```json
"permission": {
  "task": {
    "*": "deny",
    "docs-planner": "allow",
    "docs-writer-fast": "allow",
    "docs-reviewer": "allow",
    "agent-architect": "allow"
  },
  "write": "allow",
  "edit": "allow"
}
```

---

### Tier 2 — Planning Agents

---

#### `planner`

> **Produces structured implementation plans**

- **Model:** `claude-sonnet-4-6`
- **Mode:** Subagent

`planner` analyzes the task and repository context and produces a structured execution plan. It **never edits files**.

The plan includes: objective, scope, assumptions, constraints, likely affected files, step-by-step instructions, test strategy, acceptance criteria, risks, rollback notes, and escalation conditions.

**Permission boundaries:**

```json
"permission": {
  "write": "deny",
  "edit": "deny",
  "bash": "ask"
}
```

`bash` is set to `ask` (not `allow`) because the planner may need to inspect the repository to understand structure, but should not run arbitrary commands without confirmation.

---

#### `docs-planner`

> **Plans complex documentation**

- **Model:** `claude-sonnet-4-6`
- **Mode:** Subagent

`docs-planner` researches the codebase, understands its structure and behavior, and produces a compact docs execution plan (audience, goal, exact files/sections to update, concrete changes, examples to include, and how to verify accuracy).

**When triggered vs. `docs-writer-fast` directly:**
- Triggered when the task involves reading files to understand behavior, workflows, architecture, onboarding, migrations, or feature usage
- `docs-writer-fast` is used directly only for narrowly scoped, self-evident changes

**Permission boundaries:**

```json
"permission": {
  "write": "deny",
  "edit": "deny"
}
```

---

### Tier 3 — Execution Agents

---

#### `implementer-small`

> **Cheap execution agent for trivial tasks**

- **Model:** `gpt-5.1-codex-mini`
- **Mode:** Subagent

`implementer-small` is the cost-optimized path for small, localized edits. It uses the cheapest capable code model to minimize cost on work that doesn't require heavy reasoning.

**Rules:**
- Only performs edits limited to one or two closely related files
- Prefers tiny diffs
- Does not modify APIs, schemas, or security boundaries

**Self-escalation:** If scope unexpectedly expands during implementation, `implementer-small` stops and escalates to `@implementer` rather than proceeding beyond its mandate.

**Output:** Implementation summary result JSON (`result_type: "implementation_summary"`).

**Permission boundaries:**

```json
"permission": {
  "write": "allow",
  "edit": "allow",
  "bash": "allow"
}
```

---

#### `implementer`

> **Primary implementation agent**

- **Model:** `gpt-5.3-codex`
- **Mode:** Subagent

`implementer` is the main execution engine for non-trivial coding tasks. It requires an implementation plan in its assigned handoff artifact and follows it strictly.

Before editing, it restates: the objective, files it expects to modify, and its execution plan. If the plan is invalid or contradictory, it escalates back to `@planner` rather than improvising.

**Output:** Implementation summary result JSON (`result_type: "implementation_summary"`).

**Permission boundaries:**

```json
"permission": {
  "write": "allow",
  "edit": "allow",
  "bash": "allow"
}
```

> **Note:** `implementer` uses `gpt-5.3-codex`, not a Claude model. This is intentional — Codex models are optimized for code generation and execution.

---

#### `docs-writer-fast`

> **Documentation writer for explicit, plan-backed edits**

- **Model:** `claude-haiku-4-5`
- **Mode:** Subagent

`docs-writer-fast` executes narrowly scoped documentation updates (typos, small rewrites, formatting, concrete examples) and also carries out larger edits when guided by a docs plan. It keeps diffs tight and follows the plan strictly; if the plan is missing, contradictory, or the scope expands, it escalates to `@docs-planner`.

**Permission boundaries:**

```json
"permission": {
  "write": "allow",
  "edit": "allow"
}
```

---

### Tier 4 — Validation & Design Agents

---

#### `code-reviewer`

> **Reviews implementation for quality and safety**

- **Model:** `claude-sonnet-4-6`
- **Mode:** Subagent

`code-reviewer` receives an implementation summary from an implementer and performs a structured review against four axes: correctness, security, maintainability, and test adequacy.

**Output:** A single JSON object persisted as a result artifact, matching `.opencode/schemas/result.schema.json` (use `result_type: "review_result"` and `status: "approve" | "needs_changes"`).

**Permission boundaries:**

```json
"permission": {
  "write": "deny",
  "edit": "deny",
  "bash": "ask"
}
```

`code-reviewer` never modifies files. It judges; it does not act.

---

#### `docs-reviewer`

> **Reviews documentation quality**

- **Model:** `claude-sonnet-4-6`
- **Mode:** Subagent

`docs-reviewer` reviews completed documentation for accuracy, clarity, and structure.

**Output:** A single JSON object persisted as a result artifact, matching `.opencode/schemas/result.schema.json` (use `result_type: "review_result"` and `status: "approve" | "needs_changes"`).

**Permission boundaries:**

```json
"permission": {
  "write": "deny",
  "edit": "deny"
}
```

---

#### `agent-architect`

> **Designs AGENTS.md and multi-agent workflow documentation**

- **Model:** `claude-sonnet-4-6`
- **Mode:** Subagent

`agent-architect` is a specialized design agent for documentation about agent systems themselves — for example, an `AGENTS.md` you maintain in your own repo. It defines agent roles, delegation patterns, and interaction guidelines.

`agent-architect` is **write-denied** (design only). It produces design output; it does not write files directly.

**When triggered:** Only via the `docs` routing agent, when the task involves AGENTS.md or multi-agent workflow design.

**Permission boundaries:**

```json
"permission": {
  "write": "deny",
  "edit": "deny"
}
```

---

## Permission Model

Agent permissions are controlled along three axes in your configuration file (for example: `opencode.mixed.json`, `opencode.openai.json`, or a local `opencode.json`):

### 1. `task` — Subagent delegation allow-list

Primary routing agents (`coding-boss`, `docs`) use explicit task allow-lists:

```json
"task": {
  "*": "deny",
  "planner": "allow",
  "implementer": "allow"
}
```

This means `coding-boss` **can only delegate to the four agents listed in its allow-list**. It cannot spontaneously invoke any other agent. This is a hard architectural boundary — not a convention.

### 2. `write` / `edit` — File system access

| Agent | write | edit |
|-------|-------|------|
| `coding-boss` | allow* | allow* |
| `planner` | deny | deny |
| `implementer-small` | allow | allow |
| `implementer` | allow | allow |
| `code-reviewer` | deny | deny |
| `docs` | allow* | allow* |
| `docs-planner` | deny | deny |
| `docs-writer-fast` | allow | allow |
| `docs-reviewer` | deny | deny |
| `agent-architect` | deny | deny |

\* Intended for maintaining `.opencode/sessions/<session_id>/{handoffs,results}/*` only.

**Why routing, planning, and review agents are (mostly) write-denied:**

Planning and review agents exist to make decisions, not to implement. Write-denial prevents "fixing while planning" and blurring the line between review and implementation. Routing/orchestration agents are the exception: they maintain session artifacts under `.opencode/sessions/<session_id>/{handoffs,results}/*` but should not touch the rest of the repo.

### 3. `bash` — Shell access

| Agent | bash |
|-------|------|
| `planner` | ask (requires confirmation) |
| `implementer-small` | allow |
| `implementer` | allow |
| `code-reviewer` | ask |

Agents with `bash: ask` (for example: `planner`, `code-reviewer`) may legitimately need to inspect repository structure (e.g., `ls`, `find`, `grep`), but should not run build systems, tests, or mutation commands without user confirmation. Execution agents use `bash: allow` because running tests and build tools is part of their normal workflow.

---

## Workflow Walkthroughs

### Coding Workflow

```
User
  │
  ▼
coding-boss  [claude-sonnet-4-6]
  │  Routes task; maintains `.opencode/sessions/<session_id>/{handoffs,results}/*`
  │
  ├─ NOT implementation-ready
  │     │
  │     ▼
  │   planner  [claude-sonnet-4-6]
  │     │  Returns: implementation plan content
  │     │  No files touched
  │     │
  │     ▼
  │   implementer  [gpt-5.3-codex]
  │     │  Follows plan strictly
  │     │  Returns: result JSON artifact (`implementation_summary`)
  │     │
  │     ▼
  │   code-reviewer  [claude-sonnet-4-6]
  │     │  Returns: result JSON artifact (`review_result`)
  │     ▼
  │   coding-boss  [claude-sonnet-4-6]
  │       Records result artifacts; decides next step / returns final result
  │
  ├─ Implementation-ready, trivial
  │     │
  │     ▼
  │   implementer-small  [gpt-5.1-codex-mini]
  │       Returns: result JSON artifact (`implementation_summary`)
  │       (self-escalates to implementer if scope expands)
  │
  │     ▼
  │   coding-boss  [claude-sonnet-4-6]
  │       Records result artifact; decides next step / returns final result
  │
  └─ Implementation-ready, non-trivial
        │
        ▼
      implementer  [gpt-5.3-codex]
        │  Returns: result JSON artifact (`implementation_summary`)
        │
        ▼
      code-reviewer  [claude-sonnet-4-6]
        │  Returns: result JSON artifact (`review_result`)
        ▼
      coding-boss  [claude-sonnet-4-6]
          Records result artifacts; decides next step / returns final result
```

### Documentation Workflow

```
User
  │
  ▼
docs  [claude-haiku-4-5]
  │  Routes task; maintains `.opencode/sessions/<session_id>/{handoffs,results}/*`
  │
  ├─ AGENTS.md / multi-agent workflow design
  │     │
  │     ▼
  │   agent-architect  [claude-sonnet-4-6]
  │       Design output (write-denied)
  │
  ├─ Narrow, explicit, no synthesis needed
  │     │
  │     ▼
  │   docs-writer-fast  [claude-haiku-4-5]
  │       Returns: result JSON artifact (`docs_result`)
  │
  │     ▼
  │   docs  [claude-haiku-4-5]
  │       Records result artifact; decides next step / returns final result
  │
  └─ Complex, multi-file, or requires codebase synthesis
        │
        ▼
      docs-planner  [claude-sonnet-4-6]
        │  Returns: docs plan content
        │  No files touched
        │
        ▼
      docs-writer-fast  [claude-haiku-4-5]
        │  Returns: result JSON artifact (`docs_result`)
        │
        ▼
      docs-reviewer  [claude-sonnet-4-6]   (optional)
        │  Returns: result JSON artifact (`review_result`)
        ▼
      docs  [claude-haiku-4-5]
          Records result artifacts; decides next step / returns final result
```

### Concrete Routing Examples

The examples below are concrete, paste-ready delegation flows as encoded in `opencode.openai.json`. If you use `opencode.mixed.json` (or a local `opencode.json`), agent names/models may differ; treat these as routing-shape examples, not a guarantee of exact model IDs.

**Documentation tasks:**

```
User prompt: "Fix a typo in README.md: change 'workfow' to 'workflow'."
docs  [openai/gpt-5.2]
  → docs-writer-fast  [openai/gpt-5.2]
```

```
User prompt: "Add a new 'Getting Started' section describing how to run tests and lint."
docs  [openai/gpt-5.2]
  → docs-planner  [openai/gpt-5.4]
      Returns: docs plan content
      Next agent: @docs-writer-fast
  → docs-writer-fast  [openai/gpt-5.2]
  → docs-reviewer  [openai/gpt-5.4]    (optional)
```

```
User prompt: "Create/refresh an AGENTS.md to document our multi-agent workflow."
docs  [openai/gpt-5.2]
  → agent-architect  [openai/gpt-5.4]
```

**Coding tasks:**

```
User prompt: "Rename a local variable in src/foo.ts and update its references."
coding-boss  [openai/gpt-5.2]
  → implementer-small  [openai/gpt-5.1-codex-mini]
```

```
User prompt: "Add rate limiting to the API (needs design decisions + tests)."
coding-boss  [openai/gpt-5.2]
  → planner  [openai/gpt-5.4]
      Returns: implementation plan content
      Next agent: @implementer
  → implementer  [openai/gpt-5.3-codex]
  → code-reviewer  [openai/gpt-5.4]
```

### Handoff Artifacts (JSON Contract)

This repo's cross-agent contract is an on-disk, schema-validated JSON handoff artifact written by the routing/orchestration agent (for example: `coding-boss`, `docs`) and stored under a session folder: `.opencode/sessions/<session_id>/handoffs/`.

The `=== HANDOVER ... ===` blocks shown in prompts and examples are a human-readable convention; the machine-validated structure is the JSON artifact itself.

#### Router Output vs. Persisted Artifacts (Dual Contract)

- Persisted artifacts under `.opencode/sessions/<session_id>/{handoffs,results}/*` are the canonical machine-readable trace.
- Router chat output is the canonical user-facing experience: brief phase-boundary status updates, then a final human-readable outcome summary derived from the persisted result artifact.
- Routers should include an artifact reference for traceability, but must never reply with only an artifact path.

Concrete expected final message shape:

```
Done: Updated router prompts so users see status updates + a final outcome summary.
Trace: session <session_id>, final result .opencode/sessions/<session_id>/results/<seq>-<result_type>-<agent>.json
```

Canonical schemas:

- `.opencode/schemas/handoff.schema.json` (what a handoff artifact must contain)
- `.opencode/schemas/result.schema.json` (required structured result objects execution/review agents must return)

If you use the multi-agent/agent-flow configs, keep `.opencode/schemas/` in place (including both files above). The routing/orchestration prompts reference these schemas as the canonical contract; without them, orchestration cannot validate or persist handoffs/results reliably.

#### Session Storage Layout

- Handoffs: `.opencode/sessions/<session_id>/handoffs/`
- Results: `.opencode/sessions/<session_id>/results/`

Sessions prevent filename clashes across concurrent runs and give you a single folder to archive, share, or delete.

#### Human-Readable IDs

- `session_id` is the primary human-facing identifier. Prefer a ULID (time-sortable, globally unique, offline-generatable).
- Within a session, use a monotonic zero-padded sequence as the primary per-artifact handle: `0001`, `0002`, ...

When referencing artifacts in chat/issues, prefer: `session <session_id>, handoff <seq>`.

#### Handoff and Result Filename Patterns

- Handoff: `.opencode/sessions/<session_id>/handoffs/<seq>-<from>-to-<to>[-<slug>].json`
- Result: `.opencode/sessions/<session_id>/results/<seq>-<result_type>-<agent>[-<slug>].json`

The optional `-<slug>` is for scanability only; it must not be required for uniqueness.

#### Metadata and Traceability

ISO timestamps stay in JSON metadata fields (for example: `created_at`) for auditability/traceability, but are no longer the main filename/ID humans read daily.

At minimum, handoff artifacts include required top-level fields like `version`, `kind`, `handoff_id`, `parent_handoff_id`, `from_agent`, `to_agent`, `created_at`, `status`, and a `payload`. For `implementation_plan`, that payload is the canonical recovery record and captures not just `goal`, `why`, `files_to_modify`, `changes`, and acceptance/abort criteria, but also architecture decisions, target symbols or sections, execution order, contract impacts, edge cases, out-of-scope boundaries, and verification expectations.

`source_handoff_id` links the result JSON back to the handoff that triggered that agent execution.

Implementation results are likewise structured for replay and review. `implementation_summary` now records completed plan steps, deviations from plan, changed symbols, verification results, and remaining follow-ups so a fresh-context reviewer or recovery run does not depend on transient chat context.

#### Referencing Examples

- `.opencode/sessions/01JNZ2D0R3R7J5M2G7V8Q9K2C1/handoffs/0001-router-to-planner.json`
- `.opencode/sessions/01JNZ2D0R3R7J5M2G7V8Q9K2C1/results/0001-docs_result-docs-writer-fast.json`
- Reference: session `01JNZ2D0R3R7J5M2G7V8Q9K2C1`, handoff `0001`
- Optional slug: `.opencode/sessions/01JNZ2D0R3R7J5M2G7V8Q9K2C1/handoffs/0002-docs-router-to-docs-planner-id-convention.json`

Minimal paired example (handoff + result):

```json
{
  "version": 1,
  "kind": "docs_plan",
  "handoff_id": "01JNZ2D0R3R7J5M2G7V8Q9K2C1-0001",
  "parent_handoff_id": null,
  "from_agent": "docs",
  "to_agent": "docs-writer-fast",
  "created_at": "2026-03-17T00:00:00Z",
  "status": "pending",
  "payload": {
    "goal": "Update README to require persisted result artifacts",
    "why": "Make end-to-end workflow produce structured result JSON",
    "files_to_modify": ["README.md"],
    "files_to_inspect_only": [".opencode/schemas/result.schema.json"],
    "do_not_modify": [".opencode/schemas/result.schema.json"],
    "inputs_already_verified": ["result.schema.json defines required fields"],
    "changes": ["README.md: require result JSON artifacts"],
    "tests": ["none"],
    "done_when": ["README requires persisted result artifacts"],
    "abort_if": ["target sections cannot be identified"],
    "examples": []
  }
}
```

```json
{
  "version": 1,
  "result_type": "docs_result",
  "agent": "docs-writer-fast",
  "source_handoff_id": "01JNZ2D0R3R7J5M2G7V8Q9K2C1-0001",
  "created_at": "2026-03-17T00:05:00Z",
  "status": "done",
  "summary": "Updated workflow docs to require result artifacts",
  "files_changed": ["README.md"],
  "tests_run": ["none"]
}
```

---

## Getting Started

### Using the Configuration

Choose a configuration file and point OpenCode AI at it (many setups keep the active config at `opencode.json` locally):

This repository includes `opencode.mixed.json` and `opencode.openai.json` (OpenAI-model variant).

```bash
# Validate the configuration
cat opencode.mixed.json | jq . > /dev/null && echo "Valid JSON"   # or: opencode.openai.json
```

The `default_agent` field is set to `docs`:

```json
{
  "default_agent": "docs"
}
```

This means any task submitted without an explicit agent selection routes through `docs` automatically (for code changes, call `coding-boss` explicitly).

**Entry points:**
- `coding-boss` — for all code changes, bug fixes, refactors, and implementations
- `docs` — for all documentation tasks

### Iteration-1 helper CLI

This repository also includes an iteration-1 helper script at `scripts/opencode-helper` for bootstrapping a local `opencode.json` and local schema install.

```bash
# Show helper commands
sh scripts/opencode-helper help

# List bundled presets
sh scripts/opencode-helper preset list

# Initialize a project with the OpenAI preset
sh scripts/opencode-helper init --preset openai --project-root "$PWD"

# Validate helper-managed setup state
sh scripts/opencode-helper validate --project-root "$PWD"

# Equivalent output path spellings resolve to the same target
sh scripts/opencode-helper validate --project-root "$PWD" --output ./opencode.json

# Show helper version and bundled asset provenance
sh scripts/opencode-helper version
```

Generated project-local layout:

```text
<project-root>/
  opencode.json
  .opencode/
    opencode-helper-manifest.tsv
    install-manifest.tsv
    schemas/
      handoff.schema.json
      result.schema.json
```

### Helper release bundle assets (US-016)

Tagged helper releases (`v*`) now publish three deterministic GitHub Release assets:

- `opencode-helper-<tag>.tar.gz`
- `opencode-helper-<tag>-manifest.json`
- `opencode-helper-<tag>-checksums.txt`

Bundle archive layout is rooted at `opencode-helper-<tag>/` and preserves helper runtime-relative paths:

```text
opencode-helper-<tag>/
  scripts/opencode-helper
  scripts/opencode-helper-install
  opencode.openai.json
  opencode.mixed.json
  .opencode/schemas/handoff.schema.json
  .opencode/schemas/result.schema.json
  release-manifest.json
```

The external `*-manifest.json` matches the in-bundle `release-manifest.json` bytes. `*-checksums.txt` contains SHA-256 lines for the tarball and manifest in `<sha256>  <filename>` format.

`scripts/opencode-helper-install` now installs from the latest published GitHub release by default. It downloads `opencode-helper-<tag>.tar.gz` together with `opencode-helper-<tag>-checksums.txt`, verifies the tarball SHA-256 before extraction, and aborts without activating the bundle if verification fails.

### Calling Agents

```bash
# Code task: coding-boss handles routing automatically
# It will decide whether to plan first or go straight to implementation

# Documentation task: docs applies planner-first policy
# It will decide whether to plan, write directly, or involve agent-architect
```

You do not need to call `planner`, `implementer`, or `code-reviewer` directly. The routing agents manage the pipeline.

---

## Customization

### Adding a New Agent

1. Add the agent definition to your configuration file (for example: `opencode.mixed.json`, `opencode.openai.json`, or your local `opencode.json`) under `"agent"`:

```json
"my-new-agent": {
  "description": "What this agent does",
  "mode": "subagent",
  "model": "opencode/claude-sonnet-4-6",
  "prompt": "Your system prompt here.",
  "permission": {
    "write": "deny",
    "edit": "deny"
  }
}
```

2. **Update the task allow-list** of any routing agent that should be able to delegate to it. If you add a new implementation agent but don't add it to `coding-boss`'s `task` allow-list, `coding-boss` will never route work to it:

```json
"coding-boss": {
  "permission": {
    "task": {
      "*": "deny",
      "planner": "allow",
      "implementer-small": "allow",
      "implementer": "allow",
      "my-new-agent": "allow",   // ← add here
      "code-reviewer": "allow"
    }
  }
}
```

### Adjusting Model Tiers

To swap a model, update the `model` field for the relevant agent:

```json
"implementer": {
  "model": "opencode/gpt-5.3-codex"  // change to another model here
}
```

**Cost implications:**
- Upgrading `implementer-small` from `gpt-5.1-codex-mini` to a more expensive model increases cost for every trivial task
- Downgrading `planner` or `code-reviewer` from sonnet risks lower-quality plans and missed issues
- `docs` routing agent uses haiku deliberately — its routing logic is simple enough that spending sonnet tokens here is wasteful

### Modifying Routing Logic

The routing logic lives in the `prompt` field of `coding-boss` and `docs`. To change how tasks are classified:

1. Edit the relevant `prompt` in your configuration file
2. Preserve the JSON contract formats (`.opencode/schemas/handoff.schema.json`, `.opencode/schemas/result.schema.json`) — downstream agents validate these; `=== HANDOVER ... ===` blocks are a human-readable convention
3. Preserve the task allow-list entries for any agent you want to remain routable
4. Test with representative tasks across each routing branch

The orchestrator owns persistence/recording of both artifact types under `.opencode/sessions/<session_id>/{handoffs,results}/`: handoff artifacts it writes, and result artifacts returned by subagents.

---

## Security & Cost Considerations

**Security:**
- Routing/orchestration agents can write/edit, but should only touch `.opencode/sessions/<session_id>/{handoffs,results}/*` (and that directory is typically gitignored)
- Never grant `write` or `edit` permissions to planner or reviewer agents; doing so undermines the separation of decision and action
- `coding-boss` uses an explicit task allow-list (`"*": "deny"`) — it cannot spontaneously delegate to arbitrary agents
- Review `bash` permissions carefully when adding new agents; `ask` is safer than `allow` for agents that should not run commands autonomously

**Cost:**
- Use `implementer-small` for genuinely trivial tasks — it uses the cheapest code model in the stack
- `docs` routing agent uses haiku because routing classification requires no heavy reasoning
- Avoid routing all tasks to `implementer` (codex) when `implementer-small` (codex-mini) would suffice
- Planning agents use sonnet because a bad plan multiplies cost downstream — this is the right place to spend tokens

---

## License

AGPL-3.0 — see the [LICENSE](LICENSE) file for details.

---

## Links

- [OpenCode AI](https://opencode.ai)
- [OpenCode AI Configuration Schema](https://opencode.ai/config.json)
- [Anthropic Claude Models](https://docs.anthropic.com/en/docs/about-claude/models)
