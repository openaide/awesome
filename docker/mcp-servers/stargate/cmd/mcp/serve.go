package mcp

import (
	"fmt"
	"log"
	"os"

	"github.com/openaide/stargate/proxy"
	"github.com/spf13/cobra"
)

var config proxy.ProxyConfig

func serve() {
	if err := proxy.Serve(&config); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
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
	serveCmd.Flags().IntVar(&config.Port, "port", defaultPort, "Port to run the server")
	serveCmd.Flags().StringVar(&config.Host, "host", "localhost", "Host to bind the server")
	serveCmd.Flags().StringVar(&config.Config, "config", "", "MCP configuration file")

	rootCmd.AddCommand(serveCmd)
}

func Serve() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
