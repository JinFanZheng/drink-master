package repositories

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

// ProductRepositoryInterface 商品仓储接口
type ProductRepositoryInterface interface {
	GetMachineProducts(machineID string) ([]*models.MachineProductPrice, error)
}

// ProductRepository 商品仓储实现
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建商品仓储
func NewProductRepository(db *gorm.DB) ProductRepositoryInterface {
	return &ProductRepository{
		db: db,
	}
}

// GetMachineProducts 获取售货机商品列表
func (r *ProductRepository) GetMachineProducts(machineID string) ([]*models.MachineProductPrice, error) {
	var machineProducts []*models.MachineProductPrice

	err := r.db.Where("machine_id = ?", machineID).
		Preload("Product").
		Order("created_at ASC").
		Find(&machineProducts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get machine products: %w", err)
	}

	return machineProducts, nil
}
