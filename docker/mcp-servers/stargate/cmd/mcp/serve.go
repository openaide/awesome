package mcp

import (
	"fmt"
	"log"
	"os"

	"github.com/openaide/stargate/server"
	"github.com/spf13/cobra"
)

var port int
var host string

func serve() {
	ms := server.NewMCPServer(
		"Stargate",
		"1.0.0",
		// server.WithResourceCapabilities(true, true),
		// server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	if err := ms.Initialize(); err != nil {
		log.Fatalf("Failed to initialize MCP server: %v", err)
	}

	baseURL := fmt.Sprintf("http://%s:%v", host, port)
	addr := fmt.Sprintf(":%v", port)

	server.NewSSEServer(ms, baseURL)

	sse := server.NewSSEServer(ms, baseURL)

	log.Printf("SSE server listening on :%d", port)

	if err := sse.Start(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mcp",
	Short: "A command line tool for MCP server.",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server.",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	var defaultPort = 58080
	if v := os.Getenv("AI_MCP_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &defaultPort)
	}

	// flags
	serveCmd.Flags().IntVar(&port, "port", defaultPort, "Port to run the server")
	serveCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind the server")

	rootCmd.AddCommand(serveCmd)
}

func Serve() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
