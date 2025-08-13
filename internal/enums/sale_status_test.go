package enums

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaleStatus_GetSaleStatusDesc(t *testing.T) {
	tests := []struct {
		name     string
		status   SaleStatus
		expected string
	}{
		{
			name:     "SaleStatusOn",
			status:   SaleStatusOn,
			expected: "在售",
		},
		{
			name:     "SaleStatusOff",
			status:   SaleStatusOff,
			expected: "停售",
		},
		{
			name:     "Invalid status",
			status:   SaleStatus(99),
			expected: "未知状态",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSaleStatusDesc(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSaleStatus_String(t *testing.T) {
	status := SaleStatusOn
	assert.Equal(t, "在售", status.String())

	status = SaleStatusOff
	assert.Equal(t, "停售", status.String())
}

func TestSaleStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   SaleStatus
		expected bool
	}{
		{
			name:     "SaleStatusOn is valid",
			status:   SaleStatusOn,
			expected: true,
		},
		{
			name:     "SaleStatusOff is valid",
			status:   SaleStatusOff,
			expected: true,
		},
		{
			name:     "Invalid status below range",
			status:   SaleStatus(-1),
			expected: false,
		},
		{
			name:     "Invalid status above range",
			status:   SaleStatus(99),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSaleStatus_ToAPIString(t *testing.T) {
	tests := []struct {
		name     string
		status   SaleStatus
		expected string
	}{
		{
			name:     "SaleStatusOn to API string",
			status:   SaleStatusOn,
			expected: "On",
		},
		{
			name:     "SaleStatusOff to API string",
			status:   SaleStatusOff,
			expected: "Off",
		},
		{
			name:     "Invalid status defaults to Off",
			status:   SaleStatus(99),
			expected: "Off",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.ToAPIString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSaleStatusFromAPIString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SaleStatus
	}{
		{
			name:     "On string to SaleStatusOn",
			input:    "On",
			expected: SaleStatusOn,
		},
		{
			name:     "Off string to SaleStatusOff",
			input:    "Off",
			expected: SaleStatusOff,
		},
		{
			name:     "Invalid string defaults to SaleStatusOff",
			input:    "Invalid",
			expected: SaleStatusOff,
		},
		{
			name:     "Empty string defaults to SaleStatusOff",
			input:    "",
			expected: SaleStatusOff,
		},
		{
			name:     "Case sensitive - lowercase on",
			input:    "on",
			expected: SaleStatusOff, // Should default to Off for case sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SaleStatusFromAPIString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSaleStatus_Constants(t *testing.T) {
	assert.Equal(t, SaleStatus(0), SaleStatusOff)
	assert.Equal(t, SaleStatus(1), SaleStatusOn)
}
