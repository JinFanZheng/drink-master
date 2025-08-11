package contracts

import "time"

// 饮品相关的API契约定义

// DrinkCreateRequest 创建饮品请求
type DrinkCreateRequest struct {
	Name        string  `json:"name" validate:"required,max=100" example:"拿铁咖啡"`
	Category    string  `json:"category" validate:"required" example:"coffee"`
	Price       float64 `json:"price" validate:"required,gt=0" example:"25.5"`
	Description string  `json:"description,omitempty" validate:"max=500" example:"香浓的拿铁咖啡"`
	ImageURL    string  `json:"image_url,omitempty" validate:"url" example:"https://example.com/latte.jpg"`
}

// DrinkUpdateRequest 更新饮品请求
type DrinkUpdateRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,max=100" example:"拿铁咖啡"`
	Category    *string  `json:"category,omitempty" example:"coffee"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0" example:"25.5"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500" example:"香浓的拿铁咖啡"`
	ImageURL    *string  `json:"image_url,omitempty" validate:"omitempty,url" example:"https://example.com/latte.jpg"`
}

// DrinkResponse 饮品响应
type DrinkResponse struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"拿铁咖啡"`
	Category    string    `json:"category" example:"coffee"`
	Price       float64   `json:"price" example:"25.5"`
	Description string    `json:"description" example:"香浓的拿铁咖啡"`
	ImageURL    string    `json:"image_url" example:"https://example.com/latte.jpg"`
	UserID      uint      `json:"user_id" example:"1"`
	CreateAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdateAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// DrinkListRequest 饮品列表请求参数
type DrinkListRequest struct {
	Category string `form:"category" example:"coffee"`           // 按分类过滤
	Name     string `form:"name" example:"latte"`                // 按名称搜索
	MinPrice *float64 `form:"min_price" example:"10"`            // 最低价格
	MaxPrice *float64 `form:"max_price" example:"50"`            // 最高价格
	SortBy   string `form:"sort_by" example:"created_at"`        // 排序字段: name, price, created_at
	Order    string `form:"order" example:"desc"`                // 排序方向: asc, desc
	Limit    int    `form:"limit" validate:"min=1,max=100" example:"10"` // 每页数量
	Offset   int    `form:"offset" validate:"min=0" example:"0"`         // 偏移量
}

// DrinkListResponse 饮品列表响应
type DrinkListResponse struct {
	Drinks []DrinkResponse `json:"drinks"`
	Meta   PaginationMeta  `json:"meta"`
}

// ConsumptionLogRequest 消费记录请求
type ConsumptionLogRequest struct {
	DrinkID    uint      `json:"drink_id" validate:"required" example:"1"`
	Quantity   int       `json:"quantity" validate:"required,min=1" example:"2"`
	ConsumedAt time.Time `json:"consumed_at,omitempty" example:"2023-01-01T14:30:00Z"`
	Notes      string    `json:"notes,omitempty" validate:"max=200" example:"下午茶时光"`
}

// ConsumptionLogResponse 消费记录响应
type ConsumptionLogResponse struct {
	ID         uint          `json:"id" example:"1"`
	Drink      DrinkResponse `json:"drink"`
	Quantity   int           `json:"quantity" example:"2"`
	ConsumedAt time.Time     `json:"consumed_at" example:"2023-01-01T14:30:00Z"`
	Notes      string        `json:"notes" example:"下午茶时光"`
	CreateAt   time.Time     `json:"created_at" example:"2023-01-01T14:30:00Z"`
}

// ConsumptionLogListRequest 消费记录列表请求
type ConsumptionLogListRequest struct {
	DrinkID   *uint  `form:"drink_id" example:"1"`                    // 按饮品过滤
	StartDate string `form:"start_date" example:"2023-01-01"`         // 开始日期
	EndDate   string `form:"end_date" example:"2023-01-31"`           // 结束日期
	Limit     int    `form:"limit" validate:"min=1,max=100" example:"10"` // 每页数量
	Offset    int    `form:"offset" validate:"min=0" example:"0"`         // 偏移量
}

// ConsumptionLogListResponse 消费记录列表响应
type ConsumptionLogListResponse struct {
	Logs []ConsumptionLogResponse `json:"logs"`
	Meta PaginationMeta           `json:"meta"`
}