package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	baseUrl := "http://localhost:58080"
	client, err := client.NewSSEMCPClient(baseUrl + "/sse")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start the client
	if err := client.Start(ctx); err != nil {
		t.Fatalf("Failed to start client: %v", err)
	}

	// Initialize
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}

	initResult, err := client.Initialize(ctx, initRequest)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	t.Logf(
		"Initialized with server: %s %s\n",
		initResult.ServerInfo.Name,
		initResult.ServerInfo.Version,
	)

	// Test Ping
	if err := client.Ping(ctx); err != nil {
		t.Errorf("Ping failed: %v", err)
	}
	t.Log("Ping successful")

	// List Tools
	fmt.Println("Listing available tools...")
	listReq := mcp.ListToolsRequest{}
	tools, err := client.ListTools(ctx, listReq)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range tools.Tools {
		log.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// Call a tool
	req := mcp.CallToolRequest{}
	req.Params.Name = "time__convert_time"
	req.Params.Arguments = map[string]interface{}{
		"source_timezone": "America/Los_Angeles",
		"time":            "16:30",
		"target_timezone": "Asia/Shanghai",
	}
	result, err := client.CallTool(ctx, req)
	if err != nil {
		t.Fatalf("Failed to call: %+v %v", req, err)
	}

	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}

}
