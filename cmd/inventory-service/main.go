package main

import (
	"log" // Logging

	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/inventory/dellivery"
	inventoryEvent "github.com/purnama/Event-Driven-Logistic/internal/inventory/event"
	"github.com/purnama/Event-Driven-Logistic/internal/inventory/repository"
	"github.com/purnama/Event-Driven-Logistic/internal/inventory/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
	"github.com/purnama/Event-Driven-Logistic/pkg/middleware"
)

func main() {
	log.Println("🚀 Starting Inventory Service...")

	cfg := config.LoadConfig()
	cfg.PrintConfig()

	db := database.InitPostgres(cfg.Database.URL)

	log.Println("📦 Running database migrations...")
	if err := db.AutoMigrate(
		&repository.Product{},
		&repository.StockReservation{},
	); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed successfully")

	productRepo := repository.NewProductRepository(db)
	reservationRepo := repository.NewStockReservationRepository(db)
	inventorySvc := service.NewInventoryService(productRepo, reservationRepo)
	mqConn := broker.ConnectRabbitMQ()
	defer mqConn.Close()
	pubChan := broker.CreateChannel(mqConn)
	defer pubChan.Close()
	publisher, err := broker.NewPublisher(pubChan)
	if err != nil {
		log.Fatalf("❌ Failed to create publisher: %v", err)
	}
	consChan := broker.CreateChannel(mqConn)
	defer consChan.Close()
	consumer, err := broker.NewConsumer(consChan)
	if err != nil {
		log.Fatalf("❌ Failed to create consumer: %v", err)
	}
	invConsumer := inventoryEvent.NewInventoryConsumer(consumer, publisher, inventorySvc)
	if err := invConsumer.StartListening(); err != nil {
		log.Fatalf("❌ Failed to start inventory consumer: %v", err)
	}
	handler := dellivery.NewInventoryHandler(inventorySvc)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware()) // CORS untuk dashboard

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "inventory-service", "version": "1.0.0"})
	})

	dellivery.RegisterRoutes(router, handler)

	log.Printf("✅ Inventory Service is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
