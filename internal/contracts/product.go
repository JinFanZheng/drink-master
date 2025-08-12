package contracts

// SelectViewModel 产品选择列表视图模型 (对应VendingMachine SelectViewModel)
type SelectViewModel struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ProductResponse 产品响应结构
type ProductResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category *string `json:"category"`
}

// ProductSelectListResponse 产品选择列表响应
type ProductSelectListResponse struct {
	Products []SelectViewModel `json:"products"`
	Meta     *Meta             `json:"meta,omitempty"`
}
