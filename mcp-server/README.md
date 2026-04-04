# Passgen MCP Server

An MCP (Model Context Protocol) server for the passgen password generator API, written in Go.

## What is MCP?

The Model Context Protocol (MCP) is an open protocol that enables AI models to interact with external data sources and tools through a standardized interface. MCP servers expose tools and resources that AI assistants can use to extend their capabilities beyond their training data.

### Key MCP Concepts

1. **Tools**: Functions that AI models can call with specific parameters to perform actions or retrieve information.
2. **Resources**: Data sources that AI models can read from, such as files, databases, or APIs.
3. **Prompts**: Pre-defined templates for common interactions or queries.
4. **Transport**: The communication channel between the AI model and the MCP server (stdio, HTTP, SSE).

### Tool Discovery Handshake

When an MCP client (like Claude Desktop or Cursor) starts an MCP server, it follows this handshake process:

1. **Initialize**: The client sends an `initialize` request with protocol version and capabilities.
2. **Server Announcement**: The server responds with its name, version, and supported capabilities.
3. **Tool Listing**: The client requests the list of available tools via `tools/list`.
4. **Tool Schemas**: The server returns detailed schemas for each tool, including parameter definitions, types, and descriptions.
5. **Ready State**: Once handshake is complete, the AI model can call tools using the `tools/call` method.

### Why Stdio is the Default Transport

Stdio (standard input/output) is the default transport for MCP servers because:

1. **Security**: No network ports are opened, reducing attack surface.
2. **Simplicity**: No configuration needed for ports, authentication, or firewalls.
3. **Isolation**: Each server runs in its own process with clear boundaries.
4. **Portability**: Works consistently across different operating systems.
5. **Resource Management**: Processes can be easily started and stopped by the client.

## Features

- Exposes a `generate_password` tool with configurable parameters
- Connects to passgen REST API at configurable base URL
- Communicates over stdio (standard MCP transport)
- Docker container support

## Tool: `generate_password`

Generates random passwords using the passgen API.

### Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `length` | integer | Yes | - | Length of the password (8-128) |
| `uppercase` | boolean | No | true | Include uppercase letters |
| `lowercase` | boolean | No | true | Include lowercase letters |
| `digits` | boolean | No | true | Include digits |
| `symbols` | boolean | No | true | Include symbols |
| `count` | integer | No | 1 | Number of passwords to generate (1-100) |

### Example Usage

```json
{
  "length": 16,
  "uppercase": true,
  "lowercase": true,
  "digits": true,
  "symbols": false,
  "count": 3
}
```

## Installation

### Prerequisites

- Go 1.26 or later
- Docker (optional, for containerized deployment)

### Building from Source

```bash
cd mcp-server
go build -o passgen-mcp
```

### Building with Docker

```bash
cd mcp-server
docker build -t passgen-mcp .
```

## Configuration

### Environment Variables

- `PASSGEN_URL`: Base URL of the passgen API (default: `http://localhost:8080`)

### Claude Desktop Configuration

Add the following to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`  
**Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "passgen": {
      "command": "/path/to/passgen-mcp",
      "env": {
        "PASSGEN_URL": "http://localhost:8080"
      }
    }
  }
}
```

For Docker container:

```json
{
  "mcpServers": {
    "passgen": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-e",
        "PASSGEN_URL=http://host.docker.internal:8080",
        "passgen-mcp"
      ]
    }
  }
}
```

### Cursor Configuration

Add the following to your Cursor MCP settings:

**Settings > MCP Servers > Add New Server**

```json
{
  "name": "passgen",
  "command": "/path/to/passgen-mcp",
  "env": {
    "PASSGEN_URL": "http://localhost:8080"
  }
}
```

Or for Docker:

```json
{
  "name": "passgen",
  "command": "docker",
  "args": [
    "run",
    "--rm",
    "-i",
    "-e",
    "PASSGEN_URL=http://host.docker.internal:8080",
    "passgen-mcp"
  ]
}
```

## Using with Local Ollama

Ollama supports MCP servers through its experimental MCP integration. Here's how to configure the passgen MCP server with Ollama:

### Prerequisites

1. Install Ollama from [ollama.com](https://ollama.com)
2. Build or download the passgen MCP server binary
3. Ensure the passgen API is running (default: `http://localhost:8080`)

### Configuration

Create an Ollama MCP configuration file:

**macOS/Linux**: `~/.ollama/mcp/servers/passgen.json`
**Windows**: `%USERPROFILE%\.ollama\mcp\servers\passgen.json`

```json
{
  "mcpServers": {
    "passgen": {
      "command": "/absolute/path/to/passgen-mcp",
      "env": {
        "PASSGEN_URL": "http://localhost:8080"
      }
    }
  }
}
```

For Docker container:

```json
{
  "mcpServers": {
    "passgen": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-e",
        "PASSGEN_URL=http://host.docker.internal:8080",
        "passgen-mcp"
      ]
    }
  }
}
```

### Starting Ollama with MCP Support

1. **Method 1: Environment variable** (Ollama 0.5.0+):
   ```bash
   export OLLAMA_MCP_SERVERS='{"passgen":{"command":"/path/to/passgen-mcp","env":{"PASSGEN_URL":"http://localhost:8080"}}}'
   ollama run llama3.2
   ```

2. **Method 2: Configuration file** (Ollama 0.5.0+):
   Place the configuration file as shown above and Ollama will automatically load it.

3. **Method 3: Direct MCP server specification**:
   ```bash
   ollama run --mcp-server passgen:/path/to/passgen-mcp llama3.2
   ```

### Testing with Ollama

Once Ollama is running with MCP support, you can ask the LLM to generate passwords:

```
User: Generate a 16-character password with uppercase, lowercase, and digits but no symbols.

Assistant: I'll generate a password for you using the passgen tool.

<tool_call>
{
  "tool": "generate_password",
  "parameters": {
    "length": 16,
    "uppercase": true,
    "lowercase": true,
    "digits": true,
    "symbols": false,
    "count": 1
  }
}
</tool_call>
```

The LLM will receive the generated password and present it to you.

### Troubleshooting

1. **MCP server not detected**: Ensure the configuration file is in the correct location and Ollama version supports MCP (0.5.0+).
2. **Connection errors**: Verify the passgen API is running at the URL specified in `PASSGEN_URL`.
3. **Permission issues**: Make sure the passgen-mcp binary is executable (`chmod +x /path/to/passgen-mcp`).
4. **Docker networking**: Use `host.docker.internal` on macOS/Windows or `172.17.0.1` on Linux to connect to host services.

### Example Ollama Prompt

```
You are a helpful assistant with access to password generation tools.
When asked to generate passwords, use the generate_password tool with appropriate parameters.
Always inform the user about the security characteristics of the generated password.
```

## Development

### Project Structure

```
mcp-server/
├── main.go          # Main server implementation
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── Dockerfile       # Container definition
└── README.md        # This file
```

### Dependencies

- `github.com/mark3labs/mcp-go`: Official MCP SDK for Go

### Testing

1. Start the passgen API server:
   ```bash
   cd /home/dev/git/go-passgen
   go run main.go
   ```

2. Build and run the MCP server:
   ```bash
   cd mcp-server
   go build -o passgen-mcp
   ./passgen-mcp
   ```

3. Test with an MCP client or use the example test script.
