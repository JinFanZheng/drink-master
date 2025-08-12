package contracts

// Machine API契约定义 - 对应MobileAPI MachineController

// GetMachinePagingRequest 获取售货机分页列表请求
type GetMachinePagingRequest struct {
	Page           int    `json:"page" binding:"required,min=1"`
	PageSize       int    `json:"pageSize" binding:"required,min=1,max=100"`
	MachineOwnerID string `json:"machineOwnerId"` // 从JWT token获取
	Keyword        string `json:"keyword"`        // 搜索关键词
}

// GetMachinePagingResponse 获取售货机分页列表响应项
type GetMachinePagingResponse struct {
	ID                 string `json:"id"`
	MachineNo          string `json:"machineNo"`
	Name               string `json:"name"`
	Area               string `json:"area"`
	Address            string `json:"address"`
	BusinessStatus     string `json:"businessStatus"`
	BusinessStatusDesc string `json:"businessStatusDesc"`
	DeviceID           string `json:"deviceId"`
}

// PagingResult 分页结果
type PagingResult struct {
	Items      []GetMachinePagingResponse `json:"items"`
	TotalCount int64                      `json:"totalCount"`
	PageIndex  int                        `json:"pageIndex"`
	PageSize   int                        `json:"pageSize"`
}

// GetMachineListResponse 获取售货机列表响应项（简化版）
type GetMachineListResponse struct {
	ID                 string `json:"id"`
	MachineNo          string `json:"machineNo"`
	Name               string `json:"name"`
	BusinessStatus     string `json:"businessStatus"`
	BusinessStatusDesc string `json:"businessStatusDesc"`
}

// GetMachineByIDResponse 根据ID获取售货机详情响应
type GetMachineByIDResponse struct {
	ID                 string `json:"id"`
	MachineNo          string `json:"machineNo"`
	Name               string `json:"name"`
	Area               string `json:"area"`
	Address            string `json:"address"`
	BusinessStatus     string `json:"businessStatus"` // Open, Close, Offline
	BusinessStatusDesc string `json:"businessStatusDesc"`
	DeviceID           string `json:"deviceId"`
	ServicePhone       string `json:"servicePhone"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

// CheckDeviceExistRequest 检查设备是否存在请求
type CheckDeviceExistRequest struct {
	DeviceID string `form:"deviceId" binding:"required"`
}

// CheckDeviceExistResponse 检查设备是否存在响应
type CheckDeviceExistResponse struct {
	Exists bool `json:"exists"`
}

// MachineProductResponse 售货机商品响应
type MachineProductResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	PriceWithoutCup float64 `json:"priceWithoutCup"`
	Stock           int     `json:"stock"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
}

// ProductListResponse 商品列表响应（基于VendingMachine逻辑）
type ProductListResponse struct {
	Name     string                   `json:"name"` // "限时巨惠"
	Products []MachineProductResponse `json:"products"`
}

// OpenOrCloseBusinessRequest 开关营业状态请求
type OpenOrCloseBusinessRequest struct {
	ID string `form:"id" binding:"required"`
}

// OpenOrCloseBusinessResponse 开关营业状态响应
type OpenOrCloseBusinessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// DeviceStatusCheckResult 设备状态检查结果
type DeviceStatusCheckResult struct {
	DeviceID string `json:"deviceId"`
	Online   bool   `json:"online"`
	LastSeen string `json:"lastSeen,omitempty"`
}

// BusinessStatusEnum 营业状态枚举
const (
	BusinessStatusOpen    = "Open"
	BusinessStatusClose   = "Close"
	BusinessStatusOffline = "Offline"
)

// ProductGroupName 商品分组名称（对应VendingMachine）
const (
	ProductGroupTimeLimited = "限时巨惠"
)
