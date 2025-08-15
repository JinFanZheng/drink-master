package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
	"github.com/ddteam/drink-master/internal/repositories"
)

// getSaleStatusAPIString converts BitBool to API string
func getSaleStatusAPIString(isSale models.BitBool) string {
	if isSale.Bool() {
		return "On"
	}
	return "Off"
}

// parseNoToInt converts No field to int, returns 0 if nil or invalid
func parseNoToInt(no *string) int {
	if no == nil {
		return 0
	}
	// If No is a string like "01", "02", try to parse as int
	// For now, just return 0 as default
	return 0
}

// formatTimeToTime formats *time.Time to time.Time, returns zero time if nil
func formatTimeToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// MaterialSiloServiceInterface 物料槽服务接口
type MaterialSiloServiceInterface interface {
	GetPaging(req contracts.GetMaterialSiloPagingRequest) (*contracts.MaterialSiloPaging, error)
	UpdateStock(req contracts.UpdateMaterialSiloStockRequest) (*contracts.MaterialSiloOperationResult, error)
	UpdateProduct(req contracts.UpdateMaterialSiloProductRequest) (*contracts.MaterialSiloOperationResult, error)
	ToggleSaleStatus(req contracts.ToggleSaleMaterialSiloRequest) (*contracts.MaterialSiloOperationResult, error)
	ValidateMachineExists(machineID string) error
	ValidateProductExists(productID string) error
	ValidateMaterialSiloExists(siloID string) error
}

// MaterialSiloService 物料槽服务实现
type MaterialSiloService struct {
	materialSiloRepo repositories.MaterialSiloRepositoryInterface
	machineRepo      repositories.MachineRepositoryInterface
	productRepo      repositories.ProductRepositoryInterface
}

// NewMaterialSiloService 创建物料槽服务
func NewMaterialSiloService(db *gorm.DB) MaterialSiloServiceInterface {
	return &MaterialSiloService{
		materialSiloRepo: repositories.NewMaterialSiloRepository(db),
		machineRepo:      repositories.NewMachineRepository(db),
		productRepo:      repositories.NewProductRepository(db),
	}
}

// GetPaging 分页获取物料槽列表
func (s *MaterialSiloService) GetPaging(
	req contracts.GetMaterialSiloPagingRequest,
) (*contracts.MaterialSiloPaging, error) {
	// 验证机器是否存在
	if err := s.ValidateMachineExists(req.MachineID); err != nil {
		return nil, err
	}

	// 获取分页数据
	silos, totalCount, err := s.materialSiloRepo.GetPaging(req.MachineID, req.PageIndex, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get material silo paging: %w", err)
	}

	// 转换为响应格式
	items := make([]contracts.GetMaterialSiloPagingResponse, 0, len(silos))
	for _, silo := range silos {
		item := contracts.GetMaterialSiloPagingResponse{
			ID: silo.ID,
			MachineID: func() string {
				if silo.MachineId == nil {
					return ""
				}
				return *silo.MachineId
			}(),
			SiloNo:      parseNoToInt(silo.No),
			ProductID:   silo.ProductId,
			Stock:       silo.Stock,
			MaxCapacity: silo.Total,
			SaleStatus:  getSaleStatusAPIString(silo.IsSale),
			UpdatedAt:   formatTimeToTime(silo.UpdatedOn),
		}

		// 产品名称需要通过单独查询获取
		// 由于GORM关联已禁用，这里暂时留空
		// TODO: 实现产品名称查询逻辑

		items = append(items, item)
	}

	return &contracts.MaterialSiloPaging{
		Items:      items,
		TotalCount: totalCount,
		PageIndex:  req.PageIndex,
		PageSize:   req.PageSize,
	}, nil
}

