package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/delivery/repository"
	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
)

func main() {
	log.Println("🚀 Starting Delivery Service...")

	cfg := config.LoadConfig()
	cfg.PrintConfig()

	db := database.InitPostgres(cfg.Database.URL)

	log.Println("📦 Running database migrations...")
	err := db.AutoMigrate(&repository.Shipment{})
	if err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed successfully")

	shipmentRepo := repository.NewShipmentRepository(db)
	log.Printf("✅ Shipment repository initialized: %T", shipmentRepo)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "delivery-service",
			"version": "1.0.0",
		})
	})

	log.Printf("✅ Delivery Service is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
