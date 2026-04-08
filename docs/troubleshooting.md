# Troubleshooting

Common issues and solutions.

## Installation

### "command not found: oc"

The `oc` binary is not in your PATH.

**Solution:**

```sh
export PATH="$HOME/.local/bin:$PATH"
```

To make this permanent, add the line above to your shell profile (`~/.zshrc` or `~/.bashrc`).

### "install.sh not found"

The installer script is not available in the release.

**Solution:** Use an alternative install method:

```sh
# Method 1: Go install
go install github.com/sven1103-agent/opencode-config-cli@latest

# Method 2: Manual download
VERSION=v1.0.0-alpha.4
curl -L "https://github.com/sven1103-agent/opencode-config-cli/releases/download/${VERSION}/oc_${VERSION#v}_darwin_arm64.tar.gz" | tar xz
mv oc ~/.local/bin/
```

## Bundle Operations

### "source not found"

The source name doesn't exist.

**Solution:** List registered sources:

```sh
oc source list
```

### "preset not found"

The preset doesn't exist in the bundle.

**Solution:** List available presets:

```sh
oc preset list --sources
```

### "bundle validation failed"

The bundle doesn't comply with the V2 manifest contract.

**Solution:** Ensure your bundle includes `opencode-bundle.manifest.json` at the root with all required fields. See [Config Bundles](config-bundles.md).

### "checksum verification failed"

The downloaded bundle failed integrity check.

**Solution:**
1. Try updating the bundle
2. Re-register the source
3. Report the issue if it persists

## PATH Issues

### "oc: command not found" after fresh install

**Solution:**

1. Verify the binary exists:
   ```sh
   ls -la ~/.local/bin/oc
   ```

2. Add to PATH:
   ```sh
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   ```

## Migration

### "legacy config not found"

No V1 legacy config found in the project.

**Solution:** The project may already be using V2 format, or there is no OpenCode configuration present.

### "migration cancelled"

Migration was aborted by user.

**Solution:** Re-run the migration command and confirm when prompted.

## Getting Help

If your issue isn't listed here:

1. Check [existing issues](https://github.com/sven1103-agent/opencode-config-cli/issues)
2. Open a new issue with:
   - CLI version (`oc version`)
   - OS and shell
   - Install method
   - Command that failed
   - Error message
   - Reproduction steps
