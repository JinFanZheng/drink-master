package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
	"github.com/ddteam/drink-master/internal/repositories"
)

// OrderService 订单服务接口
type OrderService interface {
	GetMemberOrderPaging(request contracts.GetMemberOrderPagingRequest) (*contracts.OrderPagingResponse, error)
	GetByID(id string) (*contracts.GetOrderByIdResponse, error)
	GetByOrderNo(orderNo string) (*models.Order, error)
	Create(request contracts.CreateOrderRequest) (*contracts.CreateOrderResponse, error)
	Refund(request contracts.RefundOrderRequest) (*contracts.RefundOrderResponse, error)
}

// orderService 订单服务实现
type orderService struct {
	orderRepo   repositories.OrderRepository
	machineRepo repositories.MachineRepositoryInterface
	memberRepo  *repositories.MemberRepository
	deviceSvc   DeviceServiceInterface
}

// NewOrderService 创建订单服务
func NewOrderService(
	orderRepo repositories.OrderRepository,
	machineRepo repositories.MachineRepositoryInterface,
	memberRepo *repositories.MemberRepository,
	deviceSvc DeviceServiceInterface,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		machineRepo: machineRepo,
		memberRepo:  memberRepo,
		deviceSvc:   deviceSvc,
	}
}

// GetMemberOrderPaging 分页获取会员订单列表
func (s *orderService) GetMemberOrderPaging(
	request contracts.GetMemberOrderPagingRequest,
) (*contracts.OrderPagingResponse, error) {
	// 验证会员是否存在
	_, err := s.memberRepo.GetByID(request.MemberID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("会员不存在")
		}
		return nil, fmt.Errorf("查询会员信息失败: %w", err)
	}

	// 获取订单列表
	orders, total, err := s.orderRepo.GetByMemberPaging(request.MemberID, request.PageIndex, request.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	// 转换为响应格式
	orderResponses := make([]contracts.GetMemberOrderPagingResponse, len(orders))
	for i, order := range orders {
		productName := "Unknown Product" // Since Product association is disabled

		paymentStatus := ""
		switch enums.PaymentStatus(order.PaymentStatus) {
		case enums.PaymentStatusWaitPay:
			paymentStatus = contracts.PaymentStatusWaitPay
		case enums.PaymentStatusPaid:
			paymentStatus = contracts.PaymentStatusPaid
		case enums.PaymentStatusRefunded:
			paymentStatus = contracts.PaymentStatusRefunded
		case enums.PaymentStatusInvalid:
			paymentStatus = contracts.PaymentStatusCancelled
		default:
			paymentStatus = contracts.PaymentStatusWaitPay
		}

		orderNo := ""
		if order.OrderNo != nil {
			orderNo = *order.OrderNo
		}

		orderResponses[i] = contracts.GetMemberOrderPagingResponse{
			ID:                order.ID,
			OrderNo:           orderNo,
			ProductName:       productName,
			PayAmount:         decimal.NewFromFloat(order.PayAmount),
			CreatedAt:         order.CreatedOn,
			PaymentStatus:     paymentStatus,
			PaymentStatusDesc: order.GetPaymentStatusDesc(),
		}
	}

	// 计算分页信息
	totalPages := int(total) / request.PageSize
	if int(total)%request.PageSize > 0 {
		totalPages++
	}

	meta := contracts.PaginationMeta{
		Total:       total,
		Count:       len(orderResponses),
		PerPage:     request.PageSize,
		CurrentPage: request.PageIndex,
		TotalPages:  totalPages,
		HasNext:     request.PageIndex < totalPages,
		HasPrev:     request.PageIndex > 1,
		Meta: &contracts.Meta{
			Timestamp: time.Now(),
			Version:   "v1.0.0",
		},
	}

	return &contracts.OrderPagingResponse{
		Orders: orderResponses,
		Meta:   meta,
	}, nil
}

