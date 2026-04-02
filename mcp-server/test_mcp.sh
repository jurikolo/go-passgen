#!/bin/bash

# Test the MCP server by sending initialization and tool listing requests
# This simulates what an MCP client would do

echo "Testing MCP server..."

# Create a simple test that pipes JSON to the server and reads response
cat <<EOF | ./passgen-mcp 2>&1 | head -50
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "test-client",
      "version": "1.0.0"
    }
  }
}
EOF

echo "Test completed."