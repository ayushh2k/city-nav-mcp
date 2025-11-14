package main

import (
	"fmt"
	"mcp-server/internal/routes"
)

func main() {
	router := routes.SetupRouter()

	fmt.Println("Starting MCP server on http://localhost:8000")
	router.Run(":8000")
}
