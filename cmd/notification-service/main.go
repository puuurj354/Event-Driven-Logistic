package main

// ============================================================================
// Notification Service — Entry Point (main.go)
// ============================================================================
//
// Logic Overview:
// Notification Service adalah "jembatan" antara RabbitMQ dan browser:
//   1. RabbitMQ Consumer → listen ALL events (order.created, payment.success, etc.)
//   2. Service → persist event ke DB + broadcast via WebSocket Hub
//   3. WebSocket Hub → push real-time updates ke connected browsers
//   4. HTTP Server → serve HTMX dashboard + WebSocket endpoint
//
// Wire chain:
// Config → DB → Repo → Hub → Service → Consumer → HTTP Routes
// ============================================================================

import (
	"html/template" // Go HTML templates
	"log"           // Logging
	"net/http"      // HTTP handlers
	"path/filepath" // File path handling

	"github.com/gin-gonic/gin"                                                         // Gin web framework
	"github.com/purnama/Event-Driven-Logistic/internal/notification/event"             // RabbitMQ consumer
	"github.com/purnama/Event-Driven-Logistic/internal/notification/repository"        // DB models + repo
	"github.com/purnama/Event-Driven-Logistic/internal/notification/service"           // Business logic
	notifWs "github.com/purnama/Event-Driven-Logistic/internal/notification/websocket" // WebSocket hub+handler
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"                              // RabbitMQ helpers
	"github.com/purnama/Event-Driven-Logistic/pkg/config"                              // Config loader
	"github.com/purnama/Event-Driven-Logistic/pkg/database"                            // Database helper
)

func main() {
	log.Println("🚀 Starting Notification Service...") // Log startup

	// ── Step 1: Load konfigurasi dari .env ──
	cfg := config.LoadConfig() // Baca DB_URL, MQ_URL, PORT
	cfg.PrintConfig()          // Tampilkan (password di-mask)

	// ── Step 2: Koneksi ke PostgreSQL ──
	db := database.InitPostgres(cfg.Database.URL) // Koneksi ke db_notification

	// ── Step 3: AutoMigrate ──
	log.Println("📦 Running database migrations...")
	if err := db.AutoMigrate(&repository.NotificationLog{}); err != nil { // Migrate model
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Database migrations completed successfully")

	// ── Step 4: Wire dependencies ──
	repo := repository.NewNotificationRepository(db) // Layer 1: Database

	// ── Step 5: Start WebSocket Hub ──
	hub := notifWs.NewHub() // Buat Hub instance
	go hub.Run()            // Jalankan Hub sebagai goroutine permanen

	// ── Step 6: Wire Service (inject repo + hub) ──
	notifSvc := service.NewNotificationService(repo, hub) // Layer 2: Business logic

	// ── Step 7: Setup RabbitMQ Consumer ──
	mqConn := broker.ConnectRabbitMQ()            // Koneksi ke RabbitMQ
	defer mqConn.Close()                          // Tutup saat shutdown
	consChan := broker.CreateChannel(mqConn)      // Buat channel untuk consume
	defer consChan.Close()                        // Tutup saat shutdown
	consumer, err := broker.NewConsumer(consChan) // Buat consumer
	if err != nil {
		log.Fatalf("❌ Failed to create consumer: %v", err)
	}
	notifConsumer := event.NewNotificationConsumer(consumer, notifSvc) // Wire consumer
	if err := notifConsumer.StartListening(); err != nil {             // Mulai listen ALL events
		log.Fatalf("❌ Failed to start notification consumer: %v", err)
	}

	// ── Step 8: Setup Gin router ──
	router := gin.Default()

	// Load HTML templates from project root /templates/
	templatePath := filepath.Join("..", "..", "templates", "*.html") // Relative dari cmd/notification-service/
	tmpl := template.Must(template.ParseGlob(templatePath))          // Parse main templates
	router.SetHTMLTemplate(tmpl)                                     // Set ke Gin

	// Serve static assets
	router.Static("/static", filepath.Join("..", "..", "templates", "static")) // CSS, JS, images

	// ── Health check ──
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":     "healthy",
			"service":    "notification-service",
			"version":    "1.0.0",
			"ws_clients": hub.ClientCount(), // Jumlah WS client yang terkoneksi
		})
	})

	// ── WebSocket endpoint ──
	router.GET("/ws", notifWs.ServeWs(hub)) // Upgrade HTTP → WebSocket

	// ── Dashboard (HTMX) ──
	router.GET("/", func(c *gin.Context) {
		// Ambil 50 event terbaru untuk initial render
		logs, err := notifSvc.GetRecentLogs(50) // Query DB
		if err != nil {
			log.Printf("⚠️ Gagal ambil logs: %v", err) // Non-fatal, render kosong
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Logs": logs, // Data untuk template
		})
	})

	// ── API: Event logs per order ──
	router.GET("/api/logs/:order_id", func(c *gin.Context) {
		orderID := c.Param("order_id")                  // Ambil order_id dari URL
		logs, err := notifSvc.GetLogsByOrderID(orderID) // Query DB
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()}) // 500 Error
			return
		}
		c.JSON(200, gin.H{"logs": logs}) // 200 OK
	})

	// ── Step 9: Start HTTP server ──
	log.Printf("✅ Notification Service is running on port %s", cfg.Server.Port)
	log.Printf("🌐 Dashboard: http://localhost:%s", cfg.Server.Port)
	log.Printf("🔌 WebSocket: ws://localhost:%s/ws", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
