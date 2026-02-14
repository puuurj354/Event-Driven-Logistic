package delivery

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/purnama/Event-Driven-Logistic/internal/order/event"
	"github.com/purnama/Event-Driven-Logistic/internal/order/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
	"github.com/purnama/Event-Driven-Logistic/pkg/response"
)

type OrderHandler struct {
	svc       service.OrderService
	publisher *event.OrderPublisher
}

func NewOrderHandler(svc service.OrderService, publisher *event.OrderPublisher) *OrderHandler {
	return &OrderHandler{
		svc:       svc,
		publisher: publisher,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	order, err := h.svc.CreateOrder(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal membuat order: "+err.Error())
		return
	}

	if h.publisher != nil {
		go func() {
			payload := broker.OrderCreatedPayload{
				OrderID:    order.ID.String(), // UUID → string
				CustomerID: order.CustomerID,  // ID customer
				ItemName:   order.ItemName,    // Nama item
				Quantity:   order.Quantity,    // Jumlah
				TotalPrice: order.TotalPrice,  // Total harga
			}
			if err := h.publisher.PublishOrderCreated(payload); err != nil {
				log.Printf("⚠️ Failed to publish order.created event: %v", err)
			}
		}()
	}

	// Kirim response sukses
	c.JSON(http.StatusCreated, response.Response{
		Status:  "success",
		Message: "Order berhasil dibuat",
		Data:    order,
	})
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Format ID tidak valid, gunakan UUID")
		return
	}

	order, err := h.svc.GetOrderByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Order tidak ditemukan")
		return
	}

	response.Success(c, "Order ditemukan", order)
}

func (h *OrderHandler) GetOrdersByCustomerID(c *gin.Context) {
	customerID := c.Param("user_id")

	if customerID == "" {
		response.Error(c, http.StatusBadRequest, "user_id tidak boleh kosong")
		return
	}

	orders, err := h.svc.GetOrdersByCustomerID(customerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil riwayat order")
		return
	}

	response.Success(c, "Riwayat order ditemukan", orders)
}
