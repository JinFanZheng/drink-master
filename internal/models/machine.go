package models

import (
	"time"

	"github.com/ddteam/drink-master/internal/enums"
)

// Machine represents the vending machine entity - matches production DB structure
type Machine struct {
	ID              string               `json:"id" gorm:"primaryKey;type:varchar(36);column:Id"`
	MachineOwnerId  *string              `json:"machineOwnerId" gorm:"type:varchar(36);column:MachineOwnerId"`
	MachineNo       *string              `json:"machineNo" gorm:"type:varchar(32);column:MachineNo"`
	Name            *string              `json:"name" gorm:"type:varchar(32);column:Name"`
	Area            *string              `json:"area" gorm:"type:varchar(64);column:Area"`
	Address         *string              `json:"address" gorm:"type:varchar(128);column:Address"`
	ServicePhone    *string              `json:"servicePhone" gorm:"type:varchar(11);column:ServicePhone"`
	BusinessStatus  enums.BusinessStatus `json:"businessStatus" gorm:"type:int;column:BusinessStatus"`
	SubscribeTime   *time.Time           `json:"subscribeTime" gorm:"column:SubscribeTime"`
	UnSubscribeTime *time.Time           `json:"unSubscribeTime" gorm:"column:UnSubscribeTime"`
	DeviceId        *string              `json:"deviceId" gorm:"type:varchar(255);column:DeviceId"`
	DeviceName      *string              `json:"deviceName" gorm:"type:varchar(255);column:DeviceName"`
	DeviceSn        *string              `json:"deviceSn" gorm:"type:varchar(255);column:DeviceSn"`
	BindDeviceTime  *time.Time           `json:"bindDeviceTime" gorm:"column:BindDeviceTime"`
	IsDebugMode     BitBool              `json:"isDebugMode" gorm:"column:IsDebugMode"`
	Version         int64                `json:"version" gorm:"column:Version"`
	CreatedOn       time.Time            `json:"createdOn" gorm:"column:CreatedOn"`
	UpdatedOn       *time.Time           `json:"updatedOn" gorm:"column:UpdatedOn"`

	// Relations - disabled due to field mapping complexities in production
	// MachineOwner       *MachineOwner         `json:"machineOwner,omitempty" gorm:"foreignKey:MachineOwnerId;references:Id"`
	// MachineProductList []MachineProductPrice `json:"machineProductList,omitempty" gorm:"foreignKey:MachineId"`
	// Orders             []Order               `json:"orders,omitempty" gorm:"foreignKey:MachineId"`
}

// GetBusinessStatusDesc returns the description of the business status
func (m *Machine) GetBusinessStatusDesc() string {
	return enums.GetBusinessStatusDesc(m.BusinessStatus)
}

// TableName returns the table name for Machine
func (Machine) TableName() string {
	return "machines"
}
