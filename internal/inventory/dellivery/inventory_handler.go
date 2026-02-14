package dellivery

// ============================================================================
// Inventory HTTP Handler — Presentation Layer
// ============================================================================
//
// Logic Overview:
// Menangani request HTTP terkait produk dan stok.
// Endpoint:
//   GET /products     → Daftar semua produk (dengan pagination)
//   GET /products/:id → Detail satu produk
// ============================================================================

import (
	"net/http" // HTTP status codes
	"strconv"  // Konversi string ke integer

	"github.com/gin-gonic/gin"                                            // Gin web framework
	"github.com/purnama/Event-Driven-Logistic/internal/inventory/service" // Service layer
	"github.com/purnama/Event-Driven-Logistic/pkg/response"               // Standar response format
)

// InventoryHandler menangani semua HTTP request terkait Inventory/Product.
type InventoryHandler struct {
	svc service.InventoryService // Dependency: service layer
}

// NewInventoryHandler membuat instance baru InventoryHandler.
// Parameter: svc — InventoryService interface (dependency injection).
// Return: pointer ke InventoryHandler.
func NewInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc} // Inject service ke handler
}

// ListProducts menangani GET /products — daftar produk.
//
// Query Params (opsional):
//
//	?limit=20  — jumlah item per halaman (default: 20)
//	?offset=0  — mulai dari item ke-berapa (default: 0)
//
// Response: 200 OK dengan array produk.
func (h *InventoryHandler) ListProducts(c *gin.Context) {
	// Parse query params dengan default value
	limitStr := c.DefaultQuery("limit", "20")  // Default 20 item per halaman
	offsetStr := c.DefaultQuery("offset", "0") // Default mulai dari awal

	// Konversi string ke integer
	limit, err := strconv.Atoi(limitStr) // Parse limit
	if err != nil {
		limit = 20 // Fallback ke default jika parsing gagal
	}

	offset, err := strconv.Atoi(offsetStr) // Parse offset
	if err != nil {
		offset = 0 // Fallback ke default jika parsing gagal
	}

	// Panggil service
	products, err := h.svc.ListProducts(limit, offset) // Delegasi ke service
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil daftar produk") // 500
		return                                                                             // Stop execution
	}

	// Kirim response sukses
	response.Success(c, "Daftar produk berhasil diambil", products) // 200 OK
}

// GetProductByID menangani GET /products/:id — detail produk.
//
// URL Params:
//
//	:id — ID produk (integer)
func (h *InventoryHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id") // Ambil ":id" dari URL

	// Parse string ke uint
	id, err := strconv.ParseUint(idStr, 10, 32) // Konversi string → uint
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Format ID tidak valid") // 400
		return                                                            // Stop execution
	}

	// Panggil service
	product, err := h.svc.GetProductByID(uint(id)) // Delegasi ke service (cast ke uint)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Produk tidak ditemukan") // 404
		return                                                           // Stop execution
	}

	// Kirim response sukses
	response.Success(c, "Produk ditemukan", product) // 200 OK
}
