package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/openaide/stargate/server"
)

func main() {
	var port int
	var host string
	flag.IntVar(&port, "port", 58080, "Port to listen on")
	flag.StringVar(&host, "host", "localhost", "Host")

	flag.Parse()

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
