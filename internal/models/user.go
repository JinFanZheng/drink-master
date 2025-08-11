package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // json:"-" 防止密码被序列化
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Drinks         []Drink         `gorm:"foreignKey:UserID" json:"drinks,omitempty"`
	ConsumptionLogs []ConsumptionLog `gorm:"foreignKey:UserID" json:"consumption_logs,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM钩子：创建前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.hashPassword()
}

// BeforeUpdate GORM钩子：更新前加密密码（如果密码有变更）
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("password") {
		return u.hashPassword()
	}
	return nil
}

// hashPassword 加密密码
func (u *User) hashPassword() error {
	if u.Password == "" {
		return nil
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetDisplayName 获取显示名称
func (u *User) GetDisplayName() string {
	if u.Username != "" {
		return u.Username
	}
	return u.Email
}