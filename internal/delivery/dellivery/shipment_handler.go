package dellivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/purnama/Event-Driven-Logistic/internal/delivery/repository"
	"github.com/purnama/Event-Driven-Logistic/internal/delivery/service"
	"github.com/purnama/Event-Driven-Logistic/pkg/response"
)

type ShipmentHandler struct {
	svc service.ShipmentService
}

func NewShipmentHandler(svc service.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{svc: svc}
}

func (h *ShipmentHandler) GetShipmentByOrderID(c *gin.Context) {
	orderID := c.Param("order_id")

	if orderID == "" {
		response.Error(c, http.StatusBadRequest, "order_id tidak boleh kosong")
		return
	}

	shipment, err := h.svc.GetShipmentByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Shipment tidak ditemukan")
		return
	}

	response.Success(c, "Shipment ditemukan", shipment)
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *ShipmentHandler) UpdateShipmentStatus(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Format ID tidak valid")
		return
	}

	var req UpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	status := repository.ShipmentStatus(req.Status)

	if err := h.svc.UpdateShipmentStatus(uint(id), status); err != nil {
		response.Error(c, http.StatusBadRequest, "Gagal update status: "+err.Error())
		return
	}

	response.Success(c, "Status shipment berhasil diperbarui", gin.H{
		"id":     id,
		"status": req.Status,
	})
}
