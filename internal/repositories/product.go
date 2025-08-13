package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

// ProductRepositoryInterface 商品仓储接口
type ProductRepositoryInterface interface {
	GetByID(id string) (*models.Product, error)
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

// GetByID 根据ID获取产品
func (r *ProductRepository) GetByID(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	return &product, nil
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
