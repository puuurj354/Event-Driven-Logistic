package main


import (
	"log" 

	"github.com/gin-gonic/gin"                                           
	"github.com/purnama/Event-Driven-Logistic/internal/order/delivery"   
	"github.com/purnama/Event-Driven-Logistic/internal/order/event"      
	"github.com/purnama/Event-Driven-Logistic/internal/order/repository" 
	"github.com/purnama/Event-Driven-Logistic/internal/order/service"    
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"                
	"github.com/purnama/Event-Driven-Logistic/pkg/config"                
	"github.com/purnama/Event-Driven-Logistic/pkg/database"              
)

func main() {
	log.Println("🚀 Starting Order Service...")

	cfg := config.LoadConfig()
	cfg.PrintConfig()

	db := database.InitPostgres(cfg.Database.URL)


	log.Println("📦 Running database migrations...")
	if err := db.AutoMigrate(&repository.Order{}); err != nil { 
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed successfully")

	repo := repository.NewOrderRepository(db) 
	svc := service.NewOrderService(repo)      

	mqConn := broker.ConnectRabbitMQ()             
	defer mqConn.Close()                           
	pubChan := broker.CreateChannel(mqConn)        
	defer pubChan.Close()                          
	publisher, err := broker.NewPublisher(pubChan) 
	if err != nil {
		log.Fatalf("❌ Failed to create publisher: %v", err)
	}
	orderPublisher := event.NewOrderPublisher(publisher) 

	consChan := broker.CreateChannel(mqConn)      
	defer consChan.Close()                        
	consumer, err := broker.NewConsumer(consChan) 
	if err != nil {
		log.Fatalf("❌ Failed to create consumer: %v", err)
	}
	orderConsumer := event.NewOrderConsumer(consumer, repo) 
	if err := orderConsumer.StartListening(); err != nil {  
		log.Fatalf("❌ Failed to start order consumer: %v", err)
	}

	handler := delivery.NewOrderHandler(svc, orderPublisher) 

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "order-service", "version": "1.0.0"})
	})

	delivery.RegisterRoutes(router, handler) 


	log.Printf("✅ Order Service is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
