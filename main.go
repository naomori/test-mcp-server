package main

import (
	"log"
	"time"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
)

// TimeArgs defines the arguments for the time tool
type TimeArgs struct {
	Format string `json:"format" jsonschema:"description=The time format to use"`
}

// EchoArgs defines the arguments for the echo tool
type EchoArgs struct {
	Message string `json:"message" jsonschema:"required,description=The message to echo back"`
}

func main() {
	// Create an HTTP transport that listens on /mcp endpoint
	transport := http.NewHTTPTransport("/mcp").WithAddr(":8080")

	// Create a new server with the transport
	server := mcp.NewServer(
		transport,
		mcp.WithName("simple-mcp-server"),
		mcp.WithInstructions("A simple MCP server implementation"),
		mcp.WithVersion("1.0.0"),
	)

	// Register a time tool
	err := server.RegisterTool("time", "Returns the current time in the specified format", func(args TimeArgs) (*mcp.ToolResponse, error) {
		format := args.Format
		if format == "" {
			format = time.RFC3339
		}
		log.Printf("Time tool called with format: %s", format)
		return mcp.NewToolResponse(mcp.NewTextContent(time.Now().Format(format))), nil
	})
	if err != nil {
		log.Fatalf("Failed to register time tool: %v", err)
	}

	// Register an echo tool
	err = server.RegisterTool("echo", "Echoes back the provided message", func(args EchoArgs) (*mcp.ToolResponse, error) {
		log.Printf("Echo tool called with message: %s", args.Message)
		return mcp.NewToolResponse(mcp.NewTextContent(args.Message)), nil
	})
	if err != nil {
		log.Fatalf("Failed to register echo tool: %v", err)
	}

	// Start the server
	log.Println("Starting HTTP server on :8080...")
	if err := server.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}