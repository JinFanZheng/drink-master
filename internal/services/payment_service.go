package services

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/repositories"
	"gorm.io/gorm"
)

// PaymentServiceInterface 支付服务接口
type PaymentServiceInterface interface {
	WeChatPay(req contracts.WeChatPayRequest) (*contracts.WeChatPayResponse, error)
	TranQuery(req contracts.TranQueryRequest) (*contracts.TranQueryResponse, error)
	GetPaymentAccount(machineID string) (*contracts.PaymentAccount, error)
	PayOrder(req contracts.PayOrderRequest) error
	InvalidOrder(req contracts.InvalidOrderRequest) error
	ProcessPaymentCallback(req contracts.PaymentCallbackRequest) (*contracts.PaymentCallbackResponse, error)
}

// paymentService 支付服务实现
type paymentService struct {
	orderRepo   repositories.OrderRepository
	machineRepo repositories.MachineRepositoryInterface
	httpClient  *http.Client
}

// NewPaymentService 创建支付服务
func NewPaymentService(db *gorm.DB) PaymentServiceInterface {
	return &paymentService{
		orderRepo:   repositories.NewOrderRepository(db),
		machineRepo: repositories.NewMachineRepository(db),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WeChatPay 发起微信支付
func (s *paymentService) WeChatPay(req contracts.WeChatPayRequest) (*contracts.WeChatPayResponse, error) {
	// 在实际生产环境中，这里应该调用真实的微信支付接口
	// 目前使用mock实现，返回模拟的支付信息

	if os.Getenv("MOCK_MODE") == "true" {
		return s.mockWeChatPay(req)
	}

	// TODO: 实现真实的微信支付调用
	return s.mockWeChatPay(req)
}

// TranQuery 查询支付状态
func (s *paymentService) TranQuery(req contracts.TranQueryRequest) (*contracts.TranQueryResponse, error) {
	// 在实际生产环境中，这里应该调用真实的支付查询接口
	// 目前使用mock实现

	if os.Getenv("MOCK_MODE") == "true" {
		return s.mockTranQuery(req)
	}

	// TODO: 实现真实的支付查询调用
	return s.mockTranQuery(req)
}

// GetPaymentAccount 获取机器支付账户信息
func (s *paymentService) GetPaymentAccount(machineID string) (*contracts.PaymentAccount, error) {
	machine, err := s.machineRepo.GetByID(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine: %w", err)
	}

	if machine == nil {
		return nil, errors.New("machine not found")
	}

	// 构建支付账户信息
	// TODO: 从实际的机器配置或数据库中获取真实的支付账户信息
	account := &contracts.PaymentAccount{
		ReceivingAccount:     getEnvOrDefault("WECHAT_PAY_MERCHANT_ID", "test_merchant_001"),
		ReceivingKey:         getEnvOrDefault("WECHAT_PAY_API_KEY", "test_key_123"),
		ReceivingOrderPrefix: fmt.Sprintf("VM_%s_", machine.MachineNo),
	}

	return account, nil
}

// PayOrder 支付订单
func (s *paymentService) PayOrder(req contracts.PayOrderRequest) error {
	order, err := s.orderRepo.GetByID(req.ID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return errors.New("order not found")
	}

	// 检查订单状态
	if order.PaymentStatus != int(enums.PaymentStatusWaitPay) {
		return errors.New("order is not in wait pay status")
	}

	// 更新订单状态为已支付
	order.PaymentStatus = int(enums.PaymentStatusPaid)
	order.ChannelOrderNo = &req.ChannelOrderNo
	order.PaymentTime = &req.PaidAt

	err = s.orderRepo.Update(order)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// InvalidOrder 作废订单
func (s *paymentService) InvalidOrder(req contracts.InvalidOrderRequest) error {
	order, err := s.orderRepo.GetByID(req.ID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order == nil {
		return errors.New("order not found")
	}

	// 更新订单状态为已取消
	order.PaymentStatus = int(enums.PaymentStatusInvalid)

	err = s.orderRepo.Update(order)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

// ProcessPaymentCallback 处理支付回调
func (s *paymentService) ProcessPaymentCallback(
	req contracts.PaymentCallbackRequest,
) (*contracts.PaymentCallbackResponse, error) {
	// 根据订单号查找订单
	order, err := s.orderRepo.GetByOrderNo(req.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by order no: %w", err)
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	// 验证回调签名（实际项目中需要实现签名验证）
	if !s.verifyCallbackSignature(req) {
		return nil, errors.New("invalid callback signature")
	}

	// 根据支付状态处理订单
	switch req.Status {
	case contracts.PaymentStatusSuccess:
		return s.handleSuccessCallback(req)
	case contracts.PaymentStatusFailure, contracts.PaymentStatusCancel:
		return s.handleFailureCallback(req)
	default:
		return &contracts.PaymentCallbackResponse{
			Processed: true,
			Message:   "回调已接收，状态待处理",
		}, nil
	}
}

// mock和辅助函数
func (s *paymentService) mockWeChatPay(req contracts.WeChatPayRequest) (*contracts.WeChatPayResponse, error) {
	// 模拟微信支付响应
	return &contracts.WeChatPayResponse{
		IsSuccess: true,
		AppId:     "mock_app_id",
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  "mock_nonce_str",
		Package:   "prepay_id=mock_prepay_id",
		SignType:  "RSA",
		PaySign:   "mock_pay_sign",
		Message:   "支付信息获取成功",
	}, nil
}

func (s *paymentService) mockTranQuery(req contracts.TranQueryRequest) (*contracts.TranQueryResponse, error) {
	// 模拟支付查询响应
	return &contracts.TranQueryResponse{
		IsSuccess:     true,
		PaymentStatus: contracts.PaymentStatusSuccess,
		TransactionId: "mock_transaction_id_" + req.OrderNo,
		PaymentTime:   time.Now(),
		Message:       "查询成功",
	}, nil
}

// verifyCallbackSignature 验证回调签名
func (s *paymentService) verifyCallbackSignature(req contracts.PaymentCallbackRequest) bool {
	// TODO: 实现真实的签名验证逻辑
	// 这里应该根据支付提供商的签名算法进行验证
	return req.Signature != ""
}

// handleSuccessCallback 处理成功支付回调
func (s *paymentService) handleSuccessCallback(
	req contracts.PaymentCallbackRequest,
) (*contracts.PaymentCallbackResponse, error) {
	payOrderReq := contracts.PayOrderRequest{
		ID:             req.OrderNo, // 注意：这里需要根据OrderNo获取内部OrderID
		ChannelOrderNo: req.TransactionId,
		PaidAt:         req.PaidAt,
	}

	// 首先通过订单号获取订单ID
	order, err := s.orderRepo.GetByOrderNo(req.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	payOrderReq.ID = order.ID

	if err := s.PayOrder(payOrderReq); err != nil {
		return nil, fmt.Errorf("failed to pay order: %w", err)
	}

	return &contracts.PaymentCallbackResponse{
		Processed: true,
		Message:   "回调处理成功",
	}, nil
}

// handleFailureCallback 处理失败支付回调
func (s *paymentService) handleFailureCallback(
	req contracts.PaymentCallbackRequest,
) (*contracts.PaymentCallbackResponse, error) {
	// 首先通过订单号获取订单ID
	order, err := s.orderRepo.GetByOrderNo(req.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	invalidReq := contracts.InvalidOrderRequest{ID: order.ID}
	if err := s.InvalidOrder(invalidReq); err != nil {
		return nil, fmt.Errorf("failed to invalid order: %w", err)
	}

	return &contracts.PaymentCallbackResponse{
		Processed: true,
		Message:   "回调处理成功",
	}, nil
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
