package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	query := strings.Join(os.Args[1:], " ")
	if query == "" {
		fmt.Println("Please provide a query")
		os.Exit(1)
	}
	fmt.Println("query: ", query)
	out := filepath.Join("./outputs", "gptr")

	ctx := context.Background()

	fmt.Println("Building gptr image...")
	err := BuildGPTRImage(ctx)
	if err != nil {
		fmt.Println("Error building image: ", err)
		os.Exit(1)
	}

	fmt.Println("Running gptr...")
	err = RunGPTRContainer(ctx, query, out)
	if err != nil {
		fmt.Println("Error running gptr: ", err)
		os.Exit(1)
	}
}
