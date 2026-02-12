package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/order/repository"
)

type OrderHandler struct {
	repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order repository.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat order"})
		return
	}
	c.JSON(http.StatusCreated, order)
}
