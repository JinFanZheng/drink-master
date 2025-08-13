package contracts

import "time"

// MaterialSilo API契约定义 - 对应VendingMachine.MobileAPI MaterialSiloController

// GetMaterialSiloPagingRequest 获取物料槽分页列表请求
type GetMaterialSiloPagingRequest struct {
	MachineID string `json:"machineId" binding:"required"`
	PageIndex int    `json:"pageIndex" binding:"required,min=1"`
	PageSize  int    `json:"pageSize" binding:"required,min=1,max=100"`
}

// GetMaterialSiloPagingResponse 获取物料槽分页列表响应项
type GetMaterialSiloPagingResponse struct {
	ID          string    `json:"id"`
	MachineID   string    `json:"machineId"`
	SiloNo      int       `json:"siloNo"`      // 物料槽编号
	ProductID   *string   `json:"productId"`   // 产品ID（可能为空）
	ProductName *string   `json:"productName"` // 产品名称（可能为空）
	Stock       int       `json:"stock"`       // 当前库存
	MaxCapacity int       `json:"maxCapacity"` // 最大容量
	SaleStatus  string    `json:"saleStatus"`  // 销售状态 (On/Off)
	UpdatedAt   time.Time `json:"updatedAt"`
}

// UpdateMaterialSiloStockRequest 更新料仓库存请求
type UpdateMaterialSiloStockRequest struct {
	ID    string `json:"id" binding:"required"`
	Stock int    `json:"stock" binding:"min=0"`
}

// UpdateMaterialSiloProductRequest 更新料仓产品请求
type UpdateMaterialSiloProductRequest struct {
	ID        string `json:"id" binding:"required"`
	ProductID string `json:"productId" binding:"required"`
}

// ToggleSaleMaterialSiloRequest 切换销售状态请求
type ToggleSaleMaterialSiloRequest struct {
	ID         string `json:"id" binding:"required"`
	SaleStatus string `json:"saleStatus" binding:"required,oneof=On Off"`
}

// MaterialSiloPaging 物料槽分页结果
type MaterialSiloPaging struct {
	Items      []GetMaterialSiloPagingResponse `json:"items"`
	TotalCount int64                           `json:"totalCount"`
	PageIndex  int                             `json:"pageIndex"`
	PageSize   int                             `json:"pageSize"`
}

// MaterialSiloOperationResult 物料槽操作结果
type MaterialSiloOperationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// 销售状态枚举常量
const (
	SaleStatusOn  = "On"
	SaleStatusOff = "Off"
)
