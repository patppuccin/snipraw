---
title: Getting Started
---

# Getting started

Snipraw is a self-hosted code snippet server that serves files directly from your filesystem. No database, no syncing, no accounts. Point it at a directory and your snippets are instantly browsable.

It is built for developers who already have a folder of scripts, configs, and code fragments and just want a fast way to browse and share them without reaching for a SaaS tool.

## How it works

Snipraw reads a directory you give it. Each subdirectory becomes a project, and each file inside becomes a snippet.

Code & Markdown files render with syntax highlighting. Everything else serves as plain text.

## What snipraw is not

- Not a Gist replacement with versioning and social features
- Not a mult-user collaboration tool
- Not a pastebin with expiring links
- Not a note-taking app

If you need any of those things, snipraw is the wrong tool. If you have a `~/snippets` folder that you wish had a UI, it is the right one.

## Where to go next

If you are setting up for the first time, start with the [installation](/docs/guide/installation) guide. It covers downloading the binary, pointing it at a directory, and getting the server running.

For configuration options, check out the [configurations](/docs/guide/configuration) page. For deployment behind a reverse proxy or as a system service, visit the [deploy](/docs/deploy/) section.
