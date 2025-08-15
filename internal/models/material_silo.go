package models

import (
	"time"
)

// MaterialSilo represents the material silo entity - matches production DB structure
type MaterialSilo struct {
	ID         string     `json:"id" gorm:"primaryKey;type:varchar(36);column:Id"`
	MachineId  *string    `json:"machineId" gorm:"type:varchar(36);column:MachineId"`
	No         *string    `json:"no" gorm:"type:varchar(16);column:No"`
	Type       int        `json:"type" gorm:"type:int;column:Type"`
	ProductId  *string    `json:"productId" gorm:"type:varchar(255);column:ProductId"`
	IsSale     BitBool    `json:"isSale" gorm:"column:IsSale"`
	Total      int        `json:"total" gorm:"type:int;column:Total"`
	Stock      int        `json:"stock" gorm:"type:int;column:Stock"`
	SingleFeed int        `json:"singleFeed" gorm:"type:int;column:SingleFeed"`
	Version    int64      `json:"version" gorm:"column:Version"`
	CreatedOn  time.Time  `json:"createdOn" gorm:"column:CreatedOn"`
	UpdatedOn  *time.Time `json:"updatedOn" gorm:"column:UpdatedOn"`

	// Relations - disabled due to field mapping complexities
	// Machine *Machine `json:"machine,omitempty" gorm:"foreignKey:MachineId;references:Id"`
	// Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id"`
}

// IsStockEmpty checks if the stock is empty
func (ms *MaterialSilo) IsStockEmpty() bool {
	return ms.Stock == 0
}

// IsStockLow checks if the stock is low (below 10% of total capacity)
func (ms *MaterialSilo) IsStockLow() bool {
	if ms.Total == 0 {
		return false
	}
	threshold := float64(ms.Total) * 0.1
	return float64(ms.Stock) < threshold
}

// IsStockFull checks if the stock is at total capacity
func (ms *MaterialSilo) IsStockFull() bool {
	return ms.Stock >= ms.Total
}

// GetStockPercentage returns the stock percentage
func (ms *MaterialSilo) GetStockPercentage() float64 {
	if ms.Total == 0 {
		return 0
	}
	return (float64(ms.Stock) / float64(ms.Total)) * 100
}

// CanSale checks if the silo can be sold (has product, has stock, and sale is enabled)
func (ms *MaterialSilo) CanSale() bool {
	return ms.ProductId != nil &&
		ms.Stock > 0 &&
		ms.IsSale.Bool()
}

// UpdateStock updates the stock with validation
func (ms *MaterialSilo) UpdateStock(newStock int) error {
	if newStock < 0 {
		return ErrInvalidStock
	}
	if newStock > ms.Total {
		return ErrStockExceedsCapacity
	}
	ms.Stock = newStock
	return nil
}

// TableName returns the table name for MaterialSilo
func (MaterialSilo) TableName() string {
	return "material_silos"
}
