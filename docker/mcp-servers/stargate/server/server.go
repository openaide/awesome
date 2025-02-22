// Package server provides MCP (Model Control Protocol) server implementations.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// resourceEntry holds both a resource and its handler
type resourceEntry struct {
	// resource mcp.Resource
	// handler  ResourceHandlerFunc
}

// resourceTemplateEntry holds both a template and its handler
type resourceTemplateEntry struct {
	// template mcp.ResourceTemplate
	// handler  ResourceTemplateHandlerFunc
}

// ServerOption is a function that configures an MCPServer.
type ServerOption func(*MCPServer)

// ResourceHandlerFunc is a function that returns resource contents.
type ResourceHandlerFunc func(ctx context.Context, request mcp.ReadResourceRequest) ([]any, error)

// ResourceTemplateHandlerFunc is a function that returns a resource template.
type ResourceTemplateHandlerFunc func(ctx context.Context, request mcp.ReadResourceRequest) ([]any, error)

// PromptHandlerFunc handles prompt requests with given arguments.
type PromptHandlerFunc func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)

// ToolHandlerFunc handles tool calls with given arguments.
type ToolHandlerFunc func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)

// ServerTool combines a Tool with its ToolHandlerFunc.
type ServerTool struct {
	Tool    mcp.Tool
	Handler ToolHandlerFunc
}

// NotificationContext provides client identification for notifications
type NotificationContext struct {
	ClientID  string
	SessionID string
}

// ServerNotification combines the notification with client context
type ServerNotification struct {
	Context      NotificationContext
	Notification mcp.JSONRPCNotification
}

// NotificationHandlerFunc handles incoming notifications.
type NotificationHandlerFunc func(ctx context.Context, notification mcp.JSONRPCNotification)

// MCPServer implements a Model Control Protocol server that can handle various types of requests
// including resources, prompts, and tools.
type MCPServer struct {
	name                 string
	version              string
	resources            map[string]resourceEntry
	resourceTemplates    map[string]resourceTemplateEntry
	prompts              map[string]mcp.Prompt
	promptHandlers       map[string]PromptHandlerFunc
	tools                map[string]ServerTool
	notificationHandlers map[string]NotificationHandlerFunc
	capabilities         serverCapabilities
	notifications        chan ServerNotification
	currentClient        NotificationContext
	initialized          bool

	proxy *McpProxy
}

// serverKey is the context key for storing the server instance
type serverKey struct{}

// ServerFromContext retrieves the MCPServer instance from a context
func ServerFromContext(ctx context.Context) *MCPServer {
	if srv, ok := ctx.Value(serverKey{}).(*MCPServer); ok {
		return srv
	}
	return nil
}

// WithContext sets the current client context and returns the provided context
func (s *MCPServer) WithContext(
	ctx context.Context,
	notifCtx NotificationContext,
) context.Context {
	s.currentClient = notifCtx
	return ctx
}

// SendNotificationToClient sends a notification to the current client
func (s *MCPServer) SendNotificationToClient(
	method string,
	params map[string]any,
) error {
	if s.notifications == nil {
		return fmt.Errorf("notification channel not initialized")
	}

	notification := mcp.JSONRPCNotification{
		JSONRPC: mcp.JSONRPC_VERSION,
		Notification: mcp.Notification{
			Method: method,
			Params: mcp.NotificationParams{
				AdditionalFields: params,
			},
		},
	}

	select {
	case s.notifications <- ServerNotification{
		Context:      s.currentClient,
		Notification: notification,
	}:
		return nil
	default:
		return fmt.Errorf("notification channel full or blocked")
	}
}

// serverCapabilities defines the supported features of the MCP server
type serverCapabilities struct {
	resources *resourceCapabilities
	prompts   *promptCapabilities
	logging   bool
}

// resourceCapabilities defines the supported resource-related features
type resourceCapabilities struct {
	subscribe   bool
	listChanged bool
}

// promptCapabilities defines the supported prompt-related features
type promptCapabilities struct {
	listChanged bool
}

// WithResourceCapabilities configures resource-related server capabilities
func WithResourceCapabilities(subscribe, listChanged bool) ServerOption {
	return func(s *MCPServer) {
		s.capabilities.resources = &resourceCapabilities{
			subscribe:   subscribe,
			listChanged: listChanged,
		}
	}
}

// WithPromptCapabilities configures prompt-related server capabilities
func WithPromptCapabilities(listChanged bool) ServerOption {
	return func(s *MCPServer) {
		s.capabilities.prompts = &promptCapabilities{
			listChanged: listChanged,
		}
	}
}

// WithLogging enables logging capabilities for the server
func WithLogging() ServerOption {
	return func(s *MCPServer) {
		s.capabilities.logging = true
	}
}

