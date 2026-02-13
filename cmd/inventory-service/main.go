package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/inventory/repository"
	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
)

func main() {
	log.Println("🚀 Starting Inventory Service...")

	cfg := config.LoadConfig()
	cfg.PrintConfig()

	db := database.InitPostgres(cfg.Database.URL)

	log.Println("📦 Running database migrations...")
	err := db.AutoMigrate(
		&repository.Product{},
		&repository.StockReservation{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed successfully")

	productRepo := repository.NewProductRepository(db)
	reservationRepo := repository.NewStockReservationRepository(db)
	log.Printf("✅ Product repository initialized: %T", productRepo)
	log.Printf("✅ Stock reservation repository initialized: %T", reservationRepo)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "inventory-service",
			"version": "1.0.0",
		})
	})

	log.Printf("✅ Inventory Service is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
