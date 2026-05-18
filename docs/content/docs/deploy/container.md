---
title: Container
---

# Container

Snipraw publishes a Docker image to GitHub Container Registry.

## Docker

```bash
docker run -d
  --name snipraw \
  -p 8245:8245 \
  -v ~/snippets:/snippets \
  ghcr.io/patppuccin/snipraw:latest \
  --dir /snippets
```

Open `http://localhost:8245` to verify it is running.

## Docker Compose

```yaml
services:
  snipraw:
    image: ghcr.io/patppuccin/snipraw:latest
    container_name: snipraw-server
    restart: unless-stopped
    ports:
      - "8245:8245"
    volumes:
      - ~/snippets:/snippets
```

Start it:

```bash
docker compose up -d
```

## Passing a config file

Mount your config alongside the snippets directory:

```yaml
services:
  snipraw:
    image: ghcr.io/patppuccin/snipraw:latest
    container_name: snipraw-server
    restart: unless-stopped
    ports:
      - "8245:8245"
    volumes:
      - ~/snippets:/snippets
      - ./config.yaml:/config.yaml
```

::: tip
Pin to a specific version tag in production so updates don't catch you off guard.
:::
