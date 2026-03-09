# GEMINI.md - Context for gist-hub-mcp

## Project Overview: gist-hub-mcp

**gist-hub-mcp** is an advanced knowledge management hub that uses GitHub Gist as its primary backend. It extends the concept of a BBS (agent-hub) into a persistent "Universal Library" for AI agents across multiple projects.

### Core Vision
- **Discussion in BBS, Knowledge in Gist**: Clear role separation between real-time communication (agent-hub) and long-term wisdom (Gist).
- **Semantic Library**: Implementing RAG-like capabilities to search through cross-project Gists.
- **Agent Collaboration Wiki**: Utilizing VSCode GistPad for human-agent co-authoring of specs and guidelines.

### Key Technologies
- **Language**: Go 1.24+
- **Backend**: GitHub Gist API (via `internal/github`)
- **Memory**: SQLite Cache for performance and offline access
- **Protocol**: MCP (Model Context Protocol)

## Initial Implementation Plan (Phase 1)
1. **GitHub Gist Auth & CRUD**: Basic integration with personal access tokens.
2. **Wiki Mode**: Recursive index-based document discovery.
3. **Agent Hub Integration**: Send notifications to BBS when a Gist is updated.

### Phase 1: Detailed Design (v0.1)

#### Data Structures (internal/github)
- `Gist`: ID, Description, Public, Files (map), HTMLURL, UpdatedAt.
- `GistFile`: Filename, Type, Language, RawURL, Size, Content.

#### Interface (internal/github)
- `GetToken(host string) (string, error)`: Using `github.com/cli/go-gh/v2/pkg/auth`.
- `Client`: `ListGists`, `GetGist`, `CreateGist`, `UpdateGist`, `DeleteGist`.

#### Target Implementation
- `internal/github/auth.go`: Auth implementation using gh CLI config.
- `internal/github/client.go`: REST API client implementation.

#### Gist Identification & Filtering
- **Mandatory PREFIX**: `gist-hub:` (Default).
- **Description**: All Gists managed by this tool MUST start with this prefix in their description.
- **Filtering**: CLI and MCP server should filter Gists based on this prefix by default.

## Reference Project
The stable architecture of **agent-hub-mcp** is stored in `var/agent-hub-mcp/` as a reference library.

---
*Drafted by Gemini-CLI (Architect) - 2026-03-09*
