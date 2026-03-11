# Mission: Upgrade agent-hub-mcp to Streamable HTTP

## 1. Context
Current SSE implementation in agent-hub-mcp is legacy and incompatible with some agents like Antigravity.
We need to upgrade the server to the latest Model Context Protocol (v2025-03-26) specification.

## 2. Goals
- Implement **Streamable HTTP (v2025-03-26)** on a new endpoint (e.g., `/mcp`).
- Ensure **CORS** headers are correctly set for PWA/Browser access.
- Maintain backward compatibility with the existing `/sse` endpoint.

## 3. Work Directory
Move to: `/home/kench/workspace/agent-hub-mcp`

## 4. Team
- **Architect**: Gemini-CLI
- **Implementer**: OpenCode (to be started after relocation)

---
*Created by Gemini-CLI (Architect) via gist-hub-mcp CLI*
