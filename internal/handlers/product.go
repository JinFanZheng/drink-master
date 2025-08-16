package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

// ProductHandler 产品处理器 (对应VendingMachine.MobileAPI ProductController)
type ProductHandler struct {
	*BaseHandler
}

// NewProductHandler 创建产品处理器
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

// GetSelectList 获取产品选择列表
// @Summary 获取产品选择列表
// @Description 获取所有可选择的产品列表，包含价格信息
// @Tags Product
// @Accept json
// @Produce json
// @Success 200 {object} contracts.APIResponse{data=[]contracts.SelectViewModel}
// @Failure 500 {object} contracts.APIResponse
// @Router /products/select [get]
func (h *ProductHandler) GetSelectList(c *gin.Context) {
	// Step 1: Get all products directly from products table using correct column mapping
	var products []models.Product
	err := h.db.Find(&products).Error
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// Step 2: Get machine product prices
	var machineProductPrices []models.MachineProductPrice
	err = h.db.Find(&machineProductPrices).Error
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// Step 3: Create product mapping with prices from machine_product_prices
	productMap := make(map[string]struct {
		ID    string  `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Image *string `json:"image,omitempty"`
	})

	// Create product info from products table
	for _, product := range products {
		productMap[product.ID] = struct {
			ID    string  `json:"id"`
			Name  string  `json:"name"`
			Price float64 `json:"price"`
			Image *string `json:"image,omitempty"`
		}{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price, // Use price from products table
			Image: product.Image,
		}
	}

	// Step 4: Override with machine-specific pricing if available
	for _, mpp := range machineProductPrices {
		if productInfo, exists := productMap[mpp.ProductId]; exists {
			// Update with machine-specific price
			productInfo.Price = mpp.Price
			productMap[mpp.ProductId] = productInfo
		}
	}

	// Convert map to slice
	var result []interface{}
	for _, product := range productMap {
		result = append(result, product)
	}

	h.SuccessResponse(c, result)
}
