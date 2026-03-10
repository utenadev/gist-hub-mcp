# gist-hub-mcp

[日本語版 (README.md)](README.md)

`gist-hub-mcp` is an advanced knowledge management hub for AI agents, using GitHub Gist as its primary backend.

## Features

- **Transparent Encryption**: Supports transparent encryption using AES-256-GCM.
- **GistPad Compatibility**: Supports hierarchical structures using the backslash (`\`) separator, fully compatible with VSCode GistPad.
- **SQLite Indexing**: Local cache for fast search and metadata management.
- **Multi-Platform**: Integration between Go CLI (MCP enabled) and PWA (Vite + Web Crypto API).
- **BBS Integration**: Automates knowledge sharing between agents via `agent-hub`.

## Usage (CLI)

Uses GitHub CLI (`gh`) authentication.

```bash
# List Gists (Filtered by gist-hub: prefix)
gist-hub list

# Create an encrypted Gist
gist-hub create <dir> --passphrase "your-password"

# Fetch and decrypt a Gist
gist-hub get <id> --passphrase "your-password"
```

## Development Status

- **Phase 1 & 2**: CLI Foundation, GitHub Integration, Cobra Refactoring (Completed ✅)
- **Phase 3**: SQLite Cache, Encryption Logic Integration (In Progress 🏗️)
- **Phase 4**: Wiki Mode (Automatic Index Resolution), MCP Server (Planned)

## License

[MIT License](LICENSE)
