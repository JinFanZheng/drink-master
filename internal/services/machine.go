package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/repositories"
)

// MachineServiceInterface 售货机服务接口
type MachineServiceInterface interface {
	GetMachinePaging(req contracts.GetMachinePagingRequest) (*contracts.PagingResult, error)
	GetMachineList(machineOwnerID string) ([]*contracts.GetMachineListResponse, error)
	GetMachineByID(id string) (*contracts.GetMachineByIDResponse, error)
	GetProductList(machineID string) ([]contracts.ProductListResponse, error)
	OpenOrCloseBusiness(machineID string, ownerID string) (*contracts.OpenOrCloseBusinessResponse, error)
	CheckDeviceExist(deviceID string) (bool, error)
	ValidateMachineOwnership(machineID string, ownerID string) error
}

// MachineService 售货机服务实现
type MachineService struct {
	machineRepo   repositories.MachineRepositoryInterface
	productRepo   repositories.ProductRepositoryInterface
	deviceService DeviceServiceInterface
}

// NewMachineService 创建售货机服务
func NewMachineService(db *gorm.DB) MachineServiceInterface {
	return &MachineService{
		machineRepo:   repositories.NewMachineRepository(db),
		productRepo:   repositories.NewProductRepository(db),
		deviceService: NewDeviceService(),
	}
}

