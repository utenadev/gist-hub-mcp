# gist-hub-mcp

[日本語版 (README.md)](README.md)

`gist-hub-mcp` is an advanced knowledge management hub for AI agents, using GitHub Gist as its primary backend.

## Features

- **Universal Library**: Use GitHub Gist as a persistent knowledge base.
- **Wiki Mode**: Compatible with VSCode GistPad, supporting hierarchical structures using `_index.md` and `_category.json`.
- **SQLite Cache**: Avoid API rate limits and achieve fast full-text search.
- **MCP Compliant**: Provide Gist knowledge to any AI agent via the Model Context Protocol (MCP).
- **BBS Integration**: Integrate with `agent-hub` to notify knowledge updates in real-time.

## Usage (CLI)

Currently in Phase 1 of development. It uses GitHub CLI (`gh`) authentication.

```bash
# List Gists (Filtered by gist-hub: prefix)
gist-hub list

# Get a specific Gist
gist-hub get <gist_id>
```

## License

[MIT License](LICENSE)
