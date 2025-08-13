package handlers

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/repositories"
	"github.com/ddteam/drink-master/internal/services"
)

// PaymentHandler 支付处理器 (对应MobileAPI PaymentController)
type PaymentHandler struct {
	*BaseHandler
	paymentService services.PaymentServiceInterface
	orderService   services.OrderService
	machineService services.MachineServiceInterface
}

// NewPaymentHandler 创建支付处理器
func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{
		BaseHandler:    NewBaseHandler(db),
		paymentService: services.NewPaymentService(db),
		orderService: services.NewOrderService(
			repositories.NewOrderRepository(db),
			repositories.NewMachineRepository(db),
			repositories.NewMemberRepository(db),
			services.NewDeviceService(),
		),
		machineService: services.NewMachineService(db),
	}
}

// Get 获取支付信息
// GET /api/Payment/Get?orderId=xxx
func (h *PaymentHandler) Get(c *gin.Context) {
	// 验证请求参数
	var req contracts.GetPaymentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 获取订单信息
	order, err := h.orderService.GetByID(req.OrderID)
	if err != nil {
		h.NotFoundResponse(c, "订单不存在")
		return
	}

	// 检查订单状态
	if order.PaymentStatus != contracts.PaymentStatusWaitPay {
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodePaymentOrderAlreadyPaid, "订单已支付")
		return
	}

	// 检查订单金额是否为0（免支付）
	if order.PayAmount.IsZero() {
		h.SuccessResponse(c, contracts.GetPaymentResponse{
			Code:    200,
			Message: "免支付",
		})
		return
	}

	// 获取机器收款账户
	account, err := h.paymentService.GetPaymentAccount(order.MachineID)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// 构建微信支付请求
	notifyUrl := os.Getenv("PAYMENT_NOTIFY_URL")
	if notifyUrl == "" {
		notifyUrl = "http://vm-mobile-app/api/Callback/PaymentResult"
	}

	// 获取用户OpenId（这里需要从JWT token或其他方式获取）
	openId := h.getMemberOpenId(c)
	if openId == "" {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	wechatPayReq := contracts.WeChatPayRequest{
		Ext1:        account.ReceivingAccount,
		Ext2:        account.ReceivingKey,
		Ext3:        account.ReceivingOrderPrefix,
		NotifyUrl:   notifyUrl,
		ChannelCode: contracts.ChannelCodeFuiouMerchant,
		OrderNo:     order.OrderNo,
		OpenId:      openId,
		Attach:      "",
		OrderInfo:   fmt.Sprintf("%s(%s)", order.ProductName, h.getHasCupText(order.HasCup)),
		TransAmt:    safeInt64ToInt32(order.PayAmount.Mul(decimal.NewFromInt(100)).IntPart()), // 元转分
	}

	// 调用微信支付
	authInfo, err := h.paymentService.WeChatPay(wechatPayReq)
	if err != nil || !authInfo.IsSuccess {
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodePaymentFailed, "支付失败")
		return
	}

	h.SuccessResponse(c, contracts.GetPaymentResponse{
		Code:    200,
		Message: "获取支付信息成功",
		Data:    authInfo,
	})
}

// Query 查询支付状态
// GET /api/Payment/Query?orderId=xxx
func (h *PaymentHandler) Query(c *gin.Context) {
	// 验证请求参数
	var req contracts.QueryPaymentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 获取订单信息
	order, err := h.orderService.GetByID(req.OrderID)
	if err != nil {
		h.NotFoundResponse(c, "订单不存在")
		return
	}

	// 检查订单状态并处理
	h.handleQueryByOrderStatus(c, order, req.OrderID)
}

// handleQueryByOrderStatus 根据订单状态处理查询请求
func (h *PaymentHandler) handleQueryByOrderStatus(
	c *gin.Context, order *contracts.GetOrderByIdResponse, orderID string,
) {
	// 如果订单不是待支付状态，直接返回结果
	if order.PaymentStatus != contracts.PaymentStatusWaitPay {
		message := h.getOrderStatusMessage(order.PaymentStatus)
		h.SuccessResponse(c, contracts.QueryPaymentResponse{Message: message})
		return
	}

	// 处理免支付订单
	if order.PayAmount.IsZero() {
		h.handleFreePaymentOrder(c, orderID)
		return
	}

	// 查询第三方支付状态
	h.queryThirdPartyPaymentStatus(c, order, orderID)
}

// getOrderStatusMessage 获取订单状态对应的消息
func (h *PaymentHandler) getOrderStatusMessage(status string) string {
	switch status {
	case contracts.PaymentStatusRefunded:
		return "已退款"
	case contracts.PaymentStatusCancelled:
		return "已取消"
	default:
		return "支付成功"
	}
}

// handleFreePaymentOrder 处理免支付订单
func (h *PaymentHandler) handleFreePaymentOrder(c *gin.Context, orderID string) {
	payReq := contracts.PayOrderRequest{
		ID:             orderID,
		ChannelOrderNo: contracts.FreePaymentChannelOrderNo,
		PaidAt:         time.Now(),
	}
	if payErr := h.paymentService.PayOrder(payReq); payErr != nil {
		h.InternalErrorResponse(c, payErr)
		return
	}
	h.SuccessResponse(c, contracts.QueryPaymentResponse{Message: "支付成功"})
}

