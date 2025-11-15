package main

import (
	"fmt"
	"log"
	"mcp-server/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, loading from system env")
	}
	router := routes.SetupRouter()

	fmt.Println("Starting MCP server on http://localhost:8000")
	router.Run(":8000")
}
