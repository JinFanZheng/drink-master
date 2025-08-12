package models

import (
	"time"
)

// Product represents the product entity
type Product struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string     `json:"name" gorm:"type:varchar(200);not null"`
	Description *string    `json:"description" gorm:"type:text"`
	Category    *string    `json:"category" gorm:"type:varchar(100)"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deletedAt" gorm:"index"`

	// Relations
	MachineProductPrices []MachineProductPrice `json:"machineProductPrices,omitempty" gorm:"foreignKey:ProductId"`
	Orders               []Order               `json:"orders,omitempty" gorm:"foreignKey:ProductId"`
}

// TableName returns the table name for Product
func (Product) TableName() string {
	return "products"
}

// MachineProductPrice represents the pricing information for products in specific machines
type MachineProductPrice struct {
	ID              string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MachineId       string     `json:"machineId" gorm:"type:varchar(36);not null"`
	ProductId       string     `json:"productId" gorm:"type:varchar(36);not null"`
	Price           float64    `json:"price" gorm:"type:decimal(10,2);not null"`
	PriceWithoutCup float64    `json:"priceWithoutCup" gorm:"type:decimal(10,2);not null"`
	Stock           int        `json:"stock" gorm:"default:0"`
	CreatedAt       time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `json:"deletedAt" gorm:"index"`

	// Relations
	Machine *Machine `json:"machine,omitempty" gorm:"foreignKey:MachineId;references:ID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductId;references:ID"`
}

// TableName returns the table name for MachineProductPrice
func (MachineProductPrice) TableName() string {
	return "machine_product_prices"
}