// GetByID 根据ID获取订单详情
func (s *orderService) GetByID(id string) (*contracts.GetOrderByIdResponse, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单详情失败: %w", err)
	}

	paymentStatus := ""
	switch enums.PaymentStatus(order.PaymentStatus) {
	case enums.PaymentStatusWaitPay:
		paymentStatus = contracts.PaymentStatusWaitPay
	case enums.PaymentStatusPaid:
		paymentStatus = contracts.PaymentStatusPaid
	case enums.PaymentStatusRefunded:
		paymentStatus = contracts.PaymentStatusRefunded
	case enums.PaymentStatusInvalid:
		paymentStatus = contracts.PaymentStatusCancelled
	default:
		paymentStatus = contracts.PaymentStatusWaitPay
	}

	makeStatus := ""
	switch enums.MakeStatus(order.MakeStatus) {
	case enums.MakeStatusWaitMake:
		makeStatus = contracts.MakeStatusWaitMake
	case enums.MakeStatusMaking:
		makeStatus = contracts.MakeStatusMaking
	case enums.MakeStatusMade:
		makeStatus = contracts.MakeStatusMade
	case enums.MakeStatusMakeFail:
		makeStatus = contracts.MakeStatusFailed
	default:
		makeStatus = contracts.MakeStatusWaitMake
	}

	// 构建响应
	orderNo := ""
	machineId := ""
	productId := ""

	if order.OrderNo != nil {
		orderNo = *order.OrderNo
	}
	if order.MachineId != nil {
		machineId = *order.MachineId
	}
	if order.ProductId != nil {
		productId = *order.ProductId
	}

	response := &contracts.GetOrderByIdResponse{
		ID:                order.ID,
		OrderNo:           orderNo,
		MachineID:         machineId,
		ProductID:         productId,
		PayAmount:         decimal.NewFromFloat(order.PayAmount),
		PaymentStatus:     paymentStatus,
		PaymentStatusDesc: order.GetPaymentStatusDesc(),
		MakeStatus:        makeStatus,
		MakeStatusDesc:    order.GetMakeStatusDesc(),
		CreatedAt:         order.CreatedOn,
		PaymentTime:       order.PaymentTime,
		HasCup:            order.HasCup.Bool(),
		RefundAmount:      decimal.NewFromFloat(order.RefundAmount),
		RefundReason:      order.RefundReason,
	}

	// Machine and Product associations are disabled, set default names
	response.MachineName = "Unknown Machine"
	response.ProductName = "Unknown Product"

	return response, nil
}

// Create 创建订单
func (s *orderService) Create(request contracts.CreateOrderRequest) (*contracts.CreateOrderResponse, error) {
	// 验证会员是否存在
	_, err := s.memberRepo.GetByID(request.MemberID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("会员不存在")
		}
		return nil, fmt.Errorf("查询会员信息失败: %w", err)
	}

	// 验证机器是否存在
	machine, err := s.machineRepo.GetByID(request.MachineID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("机器不存在")
		}
		return nil, fmt.Errorf("查询机器信息失败: %w", err)
	}

	// 检查设备是否在线 - Use MachineNo as device identifier
	deviceId := machine.MachineNo
	online, err := s.deviceSvc.CheckDeviceOnline(func() string { if deviceId == nil { return "" }; return *deviceId }())
	if err != nil {
		return nil, fmt.Errorf("检查设备状态失败: %w", err)
	}
	if !online {
		return nil, fmt.Errorf("机器不在线，下单失败")
	}

	// 生成订单号
	orderNo := s.generateOrderNo()

	// 创建订单
	order := &models.Order{
		ID:            uuid.New().String(),
		MemberId:      &request.MemberID,
		MachineId:     &request.MachineID,
		ProductId:     &request.ProductID,
		OrderNo:       &orderNo,
		HasCup:        models.BitBool(func() int8 { if request.HasCup { return 1 }; return 0 }()),
		TotalAmount:   request.PayAmount.InexactFloat64(),
		PayAmount:     request.PayAmount.InexactFloat64(),
		PaymentStatus: int(enums.PaymentStatusWaitPay),
		MakeStatus:    int(enums.MakeStatusWaitMake),
		RefundAmount:  0,
	}

	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	responseOrderNo := ""
	if order.OrderNo != nil {
		responseOrderNo = *order.OrderNo
	}

	return &contracts.CreateOrderResponse{
		OrderID: order.ID,
		OrderNo: responseOrderNo,
		Message: "订单创建成功",
	}, nil
}

// Refund 退款订单
func (s *orderService) Refund(request contracts.RefundOrderRequest) (*contracts.RefundOrderResponse, error) {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(request.OrderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("获取订单信息失败: %w", err)
	}

	// 检查订单状态
	if order.PaymentStatus != int(enums.PaymentStatusPaid) {
		return nil, fmt.Errorf("订单状态不允许退款")
	}

	if order.PaymentStatus == int(enums.PaymentStatusRefunded) {
		return nil, fmt.Errorf("订单已经退款")
	}

	// 只有机主可以退款
	if !request.IsMachineOwner {
		return nil, fmt.Errorf("您不是机主，无法退款")
	}

	// 更新订单状态
	now := time.Now()
	order.PaymentStatus = int(enums.PaymentStatusRefunded)
	order.RefundTime = &now
	order.RefundAmount = order.PayAmount
	order.RefundReason = &request.Reason

	err = s.orderRepo.Update(order)
	if err != nil {
		return nil, fmt.Errorf("更新订单状态失败: %w", err)
	}

	return &contracts.RefundOrderResponse{
		OrderID:      order.ID,
		RefundAmount: decimal.NewFromFloat(order.RefundAmount),
		Message:      "退款成功",
	}, nil
}

// GetByOrderNo 根据订单号获取订单
func (s *orderService) GetByOrderNo(orderNo string) (*models.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("根据订单号获取订单失败: %w", err)
	}
	return order, nil
}

// generateOrderNo 生成订单号
func (s *orderService) generateOrderNo() string {
	return fmt.Sprintf("ORD%s", time.Now().Format("20060102150405"))
}
