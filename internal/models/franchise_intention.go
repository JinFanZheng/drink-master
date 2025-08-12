package models

import (
	"time"

	"gorm.io/gorm"
)

// FranchiseIntention 加盟意向模型
type FranchiseIntention struct {
	ID               string         `json:"id" gorm:"type:varchar(36);primaryKey;not null"`
	MemberID         string         `json:"memberId" gorm:"type:varchar(36);not null;index"`
	ContactName      string         `json:"contactName" gorm:"type:varchar(50);not null"`
	ContactPhone     string         `json:"contactPhone" gorm:"type:varchar(20);not null"`
	IntendedLocation string         `json:"intendedLocation" gorm:"type:text;not null"`
	Remarks          string         `json:"remarks" gorm:"type:text"`
	Status           string         `json:"status" gorm:"type:varchar(20);not null;default:'Pending'"`
	CreatedAt        time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Member Member `json:"member,omitempty" gorm:"foreignKey:MemberID;references:ID"`
}

// TableName 指定表名
func (FranchiseIntention) TableName() string {
	return "franchise_intentions"
}

// BeforeCreate GORM钩子，创建前验证
func (fi *FranchiseIntention) BeforeCreate(tx *gorm.DB) error {
	// ID需要在创建前由调用方设置
	if fi.ID == "" {
		return nil // 让GORM处理或在service层生成
	}
	return nil
}

// BeforeUpdate GORM钩子，更新前验证
func (fi *FranchiseIntention) BeforeUpdate(tx *gorm.DB) error {
	// 确保状态值有效
	validStatuses := []string{"Pending", "Approved", "Rejected"}
	statusValid := false
	for _, status := range validStatuses {
		if fi.Status == status {
			statusValid = true
			break
		}
	}
	if !statusValid {
		fi.Status = "Pending" // 默认状态
	}
	return nil
}

// IsApproved 检查加盟意向是否已被批准
func (fi *FranchiseIntention) IsApproved() bool {
	return fi.Status == "Approved"
}

// IsRejected 检查加盟意向是否已被拒绝
func (fi *FranchiseIntention) IsRejected() bool {
	return fi.Status == "Rejected"
}

// IsPending 检查加盟意向是否待处理
func (fi *FranchiseIntention) IsPending() bool {
	return fi.Status == "Pending"
}
