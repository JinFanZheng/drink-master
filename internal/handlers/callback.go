package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/services"
)

// CallbackHandler 回调处理器
type CallbackHandler struct {
	orderService   services.OrderService
	paymentService services.PaymentServiceInterface
	logger         *logrus.Logger
}

// NewCallbackHandler 创建回调处理器
func NewCallbackHandler(
	orderService services.OrderService,
	paymentService services.PaymentServiceInterface,
	logger *logrus.Logger,
) *CallbackHandler {
	return &CallbackHandler{
		orderService:   orderService,
		paymentService: paymentService,
		logger:         logger,
	}
}

// PaymentResult 支付结果回调
// @Summary 支付结果回调接口
// @Description 第三方支付平台回调支付结果的接口
// @Tags Callback
// @Accept json
// @Produce plain
// @Param request body contracts.PaymentCallbackResultRequest true "支付回调请求"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "参数错误"
// @Router /Callback/PaymentResult [post]
func (h *CallbackHandler) PaymentResult(c *gin.Context) {
	var request contracts.PaymentCallbackResultRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("支付回调参数解析失败")
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	h.logger.WithField("request", request).Info("支付结果回调")

	// 根据订单号查找订单
	order, err := h.orderService.GetByOrderNo(request.OrderNo)
	if err != nil {
		h.logger.WithError(err).Error("查询订单失败")
		c.String(http.StatusOK, "查询订单失败")
		return
	}

	if order == nil {
		h.logger.Warn("订单不存在")
		c.String(http.StatusOK, "订单不存在")
		return
	}

	// 如果不是未支付状态，则不处理
	if order.PaymentStatus != int(enums.PaymentStatusWaitPay) {
		h.logger.Info("订单已处理")
		c.String(http.StatusOK, "ok")
		return
	}

	// 处理支付成功
	payOrderReq := contracts.PayOrderRequest{
		ID:             order.ID,
		ChannelOrderNo: request.ChannelOrderNo,
		PaidAt:         request.PaymentTime,
	}

	err = h.paymentService.PayOrder(payOrderReq)
	if err != nil {
		h.logger.WithError(err).WithField("request", request).Error("处理支付回调异常")
		// 这里即使失败也返回ok，避免第三方重复回调
		c.String(http.StatusOK, "ok")
		return
	}

	c.String(http.StatusOK, "ok")
}
