package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestCalculatorTool(t *testing.T) {
	// Create a new MCP server for testing
	s := server.NewMCPServer(
		"Calculator Test",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add the calculator tool
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	// Add the calculator handler with our test handler
	calculatorHandler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		x, err := request.RequireFloat("x")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultError("cannot divide by zero"), nil
			}
			result = x / y
		}

		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	}

	s.AddTool(calculatorTool, calculatorHandler)

	// Test cases
	testCases := []struct {
		name           string
		operation      string
		x              float64
		y              float64
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Addition",
			operation:      "add",
			x:              5,
			y:              3,
			expectedResult: "8.00",
			expectError:    false,
		},
		{
			name:           "Subtraction",
			operation:      "subtract",
			x:              10,
			y:              4,
			expectedResult: "6.00",
			expectError:    false,
		},
		{
			name:           "Multiplication",
			operation:      "multiply",
			x:              6,
			y:              7,
			expectedResult: "42.00",
			expectError:    false,
		},
		{
			name:           "Division",
			operation:      "divide",
			x:              20,
			y:              5,
			expectedResult: "4.00",
			expectError:    false,
		},
		{
			name:           "Division by zero",
			operation:      "divide",
			x:              10,
			y:              0,
			expectedResult: "",
			expectError:    true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock request
			request := mcp.CallToolRequest{}
			
			// Set the arguments using reflection
			args := map[string]interface{}{
				"operation": tc.operation,
				"x":         tc.x,
				"y":         tc.y,
			}
			
			// Set the arguments directly in the request's params
			request.Params.Name = "calculate"
			request.Params.Arguments = args

			// Call the handler directly
			result, err := calculatorHandler(context.Background(), request)
			assert.NoError(t, err, "Handler should not return an error")

			if tc.expectError {
				assert.True(t, result.IsError, "Expected an error result")
			} else {
				assert.False(t, result.IsError, "Did not expect an error result")
				// Check the content
				if len(result.Content) > 0 {
					textContent, ok := result.Content[0].(mcp.TextContent)
					if assert.True(t, ok, "Expected TextContent") {
						assert.Equal(t, tc.expectedResult, textContent.Text, "Result text should match expected")
					}
				} else {
					t.Error("Expected content in result")
				}
			}
		})
	}
}