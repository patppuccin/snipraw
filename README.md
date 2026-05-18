# snipraw

[![Go](https://img.shields.io/badge/go-1.26+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)
[![License](https://img.shields.io/badge/license-Apache%202.0-6fa783?style=flat)](LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/patppuccin/snipraw/release.yml?style=flat&label=build)](https://github.com/patppuccin/snipraw/actions)

Snipraw is a self-hosted code snippet server that serves files directly from your filesystem. No database, no syncing, no accounts. Point it at a directory and your snippets are instantly browsable.

Snipraw does these things:

- Serves files from a directory as browsable snippets
- Renders markdown and highlights code via Chroma
- Exposes a `/raw` endpoint for every file
- Ships as a single binary with no runtime dependencies

Find the full documentation at [snipraw.patppuccin.com](https://snipraw.patppuccin.com)

## Quick start

`snipraw` is distributed as a single standalone binary and runs on **Linux, macOS, and Windows** and no runtime dependencies or package managers required.

You can install it using the convenience installer or by downloading a prebuilt release.

#### Convenience Script (Recommended)

The install script detects your platform and architecture, downloads the **latest release**, and installs the binary into your system `PATH`.

**Linux & macOS**

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/patppuccin/snipraw/main/scripts/install.sh)
```

**Windows (PowerShell 5.1+)**

```powershell
powershell -ExecutionPolicy Bypass -NoProfile -c "irm https://raw.githubusercontent.com/patppuccin/snipraw/main/scripts/install.ps1 | iex"
```

#### Install via Prebuilt Binary

You can also install `snipraw` manually.

1. Download the appropriate binary for your **OS and architecture** from the
   [GitHub Releases](https://github.com/patppuccin/snipraw/releases) page.
2. Extract the archive.
3. Move the binary to a directory in your `PATH`.

Example (Linux / macOS):

```bash
chmod +x snipraw
sudo mv snipraw /usr/local/bin/
```

Example (Windows):

Move `snipraw.exe` to a directory included in your `PATH`.

> [!TIP]
> Run the following command to confirm the installation:
>
> ```sh
> snipraw --version
> ```
>
> If the command prints the version, the installation succeeded.

Open `http://localhost:8245` in your browser.

## Disclaimer

Snipraw is not a vibe-coded application. LLMs were used during development as a thinking partner to explore architecture ideas, critique design decisions, and research shell edge cases.

The code, design decisions, and trade-offs are the author's own.

## Contributing

Contributions are not being accepted at this time while the project structure and roadmap are still evolving.

Issues are welcome for:

- bug reports
- feature ideas
- feedback

Pull requests may be opened once the project stabilizes.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) for details.