// NewMCPServer creates a new MCP server instance with the given name, version and options
func NewMCPServer(
	name, version string,
	opts ...ServerOption,
) *MCPServer {
	s := &MCPServer{
		resources:            make(map[string]resourceEntry),
		resourceTemplates:    make(map[string]resourceTemplateEntry),
		prompts:              make(map[string]mcp.Prompt),
		promptHandlers:       make(map[string]PromptHandlerFunc),
		tools:                make(map[string]ServerTool),
		name:                 name,
		version:              version,
		notificationHandlers: make(map[string]NotificationHandlerFunc),
		notifications:        make(chan ServerNotification, 100),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// HandleMessage processes an incoming JSON-RPC message and returns an appropriate response
func (s *MCPServer) HandleMessage(
	ctx context.Context,
	message json.RawMessage,
) mcp.JSONRPCMessage {
	// Add server to context
	ctx = context.WithValue(ctx, serverKey{}, s)

	var baseMessage struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		ID      any    `json:"id,omitempty"`
	}

	if err := json.Unmarshal(message, &baseMessage); err != nil {
		return createErrorResponse(
			nil,
			mcp.PARSE_ERROR,
			"Failed to parse message",
		)
	}

	// Check for valid JSONRPC version
	if baseMessage.JSONRPC != mcp.JSONRPC_VERSION {
		return createErrorResponse(
			baseMessage.ID,
			mcp.INVALID_REQUEST,
			"Invalid JSON-RPC version",
		)
	}

	if baseMessage.ID == nil {
		var notification mcp.JSONRPCNotification
		if err := json.Unmarshal(message, &notification); err != nil {
			return createErrorResponse(
				nil,
				mcp.PARSE_ERROR,
				"Failed to parse notification",
			)
		}
		s.handleNotification(ctx, notification)
		return nil // Return nil for notifications
	}

	switch baseMessage.Method {
	case "initialize":
		var request mcp.InitializeRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid initialize request",
			)
		}
		return s.handleInitialize(ctx, baseMessage.ID, request)
	case "ping":
		var request mcp.PingRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid ping request",
			)
		}
		return s.handlePing(ctx, baseMessage.ID, request)
	case "resources/list":
		if s.capabilities.resources == nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.METHOD_NOT_FOUND,
				"Resources not supported",
			)
		}
		var request mcp.ListResourcesRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid list resources request",
			)
		}
		return s.handleListResources(ctx, baseMessage.ID, request)
	case "resources/templates/list":
		if s.capabilities.resources == nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.METHOD_NOT_FOUND,
				"Resources not supported",
			)
		}
		var request mcp.ListResourceTemplatesRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid list resource templates request",
			)
		}
		return s.handleListResourceTemplates(ctx, baseMessage.ID, request)
	case "resources/read":
		if s.capabilities.resources == nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.METHOD_NOT_FOUND,
				"Resources not supported",
			)
		}
		var request mcp.ReadResourceRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid read resource request",
			)
		}
		return s.handleReadResource(ctx, baseMessage.ID, request)
	case "prompts/list":
		if s.capabilities.prompts == nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.METHOD_NOT_FOUND,
				"Prompts not supported",
			)
		}
		var request mcp.ListPromptsRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid list prompts request",
			)
		}
		return s.handleListPrompts(ctx, baseMessage.ID, request)
	case "prompts/get":
		if s.capabilities.prompts == nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.METHOD_NOT_FOUND,
				"Prompts not supported",
			)
		}
		var request mcp.GetPromptRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid get prompt request",
			)
		}
		return s.handleGetPrompt(ctx, baseMessage.ID, request)
	case "tools/list":
		// if len(s.tools) == 0 {
		// 	return createErrorResponse(
		// 		baseMessage.ID,
		// 		mcp.METHOD_NOT_FOUND,
		// 		"Tools not supported",
		// 	)
		// }
		var request mcp.ListToolsRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid list tools request",
			)
		}
		return s.handleListTools(ctx, baseMessage.ID, request)
	case "tools/call":
		// if len(s.tools) == 0 {
		// 	return createErrorResponse(
		// 		baseMessage.ID,
		// 		mcp.METHOD_NOT_FOUND,
		// 		"Tools not supported",
		// 	)
		// }
		var request mcp.CallToolRequest
		if err := json.Unmarshal(message, &request); err != nil {
			return createErrorResponse(
				baseMessage.ID,
				mcp.INVALID_REQUEST,
				"Invalid call tool request",
			)
		}
		return s.handleToolCall(ctx, baseMessage.ID, request)
	default:
		return createErrorResponse(
			baseMessage.ID,
			mcp.METHOD_NOT_FOUND,
			fmt.Sprintf("Method %s not found", baseMessage.Method),
		)
	}
}

func (s *MCPServer) handleInitialize(
	_ context.Context,
	id any,
	request mcp.InitializeRequest,
) mcp.JSONRPCMessage {
	log.Printf("Initializing server: %s %s", s.name, s.version)

	capabilities := mcp.ServerCapabilities{}
	capabilities.Tools = &struct {
		ListChanged bool `json:"listChanged,omitempty"`
	}{
		ListChanged: true,
	}

	result := mcp.InitializeResult{
		ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
		ServerInfo: mcp.Implementation{
			Name:    s.name,
			Version: s.version,
		},
		Capabilities: capabilities,
	}

	s.initialized = true
	return createResponse(id, result)
}

