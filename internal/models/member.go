package models

import (
	"time"
)

// Member represents the member/user entity - matches production DB structure
type Member struct {
	ID             string     `json:"id" gorm:"primaryKey;type:varchar(36);column:Id"`
	Nickname       *string    `json:"nickname" gorm:"type:varchar(32);column:Nickname"`
	Avatar         *string    `json:"avatar" gorm:"type:varchar(255);column:Avatar"`
	WeChatOpenId   *string    `json:"weChatOpenId" gorm:"type:varchar(36);column:WeChatOpenId"`
	Role           int        `json:"role" gorm:"type:int;column:Role"`
	MachineOwnerId *string    `json:"machineOwnerId" gorm:"type:varchar(36);column:MachineOwnerId"`
	IsAdmin        bool       `json:"isAdmin" gorm:"column:IsAdmin"`
	Version        int64      `json:"version" gorm:"column:Version"`
	CreatedOn      time.Time  `json:"createdOn" gorm:"column:CreatedOn"`
	UpdatedOn      *time.Time `json:"updatedOn" gorm:"column:UpdatedOn"`

	// Relations - disabled due to field mapping complexities
	// MachineOwner *MachineOwner `json:"machineOwner,omitempty" gorm:"foreignKey:MachineOwnerId;references:Id"`
	// Orders       []Order       `json:"orders,omitempty" gorm:"foreignKey:MemberId"`
}

// TableName returns the table name for Member
func (Member) TableName() string {
	return "members"
}
