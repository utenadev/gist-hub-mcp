# Blueprint: Universal Knowledge Platform (Semantic Search & GraphRAG)

## 1. Vision
Evolve `gist-hub-mcp` from a Gist-specific tool into a **Universal Personal Knowledge Platform**. By decoupling data ingestion from intelligence analysis, we enable secure, semantic, and relational exploration of any knowledge source.

## 2. Hybrid Architecture
- **Importers (Go)**: Handle data acquisition, decryption, and syncing. Target sources: GitHub Gists, Agent conversation logs, local memos.
- **Contract (SQLite)**: A unified data store acting as the "Single Source of Truth."
- **Analyzers (Python)**: Intelligence modules implementing a common interface for indexing and searching.

### Unified Schema (`search_source`)
```sql
CREATE TABLE IF NOT EXISTS search_source (
    uid TEXT PRIMARY KEY,    -- e.g., "gist:ID" or "log:DATE"
    source_type TEXT,        -- "gist", "agent_log", "memo"
    title TEXT,
    content TEXT,            -- Decrypted plain text
    created_at TEXT,
    updated_at TEXT,
    metadata JSON            -- Source-specific extra info
);
```

## 3. Dual Intelligence Engines
### A. Semantic Search (Lightweight & Fast)
- **Model**: `static-embedding-japanese` (CPU-optimized, high speed).
- **Storage**: `sqlite-vec` for vector similarity search (1024-dim, reduced if needed).
- **Use Case**: Quick, fuzzy recall ("What did I talk about regarding authentication?").

### B. Nano-GraphRAG (Deep Insights)
- **Logic**: Extract entities and relationships using LLM chunks.
- **Storage**: `entities`, `edges`, and `communities` tables.
- **Use Case**: Structural understanding and cross-project knowledge discovery ("How does Project A relate to the solution used in Project B?").

## 4. Security Boundary
- **Plaintext Exposure**: Decrypted data exists only within the local SQLite database and memory during Python analysis.
- **Cloud Integrity**: Only encrypted blobs are stored on GitHub Gist.

## 5. Workflow
1. **Sync (Go)**: Fetch and decrypt from Gist -> Write to `search_source`.
2. **Index (Python)**: Trigger `analyzer.index()` for new entries.
3. **Query (Python/MCP)**: Perform hybrid search (Semantic + Graph) and return unified results.

---
*Organized by Gemini-CLI (Architect) - 2026-03-12*
*Based on the "Pluggable Importers & Analyzers" vision.*
