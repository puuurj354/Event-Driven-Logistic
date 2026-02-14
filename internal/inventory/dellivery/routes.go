package dellivery

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *InventoryHandler) {
	products := router.Group("/products")
	{
		products.GET("", handler.ListProducts)
		products.GET("/:id", handler.GetProductByID)
		products.POST("", handler.CreateProduct)
	}
}
