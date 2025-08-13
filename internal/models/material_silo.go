package models

import (
	"time"

	"github.com/ddteam/drink-master/internal/enums"
)

// MaterialSilo represents the material silo entity (物料槽)
type MaterialSilo struct {
	ID          string           `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MachineID   string           `json:"machineId" gorm:"type:varchar(36);not null;index"`
	SiloNo      int              `json:"siloNo" gorm:"not null"`                                // 物料槽编号
	ProductID   *string          `json:"productId" gorm:"type:varchar(36);index"`               // 产品ID（可能为空）
	Stock       int              `json:"stock" gorm:"default:0;check:stock >= 0"`               // 当前库存
	MaxCapacity int              `json:"maxCapacity" gorm:"default:100;check:max_capacity > 0"` // 最大容量
	SaleStatus  enums.SaleStatus `json:"saleStatus" gorm:"type:int;not null;default:0"`         // 销售状态
	CreatedAt   time.Time        `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time       `json:"deletedAt" gorm:"index"`

	// Relations
	Machine *Machine `json:"machine,omitempty" gorm:"foreignKey:MachineID;references:ID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
}

// GetSaleStatusDesc returns the description of the sale status
func (ms *MaterialSilo) GetSaleStatusDesc() string {
	return enums.GetSaleStatusDesc(ms.SaleStatus)
}

// GetSaleStatusAPIString returns the API string representation of sale status
func (ms *MaterialSilo) GetSaleStatusAPIString() string {
	return ms.SaleStatus.ToAPIString()
}

// IsStockLow checks if the stock is low (below 10% of max capacity)
func (ms *MaterialSilo) IsStockLow() bool {
	if ms.MaxCapacity == 0 {
		return false
	}
	threshold := float64(ms.MaxCapacity) * 0.1
	return float64(ms.Stock) < threshold
}

// IsStockEmpty checks if the stock is empty
func (ms *MaterialSilo) IsStockEmpty() bool {
	return ms.Stock == 0
}

// IsStockFull checks if the stock is at max capacity
func (ms *MaterialSilo) IsStockFull() bool {
	return ms.Stock >= ms.MaxCapacity
}

// GetStockPercentage returns the stock percentage
func (ms *MaterialSilo) GetStockPercentage() float64 {
	if ms.MaxCapacity == 0 {
		return 0
	}
	return (float64(ms.Stock) / float64(ms.MaxCapacity)) * 100
}

// CanSale checks if the silo can be sold (has product, has stock, and sale status is on)
func (ms *MaterialSilo) CanSale() bool {
	return ms.ProductID != nil &&
		ms.Stock > 0 &&
		ms.SaleStatus == enums.SaleStatusOn
}

// UpdateStock updates the stock with validation
func (ms *MaterialSilo) UpdateStock(newStock int) error {
	if newStock < 0 {
		return ErrInvalidStock
	}
	if newStock > ms.MaxCapacity {
		return ErrStockExceedsCapacity
	}
	ms.Stock = newStock
	return nil
}

// TableName returns the table name for MaterialSilo
func (MaterialSilo) TableName() string {
	return "material_silos"
}