// GetMachinePaging 分页获取售货机列表（需要机主权限）
func (s *MachineService) GetMachinePaging(req contracts.GetMachinePagingRequest) (*contracts.PagingResult, error) {
	if req.MachineOwnerID == "" {
		return nil, errors.New("machine owner id is required")
	}

	machines, totalCount, err := s.machineRepo.GetPaging(
		req.MachineOwnerID,
		req.Keyword,
		req.Page,
		req.PageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine paging: %w", err)
	}

	// 转换为响应格式
	items := make([]contracts.GetMachinePagingResponse, len(machines))
	for i, machine := range machines {
		deviceID := ""
		if machine.DeviceId != nil {
			deviceID = *machine.DeviceId
		}

		items[i] = contracts.GetMachinePagingResponse{
			ID:             machine.ID,
			MachineNo:      machine.MachineNo,
			Name:           machine.Name,
			Area:           machine.Area,
			Address:        machine.Address,
			BusinessStatus: machine.BusinessStatus,
			DeviceID:       deviceID,
		}
	}

	return &contracts.PagingResult{
		Items:      items,
		TotalCount: totalCount,
		PageIndex:  req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// GetMachineList 获取售货机列表（需要机主权限）
func (s *MachineService) GetMachineList(machineOwnerID string) ([]*contracts.GetMachineListResponse, error) {
	if machineOwnerID == "" {
		return nil, errors.New("machine owner id is required")
	}

	machines, err := s.machineRepo.GetList(machineOwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine list: %w", err)
	}

	// 转换为响应格式
	result := make([]*contracts.GetMachineListResponse, len(machines))
	for i, machine := range machines {
		result[i] = &contracts.GetMachineListResponse{
			ID:             machine.ID,
			MachineNo:      machine.MachineNo,
			Name:           machine.Name,
			BusinessStatus: machine.BusinessStatus,
		}
	}

	return result, nil
}

// GetMachineByID 根据ID获取售货机详情（公开接口）
func (s *MachineService) GetMachineByID(id string) (*contracts.GetMachineByIDResponse, error) {
	machine, err := s.machineRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine by id: %w", err)
	}

	if machine == nil {
		return nil, nil
	}

	// 检查设备在线状态并更新BusinessStatus
	businessStatus := machine.BusinessStatus
	if machine.DeviceId != nil && *machine.DeviceId != "" {
		online, err := s.deviceService.CheckDeviceOnline(*machine.DeviceId)
		if err == nil && !online {
			businessStatus = contracts.BusinessStatusOffline
		}
	}

	deviceID := ""
	if machine.DeviceId != nil {
		deviceID = *machine.DeviceId
	}

	servicePhone := ""
	if machine.ServicePhone != nil {
		servicePhone = *machine.ServicePhone
	}

	return &contracts.GetMachineByIDResponse{
		ID:             machine.ID,
		MachineNo:      machine.MachineNo,
		Name:           machine.Name,
		Area:           machine.Area,
		Address:        machine.Address,
		BusinessStatus: businessStatus,
		DeviceID:       deviceID,
		ServicePhone:   servicePhone,
		CreatedAt:      machine.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      machine.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetProductList 获取售货机商品列表（核心接口）
func (s *MachineService) GetProductList(machineID string) ([]contracts.ProductListResponse, error) {
	machineProducts, err := s.productRepo.GetMachineProducts(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine products: %w", err)
	}

	if len(machineProducts) == 0 {
		return []contracts.ProductListResponse{}, nil
	}

	// 转换为VendingMachine格式（包装为"限时巨惠"分组）
	products := make([]contracts.MachineProductResponse, len(machineProducts))
	for i, mp := range machineProducts {
		category := ""
		description := ""
		if mp.Product != nil {
			if mp.Product.Category != nil {
				category = *mp.Product.Category
			}
			if mp.Product.Description != nil {
				description = *mp.Product.Description
			}
		}

		productName := "Unknown Product"
		if mp.Product != nil {
			productName = mp.Product.Name
		}

		products[i] = contracts.MachineProductResponse{
			ID:              mp.ID,
			Name:            productName,
			Price:           mp.Price,
			PriceWithoutCup: mp.PriceWithoutCup,
			Stock:           mp.Stock,
			Category:        category,
			Description:     description,
		}
	}

	// 基于VendingMachine逻辑，包装为分组格式
	result := []contracts.ProductListResponse{
		{
			Name:     contracts.ProductGroupTimeLimited,
			Products: products,
		},
	}

	return result, nil
}

// OpenOrCloseBusiness 开关营业状态（机主权限）
func (s *MachineService) OpenOrCloseBusiness(
	machineID string, ownerID string,
) (*contracts.OpenOrCloseBusinessResponse, error) {
	// 验证机主权限
	err := s.ValidateMachineOwnership(machineID, ownerID)
	if err != nil {
		return nil, err
	}

	// 获取当前机器状态
	machine, err := s.machineRepo.GetByID(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine: %w", err)
	}

	if machine == nil {
		return nil, errors.New("machine not found")
	}

	// 切换营业状态
	newStatus := contracts.BusinessStatusClose
	if machine.BusinessStatus == contracts.BusinessStatusClose {
		newStatus = contracts.BusinessStatusOpen
	}

	err = s.machineRepo.UpdateBusinessStatus(machineID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update business status: %w", err)
	}

	message := fmt.Sprintf("售货机已%s", map[string]string{
		contracts.BusinessStatusOpen:  "开启营业",
		contracts.BusinessStatusClose: "关闭营业",
	}[newStatus])

	return &contracts.OpenOrCloseBusinessResponse{
		Status:  newStatus,
		Message: message,
	}, nil
}

// CheckDeviceExist 检查设备是否存在
func (s *MachineService) CheckDeviceExist(deviceID string) (bool, error) {
	if deviceID == "" {
		return false, nil
	}

	return s.machineRepo.CheckDeviceExists(deviceID)
}

// ValidateMachineOwnership 验证机器所有权
func (s *MachineService) ValidateMachineOwnership(machineID string, ownerID string) error {
	machine, err := s.machineRepo.GetByID(machineID)
	if err != nil {
		return fmt.Errorf("failed to get machine: %w", err)
	}

	if machine == nil {
		return errors.New("machine not found")
	}

	if machine.MachineOwnerId != ownerID {
		return errors.New("permission denied: not machine owner")
	}

	return nil
}
