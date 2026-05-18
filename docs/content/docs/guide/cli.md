---
title: Command Line
---

# Command Line

Snipraw's CLI is a simple one level command with no subcommands. Use the `--help` flag to see the full list of flags amd options.

## Available flags and options

```sh
snipraw --help
```

Available flags and options are listed below.

| Flag           | Default     | Description                      |
| -------------- | ----------- | -------------------------------- |
| `--host`       | `127.0.0.1` | Host to bind to                  |
| `--port`       | `8245`      | Port to bind to                  |
| `--dir`        | —           | Directory to serve snippets from |
| `--log-level`  | `info`      | log level                        |
| `--config`     | —           | Path to config file              |
| `--help`, `-h` | false       | help for snipraw                 |
| `--help`, `-v` | false       | version for snipraw              |

## Examples

```bash
# Run with a specific directory
snipraw --dir ~/snippets

# Run on a different port
snipraw --dir ~/snippets --port 8080

# Run with a custom config file
snipraw --config /etc/snipraw/config.toml

```
