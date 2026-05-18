// scripts/ci-changelog-gen.mjs
// Fetches GitHub releases and writes docs/content/changelog.md.
// Usage: node scripts/ci-changelog-gen.mjs

import { writeFileSync, mkdirSync } from "fs";
import { dirname } from "path";

const REPO = "patppuccin/snipraw";
const OUTPUT = "docs/content/changelog.md";

const log_inf = (msg) => console.log(`\x1b[34mINF\x1b[0m ${msg}`);
const log_wrn = (msg) => console.log(`\x1b[33mWRN\x1b[0m ${msg}`);
const log_err = (msg) => {
  console.error(`\x1b[31mERR\x1b[0m ${msg}`);
  process.exit(1);
};

function buildHeaders() {
  const headers = {
    Accept: "application/vnd.github+json",
    "X-GitHub-Api-Version": "2022-11-28",
  };
  if (process.env.GITHUB_TOKEN) {
    headers["Authorization"] = `token ${process.env.GITHUB_TOKEN}`;
  }
  return headers;
}

async function fetchJSON(url, headers) {
  const res = await fetch(url, { headers });
  if (!res.ok)
    log_err(`Failed to fetch ${url}: ${res.status} ${res.statusText}`);
  return res.json();
}

function formatDate(published_at) {
  return new Date(published_at).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

function getBadge(index, prerelease) {
  if (index === 0) return ' <Badge type="tip" text="latest" />';
  if (prerelease) return ' <Badge type="danger" text="deprecated" />';
  return "";
}

async function fetchCommits(base, head, headers) {
  const url = `https://api.github.com/repos/${REPO}/compare/${base}...${head}`;
  const res = await fetch(url, { headers });
  if (!res.ok) {
    log_wrn(`Failed to compare ${base}...${head}`);
    return null;
  }
  return res.json();
}

function buildFrontmatter() {
  return `---
title: Changelog
---

# Changelog

All notable changes to snipraw are documented here.
Follows [Keep a Changelog](https://keepachangelog.com) and [Semantic Versioning](https://semver.org).

## Releases

`;
}

async function buildChangelog(releases, headers) {
  let md = buildFrontmatter();

  for (let i = 0; i < releases.length; i++) {
    const release = releases[i];
    const prevRelease = releases[i + 1];
    const date = formatDate(release.published_at);
    const badge = getBadge(i, release.prerelease);

    md += `### ${release.tag_name} (${date})${badge} \n\n`;
    md += `View this release on [GitHub](${release.html_url}).\n\n`;

    if (prevRelease) {
      const compare = await fetchCommits(
        prevRelease.tag_name,
        release.tag_name,
        headers,
      );

      if (compare?.commits?.length) {
        for (const commit of compare.commits) {
          const sha = commit.sha.slice(0, 7);
          const message = commit.commit.message.split("\n")[0];
          md += `- \`${sha}\` ${message}\n`;
        }
      } else {
        md += "_No commits found._\n";
      }
    } else {
      md += "_Initial release._\n";
    }

    md += "\n";
  }

  return md;
}

async function main() {
  console.log();
  console.log("\x1b[32mSnipraw Changelog Generator\x1b[0m");
  console.log(`Repo:   https://github.com/${REPO}`);
  console.log(`Output: ${OUTPUT}`);
  console.log();

  const headers = buildHeaders();
  const releases = await fetchJSON(
    `https://api.github.com/repos/${REPO}/releases`,
    headers,
  );

  if (!releases.length) {
    log_wrn("No releases found");
    process.exit(0);
  }

  log_inf(`Fetched ${releases.length} releases`);

  const md = await buildChangelog(releases, headers);

  mkdirSync(dirname(OUTPUT), { recursive: true });
  writeFileSync(OUTPUT, md);

  log_inf(`Changelog written to ${OUTPUT} with ${releases.length} releases`);
}

main().catch((err) => log_err(err.message));
