package dellivery

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/payment/event"
	"github.com/purnama/Event-Driven-Logistic/internal/payment/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/broker"
	"github.com/purnama/Event-Driven-Logistic/pkg/response"
)

type PaymentHandler struct {
	svc       service.PaymentService
	publisher *event.PaymentPublisher
}

func NewPaymentHandler(svc service.PaymentService, publisher *event.PaymentPublisher) *PaymentHandler {
	return &PaymentHandler{
		svc:       svc,
		publisher: publisher,
	}
}

type ConfirmPaymentRequest struct {
	PaymentID uint `json:"payment_id" binding:"required"`
}

func (h *PaymentHandler) ConfirmPayment(c *gin.Context) {
	var req ConfirmPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	payment, err := h.svc.ConfirmPayment(req.PaymentID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal konfirmasi pembayaran: "+err.Error())
		return
	}

	if h.publisher != nil {
		go func() {
			payload := broker.PaymentSuccessPayload{
				OrderID:   payment.OrderID.String(),
				PaymentID: payment.ID,
				Amount:    payment.Amount,
			}
			if err := h.publisher.PublishPaymentSuccess(payload); err != nil {
				log.Printf("⚠️ Failed to publish payment.success event: %v", err)
			}
		}()
	}

	response.Success(c, "Pembayaran berhasil dikonfirmasi", payment)
}

func (h *PaymentHandler) GetPaymentByOrderID(c *gin.Context) {
	orderID := c.Param("order_id")

	if orderID == "" {
		response.Error(c, http.StatusBadRequest, "order_id tidak boleh kosong")
		return
	}
	payment, err := h.svc.GetPaymentByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Pembayaran tidak ditemukan")
		return
	}

	response.Success(c, "Pembayaran ditemukan", payment)
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Format ID tidak valid")
		return
	}

	_ = id
	response.Error(c, http.StatusNotImplemented, "Endpoint belum diimplementasi")
}
