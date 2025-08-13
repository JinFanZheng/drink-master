package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/services"
)

// MaterialSiloHandler 物料槽处理器 (对应VendingMachine.MobileAPI MaterialSiloController)
type MaterialSiloHandler struct {
	*BaseHandler
	materialSiloService services.MaterialSiloServiceInterface
}

// NewMaterialSiloHandler 创建物料槽处理器
func NewMaterialSiloHandler(db *gorm.DB) *MaterialSiloHandler {
	return &MaterialSiloHandler{
		BaseHandler:         NewBaseHandler(db),
		materialSiloService: services.NewMaterialSiloService(db),
	}
}

// GetPaging 获取物料槽分页列表
// POST /api/MaterialSilo/GetPaging
func (h *MaterialSiloHandler) GetPaging(c *gin.Context) {
	var req contracts.GetMaterialSiloPagingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.materialSiloService.GetPaging(req)
	if err != nil {
		if err.Error() == "机器不存在" {
			h.NotFoundResponse(c, "machine not found")
			return
		}
		h.InternalErrorResponse(c, err)
		return
	}

	h.SuccessResponse(c, result)
}

// UpdateStock 更新料仓库存
// POST /api/MaterialSilo/UpdateStock
func (h *MaterialSiloHandler) UpdateStock(c *gin.Context) {
	var req contracts.UpdateMaterialSiloStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.materialSiloService.UpdateStock(req)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	if !result.Success {
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, result.Message)
		return
	}

	h.SuccessResponseWithMessage(c, result.Data, result.Message)
}

// UpdateProduct 更新料仓产品
// POST /api/MaterialSilo/UpdateProduct
func (h *MaterialSiloHandler) UpdateProduct(c *gin.Context) {
	var req contracts.UpdateMaterialSiloProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.materialSiloService.UpdateProduct(req)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	if !result.Success {
		if result.Message == "物料槽不存在" {
			h.NotFoundResponse(c, "material silo not found")
			return
		}
		if result.Message == "产品不存在" {
			h.NotFoundResponse(c, "product not found")
			return
		}
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, result.Message)
		return
	}

	h.SuccessResponseWithMessage(c, result.Data, result.Message)
}

// ToggleSaleStatus 切换销售状态
// POST /api/MaterialSilo/ToggleSaleStatus
func (h *MaterialSiloHandler) ToggleSaleStatus(c *gin.Context) {
	var req contracts.ToggleSaleMaterialSiloRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.materialSiloService.ToggleSaleStatus(req)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	if !result.Success {
		if result.Message == "物料槽不存在" {
			h.NotFoundResponse(c, "material silo not found")
			return
		}
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, result.Message)
		return
	}

	h.SuccessResponseWithMessage(c, result.Data, result.Message)
}
