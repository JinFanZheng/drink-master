package contracts

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColumnModel(t *testing.T) {
	// 测试创建ColumnModel
	column := ColumnModel{
		Label: "Test Machine",
		Value: decimal.NewFromFloat(123.45),
	}

	assert.Equal(t, "Test Machine", column.Label)
	assert.True(t, column.Value.Equal(decimal.NewFromFloat(123.45)))

	// 测试JSON序列化
	data, err := json.Marshal(column)
	require.NoError(t, err)

	var unmarshaled ColumnModel
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, column.Label, unmarshaled.Label)
	assert.True(t, column.Value.Equal(unmarshaled.Value))
}

func TestColumnModel_ZeroValue(t *testing.T) {
	// 测试零值情况
	column := ColumnModel{
		Label: "Empty Machine",
		Value: decimal.NewFromInt(0),
	}

	assert.Equal(t, "Empty Machine", column.Label)
	assert.True(t, column.Value.IsZero())

	// 测试JSON序列化
	data, err := json.Marshal(column)
	require.NoError(t, err)

	expected := `{"label":"Empty Machine","value":"0"}`
	assert.JSONEq(t, expected, string(data))
}

func TestGetSalesRequest(t *testing.T) {
	// 测试无日期的请求
	req := GetSalesRequest{}
	assert.Nil(t, req.DateTime)

	// 测试有日期的请求
	now := time.Now()
	req = GetSalesRequest{
		DateTime: &now,
	}
	assert.NotNil(t, req.DateTime)
	assert.Equal(t, now, *req.DateTime)
}

func TestSalesResponse(t *testing.T) {
	date := time.Date(2025, 8, 11, 0, 0, 0, 0, time.UTC)
	sales := []ColumnModel{
		{
			Label: "Machine 1",
			Value: decimal.NewFromFloat(15.50),
		},
		{
			Label: "Machine 2",
			Value: decimal.NewFromFloat(25.00),
		},
	}
	total := decimal.NewFromFloat(40.50)

	response := SalesResponse{
		Date:  date,
		Sales: sales,
		Total: total,
	}

	assert.Equal(t, date, response.Date)
	assert.Len(t, response.Sales, 2)
	assert.True(t, response.Total.Equal(total))

	// 测试JSON序列化
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var unmarshaled SalesResponse
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, response.Date.Unix(), unmarshaled.Date.Unix())
	assert.Len(t, unmarshaled.Sales, 2)
	assert.True(t, response.Total.Equal(unmarshaled.Total))

	// 验证销售数据
	for i, sale := range response.Sales {
		assert.Equal(t, sale.Label, unmarshaled.Sales[i].Label)
		assert.True(t, sale.Value.Equal(unmarshaled.Sales[i].Value))
	}
}

func TestMachineOwnerSalesResponse(t *testing.T) {
	sales := []ColumnModel{
		{
			Label: "Test Machine",
			Value: decimal.NewFromFloat(100.00),
		},
	}

	meta := &Meta{
		Timestamp: time.Now(),
		RequestID: "test-request-id",
	}

	response := MachineOwnerSalesResponse{
		Success: true,
		Data:    sales,
		Meta:    meta,
	}

	assert.True(t, response.Success)
	assert.Len(t, response.Data, 1)
	assert.NotNil(t, response.Meta)
	assert.Nil(t, response.Error)

	// 测试JSON序列化
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var unmarshaled MachineOwnerSalesResponse
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, response.Success, unmarshaled.Success)
	assert.Len(t, unmarshaled.Data, 1)
	assert.Equal(t, response.Data[0].Label, unmarshaled.Data[0].Label)
	assert.True(t, response.Data[0].Value.Equal(unmarshaled.Data[0].Value))
}

func TestMachineOwnerSalesResponse_WithError(t *testing.T) {
	apiError := &APIError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "服务器内部错误",
	}

	response := MachineOwnerSalesResponse{
		Success: false,
		Error:   apiError,
	}

	assert.False(t, response.Success)
	assert.Nil(t, response.Data)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Error.Code)
	assert.Equal(t, "服务器内部错误", response.Error.Message)

	// 测试JSON序列化
	data, err := json.Marshal(response)
	require.NoError(t, err)

	var unmarshaled MachineOwnerSalesResponse
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, response.Success, unmarshaled.Success)
	assert.Equal(t, response.Error.Code, unmarshaled.Error.Code)
	assert.Equal(t, response.Error.Message, unmarshaled.Error.Message)
}

func TestDecimalJSONSerialization(t *testing.T) {
	// 测试不同精度的decimal序列化
	testCases := []struct {
		name     string
		value    decimal.Decimal
		expected string
	}{
		{
			name:     "整数",
			value:    decimal.NewFromInt(100),
			expected: "100",
		},
		{
			name:     "小数",
			value:    decimal.NewFromFloat(123.45),
			expected: "123.45",
		},
		{
			name:     "零值",
			value:    decimal.NewFromInt(0),
			expected: "0",
		},
		{
			name:     "高精度小数",
			value:    func() decimal.Decimal { d, _ := decimal.NewFromString("123.456789"); return d }(),
			expected: "123.456789",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			column := ColumnModel{
				Label: "Test",
				Value: tc.value,
			}

			data, err := json.Marshal(column)
			require.NoError(t, err)

			var unmarshaled ColumnModel
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.True(t, tc.value.Equal(unmarshaled.Value),
				"Expected %s, got %s", tc.value.String(), unmarshaled.Value.String())
		})
	}
}
