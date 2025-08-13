package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ddteam/drink-master/internal/enums"
)

func TestMaterialSilo_GetSaleStatusDesc(t *testing.T) {
	silo := &MaterialSilo{SaleStatus: enums.SaleStatusOn}
	assert.Equal(t, "在售", silo.GetSaleStatusDesc())

	silo.SaleStatus = enums.SaleStatusOff
	assert.Equal(t, "停售", silo.GetSaleStatusDesc())
}

func TestMaterialSilo_GetSaleStatusAPIString(t *testing.T) {
	silo := &MaterialSilo{SaleStatus: enums.SaleStatusOn}
	assert.Equal(t, "On", silo.GetSaleStatusAPIString())

	silo.SaleStatus = enums.SaleStatusOff
	assert.Equal(t, "Off", silo.GetSaleStatusAPIString())
}

func TestMaterialSilo_IsStockLow(t *testing.T) {
	tests := []struct {
		name        string
		stock       int
		maxCapacity int
		expected    bool
	}{
		{
			name:        "stock is low (5% of capacity)",
			stock:       5,
			maxCapacity: 100,
			expected:    true,
		},
		{
			name:        "stock is not low (50% of capacity)",
			stock:       50,
			maxCapacity: 100,
			expected:    false,
		},
		{
			name:        "stock is exactly 10% of capacity",
			stock:       10,
			maxCapacity: 100,
			expected:    false,
		},
		{
			name:        "stock is slightly below 10% threshold",
			stock:       9,
			maxCapacity: 100,
			expected:    true,
		},
		{
			name:        "max capacity is zero",
			stock:       10,
			maxCapacity: 0,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			silo := &MaterialSilo{
				Stock:       tt.stock,
				MaxCapacity: tt.maxCapacity,
			}
			assert.Equal(t, tt.expected, silo.IsStockLow())
		})
	}
}

func TestMaterialSilo_IsStockEmpty(t *testing.T) {
	silo := &MaterialSilo{Stock: 0}
	assert.True(t, silo.IsStockEmpty())

	silo.Stock = 1
	assert.False(t, silo.IsStockEmpty())
}

func TestMaterialSilo_IsStockFull(t *testing.T) {
	silo := &MaterialSilo{Stock: 100, MaxCapacity: 100}
	assert.True(t, silo.IsStockFull())

	silo.Stock = 99
	assert.False(t, silo.IsStockFull())

	// Stock can exceed capacity in edge cases
	silo.Stock = 101
	assert.True(t, silo.IsStockFull())
}

func TestMaterialSilo_GetStockPercentage(t *testing.T) {
	tests := []struct {
		name        string
		stock       int
		maxCapacity int
		expected    float64
	}{
		{
			name:        "50% stock",
			stock:       50,
			maxCapacity: 100,
			expected:    50.0,
		},
		{
			name:        "100% stock",
			stock:       100,
			maxCapacity: 100,
			expected:    100.0,
		},
		{
			name:        "0% stock",
			stock:       0,
			maxCapacity: 100,
			expected:    0.0,
		},
		{
			name:        "zero max capacity",
			stock:       50,
			maxCapacity: 0,
			expected:    0.0,
		},
		{
			name:        "stock exceeds capacity",
			stock:       150,
			maxCapacity: 100,
			expected:    150.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			silo := &MaterialSilo{
				Stock:       tt.stock,
				MaxCapacity: tt.maxCapacity,
			}
			assert.Equal(t, tt.expected, silo.GetStockPercentage())
		})
	}
}

func TestMaterialSilo_CanSale(t *testing.T) {
	productID := "product_123"

	tests := []struct {
		name       string
		productID  *string
		stock      int
		saleStatus enums.SaleStatus
		expected   bool
	}{
		{
			name:       "can sale - has product, stock, and sale status on",
			productID:  &productID,
			stock:      10,
			saleStatus: enums.SaleStatusOn,
			expected:   true,
		},
		{
			name:       "cannot sale - no product",
			productID:  nil,
			stock:      10,
			saleStatus: enums.SaleStatusOn,
			expected:   false,
		},
		{
			name:       "cannot sale - no stock",
			productID:  &productID,
			stock:      0,
			saleStatus: enums.SaleStatusOn,
			expected:   false,
		},
		{
			name:       "cannot sale - sale status off",
			productID:  &productID,
			stock:      10,
			saleStatus: enums.SaleStatusOff,
			expected:   false,
		},
		{
			name:       "cannot sale - multiple conditions fail",
			productID:  nil,
			stock:      0,
			saleStatus: enums.SaleStatusOff,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			silo := &MaterialSilo{
				ProductID:  tt.productID,
				Stock:      tt.stock,
				SaleStatus: tt.saleStatus,
			}
			assert.Equal(t, tt.expected, silo.CanSale())
		})
	}
}

func TestMaterialSilo_UpdateStock(t *testing.T) {
	tests := []struct {
		name        string
		maxCapacity int
		newStock    int
		expectError bool
		expectedErr error
	}{
		{
			name:        "valid stock update",
			maxCapacity: 100,
			newStock:    50,
			expectError: false,
		},
		{
			name:        "update to zero stock",
			maxCapacity: 100,
			newStock:    0,
			expectError: false,
		},
		{
			name:        "update to max capacity",
			maxCapacity: 100,
			newStock:    100,
			expectError: false,
		},
		{
			name:        "negative stock",
			maxCapacity: 100,
			newStock:    -1,
			expectError: true,
			expectedErr: ErrInvalidStock,
		},
		{
			name:        "stock exceeds capacity",
			maxCapacity: 100,
			newStock:    101,
			expectError: true,
			expectedErr: ErrStockExceedsCapacity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			silo := &MaterialSilo{MaxCapacity: tt.maxCapacity}
			err := silo.UpdateStock(tt.newStock)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newStock, silo.Stock)
			}
		})
	}
}

func TestMaterialSilo_TableName(t *testing.T) {
	silo := &MaterialSilo{}
	assert.Equal(t, "material_silos", silo.TableName())
}

func TestMaterialSilo_ModelStructure(t *testing.T) {
	productID := "product_123"
	now := time.Now()

	silo := MaterialSilo{
		ID:          "silo_123",
		MachineID:   "machine_123",
		SiloNo:      1,
		ProductID:   &productID,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  enums.SaleStatusOn,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
	}

	assert.Equal(t, "silo_123", silo.ID)
	assert.Equal(t, "machine_123", silo.MachineID)
	assert.Equal(t, 1, silo.SiloNo)
	assert.NotNil(t, silo.ProductID)
	assert.Equal(t, "product_123", *silo.ProductID)
	assert.Equal(t, 50, silo.Stock)
	assert.Equal(t, 100, silo.MaxCapacity)
	assert.Equal(t, enums.SaleStatusOn, silo.SaleStatus)
	assert.Equal(t, now, silo.CreatedAt)
	assert.Equal(t, now, silo.UpdatedAt)
	assert.Nil(t, silo.DeletedAt)
}

func TestMaterialSiloErrors(t *testing.T) {
	assert.Equal(t, "invalid stock: stock cannot be negative", ErrInvalidStock.Error())
	assert.Equal(t, "stock exceeds max capacity", ErrStockExceedsCapacity.Error())
}
