// +build ignore

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

	// List available tools
	tools, err := client.ListTools(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	log.Println("Available Tools:")
	for _, tool := range tools.Tools {
		desc := ""
		if tool.Description != nil {
			desc = *tool.Description
		}
		log.Printf("Tool: %s. Description: %s", tool.Name, desc)
	}

	// Call the echo tool
	echoArgs := map[string]interface{}{
		"message": "Hello from MCP client!",
	}

	echoResponse, err := client.CallTool(context.Background(), "echo", echoArgs)
	if err != nil {
		log.Fatalf("Failed to call echo tool: %v", err)
	}

	if len(echoResponse.Content) > 0 && echoResponse.Content[0].TextContent != nil {
		log.Printf("Echo response: %s", echoResponse.Content[0].TextContent.Text)
	}

	// Call the time tool
	timeArgs := map[string]interface{}{
		"format": "2006-01-02 15:04:05",
	}

	timeResponse, err := client.CallTool(context.Background(), "time", timeArgs)
	if err != nil {
		log.Fatalf("Failed to call time tool: %v", err)
	}

	if len(timeResponse.Content) > 0 && timeResponse.Content[0].TextContent != nil {
		log.Printf("Time response: %s", timeResponse.Content[0].TextContent.Text)
	}
}