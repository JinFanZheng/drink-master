package models

import (
	"time"

	"github.com/ddteam/drink-master/internal/enums"
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
	PaymentStatus  int        `json:"paymentStatus" gorm:"type:int;not null;default:0"`
	PaymentTime    *time.Time `json:"paymentTime"`
	ChannelOrderNo *string    `json:"channelOrderNo" gorm:"type:varchar(100)"`
	MakeStatus     int        `json:"makeStatus" gorm:"type:int;not null;default:0"`
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

// GetPaymentStatusDesc 获取支付状态描述
func (o *Order) GetPaymentStatusDesc() string {
	return enums.GetPaymentStatusDesc(enums.PaymentStatus(o.PaymentStatus))
}

// GetMakeStatusDesc 获取制作状态描述
func (o *Order) GetMakeStatusDesc() string {
	return enums.GetMakeStatusDesc(enums.MakeStatus(o.MakeStatus))
}
