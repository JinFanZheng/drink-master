package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/services"
)

// MachineOwnerHandler 机主控制器 (对应VendingMachine.MobileAPI MachineOwnerController)
type MachineOwnerHandler struct {
	*BaseHandler
	machineOwnerService *services.MachineOwnerService
}

// NewMachineOwnerHandler 创建机主控制器
func NewMachineOwnerHandler(db *gorm.DB) *MachineOwnerHandler {
	return &MachineOwnerHandler{
		BaseHandler:         NewBaseHandler(db),
		machineOwnerService: services.NewMachineOwnerService(db),
	}
}

// GetSales 获取销售情况
// @Summary 获取机主销售情况
// @Description 获取指定日期的机主销售数据，默认为当天
// @Tags MachineOwner
// @Accept json
// @Produce json
// @Param dateTime query string false "查询日期 (YYYY-MM-DD格式)"
// @Success 200 {object} contracts.MachineOwnerSalesResponse
// @Failure 401 {object} contracts.APIResponse
// @Failure 403 {object} contracts.APIResponse
// @Failure 500 {object} contracts.APIResponse
// @Router /MachineOwner/GetSales [get]
// @Security Bearer
// 对应原方法: Task<List<ColumnModel>> GetSalesAsync([FromQuery] DateTime? dateTime)
func (h *MachineOwnerHandler) GetSales(c *gin.Context) {
	// 验证是否为机主
	if !h.IsMachineOwner(c) {
		h.ForbiddenResponse(c, "您不是机主，无法查看机器列表")
		return
	}

	// 获取机主ID
	machineOwnerID, exists := h.GetMachineOwnerID(c)
	if !exists || machineOwnerID == "" {
		h.UnauthorizedResponse(c, "无效的机主信息")
		return
	}

	// 解析请求参数
	var req contracts.GetSalesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 确定查询日期
	var targetDate time.Time
	if req.DateTime == nil {
		// 默认为今天
		targetDate = time.Now().Truncate(24 * time.Hour)
	} else {
		targetDate = req.DateTime.Truncate(24 * time.Hour)
	}

	// 获取销售数据
	sales, err := h.machineOwnerService.GetSales(machineOwnerID, targetDate)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// 返回响应 (基于VendingMachine API格式)
	response := contracts.MachineOwnerSalesResponse{
		Success: true,
		Data:    sales,
		Meta: &contracts.Meta{
			Timestamp: time.Now(),
			RequestID: getRequestID(c),
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetSalesStats 获取销售统计 (扩展功能，可选实现)
// @Summary 获取机主销售统计
// @Description 获取指定日期范围内的机主销售统计数据
// @Tags MachineOwner
// @Accept json
// @Produce json
// @Param startDate query string false "开始日期 (YYYY-MM-DD格式)，默认7天前"
// @Param endDate query string false "结束日期 (YYYY-MM-DD格式)，默认今天"
// @Success 200 {object} contracts.APIResponse
// @Failure 400 {object} contracts.APIResponse
// @Failure 401 {object} contracts.APIResponse
// @Failure 403 {object} contracts.APIResponse
// @Failure 500 {object} contracts.APIResponse
// @Router /MachineOwner/GetSalesStats [get]
// @Security Bearer
func (h *MachineOwnerHandler) GetSalesStats(c *gin.Context) {
	// 验证是否为机主
	if !h.IsMachineOwner(c) {
		h.ForbiddenResponse(c, "您不是机主，无法查看统计数据")
		return
	}

	// 获取机主ID
	machineOwnerID, exists := h.GetMachineOwnerID(c)
	if !exists || machineOwnerID == "" {
		h.UnauthorizedResponse(c, "无效的机主信息")
		return
	}

	// 解析日期参数
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, "开始日期格式错误，请使用YYYY-MM-DD格式")
			return
		}
	} else {
		startDate = time.Now().AddDate(0, 0, -7) // 默认7天前
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, "结束日期格式错误，请使用YYYY-MM-DD格式")
			return
		}
	} else {
		endDate = time.Now() // 默认今天
	}

	// 验证日期范围
	if endDate.Before(startDate) {
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, "结束日期不能早于开始日期")
		return
	}

	// 获取统计数据
	stats, err := h.machineOwnerService.GetSalesStats(machineOwnerID, startDate, endDate)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	h.SuccessResponse(c, stats)
}
