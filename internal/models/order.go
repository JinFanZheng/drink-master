package models

import (
	"time"
)

// Order represents the order entity
type Order struct {
	ID             string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	MemberId       string     `json:"memberId" gorm:"type:varchar(36);not null"`
	MachineId      string     `json:"machineId" gorm:"type:varchar(36);not null"`
	ProductId      string     `json:"productId" gorm:"type:varchar(36);not null"`
	OrderNo        string     `json:"orderNo" gorm:"uniqueIndex;type:varchar(100);not null"`
	HasCup         bool       `json:"hasCup" gorm:"default:true"`
	TotalAmount    float64    `json:"totalAmount" gorm:"type:decimal(10,2);not null"`
	PayAmount      float64    `json:"payAmount" gorm:"type:decimal(10,2);not null"`
	PaymentStatus string `json:"paymentStatus" gorm:"type:varchar(20);not null;default:'WaitPay'"` // WaitPay, Paid, Invalid, Refunded
	PaymentTime    *time.Time `json:"paymentTime"`
	ChannelOrderNo *string `json:"channelOrderNo" gorm:"type:varchar(100)"` // Third-party payment order number
	MakeStatus string `json:"makeStatus" gorm:"type:varchar(20);not null;default:'WaitMake'"` // WaitMake, Making, Made, MakeFail
	RefundTime     *time.Time `json:"refundTime"`
	RefundAmount   float64    `json:"refundAmount" gorm:"type:decimal(10,2);default:0"`
	RefundReason   *string    `json:"refundReason" gorm:"type:text"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"index"`

	// Relations
	Member  *Member  `json:"member,omitempty" gorm:"foreignKey:MemberId;references:ID"`
	Machine *Machine `json:"machine,omitempty" gorm:"foreignKey:MachineId;references:ID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductId;references:ID"`
}

// TableName returns the table name for Order
func (Order) TableName() string {
	return "orders"
}
