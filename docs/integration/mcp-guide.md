# MCP Guide — Connecting AI Agents to EduTrack

The EduTrack MCP server lets AI agents (Claude, custom agents, IDEs) read
routes, progress, coverage, gaps and resources over the Model Context
Protocol.

## Running the server

```
vedo-edutrack mcp
```

The server speaks JSON-RPC 2.0 over **stdio**: requests arrive on stdin,
responses go to stdout, and all logs go to stderr. Run it with an API key to
enable authentication:

```
VEDO_MCP_API_KEY=<secret> vedo-edutrack mcp
```

## Supported tools

| Tool            | Arguments                              | Returns                          |
|-----------------|----------------------------------------|----------------------------------|
| `get_route`     | `learner_id`, `goal_topic_id`          | Route steps + horizons           |
| `get_progress`  | `learner_id`                           | Plan-vs-actual + forecast        |
| `get_coverage`  | `learner_id`                           | FGOS coverage report             |
| `get_gaps`      | `learner_id`, `lag_module_id`          | Root-cause diagnosis             |
| `get_resources` | `module_id`                            | Bound resources                  |

## Protocol

```json
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"my-agent","version":"1.0"},"capabilities":{"apiKey":"<secret>"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_route","arguments":{"learner_id":"user-1","goal_topic_id":"math-5-11"}}}
```

## Handshake order

1. `initialize` — authenticate (API key) and receive server capabilities.
2. `tools/list` — discover tool schemas.
3. `tools/call` — invoke tools.

Without a valid API key every method except `initialize` returns error
`-32002 (Not authenticated)`.

## Claude Desktop configuration

```json
{
  "mcpServers": {
    "vedo-edutrack": {
      "command": "vedo-edutrack",
      "args": ["mcp"],
      "env": { "VEDO_MCP_API_KEY": "<secret>" }
    }
  }
}
```

## See also

- `quickstart.md` — get started in 5 minutes
- `examples/curl-examples.sh` — REST equivalents of every tool