// UpdateStock 更新物料槽库存
func (s *MaterialSiloService) UpdateStock(
	req contracts.UpdateMaterialSiloStockRequest,
) (*contracts.MaterialSiloOperationResult, error) {
	// 验证物料槽是否存在
	silo, err := s.materialSiloRepo.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get material silo: %w", err)
	}
	if silo == nil {
		return &contracts.MaterialSiloOperationResult{
			Success: false,
			Message: "物料槽不存在",
		}, nil
	}

	// 验证库存是否超出容量
	if req.Stock > silo.Total {
		return &contracts.MaterialSiloOperationResult{
			Success: false,
			Message: fmt.Sprintf("库存不能超过最大容量 %d", silo.Total),
		}, nil
	}

	// 更新库存
	err = s.materialSiloRepo.UpdateStock(req.ID, req.Stock)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	return &contracts.MaterialSiloOperationResult{
		Success: true,
		Message: "库存更新成功",
		Data:    "库存已更新",
	}, nil
}

// UpdateProduct 更新物料槽产品
func (s *MaterialSiloService) UpdateProduct(
	req contracts.UpdateMaterialSiloProductRequest,
) (*contracts.MaterialSiloOperationResult, error) {
	// 验证物料槽是否存在
	if err := s.ValidateMaterialSiloExists(req.ID); err != nil {
		return &contracts.MaterialSiloOperationResult{
			Success: false,
			Message: "物料槽不存在",
		}, nil
	}

	// 验证产品是否存在
	if err := s.ValidateProductExists(req.ProductID); err != nil {
		return &contracts.MaterialSiloOperationResult{
			Success: false,
			Message: "产品不存在",
		}, nil
	}

	// 更新产品
	err := s.materialSiloRepo.UpdateProduct(req.ID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &contracts.MaterialSiloOperationResult{
		Success: true,
		Message: "产品更新成功",
		Data:    "产品已更新",
	}, nil
}

// ToggleSaleStatus 切换销售状态
func (s *MaterialSiloService) ToggleSaleStatus(
	req contracts.ToggleSaleMaterialSiloRequest,
) (*contracts.MaterialSiloOperationResult, error) {
	// 验证物料槽是否存在
	silo, err := s.materialSiloRepo.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get material silo: %w", err)
	}
	if silo == nil {
		return &contracts.MaterialSiloOperationResult{
			Success: false,
			Message: "物料槽不存在",
		}, nil
	}

	// 转换销售状态
	saleStatus := enums.SaleStatusFromAPIString(req.SaleStatus)

	// 如果要开启销售，需要检查是否有产品和库存
	if saleStatus == enums.SaleStatusOn {
		if silo.ProductId == nil {
			return &contracts.MaterialSiloOperationResult{
				Success: false,
				Message: "开启销售前需要先设置产品",
			}, nil
		}
		if silo.Stock <= 0 {
			return &contracts.MaterialSiloOperationResult{
				Success: false,
				Message: "开启销售前需要先补充库存",
			}, nil
		}
	}

	// 更新销售状态
	err = s.materialSiloRepo.UpdateSaleStatus(req.ID, saleStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update sale status: %w", err)
	}

	statusDesc := "停售"
	if saleStatus == enums.SaleStatusOn {
		statusDesc = "在售"
	}

	return &contracts.MaterialSiloOperationResult{
		Success: true,
		Message: fmt.Sprintf("销售状态已切换为%s", statusDesc),
		Data:    statusDesc,
	}, nil
}

// ValidateMachineExists 验证机器是否存在
func (s *MaterialSiloService) ValidateMachineExists(machineID string) error {
	machine, err := s.machineRepo.GetByID(machineID)
	if err != nil {
		return fmt.Errorf("failed to check machine existence: %w", err)
	}
	if machine == nil {
		return errors.New("机器不存在")
	}
	return nil
}

// ValidateProductExists 验证产品是否存在
func (s *MaterialSiloService) ValidateProductExists(productID string) error {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("failed to check product existence: %w", err)
	}
	if product == nil {
		return errors.New("产品不存在")
	}
	return nil
}

// ValidateMaterialSiloExists 验证物料槽是否存在
func (s *MaterialSiloService) ValidateMaterialSiloExists(siloID string) error {
	silo, err := s.materialSiloRepo.GetByID(siloID)
	if err != nil {
		return fmt.Errorf("failed to check material silo existence: %w", err)
	}
	if silo == nil {
		return errors.New("物料槽不存在")
	}
	return nil
}