func (s *MCPServer) handlePing(
	_ context.Context,
	id any,
	request mcp.PingRequest,
) mcp.JSONRPCMessage {
	return createResponse(id, mcp.EmptyResult{})
}

func (s *MCPServer) handleListResources(
	_ context.Context,
	id any,
	request mcp.ListResourcesRequest,
) mcp.JSONRPCMessage {
	resources := make([]mcp.Resource, 0, len(s.resources))
	result := mcp.ListResourcesResult{
		Resources: resources,
	}
	if request.Params.Cursor != "" {
		result.NextCursor = "" // Handle pagination if needed
	}
	return createResponse(id, result)
}

func (s *MCPServer) handleListResourceTemplates(
	_ context.Context,
	id any,
	request mcp.ListResourceTemplatesRequest,
) mcp.JSONRPCMessage {
	templates := make([]mcp.ResourceTemplate, 0, len(s.resourceTemplates))
	result := mcp.ListResourceTemplatesResult{
		ResourceTemplates: templates,
	}
	if request.Params.Cursor != "" {
		result.NextCursor = "" // Handle pagination if needed
	}
	return createResponse(id, result)
}

func (s *MCPServer) handleReadResource(
	_ context.Context,
	id any,
	request mcp.ReadResourceRequest,
) mcp.JSONRPCMessage {
	return createErrorResponse(
		id,
		mcp.INVALID_PARAMS,
		fmt.Sprintf(
			"No handler found for resource URI: %s",
			request.Params.URI,
		),
	)
}

func (s *MCPServer) handleListPrompts(
	_ context.Context,
	id any,
	request mcp.ListPromptsRequest,
) mcp.JSONRPCMessage {
	prompts := make([]mcp.Prompt, 0, len(s.prompts))
	result := mcp.ListPromptsResult{
		Prompts: prompts,
	}
	if request.Params.Cursor != "" {
		result.NextCursor = "" // Handle pagination if needed
	}
	return createResponse(id, result)
}

func (s *MCPServer) handleGetPrompt(
	_ context.Context,
	id any,
	request mcp.GetPromptRequest,
) mcp.JSONRPCMessage {
	return createErrorResponse(
		id,
		mcp.INVALID_PARAMS,
		fmt.Sprintf("Prompt not found: %s", request.Params.Name),
	)
}

func (s *MCPServer) handleListTools(
	_ context.Context,
	id any,
	request mcp.ListToolsRequest,
) mcp.JSONRPCMessage {
	log.Println("Listing tools...")

	tools := make([]mcp.Tool, 0, len(s.tools))
	funcMap, err := s.proxy.ListTools()
	if err != nil {
		return createErrorResponse(id, mcp.INTERNAL_ERROR, err.Error())
	}
	for serverName, funcs := range funcMap {
		for _, v := range funcs {
			name := fmt.Sprintf("%s__%s", serverName, v.Name)
			params := v.Parameters
			tool := mcp.Tool{
				Name:        name,
				Description: v.Description,
				InputSchema: mcp.ToolInputSchema{
					Type:       params["type"].(string),
					Properties: params["properties"].(map[string]any),
					Required:   params["required"].([]string),
				},
			}
			tools = append(tools, tool)
		}
	}

	result := mcp.ListToolsResult{
		Tools: tools,
	}
	if request.Params.Cursor != "" {
		result.NextCursor = "" // Handle pagination if needed
	}
	return createResponse(id, result)
}

func (s *MCPServer) handleToolCall(
	ctx context.Context,
	id any,
	request mcp.CallToolRequest,
) mcp.JSONRPCMessage {
	log.Printf("Calling tool: %s", request.Params.Name)

	parts := strings.SplitN(request.Params.Name, "__", 2)
	if len(parts) != 2 {
		return createErrorResponse(
			id,
			mcp.INVALID_PARAMS,
			fmt.Sprintf("Invalid tool name format: %s", request.Params.Name),
		)
	}

	resp, err := s.proxy.CallTool(parts[0], parts[1], request.Params.Arguments)
	if err != nil {
		return createErrorResponse(id, mcp.INTERNAL_ERROR, err.Error())
	}

	result := &mcp.CallToolResult{
		Content: []any{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Echo: %s", resp),
			},
		},
	}

	return createResponse(id, result)
}

func (s *MCPServer) handleNotification(
	_ context.Context,
	notification mcp.JSONRPCNotification,
) mcp.JSONRPCMessage {
	return nil
}

func createResponse(id any, result any) mcp.JSONRPCMessage {
	return mcp.JSONRPCResponse{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      id,
		Result:  result,
	}
}

func createErrorResponse(
	id any,
	code int,
	message string,
) mcp.JSONRPCMessage {
	return mcp.JSONRPCError{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      id,
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    any    `json:"data,omitempty"`
		}{
			Code:    code,
			Message: message,
		},
	}
}
