---
name: fp-go-mcp
description: Use this skill when working with the fp-go MCP (Model Context Protocol) server located in github.com/IBM/fp-go/v2/gen. Trigger on mentions of MCP server, fp-go MCP tools, list_skills, use_skill, search_examples, get_example, configuring fp-go for Claude Desktop, or when the user needs to access fp-go examples and skills programmatically. This skill explains how to configure and use the MCP server to access fp-go documentation, examples, and skills.
---

# fp-go MCP Server

## Overview

The fp-go MCP (Model Context Protocol) server provides programmatic access to fp-go skills, examples, and documentation. It's located in `github.com/IBM/fp-go/v2/gen` and can be integrated into any MCP-compatible client such as Claude Desktop.

**Prerequisite**: the working directory must contain a `go.mod` that lists `github.com/IBM/fp-go/v2/gen` as a tool dependency (added via `go get -tool`). The server is launched with `go tool gen mcp` and uses stdio transport, so no global installation is needed.

The fp-go MCP server exposes four tools:

1. **`list_skills`** — List all available fp-go skills
2. **`use_skill`** — Retrieve the full content of a specific skill
3. **`search_examples`** — Search for Go examples using full-text search
4. **`get_example`** — Retrieve a specific example by symbol name

## Installation

Add the fp-go generator as a tool dependency in your project:

```bash
go get -tool github.com/IBM/fp-go/v2/gen
```

This makes the tool available via `go tool gen` without requiring global installation.

## Configuration

### For Claude Desktop (Anthropic)

Add to your Claude Desktop configuration file:

**macOS/Linux**: `~/Library/Application Support/Claude/claude_desktop_config.json`

**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "fp-go": {
      "command": "go",
      "args": ["tool", "gen", "mcp"],
      "env": {}
    }
  }
}
```

After configuration, restart Claude Desktop to activate the server.

### For Other MCP Clients

The server uses stdio transport and follows the MCP specification. Configure your client to:
- Execute: `go tool gen mcp`
- Use stdio for communication
- Optionally add `--verbose` flag for debugging

## Available Tools

### 1. list_skills

**Description**: List all available fp-go skills with their names, descriptions, and paths.

**Parameters**: None

**Returns**:
```json
{
  "skills": [
    {
      "name": "fp-go",
      "description": "Core fp-go patterns and best practices",
      "path": "fp-go"
    },
    {
      "name": "fp-go-pipe-flow",
      "description": "Pipe and Flow composition patterns",
      "path": "fp-go-pipe-flow"
    }
  ]
}
```

**Example Usage**:
```
User: List all available fp-go skills
Assistant: [calls list_skills tool]
```

**Use Cases**:
- Discover available skills
- Find skills for specific topics
- Get an overview of fp-go documentation

### 2. use_skill

**Description**: Retrieve the full content of a specific skill by name.

**Parameters**:
- `name` (required): The skill name from `list_skills` (e.g., "fp-go", "fp-go-pipe-flow")

**Returns**:
```json
{
  "name": "fp-go",
  "description": "Core fp-go patterns and best practices",
  "content": "# fp-go v2 — Functional Programming for Go\n\n..."
}
```

**Example Usage**:
```
User: Show me the fp-go-pipe-flow skill
Assistant: [calls use_skill with name="fp-go-pipe-flow"]
```

**Use Cases**:
- Load detailed documentation for a specific topic
- Get best practices and examples
- Reference API patterns

**Note**: The content excludes the YAML frontmatter header — only the markdown content is returned.

### 3. search_examples

**Description**: Search for Go examples using full-text search across example names, symbols, packages, documentation comments, and code.

**Parameters**:
- `query` (required): Search query (supports SQLite FTS5 syntax)
- `package_filter` (optional): Filter by package name (e.g., "option", "either")

**Returns**:
```json
{
  "examples": [
    {
      "id": "option::ExampleMap",
      "package": "github.com/IBM/fp-go/v2/option",
      "symbol": "Map",
      "name": "ExampleMap",
      "doc_comment": "// ExampleMap demonstrates mapping over an Option",
      "code": "func ExampleMap() {\n\t...\n}",
      "output": "Some(42)\n",
      "imports": "import O \"github.com/IBM/fp-go/v2/option\"",
      "file": "option/option_test.go"
    }
  ],
  "count": 1
}
```

**Search Syntax**:
- Simple terms: `"Map"` — finds examples mentioning Map
- Phrases: `"\"point free\""` — exact phrase match
- Boolean: `"Map AND Option"` — both terms required
- Wildcards: `"Trav*"` — matches Traverse, TraverseArray, etc.
- Package filter: `query="Map", package_filter="option"` — only option package

**Example Usage**:
```
User: Find examples of using Map with Option
Assistant: [calls search_examples with query="Map Option", package_filter="option"]

