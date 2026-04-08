# CLI E2E Coverage

## Purpose

This document tracks black-box end-to-end coverage for the shipped Go CLI binary (`oc`).

It is the durable inventory for what is already covered, what is verified in CI, and what is still deferred.

## Scope

- Test the built `oc` binary, not internal package functions
- Focus on deterministic, offline scenarios first
- Track both happy paths and operator-facing failure paths
- Keep coverage aligned with supported Linux and macOS release targets

## Principles

- Run tests against the compiled binary via `os/exec`
- Isolate each test with temp directories and isolated user config state
- Keep fixtures local and deterministic
- Avoid live network dependencies in CI
- Expand coverage incrementally and update this document in the same PR

## Workflow

- Workflow file: `.github/workflows/e2e-cli.yml`
- Binary handoff env var: `OC_E2E_BINARY`
- Current CI platforms:
  - `ubuntu-latest`
  - `macos-latest`
- Latest verified `main` run:
  - GitHub Actions run `23850747833`
  - `E2E CLI (ubuntu-latest)`: success
  - `E2E CLI (macos-latest)`: success

## Test Harness

- Package: `e2e/`
- Fixture root: `e2e/testdata/`
- Execution model:
  - build `oc`
  - invoke commands via subprocesses
  - set isolated `HOME` and `XDG_CONFIG_HOME`
  - verify stdout, stderr, exit status, and written files

## Covered Scenarios

| Scenario | Command(s) | Source Type | Status | CI Verified | Notes |
|---|---|---|---|---|---|
| Print version | `oc version` | n/a | Implemented | ✅ | Verified by successful `main` workflow matrix run on Linux and macOS |
| Register local directory source | `oc source add <dir>` | local directory | Implemented | ✅ | Uses fixture bundle with manifest |
| Register local archive source | `oc source add <tar.gz>` | local archive | Implemented | ✅ | Archive generated at test runtime |
| List registered sources | `oc source list` | registry | Implemented | ✅ | Covered inside the local directory flow |
| Apply preset from local directory source | `oc bundle apply <id> --preset <name>` | local directory | Implemented | ✅ | Verifies written config and provenance |
| Apply preset from local archive source | `oc bundle apply <id> --preset <name>` | local archive | Implemented | ✅ | Exercises tar extraction path |
| Show applied bundle provenance | `oc bundle status --project-root <dir>` | project provenance | Implemented | ✅ | Reads `.opencode/bundle-provenance.json` through CLI |
| Refuse overwrite without force | `oc bundle apply <id> --preset <name>` | local directory | Implemented | ✅ | Must fail when output already exists |
| Reject missing manifest source | `oc source add <dir>` | invalid local directory | Implemented | ✅ | Validates manifest presence check |
| Reject invalid tarball | `oc source add <tar.gz>` and/or `oc bundle apply <id> --preset <name>` | invalid local archive | Implemented | ✅ | Verifies archive extraction failure path |
| Reject unknown source ID | `oc bundle apply <id> --preset <name>` | registry lookup | Implemented | ✅ | Verifies user-facing error path |

## Deferred Scenarios

| Gap | Reason | Planned Follow-up |
|---|---|---|
| GitHub release source E2E | Current implementation still has remote operations marked not implemented in this branch line | Add once remote source behavior is stable end to end |
| `httptest`-backed remote source coverage | Not needed for initial local-flow milestone | Add with remote-source implementation follow-up |
| `bundle update` E2E | Current command behavior is intentionally limited and network-dependent | Add when update behavior is production-ready |
| Windows runner coverage | Windows is out of current scope | Revisit only if product scope changes |

## Fixture Strategy

- Keep a small valid local bundle fixture in `e2e/testdata/fixture-bundle/`
- Generate `.tar.gz` archives during test execution instead of committing binary archives
- Keep invalid cases lightweight and explicit

## Environment Isolation

- Each test uses a fresh temporary project root
- Each test uses a fresh temporary config home
- Subprocess env sets:
  - `HOME`
  - `XDG_CONFIG_HOME`

## Change Log

- 2026-04-01: Document created for `US-052` initial local-flow coverage planning
- 2026-04-01: Initial local-flow scenarios implemented in `e2e/` and wired into `.github/workflows/e2e-cli.yml`
- 2026-04-01: Marked initial local-flow scenarios as CI-verified after successful `main` workflow run `23850747833` on Linux and macOS
