package models

import (
	"time"
)

// Product represents the product entity - matches production DB structure
type Product struct {
	ID              string     `json:"id" gorm:"primaryKey;type:varchar(36);column:Id"`
	Name            string     `json:"name" gorm:"type:varchar(32);column:Name"`
	Image           *string    `json:"image" gorm:"type:varchar(255);column:Image"`
	Status          int        `json:"status" gorm:"type:int;column:Status"`
	Price           float64    `json:"price" gorm:"type:decimal(10,2);column:Price"`
	PriceWithoutCup float64    `json:"priceWithoutCup" gorm:"type:decimal(10,2);column:PriceWithoutCup"`
	Version         int64      `json:"version" gorm:"column:Version"`
	CreatedOn       time.Time  `json:"createdOn" gorm:"column:CreatedOn"`
	UpdatedOn       *time.Time `json:"updatedOn" gorm:"column:UpdatedOn"`

	// Relations - disabled due to field mapping complexities
	// MachineProductPrices []MachineProductPrice `json:"machineProductPrices,omitempty" gorm:"foreignKey:ProductId"`
	// Orders               []Order               `json:"orders,omitempty" gorm:"foreignKey:ProductId"`
}

// TableName returns the table name for Product
func (Product) TableName() string {
	return "products"
}

// MachineProductPrice represents the pricing information for products in specific machines - matches production DB
type MachineProductPrice struct {
	ID              string     `json:"id" gorm:"primaryKey;type:varchar(36);column:Id"`
	MachineId       string     `json:"machineId" gorm:"type:varchar(255);column:MachineId"`
	ProductId       string     `json:"productId" gorm:"type:varchar(255);column:ProductId"`
	Price           float64    `json:"price" gorm:"type:decimal(10,2);column:Price"`
	PriceWithoutCup float64    `json:"priceWithoutCup" gorm:"type:decimal(10,2);column:PriceWithoutCup"`
	Version         int64      `json:"version" gorm:"column:Version"`
	CreatedOn       time.Time  `json:"createdOn" gorm:"column:CreatedOn"`
	UpdatedOn       *time.Time `json:"updatedOn" gorm:"column:UpdatedOn"`

	// Relations - disabled due to field mapping complexities in production
	// Machine *Machine `json:"machine,omitempty" gorm:"foreignKey:MachineId;references:Id"`
	// Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id"`
}

// TableName returns the table name for MachineProductPrice
func (MachineProductPrice) TableName() string {
	return "machine_product_prices"
}
