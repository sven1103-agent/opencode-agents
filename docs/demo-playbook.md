# CLI Demo Playbook

Repeat this playbook when recording a short terminal demo for the repository landing page. It is designed to show the core value of `oc` in under 30 seconds using a deterministic happy path.

## Goal

Show that `oc` can:

- register a configuration bundle source
- inspect the registered source list
- apply a preset into a project
- show provenance for the applied bundle

## Demo Shape

- Duration: 20 to 30 seconds
- Style: terminal-only, no voice-over required
- Scope: happy path only
- Output asset target: README-friendly GIF or preview image plus the source recording

## Commands

Run the demo in a clean temporary directory so the output stays predictable.

```sh
mkdir -p demo-project
oc --version
oc source add qbicsoftware/opencode-config-bundle --name qbic
oc source list
oc bundle apply qbic --preset mixed --project-root demo-project
oc bundle status --project-root demo-project
ls demo-project
```

## Recording Notes

- Keep the terminal large enough to avoid line wrapping.
- Use a clean shell prompt with minimal noise.
- Prefer a fresh OpenCode source registry before recording so `oc source add` succeeds on the first try.
- If the source is already registered on the recording machine, clear it first or use a fresh config home.
- Keep typing speed steady and slightly faster than normal conversation speed.
- Leave a short pause after each command so the rendered GIF is readable.

## Expected Story Beat

Use this sequence when recording:

1. Show the installed CLI version.
2. Register the official bundle source as `qbic`.
3. Show that the source is now available in the local registry.
4. Apply the `mixed` preset into `demo-project`.
5. Show bundle provenance for the generated project config.
6. End by showing that `opencode.json` exists in the project directory.

## Environment Prep

Prepare the environment before starting the actual recording:

```sh
rm -rf demo-project
mkdir -p demo-project
```

If you want a fully isolated source registry for recording, prefer an ephemeral config home:

```sh
export XDG_CONFIG_HOME="$PWD/.demo-config"
rm -rf "$XDG_CONFIG_HOME"
mkdir -p "$XDG_CONFIG_HOME"
```

## Capture Workflow

Recommended tools:

- `vhs` for reproducible terminal capture and GIF rendering
- `scripts/record-demo.sh` in this repo for the current committed capture flow

Install `vhs` with Homebrew if needed:

```sh
brew install vhs
```

Example capture flow:

```sh
scripts/record-demo.sh
```

This script builds the CLI in a temporary workspace and renders `docs/demo.gif` from `docs/demo.tape`.

Keep `docs/demo.tape` as the source artifact so the README asset can be re-rendered for later releases without inventing a new script.

## Re-record Checklist

Before recording against a new CLI release:

1. Confirm the commands in this playbook still match the current CLI.
2. Confirm the selected preset name still exists in the referenced bundle.
3. Confirm the command output still fits inside the recording frame.
4. Re-run `scripts/record-demo.sh`.
5. Re-render the README asset from the updated tape.
6. Update README references if the demo asset path changes.
