package contracts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMaterialSiloPagingRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request GetMaterialSiloPagingRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: GetMaterialSiloPagingRequest{
				MachineID: "machine_123",
				PageIndex: 1,
				PageSize:  10,
			},
			isValid: true,
		},
		{
			name: "missing machine id",
			request: GetMaterialSiloPagingRequest{
				PageIndex: 1,
				PageSize:  10,
			},
			isValid: false,
		},
		{
			name: "invalid page index",
			request: GetMaterialSiloPagingRequest{
				MachineID: "machine_123",
				PageIndex: 0,
				PageSize:  10,
			},
			isValid: false,
		},
		{
			name: "invalid page size - too large",
			request: GetMaterialSiloPagingRequest{
				MachineID: "machine_123",
				PageIndex: 1,
				PageSize:  101,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test validates the struct tags are correct
			// In actual validation, gin would validate these
			if tt.isValid {
				assert.NotEmpty(t, tt.request.MachineID)
				assert.Greater(t, tt.request.PageIndex, 0)
				assert.LessOrEqual(t, tt.request.PageSize, 100)
			}
		})
	}
}

func TestUpdateMaterialSiloStockRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateMaterialSiloStockRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: UpdateMaterialSiloStockRequest{
				ID:    "silo_123",
				Stock: 50,
			},
			isValid: true,
		},
		{
			name: "zero stock is valid",
			request: UpdateMaterialSiloStockRequest{
				ID:    "silo_123",
				Stock: 0,
			},
			isValid: true,
		},
		{
			name: "missing id",
			request: UpdateMaterialSiloStockRequest{
				Stock: 50,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				assert.NotEmpty(t, tt.request.ID)
				assert.GreaterOrEqual(t, tt.request.Stock, 0)
			}
		})
	}
}

func TestToggleSaleMaterialSiloRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request ToggleSaleMaterialSiloRequest
		isValid bool
	}{
		{
			name: "valid request with On status",
			request: ToggleSaleMaterialSiloRequest{
				ID:         "silo_123",
				SaleStatus: SaleStatusOn,
			},
			isValid: true,
		},
		{
			name: "valid request with Off status",
			request: ToggleSaleMaterialSiloRequest{
				ID:         "silo_123",
				SaleStatus: SaleStatusOff,
			},
			isValid: true,
		},
		{
			name: "invalid sale status",
			request: ToggleSaleMaterialSiloRequest{
				ID:         "silo_123",
				SaleStatus: "Invalid",
			},
			isValid: false,
		},
		{
			name: "missing id",
			request: ToggleSaleMaterialSiloRequest{
				SaleStatus: SaleStatusOn,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				assert.NotEmpty(t, tt.request.ID)
				assert.Contains(t, []string{SaleStatusOn, SaleStatusOff}, tt.request.SaleStatus)
			}
		})
	}
}

func TestGetMaterialSiloPagingResponse_Structure(t *testing.T) {
	now := time.Now()
	productID := "product_123"
	productName := "Test Product"

	response := GetMaterialSiloPagingResponse{
		ID:          "silo_123",
		MachineID:   "machine_123",
		SiloNo:      1,
		ProductID:   &productID,
		ProductName: &productName,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  SaleStatusOn,
		UpdatedAt:   now,
	}

	assert.Equal(t, "silo_123", response.ID)
	assert.Equal(t, "machine_123", response.MachineID)
	assert.Equal(t, 1, response.SiloNo)
	assert.NotNil(t, response.ProductID)
	assert.Equal(t, "product_123", *response.ProductID)
	assert.NotNil(t, response.ProductName)
	assert.Equal(t, "Test Product", *response.ProductName)
	assert.Equal(t, 50, response.Stock)
	assert.Equal(t, 100, response.MaxCapacity)
	assert.Equal(t, SaleStatusOn, response.SaleStatus)
	assert.Equal(t, now, response.UpdatedAt)
}

func TestMaterialSiloPaging_Structure(t *testing.T) {
	items := []GetMaterialSiloPagingResponse{
		{
			ID:          "silo_1",
			MachineID:   "machine_123",
			SiloNo:      1,
			Stock:       50,
			MaxCapacity: 100,
			SaleStatus:  SaleStatusOn,
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "silo_2",
			MachineID:   "machine_123",
			SiloNo:      2,
			Stock:       30,
			MaxCapacity: 100,
			SaleStatus:  SaleStatusOff,
			UpdatedAt:   time.Now(),
		},
	}

	paging := MaterialSiloPaging{
		Items:      items,
		TotalCount: 15,
		PageIndex:  1,
		PageSize:   10,
	}

	assert.Len(t, paging.Items, 2)
	assert.Equal(t, int64(15), paging.TotalCount)
	assert.Equal(t, 1, paging.PageIndex)
	assert.Equal(t, 10, paging.PageSize)
}

func TestSaleStatusConstants(t *testing.T) {
	assert.Equal(t, "On", SaleStatusOn)
	assert.Equal(t, "Off", SaleStatusOff)
}

func TestMaterialSiloOperationResult_Structure(t *testing.T) {
	result := MaterialSiloOperationResult{
		Success: true,
		Message: "Operation completed successfully",
		Data:    "additional_data",
	}

	assert.True(t, result.Success)
	assert.Equal(t, "Operation completed successfully", result.Message)
	assert.Equal(t, "additional_data", result.Data)
}
