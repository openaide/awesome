package proxy

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tailscale/hujson"
)

//go:embed mcp_config.jsonc
var mcpConfigData []byte

type McpConfig struct {
	ServerConfigs map[string]*McpServerConfig `json:"mcpServers"`
}

type McpServerConfig struct {
	Server  string         `json:"-"`
	Command string         `json:"command"`
	Args    []string       `json:"args"`
	Env     map[string]any `json:"env"`
}

type ToolFunc struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

func NewMcpConfig() *McpConfig {
	return &McpConfig{
		ServerConfigs: make(map[string]*McpServerConfig),
	}
}

func (r *McpConfig) LoadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return r.Load(data)
}

func (r *McpConfig) Load(data []byte) error {
	hu, err := hujson.Standardize(data)
	if err != nil {
		return err
	}
	ex := expandWithDefault(string(hu))
	err = json.Unmarshal([]byte(ex), r)
	if err != nil {
		return fmt.Errorf("unmarshal mcp config: %v", err)
	}

	// set server name for each config
	for k, v := range r.ServerConfigs {
		v.Server = k
	}
	return nil
}

type McpClientSession struct {
	cfg *McpServerConfig

	client *client.Client
}

func (r *McpClientSession) Connect(ctx context.Context) error {
	envMap := make(map[string]string)
	// default /config/<server>.env
	root := os.Getenv("MCP_ROOT")
	b, err := os.ReadFile(filepath.Join(root, fmt.Sprintf("/config/%s.env", r.cfg.Server)))
	if err == nil {
		ex := expandWithDefault(string(b))
		for _, v := range strings.Split(ex, "\n") {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			// skip comment
			if strings.HasPrefix(v, "#") {
				continue
			}
			// key=value
			parts := strings.SplitN(v, "=", 2)
			if len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}
	}
	// custom
	for k, v := range r.cfg.Env {
		if v != nil {
			envMap[k] = fmt.Sprintf("%v", v)
		}
	}
	env := make([]string, 0)
	for k, v := range envMap {
		if v != "" {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	log.Printf("Connecting to %s: %s %v\n", r.cfg.Server, r.cfg.Command, r.cfg.Args)

	client, err := client.NewStdioMCPClient(
		r.cfg.Command,
		env,
		r.cfg.Args...,
	)
	if err != nil {
		return err
	}

	//
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "swarm-client",
		Version: "1.0.0",
	}

	result, err := client.Initialize(ctx, initRequest)
	if err != nil {
		return err
	}

	log.Printf("Initialized: %s %s\n", result.ServerInfo.Name, result.ServerInfo.Version)

	r.client = client
	return nil
}

func (r *McpClientSession) ListTools(ctx context.Context) (*mcp.ListToolsResult, error) {
	toolsRequest := mcp.ListToolsRequest{}
	return r.client.ListTools(ctx, toolsRequest)
}

func (r *McpClientSession) GetTools(ctx context.Context, server string) ([]*ToolFunc, error) {
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := r.client.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, err
	}

	funcs := make([]*ToolFunc, 0)
	for _, v := range tools.Tools {
		funcs = append(funcs, &ToolFunc{
			Name:        v.Name,
			Description: v.Description,
			Parameters: map[string]any{
				"type":       v.InputSchema.Type,
				"properties": v.InputSchema.Properties,
				"required":   v.InputSchema.Required,
			},
		})
	}
	return funcs, nil
}

func (r *McpClientSession) CallTool(ctx context.Context, tool string, args map[string]any) (*mcp.CallToolResult, error) {
	req := mcp.CallToolRequest{}
	req.Params.Name = tool
	req.Params.Arguments = args

	return r.client.CallTool(ctx, req)
}

func (r *McpClientSession) Close() error {
	// TODO remove docker instance
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}

type McpClient struct {
	ServerConfig *McpServerConfig
}

func (r *McpClient) ListTools(ctx context.Context) (*mcp.ListToolsResult, error) {
	clientSession := &McpClientSession{
		cfg: r.ServerConfig,
	}
	if err := clientSession.Connect(ctx); err != nil {
		return nil, err
	}
	defer clientSession.Close()

	return clientSession.ListTools(ctx)
}

