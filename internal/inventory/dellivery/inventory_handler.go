package dellivery



import (
	"net/http" 
	"strconv"  

	"github.com/gin-gonic/gin"                                            
	"github.com/purnama/Event-Driven-Logistic/internal/inventory/service" 
	"github.com/purnama/Event-Driven-Logistic/pkg/response"               
)


type InventoryHandler struct {
	svc service.InventoryService 
}


func NewInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc} 
}


func (h *InventoryHandler) ListProducts(c *gin.Context) {
	
	limitStr := c.DefaultQuery("limit", "20")  
	offsetStr := c.DefaultQuery("offset", "0") 

	
	limit, err := strconv.Atoi(limitStr) 
	if err != nil {
		limit = 20 
	}

	offset, err := strconv.Atoi(offsetStr) 
	if err != nil {
		offset = 0 
	}


	products, err := h.svc.ListProducts(limit, offset) 
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil daftar produk") 
		return                                                                             
	}

	
	response.Success(c, "Daftar produk berhasil diambil", products) 
}


func (h *InventoryHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id") 

	
	id, err := strconv.ParseUint(idStr, 10, 32) 
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Format ID tidak valid") 
		return                                                            
	}

	
	product, err := h.svc.GetProductByID(uint(id)) 
	if err != nil {
		response.Error(c, http.StatusNotFound, "Produk tidak ditemukan") 
		return                                                           
	}

	
	response.Success(c, "Produk ditemukan", product) 
}

func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	
	var req struct {
		Name  string `json:"name" binding:"required"`  
		Stock int    `json:"stock" binding:"required"` 
	}

	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Format request tidak valid: "+err.Error()) 
		return                                                                               
	}

	
	product, err := h.svc.CreateProduct(req.Name, req.Stock) 
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal membuat produk: "+err.Error()) 
		return                                                                                  
	}

	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,                     
		"message": "Produk berhasil dibuat", 
		"data":    product,                  
	})
}
