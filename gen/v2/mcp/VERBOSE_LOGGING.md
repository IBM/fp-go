# MCP Server Verbose Logging

## Overview

The MCP server now supports verbose logging to help with debugging and monitoring server operations. When enabled, detailed logs are written to stderr, allowing you to trace the flow of requests and tool executions.

## Usage

To enable verbose logging, use the `--verbose` or `-v` flag when starting the MCP server:

```bash
./gen mcp --verbose
```

or

```bash
./gen mcp -v
```

## What Gets Logged

When verbose mode is enabled, the following information is logged to stderr:

### Server Startup
- Verbose logging enabled confirmation
- Server creation with tool registration
- Individual tool registration (list_skills, use_skill, search_examples, get_example)
- Server startup with stdio transport
- Total number of tools registered

### Tool Execution
Each tool execution logs:
- Tool name being called
- Tool arguments/parameters
- Execution progress
- Results (number of items found/retrieved)
- Any errors encountered

### Server Shutdown
- Server stop status (successful or with error)

## Log Format

Logs are prefixed with `[MCP]` and include timestamps, microseconds, and file locations:

```
2026/06/16 13:42:29.123456 mcp/server.go:145: [MCP] Verbose logging enabled
2026/06/16 13:42:29.123789 mcp/server.go:98: [MCP] Creating MCP server with fp-go tools
2026/06/16 13:42:29.124012 mcp/server.go:107: [MCP] Registering tool: list_skills
```

## Example Output

### Server Startup
```bash
$ ./gen mcp --verbose
2026/06/16 13:42:29.123456 mcp/server.go:145: [MCP] Verbose logging enabled
2026/06/16 13:42:29.123789 mcp/server.go:98: [MCP] Creating MCP server with fp-go tools
2026/06/16 13:42:29.124012 mcp/server.go:107: [MCP] Registering tool: list_skills
2026/06/16 13:42:29.124234 mcp/server.go:116: [MCP] Registering tool: use_skill
2026/06/16 13:42:29.124456 mcp/server.go:125: [MCP] Registering tool: search_examples
2026/06/16 13:42:29.124678 mcp/server.go:134: [MCP] Registering tool: get_example
2026/06/16 13:42:29.124890 mcp/server.go:143: [MCP] Server created with 4 tools registered
2026/06/16 13:42:29.125012 mcp/server.go:152: [MCP] Starting MCP server with stdio transport...
```

### Tool Execution - list_skills
```
2026/06/16 13:42:30.567890 mcp/server.go:195: [MCP] Executing tool: list_skills
2026/06/16 13:42:30.568123 mcp/server.go:238: [MCP] list_skills completed: found 5 skills
```

### Tool Execution - use_skill
```
2026/06/16 13:42:31.234567 mcp/server.go:281: [MCP] Executing tool: use_skill (name=fp-go-logging)
2026/06/16 13:42:31.234890 mcp/server.go:319: [MCP] use_skill completed: retrieved skill 'fp-go-logging'
```

### Tool Execution - search_examples
```
2026/06/16 13:42:32.345678 mcp/server.go:330: [MCP] Executing tool: search_examples (query=Map, package_filter=option)
2026/06/16 13:42:32.456789 mcp/server.go:444: [MCP] search_examples completed: found 8 examples
```

### Tool Execution - get_example
```
2026/06/16 13:42:33.567890 mcp/server.go:456: [MCP] Executing tool: get_example (symbol=Map)
2026/06/16 13:42:33.678901 mcp/server.go:565: [MCP] get_example completed: retrieved 3 examples
```

### Error Handling
```
2026/06/16 13:42:34.789012 mcp/server.go:281: [MCP] Executing tool: use_skill (name=invalid-skill)
2026/06/16 13:42:34.789234 mcp/server.go:297: [MCP] use_skill error: skill 'invalid-skill' not found
```

## Available Tools

The MCP server provides the following tools:

1. **list_skills** - List all available fp-go skills with descriptions
2. **use_skill** - Retrieve full content of a specific skill
3. **search_examples** - Full-text search across Go examples
4. **get_example** - Retrieve specific examples by symbol name

## Benefits

- **Debugging**: Trace the flow of tool calls and identify issues
- **Monitoring**: Track tool execution and performance
- **Troubleshooting**: Identify errors and understand their context
- **Development**: Understand the MCP protocol flow during development
- **Auditing**: Keep track of which tools are being used and when

## Performance Impact

Verbose logging adds minimal overhead:
- Log statements only execute when verbose mode is enabled
- Logs are written to stderr, not blocking stdout (which is used for MCP protocol)
- No performance impact when verbose mode is disabled (default)

## Integration with MCP Clients

The verbose logging is designed to work seamlessly with MCP clients:
- All protocol communication happens on stdout (unaffected by logging)
- All verbose logs go to stderr (separate stream)
- Clients can capture stderr separately for debugging purposes
- The MCP protocol remains clean and unaffected

## Tips

1. **Redirect stderr to a file** for later analysis:
   ```bash
   ./gen mcp --verbose 2> mcp-debug.log
   ```

2. **Filter logs** using grep:
   ```bash
   ./gen mcp --verbose 2>&1 | grep "Executing tool"
   ```

3. **Combine with MCP client logging** for full visibility into the communication flow

4. **Use during development** to understand tool behavior and debug issues

5. **Disable in production** unless troubleshooting specific issues