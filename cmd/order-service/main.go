package main

import (
	"fmt"
	"log"

	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
)

func main() {
	fmt.Println("🚀 Starting Order Service...")
	fmt.Println("============================")

	// Load configuration dari .env di direktori ini (cmd/order-service/.env)
	cfg := config.LoadConfig()
	cfg.PrintConfig()

	// Connect to database
	db, err := database.ConnectDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB(db)

	// Test ping database
	if err := database.PingDB(db); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("============================")
	fmt.Printf("🌐 Order Service is running on port %s\n", cfg.Server.Port)
	fmt.Println("Press Ctrl+C to stop")

	// TODO: Initialize Gin router, handlers, and start HTTP server
	// TODO: Initialize RabbitMQ connection and event publishers

	// Keep running
	select {}
}
