package dellivery

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *PaymentHandler) {
	payments := router.Group("/payments")
	{
		payments.POST("", handler.ConfirmPayment)
		payments.GET("/:order_id", handler.GetPaymentByOrderID)
	}
}
