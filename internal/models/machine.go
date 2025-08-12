package models

import (
	"time"

	"github.com/ddteam/drink-master/internal/enums"
)

// Machine represents the vending machine entity
type Machine struct {
	ID             string               `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MachineOwnerId string               `json:"machineOwnerId" gorm:"type:varchar(36);not null"`
	MachineNo      string               `json:"machineNo" gorm:"uniqueIndex;type:varchar(50);not null"`
	Name           string               `json:"name" gorm:"type:varchar(200);not null"`
	Area           string               `json:"area" gorm:"type:varchar(100)"`
	Address        string               `json:"address" gorm:"type:text"`
	ServicePhone   *string              `json:"servicePhone" gorm:"type:varchar(20)"`
	DeviceId       *string              `json:"deviceId" gorm:"type:varchar(100)"`
	BusinessStatus enums.BusinessStatus `json:"businessStatus" gorm:"type:int;not null;default:1"`
	CreatedAt      time.Time            `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time            `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time           `json:"deletedAt" gorm:"index"`

	// Relations
	MachineOwner       *MachineOwner         `json:"machineOwner,omitempty" gorm:"foreignKey:MachineOwnerId;references:ID"`
	MachineProductList []MachineProductPrice `json:"machineProductList,omitempty" gorm:"foreignKey:MachineId"`
	Orders             []Order               `json:"orders,omitempty" gorm:"foreignKey:MachineId"`
}

// GetBusinessStatusDesc returns the description of the business status
func (m *Machine) GetBusinessStatusDesc() string {
	return enums.GetBusinessStatusDesc(m.BusinessStatus)
}

// TableName returns the table name for Machine
func (Machine) TableName() string {
	return "machines"
}
