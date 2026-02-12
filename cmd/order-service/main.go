package main

import (
	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/order/delivery"
	"github.com/purnama/Event-Driven-Logistic/internal/order/repository"
	"github.com/purnama/Event-Driven-Logistic/pkg/config"
	"github.com/purnama/Event-Driven-Logistic/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitPostgres(cfg.Database.URL)

	db.AutoMigrate(&repository.Order{})

	repo := repository.NewOrderRepository(db)
	handler := delivery.NewOrderHandler(repo)

	r := gin.Default()
	r.POST("/orders", handler.CreateOrder)

	r.Run(":" + cfg.Server.Port)
}
