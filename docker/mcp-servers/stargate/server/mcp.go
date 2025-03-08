package server

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
	"github.com/tailscale/hujson"
)

//go:embed mcp_config.jsonc
var mcpConfigData []byte

var mcpConfig = NewMcpConfig()

func init() {
	mcpConfig.Load(mcpConfigData)
}

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

func (c *McpConfig) LoadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return c.Load(data)
}

func (c *McpConfig) Load(data []byte) error {
	hu, err := hujson.Standardize(data)
	if err != nil {
		return err
	}
	ex := expandWithDefault(string(hu))
	err = json.Unmarshal([]byte(ex), &c)
	if err != nil {
		return fmt.Errorf("unmarshal mcp config: %v", err)
	}

	// set server name for each config
	for k, v := range c.ServerConfigs {
		v.Server = k
	}
	return nil
}

type McpClientSession struct {
	cfg *McpServerConfig

	client *client.StdioMCPClient
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

func (r *McpClientSession) ListTools(ctx context.Context) ([]*ToolFunc, error) {
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

func (r *McpClientSession) CallTool(ctx context.Context, tool string, args map[string]any) (string, error) {
	req := mcp.CallToolRequest{}
	req.Params.Name = tool
	req.Params.Arguments = args

	resp, err := r.client.CallTool(ctx, req)
	if err != nil {
		return "", err
	}
	for _, content := range resp.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			return textContent.Text, nil
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			return string(jsonBytes), nil
		}
	}
	return "", nil
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

func (r *McpClient) ListTools(ctx context.Context) ([]*ToolFunc, error) {
	clientSession := &McpClientSession{
		cfg: r.ServerConfig,
	}
	if err := clientSession.Connect(ctx); err != nil {
		return nil, err
	}
	defer clientSession.Close()

	return clientSession.ListTools(ctx)
}

func (r *McpClient) CallTool(ctx context.Context, tool string, args map[string]any) (string, error) {
	clientSession := &McpClientSession{
		cfg: r.ServerConfig,
	}
	if err := clientSession.Connect(ctx); err != nil {
		return "", err
	}
	defer clientSession.Close()

	return clientSession.CallTool(ctx, tool, args)
}

type McpProxy struct {
	Config *McpConfig

	cached map[string][]*ToolFunc

	sync.Mutex
}

func NewMcpProxy(cfg *McpConfig) *McpProxy {
	return &McpProxy{
		Config: cfg,
	}
}

func (r *McpProxy) ListTools() (map[string][]*ToolFunc, error) {
	r.Lock()
	defer r.Unlock()
	if r.cached != nil {
		log.Printf("Using cached tools total %v\n", len(r.cached))
		return r.cached, nil
	}

	var tools = map[string][]*ToolFunc{}
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

func (r *McpProxy) GetTools(server string) ([]*ToolFunc, error) {
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

func (r *McpProxy) CallTool(server, tool string, args map[string]any) (string, error) {
	ctx := context.Background()
	for v, cfg := range r.Config.ServerConfigs {
		if v == server {
			client := &McpClient{
				ServerConfig: cfg,
			}
			resp, err := client.CallTool(ctx, tool, args)
			if err != nil {
				return "", err
			}
			if resp != "" {
				return resp, nil
			}
		}
	}
	return "", fmt.Errorf("no such server: %s", server)
}

// Initialize the MCP server proxy
func (s *MCPServer) Initialize() error {
	// default
	config := mcpConfigData
	var cfg = NewMcpConfig()
	if err := cfg.Load(config); err != nil {

	}
	s.proxy = NewMcpProxy(cfg)
	return nil
}
