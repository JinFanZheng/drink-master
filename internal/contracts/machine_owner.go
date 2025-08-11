package contracts

import (
	"time"

	"github.com/shopspring/decimal"
)

// ColumnModel represents sales data for a machine (基于VendingMachine.MobileAPI)
type ColumnModel struct {
	Label string          `json:"label"` // 机器名称或标识
	Value decimal.Decimal `json:"value"` // 销售额
}

// GetSalesRequest 获取销售情况请求
type GetSalesRequest struct {
	DateTime *time.Time `form:"dateTime"` // 查询日期，默认为今天
}

// SalesResponse 销售情况响应
type SalesResponse struct {
	Date  time.Time       `json:"date"`
	Sales []ColumnModel   `json:"sales"`
	Total decimal.Decimal `json:"total"`
}

// MachineOwnerSalesResponse API响应格式
type MachineOwnerSalesResponse struct {
	Success bool          `json:"success"`
	Data    []ColumnModel `json:"data"`
	Meta    *Meta         `json:"meta,omitempty"`
	Error   *APIError     `json:"error,omitempty"`
}
