package contracts

import "time"

// 通用API契约定义

// APIResponse 通用API响应结构
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// Meta 响应元数据
type Meta struct {
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	RequestID string    `json:"request_id" example:"req_123456"`
	Version   string    `json:"version" example:"v1.0.0"`
	Warnings  []string  `json:"warnings,omitempty"`
}

// PaginationMeta 分页元数据
type PaginationMeta struct {
	Total       int64 `json:"total" example:"100"`         // 总记录数
	Count       int   `json:"count" example:"10"`          // 当前页记录数
	PerPage     int   `json:"per_page" example:"10"`       // 每页记录数
	CurrentPage int   `json:"current_page" example:"1"`    // 当前页码
	TotalPages  int   `json:"total_pages" example:"10"`    // 总页数
	HasNext     bool  `json:"has_next" example:"true"`     // 是否有下一页
	HasPrev     bool  `json:"has_prev" example:"false"`    // 是否有上一页
	*Meta              // 继承通用元数据
}

// APIError 错误响应结构
type APIError struct {
	Code       string                 `json:"code" example:"VALIDATION_ERROR"`
	Message    string                 `json:"message" example:"请求参数验证失败"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Timestamp  time.Time              `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Path       string                 `json:"path" example:"/api/drinks"`
	Method     string                 `json:"method" example:"POST"`
	RequestID  string                 `json:"request_id" example:"req_123456"`
}

// ValidationError 验证错误详情
type ValidationError struct {
	Field   string      `json:"field" example:"name"`
	Value   interface{} `json:"value" example:""`
	Message string      `json:"message" example:"名称不能为空"`
	Tag     string      `json:"tag" example:"required"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string                 `json:"status" example:"ok"`
	Timestamp time.Time              `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Version   string                 `json:"version" example:"1.0.0"`
	Services  map[string]ServiceHealth `json:"services"`
}

// ServiceHealth 服务健康状态
type ServiceHealth struct {
	Status      string    `json:"status" example:"ok"`         // ok, degraded, down
	Latency     string    `json:"latency,omitempty" example:"2ms"`
	LastChecked time.Time `json:"last_checked" example:"2023-01-01T00:00:00Z"`
	Error       string    `json:"error,omitempty"`
}

// DrinkCategory 饮品分类
type DrinkCategory struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"coffee"`
	DisplayName string    `json:"display_name" example:"咖啡"`
	Description string    `json:"description" example:"各类咖啡饮品"`
	Icon        string    `json:"icon,omitempty" example:"☕"`
	Color       string    `json:"color,omitempty" example:"#8B4513"`
	DrinkCount  int       `json:"drink_count" example:"25"`    // 该分类下的饮品数量
	CreateAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdateAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CategoryListResponse 分类列表响应
type CategoryListResponse struct {
	Categories []DrinkCategory `json:"categories"`
	Meta       *Meta           `json:"meta,omitempty"`
}

// MessageResponse 简单消息响应
type MessageResponse struct {
	Message string `json:"message" example:"操作成功"`
}

// IDResponse ID响应（用于创建操作）
type IDResponse struct {
	ID      uint   `json:"id" example:"1"`
	Message string `json:"message" example:"创建成功"`
}

// BulkOperationRequest 批量操作请求
type BulkOperationRequest struct {
	IDs    []uint `json:"ids" validate:"required,min=1" example:"[1,2,3]"`
	Action string `json:"action" validate:"required" example:"delete"`
}

// BulkOperationResponse 批量操作响应
type BulkOperationResponse struct {
	Successful []uint   `json:"successful" example:"[1,2]"`
	Failed     []uint   `json:"failed" example:"[3]"`
	Errors     []string `json:"errors,omitempty" example:"[\"记录不存在: ID 3\"]"`
	Total      int      `json:"total" example:"3"`
	Success    int      `json:"success" example:"2"`
	Failure    int      `json:"failure" example:"1"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query    string   `form:"q" example:"latte"`                               // 搜索关键词
	Fields   []string `form:"fields" example:"name,description"`               // 搜索字段
	Category string   `form:"category" example:"coffee"`                       // 分类过滤
	Tags     []string `form:"tags" example:"hot,sweet"`                        // 标签过滤
	Limit    int      `form:"limit" validate:"min=1,max=100" example:"10"`     // 结果数量
	Offset   int      `form:"offset" validate:"min=0" example:"0"`             // 偏移量
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Query   string          `json:"query" example:"latte"`
	Results []DrinkResponse `json:"results"`
	Meta    PaginationMeta  `json:"meta"`
}

// ExportRequest 导出请求
type ExportRequest struct {
	Format    string   `form:"format" validate:"oneof=csv json xlsx" example:"csv"`
	StartDate string   `form:"start_date" example:"2023-01-01"`
	EndDate   string   `form:"end_date" example:"2023-01-31"`
	Category  string   `form:"category" example:"coffee"`
	Fields    []string `form:"fields" example:"name,price,created_at"`
}

// ExportResponse 导出响应
type ExportResponse struct {
	FileURL   string    `json:"file_url" example:"https://example.com/exports/drinks_2023-01.csv"`
	Format    string    `json:"format" example:"csv"`
	Records   int       `json:"records" example:"150"`
	ExpiresAt time.Time `json:"expires_at" example:"2023-01-08T00:00:00Z"`
}

// 常见的HTTP状态码常量
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusInternalServerError = 500
)

// 常见的错误码常量
const (
	ErrorCodeValidation      = "VALIDATION_ERROR"
	ErrorCodeNotFound        = "NOT_FOUND"
	ErrorCodeUnauthorized    = "UNAUTHORIZED"
	ErrorCodeForbidden       = "FORBIDDEN"
	ErrorCodeConflict        = "CONFLICT"
	ErrorCodeInternalServer  = "INTERNAL_SERVER_ERROR"
	ErrorCodeDatabaseError   = "DATABASE_ERROR"
	ErrorCodeInvalidToken    = "INVALID_TOKEN"
	ErrorCodeTokenExpired    = "TOKEN_EXPIRED"
	ErrorCodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
)