package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

// MaterialSiloRepositoryInterface 物料槽仓储接口
type MaterialSiloRepositoryInterface interface {
	GetByID(id string) (*models.MaterialSilo, error)
	GetByMachineID(machineID string) ([]*models.MaterialSilo, error)
	GetPaging(machineID string, page, pageSize int) ([]*models.MaterialSilo, int64, error)
	Create(silo *models.MaterialSilo) error
	Update(silo *models.MaterialSilo) error
	UpdateStock(id string, stock int) error
	UpdateProduct(id string, productID string) error
	UpdateSaleStatus(id string, status enums.SaleStatus) error
	Delete(id string) error
	GetBySiloNo(machineID string, siloNo int) (*models.MaterialSilo, error)
	GetByMachineAndProduct(machineID string, productID string) ([]*models.MaterialSilo, error)
}

// MaterialSiloRepository 物料槽仓储实现
type MaterialSiloRepository struct {
	db *gorm.DB
}

// NewMaterialSiloRepository 创建物料槽仓储
func NewMaterialSiloRepository(db *gorm.DB) MaterialSiloRepositoryInterface {
	return &MaterialSiloRepository{
		db: db,
	}
}

// GetByID 根据ID获取物料槽
func (r *MaterialSiloRepository) GetByID(id string) (*models.MaterialSilo, error) {
	var silo models.MaterialSilo
	err := r.db.Where("id = ?", id).
		Preload("Machine").
		Preload("Product").
		First(&silo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get material silo by id: %w", err)
	}

	return &silo, nil
}

// GetByMachineID 根据机器ID获取所有物料槽
func (r *MaterialSiloRepository) GetByMachineID(machineID string) ([]*models.MaterialSilo, error) {
	var silos []*models.MaterialSilo
	err := r.db.Where("MachineId = ?", machineID).
		Order("No ASC").
		Find(&silos).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get material silos by machine id: %w", err)
	}

	return silos, nil
}

// GetPaging 分页获取物料槽列表
func (r *MaterialSiloRepository) GetPaging(
	machineID string, page, pageSize int,
) ([]*models.MaterialSilo, int64, error) {
	var silos []*models.MaterialSilo
	var totalCount int64

	// 计算总数
	err := r.db.Model(&models.MaterialSilo{}).
		Where("MachineId = ?", machineID).
		Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count material silos: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.Where("MachineId = ?", machineID).
		Order("No ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&silos).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get material silos paging: %w", err)
	}

	return silos, totalCount, nil
}

// Create 创建物料槽
func (r *MaterialSiloRepository) Create(silo *models.MaterialSilo) error {
	err := r.db.Create(silo).Error
	if err != nil {
		return fmt.Errorf("failed to create material silo: %w", err)
	}
	return nil
}

// Update 更新物料槽
func (r *MaterialSiloRepository) Update(silo *models.MaterialSilo) error {
	err := r.db.Save(silo).Error
	if err != nil {
		return fmt.Errorf("failed to update material silo: %w", err)
	}
	return nil
}

// UpdateStock 更新库存
func (r *MaterialSiloRepository) UpdateStock(id string, stock int) error {
	err := r.db.Model(&models.MaterialSilo{}).
		Where("id = ?", id).
		Update("Stock", stock).Error

	if err != nil {
		return fmt.Errorf("failed to update material silo stock: %w", err)
	}

	return nil
}

// UpdateProduct 更新产品
func (r *MaterialSiloRepository) UpdateProduct(id string, productID string) error {
	err := r.db.Model(&models.MaterialSilo{}).
		Where("id = ?", id).
		Update("ProductId", productID).Error

	if err != nil {
		return fmt.Errorf("failed to update material silo product: %w", err)
	}

	return nil
}

// UpdateSaleStatus 更新销售状态
func (r *MaterialSiloRepository) UpdateSaleStatus(id string, status enums.SaleStatus) error {
	var isSale models.BitBool
	if status == enums.SaleStatusOn {
		isSale = models.BitBool(1)
	} else {
		isSale = models.BitBool(0)
	}

	err := r.db.Model(&models.MaterialSilo{}).
		Where("id = ?", id).
		Update("IsSale", isSale).Error

	if err != nil {
		return fmt.Errorf("failed to update material silo sale status: %w", err)
	}

	return nil
}

// Delete 删除物料槽（软删除）
func (r *MaterialSiloRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.MaterialSilo{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete material silo: %w", err)
	}
	return nil
}

// GetBySiloNo 根据机器ID和槽位号获取物料槽
func (r *MaterialSiloRepository) GetBySiloNo(machineID string, siloNo int) (*models.MaterialSilo, error) {
	var silo models.MaterialSilo
	err := r.db.Where("MachineId = ? AND No = ?", machineID, siloNo).
		First(&silo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get material silo by silo no: %w", err)
	}

	return &silo, nil
}

// GetByMachineAndProduct 根据机器ID和产品ID获取物料槽列表
func (r *MaterialSiloRepository) GetByMachineAndProduct(
	machineID string, productID string,
) ([]*models.MaterialSilo, error) {
	var silos []*models.MaterialSilo
	err := r.db.Where("MachineId = ? AND ProductId = ?", machineID, productID).
		Order("No ASC").
		Find(&silos).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get material silos by machine and product: %w", err)
	}

	return silos, nil
}
