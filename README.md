# Simple MCP Server

A minimal implementation of a Model Context Protocol (MCP) server using the [metoro-io/mcp-golang](https://github.com/metoro-io/mcp-golang) library with HTTP transport.

## Usage

To run the server:

```bash
go run .
```

The server will start on port 8080 and listen for requests at the `/mcp` endpoint.

## Available Tools

This MCP server provides two simple tools:

1. **time** - Returns the current time in the specified format
   - Parameters: `format` (optional) - The time format to use (defaults to RFC3339)

2. **echo** - Echoes back the provided message
   - Parameters: `message` (required) - The message to echo back

## Example Usage

### Using the MCP client library

```go
package main

import (
	"context"
	"log"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
)

func main() {
	// Create an HTTP transport that connects to the server
	transport := http.NewHTTPClientTransport("/mcp")
	transport.WithBaseURL("http://localhost:8080")

	// Create a new client with the transport
	client := mcp.NewClient(transport)

	// Initialize the client
	if _, err := client.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// Call the echo tool
	args := map[string]interface{}{
		"message": "Hello from MCP client!",
	}

	response, err := client.CallTool(context.Background(), "echo", args)
	if err != nil {
		log.Fatalf("Failed to call echo tool: %v", err)
	}

	if len(response.Content) > 0 && response.Content[0].TextContent != nil {
		log.Printf("Echo response: %s", response.Content[0].TextContent.Text)
	}
}
```
