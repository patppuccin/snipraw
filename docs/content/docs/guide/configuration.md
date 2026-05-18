---
title: Configuration
---

# Configuration

Snipraw can be configured via a config file or flags passed directly on the command line. Flags take precedence over the config file.

## Config file

By default snipraw looks for a config file at `~/.config/snipraw/config.toml`. You can override this with the `--config` flag.

```toml
# Host and port to listen on
host = "127.0.0.1"
port = 7070

# Directory to serve snippets from
dir = "~/snippets"
```

All fields are optional. If no config file is found, snipraw falls back to defaults. The only thing required to run is a directory — either in the config file or via `--dir`.

## Reference

| Key    | Default     | Description                   |
| ------ | ----------- | ----------------------------- |
| `host` | `127.0.0.1` | Host to bind to               |
| `port` | `7070`      | Port to listen on             |
| `dir`  | —           | Directory to serve. Required. |

## Binding to all interfaces

To make snipraw accessible on your local network, bind to `0.0.0.0`:

```toml
host = "0.0.0.0"
port = 7070
dir  = "~/snippets"
```

::: warning
Binding to `0.0.0.0` exposes snipraw to anyone on your network. There is no authentication. Only do this on trusted networks.
:::

## Using a different port

If port `7070` is taken:

```bash
snipraw --port 8080 --dir ~/snippets
```

Or in the config file:

```toml
port = 8080
```

## Config file location

You can put the config file anywhere and point snipraw at it:

```bash
snipraw --config /path/to/config.toml
```

This is useful when running snipraw as a system service with a config stored in `/etc/snipraw/`.
