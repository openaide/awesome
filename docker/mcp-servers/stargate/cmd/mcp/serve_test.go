package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	baseUrl := "http://localhost:58080/sse"
	port := os.Getenv("AI_MCP_PORT")
	if port != "" {
		baseUrl = fmt.Sprintf("http://localhost:%s/sse", port)
	}

	client, err := client.NewSSEMCPClient(baseUrl)
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
	callTests := []struct {
		name       string
		toolName   string
		arguments  map[string]interface{}
		wantResult bool
	}{
		{
			name:     "Convert time",
			toolName: "time__convert_time",
			arguments: map[string]interface{}{
				"source_timezone": "America/Los_Angeles",
				"time":            "16:30",
				"target_timezone": "Asia/Shanghai",
			},
			wantResult: true,
		},
		// {
		// 	name:     "Search on DuckDuckGo",
		// 	toolName: "ddg__search",
		// 	arguments: map[string]interface{}{
		// 		"query":       "weather in sfo ca",
		// 		"max_results": 1,
		// 	},
		// 	wantResult: true,
		// },
		// {
		// 	name:     "List Docker Containers",
		// 	toolName: "docker__list-containers",
		// 	arguments: map[string]interface{}{
		// 		"container_name": "stargate",
		// 	},
		// 	wantResult: true,
		// },
		// {
		// 	name:     "Get Docker Logs",
		// 	toolName: "docker__get-logs",
		// 	arguments: map[string]interface{}{
		// 		"container_name": "stargate",
		// 	},
		// 	wantResult: true,
		// },
	}

	for _, test := range callTests {
		t.Run(test.name, func(t *testing.T) {
			req := mcp.CallToolRequest{}
			req.Params.Name = test.toolName
			req.Params.Arguments = test.arguments

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
		})
	}
}
