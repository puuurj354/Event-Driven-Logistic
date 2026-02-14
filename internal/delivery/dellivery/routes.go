package dellivery

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *ShipmentHandler) {
	shipments := router.Group("/shipments")
	{
		shipments.GET("/:order_id", handler.GetShipmentByOrderID)
		shipments.PATCH("/:id/status", handler.UpdateShipmentStatus)
	}
}
