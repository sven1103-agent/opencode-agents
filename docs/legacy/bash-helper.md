# Legacy: Bash Helper

> **Status: Deprecated** — The Bash-based helper is no longer maintained. Please use the Go CLI (`oc`) instead.

## Overview

The original iteration of the OpenCode helper was implemented as a Bash script (`scripts/opencode-helper`). This version is deprecated and receives no further updates.

## Migration

Migrate to the Go CLI:

1. [Install the Go CLI](../installation.md)
2. Register your config source:
   ```sh
   oc source add qbicsoftware/opencode-config-bundle --name qbic
   ```
3. Apply a preset:
   ```sh
   oc bundle apply qbic --preset mixed --project-root .
   ```

## Legacy Projects

Existing projects using the Bash helper will continue to work. To migrate to V2:

```sh
oc migrate legacy-config --project-root ./your-project
```

## Archive Location

The old Bash scripts are preserved in the repository at:
- `scripts/opencode-helper` — Main helper script
- `scripts/opencode-helper-install` — Installation script

These are kept for reference only and should not be used for new projects.
