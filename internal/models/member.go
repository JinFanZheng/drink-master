package models

import (
	"time"
)

// Member represents the member/user entity
type Member struct {
	ID             string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Nickname       string     `json:"nickname" gorm:"type:varchar(100);not null"`
	Avatar         string     `json:"avatar" gorm:"type:text"`
	WeChatOpenId   string     `json:"weChatOpenId" gorm:"uniqueIndex;type:varchar(100);not null"`
	Role           string     `json:"role" gorm:"type:varchar(20);not null;default:'Member'"` // Member, Owner
	MachineOwnerId *string    `json:"machineOwnerId" gorm:"type:varchar(36)"`
	IsAdmin        bool       `json:"isAdmin" gorm:"default:false"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"index"`

	// Relations
	MachineOwner *MachineOwner `json:"machineOwner,omitempty" gorm:"foreignKey:MachineOwnerId;references:ID"`
	Orders       []Order       `json:"orders,omitempty" gorm:"foreignKey:MemberId"`
}

// TableName returns the table name for Member
func (Member) TableName() string {
	return "members"
}
