package models

import (
	"time"
)

// MachineOwner represents the machine owner entity - matches production DB structure
type MachineOwner struct {
	ID                   string     `json:"id" gorm:"primaryKey;type:varchar(36);column:Id"`
	Name                 *string    `json:"name" gorm:"type:varchar(32);column:Name"`
	Mobile               *string    `json:"mobile" gorm:"type:varchar(11);column:Mobile"`
	Email                *string    `json:"email" gorm:"type:varchar(64);column:Email"`
	ReceivingAccount     *string    `json:"receivingAccount" gorm:"type:varchar(16);column:ReceivingAccount"`
	ReceivingKey         *string    `json:"receivingKey" gorm:"type:varchar(32);column:ReceivingKey"`
	ReceivingOrderPrefix *string    `json:"receivingOrderPrefix" gorm:"type:varchar(5);column:ReceivingOrderPrefix"`
	Version              int64      `json:"version" gorm:"column:Version"`
	CreatedOn            time.Time  `json:"createdOn" gorm:"column:CreatedOn"`
	UpdatedOn            *time.Time `json:"updatedOn" gorm:"column:UpdatedOn"`

	// Relations - disabled due to field mapping complexities
	// Machines []Machine `json:"machines,omitempty" gorm:"foreignKey:MachineOwnerId"`
	// Members  []Member  `json:"members,omitempty" gorm:"foreignKey:MachineOwnerId"`
}

// TableName returns the table name for MachineOwner
func (MachineOwner) TableName() string {
	return "machine_owners"
}
