package models

import (
	"gorm.io/gorm"
)

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&MachineOwner{},
		&Member{},
		&Machine{},
		&Product{},
		&MachineProductPrice{},
		&Order{},
		&FranchiseIntention{},
	)
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
	}
}
