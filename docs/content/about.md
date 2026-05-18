---
title: About
---

# About Snipraw

Snipraw started as a personal frustration. Code snippets scattered across Gists, Notion pages, and random text files — none of it browsable in one place without depending on a third-party service or setting up a database.

The goal was simple: point it at a directory, get a clean web interface, done.

## Philosophy

**No database.** Your snippets are plain files. You can edit them with any editor, version them with git, back them up with rsync. Snipraw never owns your data.

**Single binary.** No runtime dependencies, no installation wizard, no configuration ceremony. Download, run, done.

**Self-hosted.** Your snippets stay on your machine. No accounts, no telemetry, no third-party services involved.

**Boring tech.** The stack is deliberately simple — proven tools that do their job without drama.

## Built With

Snipraw is built with the following open source tools:

- [Go](https://go.dev) — the binary
- [templ](https://templ.guide) — type-safe HTML templates
- [DaisyUI](https://daisyui.com) — UI components
- [Tailwind CSS](https://tailwindcss.com) — styling
- [Chroma](https://github.com/alecthomas/chroma) — syntax highlighting
- [Zerolog](https://github.com/rs/zerolog) — structured logging
- [Cobra](https://github.com/spf13/cobra) — CLI
- [VitePress](https://vitepress.dev) — this documentation site
