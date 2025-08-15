package models

import (
	"errors"

	"gorm.io/gorm"
)

// MaterialSilo related errors
var (
	ErrInvalidStock         = errors.New("invalid stock: stock cannot be negative")
	ErrStockExceedsCapacity = errors.New("stock exceeds max capacity")
)

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	// Note: Skip auto-migration for existing production database
	// Models need to match existing database schema structure
	// Issue #54: Align GORM models with production database fields
	return nil
}

// AllModels returns a slice of all model pointers for batch operations
func AllModels() []interface{} {
	return []interface{}{
		&MachineOwner{},
		&Member{},
		&Machine{},
		&Product{},
		&MachineProductPrice{},
		&Order{},
		&FranchiseIntention{},
		&MaterialSilo{},
	}
}
