package models

import (
	"time"

	"github.com/ddteam/drink-master/internal/enums"
)

// Order represents the order entity - matches production DB structure
type Order struct {
	ID             string     `json:"id" gorm:"primaryKey;type:varchar(36);column:Id"`
	MemberId       *string    `json:"memberId" gorm:"type:varchar(36);column:MemberId"`
	MachineId      *string    `json:"machineId" gorm:"type:varchar(36);column:MachineId"`
	ProductId      *string    `json:"productId" gorm:"type:varchar(36);column:ProductId"`
	HasCup         BitBool    `json:"hasCup" gorm:"column:HasCup"`
	OrderNo        *string    `json:"orderNo" gorm:"type:varchar(32);column:OrderNo"`
	TotalAmount    float64    `json:"totalAmount" gorm:"type:decimal(10,2);column:TotalAmount"`
	PayAmount      float64    `json:"payAmount" gorm:"type:decimal(10,2);column:PayAmount"`
	PaymentStatus  int        `json:"paymentStatus" gorm:"type:int;column:PaymentStatus"`
	PaymentTime    *time.Time `json:"paymentTime" gorm:"column:PaymentTime"`
	ChannelOrderNo *string    `json:"channelOrderNo" gorm:"type:varchar(32);column:ChannelOrderNo"`
	MakeStatus     int        `json:"makeStatus" gorm:"type:int;column:MakeStatus"`
	RefundTime     *time.Time `json:"refundTime" gorm:"column:RefundTime"`
	RefundAmount   float64    `json:"refundAmount" gorm:"type:decimal(10,2);column:RefundAmount"`
	RefundReason   *string    `json:"refundReason" gorm:"type:varchar(512);column:RefundReason"`
	Version        int64      `json:"version" gorm:"column:Version"`
	CreatedOn      time.Time  `json:"createdOn" gorm:"column:CreatedOn"`
	UpdatedOn      *time.Time `json:"updatedOn" gorm:"column:UpdatedOn"`

	// Relations - disabled due to field mapping complexities
	// Member  *Member  `json:"member,omitempty" gorm:"foreignKey:MemberId;references:Id"`
	// Machine *Machine `json:"machine,omitempty" gorm:"foreignKey:MachineId;references:Id"`
	// Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductId;references:Id"`
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
