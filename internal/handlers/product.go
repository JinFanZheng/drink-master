package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
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
// GET /api/products/select
// 对应原方法: Task<List<SelectViewModel>> GetSelectListAsync()
func (h *ProductHandler) GetSelectList(c *gin.Context) {
	// 查询产品及其机器价格信息
	var machineProductPrices []models.MachineProductPrice

	// 预加载产品信息，获取所有可用产品的价格
	err := h.db.Preload("Product").
		Where("deleted_at IS NULL").
		Find(&machineProductPrices).Error

	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// 转换为SelectViewModel格式
	var selectViewModels []contracts.SelectViewModel
	productMap := make(map[string]*contracts.SelectViewModel)

	for _, mpp := range machineProductPrices {
		if mpp.Product != nil {
			// 使用产品ID作为key，避免重复产品
			if _, exists := productMap[mpp.Product.ID]; !exists {
				productMap[mpp.Product.ID] = &contracts.SelectViewModel{
					ID:    mpp.Product.ID,
					Name:  mpp.Product.Name,
					Price: mpp.Price, // 使用第一个找到的价格
				}
			}
		}
	}

	// 将map转换为slice
	for _, vm := range productMap {
		selectViewModels = append(selectViewModels, *vm)
	}

	// 返回响应
	response := contracts.ProductSelectListResponse{
		Products: selectViewModels,
	}

	h.SuccessResponse(c, response.Products)
}
