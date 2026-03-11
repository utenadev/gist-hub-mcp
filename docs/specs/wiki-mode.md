# Wiki Mode: Virtual Directory Tree Resolution Spec (v0.6)

## 1. Goal
Convert flat Gist file lists (stored in SQLite with backslashes) into a hierarchical directory tree for MCP resources and PWA display.

## 2. Input Data
From SQLite \`files\` table:
- \`original_path\`: "docs/architecture/overview.md"
- \`gist_filename\`: "docs\\architecture\\overview.md"

## 3. Core Logic (internal/wiki)

### Data Structure
\`\`\`go
type Node struct {
    Name     string
    Path     string           // Full slash-separated path (e.g., "docs/architecture")
    IsDir    bool
    FileID   string           // Gist ID if it's a file
    Children map[string]*Node // Sub-nodes
    Metadata *CategoryMeta    // From _category.json if exists
}

type CategoryMeta struct {
    Label       string \`json:"label,omitempty"\`
    Position    int    \`json:"position,omitempty"\`
    Collapsible bool   \`json:"collapsible,omitempty"\`
}
\`\`\`

### Resolution Rules
1. **Separator**: Use \`\\\` (backslash) to split \`gist_filename\` into segments.
2. **Directory Index**: 
   - If a directory contains \`_index.md\`, it is the primary content for that path.
   - Fallback to \`README.md\` if \`_index.md\` is missing.
3. **Metadata**: 
   - \`_category.json\` provides \`label\`, \`position\`, and \`collapsible\` properties.

## 4. MCP Tools (internal/mcp)
- \`wiki_list(path string)\`: List children of a specific virtual directory.
- \`wiki_read(path string)\`: Resolve path to content (using index rules).

## 5. Next Steps for Sisyphus
1. Create \`internal/wiki\` package.
2. Implement \`BuildTree(files []db.File) *Node\`.
3. Add unit tests verifying deep hierarchy resolution (e.g., \`a\\b\\c.md\`).

---
*Drafted by Gemini-CLI (Architect) - 2026-03-11*
