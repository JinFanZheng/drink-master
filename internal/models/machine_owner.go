package models

import (
	"time"
)

// MachineOwner represents the machine owner entity
type MachineOwner struct {
	ID               string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name             string     `json:"name" gorm:"type:varchar(100);not null"`
	ReceivingAccount *string    `json:"receivingAccount" gorm:"type:varchar(200)"` // Payment receiving account
	CreatedAt        time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt        *time.Time `json:"deletedAt" gorm:"index"`

	// Relations
	Machines []Machine `json:"machines,omitempty" gorm:"foreignKey:MachineOwnerId"`
	Members  []Member  `json:"members,omitempty" gorm:"foreignKey:MachineOwnerId"`
}

// TableName returns the table name for MachineOwner
func (MachineOwner) TableName() string {
	return "machine_owners"
}
