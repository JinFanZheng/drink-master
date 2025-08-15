package models

import (
	"time"

	"gorm.io/gorm"
)

// FranchiseIntention 加盟意向模型 - matches production DB structure
type FranchiseIntention struct {
	ID        string     `json:"id" gorm:"primaryKey;type:varchar(36);column:Id"`
	MemberId  *string    `json:"memberId" gorm:"type:varchar(36);column:MemberId"`
	Company   *string    `json:"company" gorm:"type:varchar(255);column:Company"`
	Name      *string    `json:"name" gorm:"type:varchar(32);column:Name"`
	Mobile    *string    `json:"mobile" gorm:"type:varchar(11);column:Mobile"`
	Area      *string    `json:"area" gorm:"type:varchar(64);column:Area"`
	Remark    *string    `json:"remark" gorm:"type:varchar(512);column:Remark"`
	IsHandled BitBool    `json:"isHandled" gorm:"column:IsHandled"`
	Version   int64      `json:"version" gorm:"column:Version"`
	CreatedOn time.Time  `json:"createdOn" gorm:"column:CreatedOn"`
	UpdatedOn *time.Time `json:"updatedOn" gorm:"column:UpdatedOn"`

	// Relations - disabled due to field mapping complexities
	// Member Member `json:"member,omitempty" gorm:"foreignKey:MemberId;references:Id"`
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
	// 数据库使用 IsHandled 字段，无需状态验证
	return nil
}

// GetHandledStatus 检查加盟意向是否已被处理
func (fi *FranchiseIntention) GetHandledStatus() bool {
	return fi.IsHandled.Bool()
}

// IsPending 检查加盟意向是否待处理
func (fi *FranchiseIntention) IsPending() bool {
	return !fi.IsHandled.Bool()
}
