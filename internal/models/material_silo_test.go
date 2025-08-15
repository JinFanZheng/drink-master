package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaterialSilo_IsStockEmpty(t *testing.T) {
	silo := &MaterialSilo{Stock: 0}
	assert.True(t, silo.IsStockEmpty())

	silo.Stock = 1
	assert.False(t, silo.IsStockEmpty())
}

func TestMaterialSilo_IsStockLow(t *testing.T) {
	tests := []struct {
		name     string
		stock    int
		total    int
		expected bool
	}{
		{
			name:     "stock is low (5% of total)",
			stock:    5,
			total:    100,
			expected: true,
		},
		{
			name:     "stock is not low (50% of total)",
			stock:    50,
			total:    100,
			expected: false,
		},
		{
			name:     "stock is exactly 10% of total",
			stock:    10,
			total:    100,
			expected: false,
		},
		{
			name:     "stock is slightly below 10% threshold",
			stock:    9,
			total:    100,
			expected: true,
		},
		{
			name:     "total is zero",
			stock:    10,
			total:    0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			silo := &MaterialSilo{
				Stock: tt.stock,
				Total: tt.total,
			}
			assert.Equal(t, tt.expected, silo.IsStockLow())
		})
	}
}

func TestMaterialSilo_IsStockFull(t *testing.T) {
	silo := &MaterialSilo{Stock: 100, Total: 100}
	assert.True(t, silo.IsStockFull())

	silo.Stock = 99
	assert.False(t, silo.IsStockFull())

	// Stock can exceed total in edge cases
	silo.Stock = 101
	assert.True(t, silo.IsStockFull())
}

func TestMaterialSilo_GetStockPercentage(t *testing.T) {
	tests := []struct {
		name     string
		stock    int
		total    int
		expected float64
	}{
		{
			name:     "50% stock",
			stock:    50,
			total:    100,
			expected: 50.0,
		},
		{
			name:     "100% stock",
			stock:    100,
			total:    100,
			expected: 100.0,
		},
		{
			name:     "0% stock",
			stock:    0,
			total:    100,
			expected: 0.0,
		},
		{
			name:     "zero total",
			stock:    50,
			total:    0,
			expected: 0.0,
		},
		{
			name:     "stock exceeds total",
			stock:    150,
			total:    100,
			expected: 150.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			silo := &MaterialSilo{
				Stock: tt.stock,
				Total: tt.total,
			}
			assert.Equal(t, tt.expected, silo.GetStockPercentage())
		})
	}
}

func TestMaterialSilo_CanSale(t *testing.T) {
	productID := "product_123"

	tests := []struct {
		name      string
		productId *string
		stock     int
		isSale    BitBool
		expected  bool
	}{
		{
			name:      "can sale - has product, stock, and sale enabled",
			productId: &productID,
			stock:     10,
			isSale:    BitBool(1),
			expected:  true,
		},
		{
			name:      "cannot sale - no product",
			productId: nil,
			stock:     10,
			isSale:    BitBool(1),
			expected:  false,
		},
		{
			name:      "cannot sale - no stock",
			productId: &productID,
			stock:     0,
			isSale:    BitBool(1),
			expected:  false,
		},
		{
			name:      "cannot sale - sale disabled",
			productId: &productID,
			stock:     10,
			isSale:    BitBool(0),
			expected:  false,
		},
		{
			name:      "cannot sale - multiple conditions fail",
			productId: nil,
			stock:     0,
			isSale:    BitBool(0),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			silo := &MaterialSilo{
				ProductId: tt.productId,
				Stock:     tt.stock,
				IsSale:    tt.isSale,
			}
			assert.Equal(t, tt.expected, silo.CanSale())
		})
	}
}

func TestMaterialSilo_UpdateStock(t *testing.T) {
	tests := []struct {
		name        string
		total       int
		newStock    int
		expectError bool
		expectedErr error
	}{
		{
			name:        "valid stock update",
			total:       100,
			newStock:    50,
			expectError: false,
		},
		{
			name:        "update to zero stock",
			total:       100,
			newStock:    0,
			expectError: false,
		},
		{
			name:        "update to total capacity",
			total:       100,
			newStock:    100,
			expectError: false,
		},
		{
			name:        "negative stock",
			total:       100,
			newStock:    -1,
			expectError: true,
			expectedErr: ErrInvalidStock,
		},
		{
			name:        "stock exceeds total",
			total:       100,
			newStock:    101,
			expectError: true,
			expectedErr: ErrStockExceedsCapacity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			silo := &MaterialSilo{Total: tt.total}
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
	productId := "product_123"
	machineId := "machine_123"
	siloNo := "01"
	now := time.Now()

	silo := MaterialSilo{
		ID:         "silo_123",
		MachineId:  &machineId,
		No:         &siloNo,
		Type:       1,
		ProductId:  &productId,
		IsSale:     BitBool(1),
		Total:      100,
		Stock:      50,
		SingleFeed: 1,
		Version:    1,
		CreatedOn:  now,
		UpdatedOn:  &now,
	}

	assert.Equal(t, "silo_123", silo.ID)
	assert.NotNil(t, silo.MachineId)
	assert.Equal(t, "machine_123", *silo.MachineId)
	assert.NotNil(t, silo.No)
	assert.Equal(t, "01", *silo.No)
	assert.Equal(t, 1, silo.Type)
	assert.NotNil(t, silo.ProductId)
	assert.Equal(t, "product_123", *silo.ProductId)
	assert.Equal(t, BitBool(1), silo.IsSale)
	assert.Equal(t, 100, silo.Total)
	assert.Equal(t, 50, silo.Stock)
	assert.Equal(t, 1, silo.SingleFeed)
	assert.Equal(t, int64(1), silo.Version)
	assert.Equal(t, now, silo.CreatedOn)
	assert.NotNil(t, silo.UpdatedOn)
}

func TestMaterialSiloErrors(t *testing.T) {
	assert.Equal(t, "invalid stock: stock cannot be negative", ErrInvalidStock.Error())
	assert.Equal(t, "stock exceeds max capacity", ErrStockExceedsCapacity.Error())
}
