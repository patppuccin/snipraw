// scripts/ci-deps-setup.mjs
// Downloads Tailwind CSS CLI and templ binaries to tools/.
// Skips download if binary already exists.

import { writeFileSync, chmodSync, mkdirSync, existsSync } from "fs";
import { unlinkSync } from "fs";
import { join } from "path";
import { platform, arch } from "process";
import { execSync } from "child_process";

const TOOLS_DIR = "tools";
const TAILWIND_BASE =
  "https://github.com/tailwindlabs/tailwindcss/releases/latest/download";
const TEMPL_REPO = "a-h/templ";

const log_inf = (msg) => console.log(`\x1b[34mINF\x1b[0m ${msg}`);
const log_wrn = (msg) => console.log(`\x1b[33mWRN\x1b[0m ${msg}`);
const log_err = (msg) => {
  console.error(`\x1b[31mERR\x1b[0m ${msg}`);
  process.exit(1);
};

const isWindows = platform === "win32";

function binName(name) {
  return isWindows ? `${name}.exe` : name;
}

// Tailwind OS/arch naming
function getTailwindFilename() {
  const os = { linux: "linux", darwin: "macos", win32: "windows" }[platform];
  const cpu = { x64: "x64", arm64: "arm64" }[arch];
  if (!os) log_err(`Unsupported platform: ${platform}`);
  if (!cpu) log_err(`Unsupported architecture: ${arch}`);
  return `tailwindcss-${os}-${cpu}${isWindows ? ".exe" : ""}`;
}

// templ OS/arch naming — matches actual release assets exactly
// e.g. templ_Linux_x86_64.tar.gz, templ_Darwin_arm64.tar.gz, templ_Windows_x86_64.tar.gz
function getTemplAssetName(version) {
  const os = { linux: "Linux", darwin: "Darwin", win32: "Windows" }[platform];
  const cpu = { x64: "x86_64", arm64: "arm64" }[arch];
  if (!os) log_err(`Unsupported platform: ${platform}`);
  if (!cpu) log_err(`Unsupported architecture: ${arch}`);
  return `templ_${os}_${cpu}.tar.gz`;
}

async function fetchJSON(url) {
  const res = await fetch(url);
  if (!res.ok)
    log_err(`Failed to fetch ${url}: ${res.status} ${res.statusText}`);
  return res.json();
}

async function download(url, destPath) {
  const res = await fetch(url);
  if (!res.ok)
    log_err(`Failed to download ${url}: ${res.status} ${res.statusText}`);
  const buffer = await res.arrayBuffer();
  writeFileSync(destPath, Buffer.from(buffer));
}

async function getLatestTemplVersion() {
  const data = await fetchJSON(
    `https://api.github.com/repos/${TEMPL_REPO}/releases/latest`,
  );
  return data.tag_name;
}

async function setupTailwind() {
  const dest = join(TOOLS_DIR, binName("tailwindcss"));

  if (existsSync(dest)) {
    log_wrn(`Tailwind CLI already exists at ${dest}, skipping`);
    return;
  }

  const filename = getTailwindFilename();
  const url = `${TAILWIND_BASE}/${filename}`;

  log_inf(`Downloading Tailwind CSS CLI (${filename})`);
  await download(url, dest);
  if (!isWindows) chmodSync(dest, 0o755);
  log_inf(`Tailwind CSS CLI installed to ${dest}`);
}

async function setupTempl() {
  const dest = join(TOOLS_DIR, binName("templ"));

  if (existsSync(dest)) {
    log_wrn(`templ already exists at ${dest}, skipping`);
    return;
  }

  const version = await getLatestTemplVersion();
  const assetName = getTemplAssetName();
  const url = `https://github.com/${TEMPL_REPO}/releases/download/${version}/${assetName}`;
  const archivePath = join(TOOLS_DIR, assetName);

  log_inf(`Downloading templ ${version} (${assetName})`);
  await download(url, archivePath);

  // all templ archives are .tar.gz including Windows
  execSync(
    `tar -xzf "${archivePath}" -C "${TOOLS_DIR}" templ${isWindows ? ".exe" : ""}`,
  );
  unlinkSync(archivePath);

  if (!isWindows) chmodSync(dest, 0o755);
  log_inf(`templ ${version} installed to ${dest}`);
}

async function main() {
  console.log();
  console.log("\x1b[32mSnipraw Dependency Setup\x1b[0m");
  console.log(`Tools directory: ${TOOLS_DIR}`);
  console.log();

  mkdirSync(TOOLS_DIR, { recursive: true });

  await setupTailwind();
  await setupTempl();

  console.log();
  log_inf("Done");
}

main().catch((err) => log_err(err.message));
