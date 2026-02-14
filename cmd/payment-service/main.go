package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/payment/dellivery"
	"github.com/purnama/Event-Driven-Logistic/internal/payment/event"
	"github.com/purnama/Event-Driven-Logistic/internal/payment/repository"
	"github.com/purnama/Event-Driven-Logistic/internal/payment/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
	"github.com/purnama/Event-Driven-Logistic/pkg/middleware"
)

func main() {
	log.Println("🚀 Starting Payment Service...")

	cfg := config.LoadConfig()
	cfg.PrintConfig()

	db := database.InitPostgres(cfg.Database.URL)
	log.Println("📦 Running database migrations...")
	if err := db.AutoMigrate(&repository.Payment{}); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed successfully")

	paymentRepo := repository.NewPaymentRepository(db)
	paymentSvc := service.NewPaymentService(paymentRepo)

	mqConn := broker.ConnectRabbitMQ()
	defer mqConn.Close()
	pubChan := broker.CreateChannel(mqConn)
	defer pubChan.Close()
	publisher, err := broker.NewPublisher(pubChan)
	if err != nil {
		log.Fatalf("❌ Failed to create publisher: %v", err)
	}
	paymentPublisher := event.NewPaymentPublisher(publisher)
	consChan := broker.CreateChannel(mqConn)
	defer consChan.Close()
	consumer, err := broker.NewConsumer(consChan)
	if err != nil {
		log.Fatalf("❌ Failed to create consumer: %v", err)
	}
	paymentConsumer := event.NewPaymentConsumer(consumer, paymentSvc)
	if err := paymentConsumer.StartListening(); err != nil {
		log.Fatalf("❌ Failed to start payment consumer: %v", err)
	}

	handler := dellivery.NewPaymentHandler(paymentSvc, paymentPublisher)
	router := gin.Default()
	router.Use(middleware.CORSMiddleware()) // CORS untuk dashboard

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "payment-service", "version": "1.0.0"})
	})

	dellivery.RegisterRoutes(router, handler)

	log.Printf("✅ Payment Service is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
