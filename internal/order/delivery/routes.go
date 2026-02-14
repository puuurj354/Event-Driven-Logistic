package delivery

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *OrderHandler) {
	orders := router.Group("/orders")
	{
		orders.POST("", handler.CreateOrder)
		orders.GET("/:id", handler.GetOrderByID)
		orders.GET("/user/:user_id", handler.GetOrdersByCustomerID)
	}
}