// queryThirdPartyPaymentStatus 查询第三方支付状态
func (h *PaymentHandler) queryThirdPartyPaymentStatus(
	c *gin.Context, order *contracts.GetOrderByIdResponse, orderID string,
) {
	// 获取机器支付账户
	account, err := h.paymentService.GetPaymentAccount(order.MachineID)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// 查询支付状态
	queryReq := contracts.TranQueryRequest{
		Ext1:          account.ReceivingAccount,
		Ext2:          account.ReceivingKey,
		Ext3:          account.ReceivingOrderPrefix,
		ChannelCode:   contracts.ChannelCodeFuiouMerchant,
		OrderNo:       order.OrderNo,
		ModeOfPayment: contracts.ModeOfPaymentWeChat,
	}

	payInfo, queryErr := h.paymentService.TranQuery(queryReq)
	if queryErr != nil || !payInfo.IsSuccess {
		h.SuccessResponse(c, contracts.QueryPaymentResponse{Message: "查询失败"})
		return
	}

	// 处理支付状态
	h.handlePaymentStatusResult(c, payInfo, orderID)
}

// handlePaymentStatusResult 处理支付状态结果
func (h *PaymentHandler) handlePaymentStatusResult(
	c *gin.Context, payInfo *contracts.TranQueryResponse, orderID string,
) {
	switch payInfo.PaymentStatus {
	case contracts.PaymentStatusPaying:
		h.SuccessResponse(c, contracts.QueryPaymentResponse{Message: "支付中"})
	case contracts.PaymentStatusSuccess:
		h.processSuccessfulPayment(c, payInfo, orderID)
	case contracts.PaymentStatusCancel, contracts.PaymentStatusFailure,
		contracts.PaymentStatusTimeout, contracts.PaymentStatusException:
		h.processFailedPayment(c, payInfo, orderID)
	default:
		h.SuccessResponse(c, contracts.QueryPaymentResponse{Message: "支付中"})
	}
}

// processSuccessfulPayment 处理成功支付
func (h *PaymentHandler) processSuccessfulPayment(
	c *gin.Context, payInfo *contracts.TranQueryResponse, orderID string,
) {
	payReq := contracts.PayOrderRequest{
		ID:             orderID,
		ChannelOrderNo: payInfo.TransactionId,
		PaidAt:         payInfo.PaymentTime,
	}
	if payErr := h.paymentService.PayOrder(payReq); payErr != nil {
		h.InternalErrorResponse(c, payErr)
		return
	}
	h.SuccessResponse(c, contracts.QueryPaymentResponse{Message: "支付成功"})
}

// processFailedPayment 处理失败支付
func (h *PaymentHandler) processFailedPayment(c *gin.Context, payInfo *contracts.TranQueryResponse, orderID string) {
	invalidReq := contracts.InvalidOrderRequest{ID: orderID}
	if invalidErr := h.paymentService.InvalidOrder(invalidReq); invalidErr != nil {
		h.InternalErrorResponse(c, invalidErr)
		return
	}
	h.SuccessResponse(c, contracts.QueryPaymentResponse{Message: payInfo.PaymentStatus})
}

// getMemberOpenId 获取会员微信OpenId
func (h *PaymentHandler) getMemberOpenId(c *gin.Context) string {
	// TODO: 从JWT token或session中获取用户的OpenId
	// 这里需要根据实际的认证机制来实现
	// 临时返回mock值用于测试
	return "mock_open_id_" + h.getMemberIDString(c)
}

// getMemberIDString 获取会员ID字符串
func (h *PaymentHandler) getMemberIDString(c *gin.Context) string {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		return ""
	}
	return memberID
}

// getHasCupText 获取是否需要杯子的文本描述
func (h *PaymentHandler) getHasCupText(hasCup bool) string {
	if hasCup {
		return "需要杯子"
	}
	return "不需要杯子"
}

// CallbackHandler 回调处理器
type CallbackHandler struct {
	*BaseHandler
	paymentService services.PaymentServiceInterface
}

// NewCallbackHandler 创建回调处理器
func NewCallbackHandler(db *gorm.DB) *CallbackHandler {
	return &CallbackHandler{
		BaseHandler:    NewBaseHandler(db),
		paymentService: services.NewPaymentService(db),
	}
}

// PaymentResult 支付结果回调
// POST /api/Callback/PaymentResult
func (h *CallbackHandler) PaymentResult(c *gin.Context) {
	var req contracts.PaymentCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"details": err.Error(),
		})
		return
	}

	// 处理支付回调
	response, err := h.paymentService.ProcessPaymentCallback(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "回调处理失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// safeInt64ToInt32 安全地将int64转换为int32，防止整数溢出
func safeInt64ToInt32(value int64) int32 {
	if value > math.MaxInt32 {
		return math.MaxInt32
	}
	if value < math.MinInt32 {
		return math.MinInt32
	}
	return int32(value)
}
