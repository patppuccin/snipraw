---
title: Deploy
---

# Deploy

Snipraw is a single binary that runs as a plain HTTP server. Deploying it is largely a matter of deciding how you want it to start and whether you want it behind a reverse proxy.

## Options

- [Container](/docs/deploy/container) — run snipraw in Docker or with Docker Compose
- [System Service](/docs/deploy/service) — run snipraw as a systemd service on Linux or a Windows service
- [Reverse Proxy](/docs/deploy/reverse-proxy) — put snipraw behind Caddy or nginx

## Choosing an approach

If you are running snipraw on a server and want it to start automatically on boot, use the system service guide. If you want to expose it publicly with TLS, pair the system service with a reverse proxy. If you prefer containers, the Docker setup handles both.

For local use on a workstation, none of this is necessary — just run `snipraw --dir ~/snippets` and add it to your shell startup if you want it persistent.
