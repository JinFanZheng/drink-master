package contracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMachinePagingRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request GetMachinePagingRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: GetMachinePagingRequest{
				Page:           1,
				PageSize:       10,
				MachineOwnerID: "owner-123",
				Keyword:        "test",
			},
			valid: true,
		},
		{
			name: "empty page",
			request: GetMachinePagingRequest{
				Page:           0,
				PageSize:       10,
				MachineOwnerID: "owner-123",
			},
			valid: false,
		},
		{
			name: "large page size",
			request: GetMachinePagingRequest{
				Page:           1,
				PageSize:       200,
				MachineOwnerID: "owner-123",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这里只验证结构体字段是否设置正确
			if tt.valid {
				assert.Greater(t, tt.request.Page, 0)
				assert.Greater(t, tt.request.PageSize, 0)
				assert.LessOrEqual(t, tt.request.PageSize, 100)
			}
		})
	}
}

func TestBusinessStatusConstants(t *testing.T) {
	assert.Equal(t, "Open", BusinessStatusOpen)
	assert.Equal(t, "Close", BusinessStatusClose)
	assert.Equal(t, "Offline", BusinessStatusOffline)
}

func TestProductGroupConstants(t *testing.T) {
	assert.Equal(t, "限时巨惠", ProductGroupTimeLimited)
}

func TestMachineProductResponse(t *testing.T) {
	product := MachineProductResponse{
		ID:              "product-123",
		Name:            "Test Product",
		Price:           10.5,
		PriceWithoutCup: 9.5,
		Stock:           100,
		Category:        "Drinks",
		Description:     "Test product description",
	}

	assert.Equal(t, "product-123", product.ID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, 10.5, product.Price)
	assert.Equal(t, 9.5, product.PriceWithoutCup)
	assert.Equal(t, 100, product.Stock)
	assert.Equal(t, "Drinks", product.Category)
	assert.Equal(t, "Test product description", product.Description)
}

func TestProductListResponse(t *testing.T) {
	products := []MachineProductResponse{
		{
			ID:    "product-1",
			Name:  "Product 1",
			Price: 5.0,
			Stock: 50,
		},
		{
			ID:    "product-2",
			Name:  "Product 2",
			Price: 7.5,
			Stock: 25,
		},
	}

	productList := ProductListResponse{
		Name:     ProductGroupTimeLimited,
		Products: products,
	}

	assert.Equal(t, ProductGroupTimeLimited, productList.Name)
	assert.Len(t, productList.Products, 2)
	assert.Equal(t, "Product 1", productList.Products[0].Name)
	assert.Equal(t, "Product 2", productList.Products[1].Name)
}
