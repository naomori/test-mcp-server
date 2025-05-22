package main

import (
	"context"
	"log"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
)

// Example client code - you can run this separately to test the server
func runClient() {
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