func (r *McpClient) CallTool(ctx context.Context, tool string, args map[string]any) (*mcp.CallToolResult, error) {
	clientSession := &McpClientSession{
		cfg: r.ServerConfig,
	}
	if err := clientSession.Connect(ctx); err != nil {
		return nil, err
	}
	defer clientSession.Close()

	return clientSession.CallTool(ctx, tool, args)
}

type McpProxy struct {
	Config *McpConfig

	cached map[string]*mcp.ListToolsResult

	sync.Mutex
}

func NewMcpProxy(cfg *McpConfig) *McpProxy {
	return &McpProxy{
		Config: cfg,
		cached: nil,
	}
}

func (r *McpProxy) ListTools() (map[string]*mcp.ListToolsResult, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.cached) != 0 {
		log.Printf("Using cached tools total %v\n", len(r.cached))
		return r.cached, nil
	}

	var tools = make(map[string]*mcp.ListToolsResult)
	ctx := context.Background()
	for v, cfg := range r.Config.ServerConfigs {
		client := &McpClient{
			ServerConfig: cfg,
		}
		funcs, err := client.ListTools(ctx)
		if err != nil {
			return nil, err
		}
		tools[v] = funcs
	}
	r.cached = tools

	return tools, nil
}

func (r *McpProxy) GetTools(server string) (*mcp.ListToolsResult, error) {
	ctx := context.Background()
	for v, cfg := range r.Config.ServerConfigs {
		if v == server {
			client := &McpClient{
				ServerConfig: cfg,
			}
			return client.ListTools(ctx)
		}
	}
	return nil, fmt.Errorf("no such server: %s", server)
}

func (r *McpProxy) CallTool(ctx context.Context, server, tool string, args map[string]any) (*mcp.CallToolResult, error) {
	for v, cfg := range r.Config.ServerConfigs {
		if v == server {
			client := &McpClient{
				ServerConfig: cfg,
			}
			return client.CallTool(ctx, tool, args)
		}
	}
	return nil, fmt.Errorf("no such server: %s", server)
}

type ProxyConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`

	Config string `json:"config"`
}

func createProxy(cfg *ProxyConfig) (*McpProxy, error) {
	var c = NewMcpConfig()

	var err error
	if cfg.Config != "" {
		err = c.LoadFile(cfg.Config)
	} else {
		err = c.Load(mcpConfigData)
	}
	if err != nil {
		return nil, fmt.Errorf("load mcp config: %v", err)
	}
	return NewMcpProxy(c), nil
}

func Serve(cfg *ProxyConfig) error {
	ms := server.NewMCPServer(
		"StarGate",
		"1.0.0",
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	proxy, err := createProxy(cfg)
	if err != nil {
		return err
	}

	result, err := proxy.ListTools()
	if err != nil {
		return err
	}

	newId := func(server, tool string) string {
		return fmt.Sprintf("%s__%s", server, tool)
	}

	for server, v := range result {
		log.Printf("Server: %s, Tools: %d\n", server, len(v.Tools))

		for _, tool := range v.Tools {
			id := newId(server, tool.Name)
			log.Printf("Adding tool: %s,  Name %s, Description: %s\n", id, tool.Name, tool.Description)

			ms.AddTool(mcp.Tool{
				Name:        id,
				Description: tool.Description,
				InputSchema: tool.InputSchema,
			}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				log.Printf("Calling tool %s request: %q %+v\n", id, server, request)
				return proxy.CallTool(ctx, server, tool.Name, request.Params.Arguments)
			})
		}
	}

	// start the server
	baseURL := fmt.Sprintf("http://%s:%v/", cfg.Host, cfg.Port)
	addr := fmt.Sprintf(":%v", cfg.Port)

	sse := server.NewSSEServer(ms, server.WithBaseURL(baseURL))

	log.Printf("SSE server listening: %s", addr)

	if err := sse.Start(addr); err != nil {
		return err
	}
	return nil
}
