package models

import (
	"time"

	"gorm.io/gorm"
)

// Drink 饮品模型
type Drink struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;index" json:"name"`
	Category    string         `gorm:"size:50;not null;index" json:"category"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Description string         `gorm:"type:text" json:"description"`
	ImageURL    string         `gorm:"size:500" json:"image_url"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	User            User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ConsumptionLogs []ConsumptionLog `gorm:"foreignKey:DrinkID" json:"consumption_logs,omitempty"`
}

// TableName 指定表名
func (Drink) TableName() string {
	return "drinks"
}

// DrinkCategory 饮品分类模型
type DrinkCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	DisplayName string         `gorm:"size:100;not null" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"size:50" json:"icon"`
	Color       string         `gorm:"size:20" json:"color"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (DrinkCategory) TableName() string {
	return "drink_categories"
}

// ConsumptionLog 消费记录模型
type ConsumptionLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	DrinkID    uint           `gorm:"not null;index" json:"drink_id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Quantity   int            `gorm:"not null;default:1" json:"quantity"`
	ConsumedAt time.Time      `gorm:"not null;index" json:"consumed_at"`
	Notes      string         `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Drink Drink `gorm:"foreignKey:DrinkID" json:"drink,omitempty"`
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (ConsumptionLog) TableName() string {
	return "consumption_logs"
}

// BeforeCreate GORM钩子：创建前设置默认消费时间
func (c *ConsumptionLog) BeforeCreate(tx *gorm.DB) error {
	if c.ConsumedAt.IsZero() {
		c.ConsumedAt = time.Now()
	}
	if c.Quantity <= 0 {
		c.Quantity = 1
	}
	return nil
}

// GetTotalAmount 计算总金额（需要关联查询Drink）
func (c *ConsumptionLog) GetTotalAmount() float64 {
	if c.Drink.Price > 0 {
		return c.Drink.Price * float64(c.Quantity)
	}
	return 0
}