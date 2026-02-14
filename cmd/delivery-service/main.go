package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/delivery/dellivery"
	deliveryEvent "github.com/purnama/Event-Driven-Logistic/internal/delivery/event"
	"github.com/purnama/Event-Driven-Logistic/internal/delivery/repository"
	"github.com/purnama/Event-Driven-Logistic/internal/delivery/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
	"github.com/purnama/Event-Driven-Logistic/pkg/middleware"
)

func main() {
	log.Println("🚀 Starting Delivery Service...")

	cfg := config.LoadConfig()
	cfg.PrintConfig()

	db := database.InitPostgres(cfg.Database.URL)
	log.Println("📦 Running database migrations...")
	if err := db.AutoMigrate(&repository.Shipment{}); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed successfully")

	shipmentRepo := repository.NewShipmentRepository(db)
	shipmentSvc := service.NewShipmentService(shipmentRepo)

	mqConn := broker.ConnectRabbitMQ()
	defer mqConn.Close()
	consChan := broker.CreateChannel(mqConn)
	defer consChan.Close()
	consumer, err := broker.NewConsumer(consChan)
	if err != nil {
		log.Fatalf("❌ Failed to create consumer: %v", err)
	}
	delConsumer := deliveryEvent.NewDeliveryConsumer(consumer, shipmentSvc)
	if err := delConsumer.StartListening(); err != nil {
		log.Fatalf("❌ Failed to start delivery consumer: %v", err)
	}

	handler := dellivery.NewShipmentHandler(shipmentSvc)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware()) // CORS untuk dashboard

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "delivery-service", "version": "1.0.0"})
	})

	dellivery.RegisterRoutes(router, handler)

	log.Printf("✅ Delivery Service is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
