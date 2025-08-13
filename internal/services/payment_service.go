package services

import (
	"bytes"
	"encoding/json"
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

	// 真实的微信支付API调用逻辑
	return s.realWeChatPay(req)
}

// TranQuery 查询支付状态
func (s *paymentService) TranQuery(req contracts.TranQueryRequest) (*contracts.TranQueryResponse, error) {
	// 在实际生产环境中，这里应该调用真实的支付查询接口
	// 目前使用mock实现

	if os.Getenv("MOCK_MODE") == "true" {
		return s.mockTranQuery(req)
	}

	// 真实的支付查询API调用逻辑
	return s.realTranQuery(req)
}

// GetPaymentAccount 获取机器的支付账户信息
func (s *paymentService) GetPaymentAccount(machineID string) (*contracts.PaymentAccount, error) {
	machine, err := s.machineRepo.GetByID(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine: %w", err)
	}

	if machine == nil {
		return nil, errors.New("machine not found")
	}

	// 在实际实现中，这些账户信息应该从机器配置或者机主配置中获取
	// 目前使用环境变量或默认值
	account := &contracts.PaymentAccount{
		ReceivingAccount:     getEnvOrDefault("WECHAT_PAY_MERCHANT_ID", "default_merchant_id"),
		ReceivingKey:         getEnvOrDefault("WECHAT_PAY_API_KEY", "default_api_key"),
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

	// 更新订单状态
	order.PaymentStatus = int(enums.PaymentStatusPaid)
	order.PaymentTime = &req.PaidAt
	order.ChannelOrderNo = &req.ChannelOrderNo

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

	// 更新订单状态为失效
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

	// 根据支付状态更新订单
	switch req.Status {
	case contracts.PaymentStatusSuccess:
		payReq := contracts.PayOrderRequest{
			ID:             order.ID,
			ChannelOrderNo: req.TransactionId,
			PaidAt:         req.PaidAt,
		}
		err = s.PayOrder(payReq)
		if err != nil {
			return nil, fmt.Errorf("failed to pay order: %w", err)
		}
	case contracts.PaymentStatusCancel, contracts.PaymentStatusFailure, contracts.PaymentStatusTimeout:
		invalidReq := contracts.InvalidOrderRequest{ID: order.ID}
		err = s.InvalidOrder(invalidReq)
		if err != nil {
			return nil, fmt.Errorf("failed to invalid order: %w", err)
		}
	}

	return &contracts.PaymentCallbackResponse{
		Processed: true,
		Message:   "回调处理成功",
	}, nil
}

// mockWeChatPay mock微信支付
func (s *paymentService) mockWeChatPay(req contracts.WeChatPayRequest) (*contracts.WeChatPayResponse, error) {
	// 模拟支付成功的响应
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

// mockTranQuery mock支付查询
func (s *paymentService) mockTranQuery(req contracts.TranQueryRequest) (*contracts.TranQueryResponse, error) {
	// 模拟支付成功的查询响应
	return &contracts.TranQueryResponse{
		IsSuccess:     true,
		PaymentStatus: contracts.PaymentStatusSuccess,
		TransactionId: "mock_transaction_id_" + req.OrderNo,
		PaymentTime:   time.Now(),
		Message:       "查询成功",
	}, nil
}

// realWeChatPay 真实的微信支付API调用
func (s *paymentService) realWeChatPay(req contracts.WeChatPayRequest) (*contracts.WeChatPayResponse, error) {
	// TODO: 实现真实的微信支付API调用
	// 这里需要根据实际的第三方支付接口文档进行实现

	paymentApiUrl := getEnvOrDefault("PAYMENT_API_URL", "https://payment-api.example.com")

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.httpClient.Post(paymentApiUrl+"/wechat/pay", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call payment api: %w", err)
	}
	defer resp.Body.Close()

	var payResp contracts.WeChatPayResponse
	err = json.NewDecoder(resp.Body).Decode(&payResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &payResp, nil
}

// realTranQuery 真实的支付查询API调用
func (s *paymentService) realTranQuery(req contracts.TranQueryRequest) (*contracts.TranQueryResponse, error) {
	// TODO: 实现真实的支付查询API调用

	paymentApiUrl := getEnvOrDefault("PAYMENT_API_URL", "https://payment-api.example.com")

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.httpClient.Post(paymentApiUrl+"/query", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call query api: %w", err)
	}
	defer resp.Body.Close()

	var queryResp contracts.TranQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &queryResp, nil
}

// verifyCallbackSignature 验证回调签名
func (s *paymentService) verifyCallbackSignature(req contracts.PaymentCallbackRequest) bool {
	// TODO: 实现实际的签名验证逻辑
	// 目前为了简化，总是返回true
	return true
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
