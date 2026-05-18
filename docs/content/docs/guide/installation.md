---
title: Installation
---

# Installation

Snipraw ships as a single binary with no runtime dependencies. Download it, make it executable, and run it.

## Download

Grab the latest release for your platform from the [GitHub releases page](https://github.com/patppuccin/snipraw/releases).

| Platform        | File                        |
| --------------- | --------------------------- |
| Linux (amd64)   | `snipraw_linux_amd64`       |
| macOS (arm64)   | `snipraw_darwin_arm64`      |
| Windows (amd64) | `snipraw_windows_amd64.exe` |

## Linux and macOS

```bash
# Download the binary (replace with your platform)
curl -Lo snipraw https://github.com/patppuccin/snipraw/releases/latest/download/snipraw_linux_amd64

# Make it executable
chmod +x snipraw

# Move to somewhere on your PATH
mv snipraw /usr/local/bin/snipraw
```

Verify the installation:

```bash
snipraw --version
```

## Windows

Download the `.exe` from the releases page and place it somewhere on your `PATH`, or run it directly from a directory.

```powershell
# Verify
.\snipraw.exe version
```

## First run

Point snipraw at a directory and it starts serving immediately:

```bash
snipraw --dir ~/snippets
```

Open `http://localhost:8245` in your browser. If the directory doesn't exist, snipraw will exit with an error — it won't create it for you.

## Updating

Download the new binary and replace the old one. Snipraw has no state to migrate.

::: tip
If you installed via a package manager or custom path, make sure the new binary lands in the same location.
:::