User: Show me examples of TraverseArray
Assistant: [calls search_examples with query="TraverseArray"]
```

**Use Cases**:
- Find examples for specific functions
- Discover usage patterns
- Learn from working code
- Find examples in a specific package

**Limits**: Returns up to 10 results, ranked by relevance.

### 4. get_example

**Description**: Retrieve a specific example by symbol name. Performs exact lookup by symbol or function name.

**Parameters**:
- `symbol` (required): Symbol name (e.g., "Map", "Type.Method") or function name (e.g., "ExampleMap")

**Returns**:
```json
{
  "examples": [
    {
      "id": "option::ExampleMap",
      "package": "github.com/IBM/fp-go/v2/option",
      "symbol": "Map",
      "name": "ExampleMap",
      "doc_comment": "// ExampleMap demonstrates mapping over an Option",
      "code": "func ExampleMap() {\n\t...\n}",
      "output": "Some(42)\n",
      "imports": "import O \"github.com/IBM/fp-go/v2/option\"",
      "file": "option/option_test.go"
    }
  ],
  "count": 1
}
```

**Example Usage**:
```
User: Get the example for Option.Map
Assistant: [calls get_example with symbol="Map"]

User: Show me the ExampleTraverseArray example
Assistant: [calls get_example with symbol="ExampleTraverseArray"]
```

**Use Cases**:
- Get a specific example by name
- Retrieve all examples for a symbol
- Access complete example code with imports and output

**Note**: Supports both exact matches and pattern matching (LIKE queries).

## Workflow Examples

### Discovering and Using Skills

```
1. User: "What fp-go skills are available?"
   → Assistant calls list_skills
   → Returns: fp-go, fp-go-pipe-flow, fp-go-http, fp-go-logging, fp-go-lens, fp-go-pr-review, fp-go-mcp

2. User: "Show me the fp-go-pipe-flow skill"
   → Assistant calls use_skill(name="fp-go-pipe-flow")
   → Returns full skill content with examples and best practices

3. User: "Now help me refactor this code using Pipe"
   → Assistant uses the loaded skill to provide guidance
```

### Finding Examples

```
1. User: "How do I use TraverseArray?"
   → Assistant calls search_examples(query="TraverseArray")
   → Returns ranked examples with code

2. User: "Show me more examples from the array package"
   → Assistant calls search_examples(query="Traverse", package_filter="array")
   → Returns array-specific examples

3. User: "Get the exact example for Array.Map"
   → Assistant calls get_example(symbol="Array.Map")
   → Returns the specific example with full code
```

### Combined Workflow

```
1. User: "I need to work with Option types"
   → Assistant calls list_skills to find relevant skills
   → Loads fp-go skill with use_skill(name="fp-go")

2. User: "Show me examples of Option.Map"
   → Assistant calls search_examples(query="Map", package_filter="option")
   → Returns examples with code and output

3. User: "How do I compose multiple Option operations?"
   → Assistant references the loaded skill content
   → Calls search_examples(query="Flow Option") for examples
   → Provides guidance based on skill + examples
```

## Troubleshooting

### Server Won't Start

**Issue**: Tool not found or not working

**Solution**: Ensure the tool is installed in your project:
```bash
go get -tool github.com/IBM/fp-go/v2/gen
```

Then run with:
```bash
go tool gen mcp
```

### No Skills Found

**Issue**: `list_skills` returns empty list

**Solution**: The skills are embedded at build time. Update the tool:
```bash
go get -tool github.com/IBM/fp-go/v2/gen@latest
```

### Search Returns No Results

**Issue**: `search_examples` finds nothing

**Solution**: 
1. Check your search query syntax (SQLite FTS5)
2. Try broader terms (e.g., "Map" instead of "Option.Map")
3. Verify the examples database is embedded (rebuild if needed)

### Verbose Logging

Enable verbose mode to see detailed execution:
```bash
go tool gen mcp --verbose
```

Logs go to stderr and include:
- Tool registration
- Tool calls with parameters
- Query execution
- Result counts

## Best Practices

### For AI Assistants

1. **Start with list_skills** — Discover available skills before loading
2. **Load skills on demand** — Use `use_skill` when needed, not preemptively
3. **Search before get** — Use `search_examples` to find relevant examples, then `get_example` for details
4. **Cache skill content** — Skills don't change during a session
5. **Use package filters** — Narrow search results with `package_filter`

### For Users

1. **Keep the server running** — Configure it in your MCP client for persistent access
2. **Use verbose mode for debugging** — Helps diagnose issues
3. **Update after changes** — Run `go get -tool github.com/IBM/fp-go/v2/gen@latest` when fp-go or skills are updated
4. **Combine tools** — Use skills for concepts, examples for code

## Security Considerations

- **Read-only operations** — All tools are marked as read-only
- **No file system access** — Uses embedded data only
- **No network access** — Operates entirely offline
- **Temporary files** — Examples database is extracted to temp, cleaned up automatically
- **Stdio transport** — No network ports or external connections

## References

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [fp-go Repository](https://github.com/IBM/fp-go)
- [fp-go Documentation](https://pkg.go.dev/github.com/IBM/fp-go/v2)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

