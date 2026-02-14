package dellivery

// ============================================================================
// Inventory Routes — URL Mapping
// ============================================================================
//
// Routes yang didaftarkan:
//   GET /products     → ListProducts (daftar + pagination)
//   GET /products/:id → GetProductByID (detail)
// ============================================================================

import "github.com/gin-gonic/gin" // Gin web framework

// RegisterRoutes mendaftarkan semua route Inventory ke Gin router.
// Parameter: router — Gin Engine, handler — InventoryHandler yang sudah di-wire.
func RegisterRoutes(router *gin.Engine, handler *InventoryHandler) {
	products := router.Group("/products") // Group semua route di bawah /products
	{
		products.GET("", handler.ListProducts)       // GET /products → daftar produk
		products.GET("/:id", handler.GetProductByID) // GET /products/:id → detail produk
	}
}
