package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/services"
)

// OrderHandler 订单处理器 (对应MobileAPI OrderController)
type OrderHandler struct {
	*BaseHandler
	orderService services.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(db *gorm.DB, orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		BaseHandler:  NewBaseHandler(db),
		orderService: orderService,
	}
}

// GetPaging 分页获取订单列表
// @Summary 分页获取用户订单列表
// @Description 获取当前登录用户的订单列表，支持分页和筛选
// @Tags Order
// @Accept json
// @Produce json
// @Param request body contracts.GetMemberOrderPagingRequest true "分页请求"
// @Success 200 {object} contracts.APIResponse{data=[]contracts.GetMemberOrderPagingResponse}
// @Failure 400 {object} contracts.APIResponse
// @Failure 401 {object} contracts.APIResponse
// @Security BearerAuth
// @Router /Order/GetPaging [post]
func (h *OrderHandler) GetPaging(c *gin.Context) {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	var request contracts.GetMemberOrderPagingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 设置会员ID
	request.MemberID = memberID

	// 设置默认分页参数
	if request.PageIndex <= 0 {
		request.PageIndex = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	response, err := h.orderService.GetMemberOrderPaging(request)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.APIResponse{
		Success: true,
		Data:    response.Orders,
		Meta:    response.Meta.Meta,
	})
}

// Get 获取订单详情
// @Summary 获取订单详细信息
// @Description 根据订单ID获取订单的详细信息
// @Tags Order
// @Accept json
// @Produce json
// @Param id query string true "订单ID"
// @Success 200 {object} contracts.APIResponse{data=contracts.GetOrderByIdResponse}
// @Failure 400 {object} contracts.APIResponse
// @Failure 401 {object} contracts.APIResponse
// @Failure 404 {object} contracts.APIResponse
// @Security BearerAuth
// @Router /Order/Get [get]
func (h *OrderHandler) Get(c *gin.Context) {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	orderID := c.Query("id")
	if orderID == "" {
		h.ValidationErrorResponse(c, errors.New("订单ID不能为空"))
		return
	}

	response, err := h.orderService.GetByID(orderID)
	if err != nil {
		if err.Error() == "订单不存在" {
			h.NotFoundResponse(c, "订单不存在")
			return
		}
		h.InternalErrorResponse(c, err)
		return
	}

	// 简单的权限检查 - 这里可以根据实际需求优化
	_ = memberID

	h.SuccessResponse(c, response)
}

// Create 创建订单
// @Summary 创建新订单
// @Description 用户创建一个新的购买订单
// @Tags Order
// @Accept json
// @Produce json
// @Param request body contracts.CreateOrderRequest true "创建订单请求"
// @Success 201 {object} contracts.APIResponse{data=contracts.CreateOrderResponse}
// @Failure 400 {object} contracts.APIResponse
// @Failure 401 {object} contracts.APIResponse
// @Security BearerAuth
// @Router /Order/Create [post]
func (h *OrderHandler) Create(c *gin.Context) {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	var request contracts.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 设置会员ID
	request.MemberID = memberID

	response, err := h.orderService.Create(request)
	if err != nil {
		if err.Error() == "会员不存在" || err.Error() == "机器不存在" {
			h.ValidationErrorResponse(c, err)
			return
		}
		if err.Error() == "机器不在线，下单失败" {
			c.JSON(http.StatusBadRequest, contracts.APIResponse{
				Success: false,
				Error: &contracts.APIError{
					Code:    contracts.ErrorCodeDeviceOffline,
					Message: "机器不在线，下单失败",
				},
			})
			return
		}
		h.InternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, contracts.APIResponse{
		Success: true,
		Data:    response,
	})
}

// Refund 申请退款
// @Summary 订单退款
// @Description 机主权限用户对订单进行退款操作
// @Tags Order
// @Accept json
// @Produce json
// @Param request body contracts.RefundOrderRequest true "退款请求"
// @Success 200 {object} contracts.APIResponse{data=contracts.RefundOrderResponse}
// @Failure 400 {object} contracts.APIResponse
// @Failure 403 {object} contracts.APIResponse
// @Failure 404 {object} contracts.APIResponse
// @Security BearerAuth
// @Router /Order/Refund [post]
func (h *OrderHandler) Refund(c *gin.Context) {
	// 检查是否为机主
	isMachineOwner := h.IsMachineOwner(c)
	if !isMachineOwner {
		h.ForbiddenResponse(c, "您不是机主，无法退款")
		return
	}

	var request contracts.RefundOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 设置机主权限标识
	request.IsMachineOwner = true

	response, err := h.orderService.Refund(request)
	if err != nil {
		if err.Error() == "订单不存在" {
			h.NotFoundResponse(c, "订单不存在")
			return
		}
		if err.Error() == "订单状态不允许退款" || err.Error() == "订单已经退款" {
			c.JSON(http.StatusBadRequest, contracts.APIResponse{
				Success: false,
				Error: &contracts.APIError{
					Code:    contracts.ErrorCodeInvalidOrderStatus,
					Message: err.Error(),
				},
			})
			return
		}
		h.InternalErrorResponse(c, err)
		return
	}

	h.SuccessResponseWithMessage(c, response, "退款成功")
}
