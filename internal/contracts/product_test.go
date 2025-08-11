package contracts

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectViewModel_JSON(t *testing.T) {
	// 创建测试数据
	vm := SelectViewModel{
		ID:    "test-id",
		Name:  "测试产品",
		Price: 9.99,
	}

	// 测试JSON序列化
	jsonData, err := json.Marshal(vm)
	assert.NoError(t, err)

	expectedJSON := `{"id":"test-id","name":"测试产品","price":9.99}`
	assert.JSONEq(t, expectedJSON, string(jsonData))

	// 测试JSON反序列化
	var deserializedVM SelectViewModel
	err = json.Unmarshal(jsonData, &deserializedVM)
	assert.NoError(t, err)
	assert.Equal(t, vm, deserializedVM)
}

func TestProductResponse_JSON(t *testing.T) {
	category := "饮料"
	response := ProductResponse{
		ID:       "prod-123",
		Name:     "产品名称",
		Price:    15.50,
		Category: &category,
	}

	// 测试JSON序列化
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	expectedJSON := `{"id":"prod-123","name":"产品名称","price":15.5,"category":"饮料"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))

	// 测试JSON反序列化
	var deserializedResponse ProductResponse
	err = json.Unmarshal(jsonData, &deserializedResponse)
	assert.NoError(t, err)
	assert.Equal(t, response, deserializedResponse)
}

func TestProductResponse_NullCategory(t *testing.T) {
	response := ProductResponse{
		ID:       "prod-456",
		Name:     "无分类产品",
		Price:    5.00,
		Category: nil,
	}

	// 测试JSON序列化
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	expectedJSON := `{"id":"prod-456","name":"无分类产品","price":5,"category":null}`
	assert.JSONEq(t, expectedJSON, string(jsonData))

	// 测试JSON反序列化
	var deserializedResponse ProductResponse
	err = json.Unmarshal(jsonData, &deserializedResponse)
	assert.NoError(t, err)
	assert.Equal(t, response, deserializedResponse)
	assert.Nil(t, deserializedResponse.Category)
}

func TestProductSelectListResponse_JSON(t *testing.T) {
	products := []SelectViewModel{
		{ID: "1", Name: "产品1", Price: 10.0},
		{ID: "2", Name: "产品2", Price: 20.0},
	}

	response := ProductSelectListResponse{
		Products: products,
		Meta:     nil,
	}

	// 测试JSON序列化
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	// 测试JSON反序列化
	var deserializedResponse ProductSelectListResponse
	err = json.Unmarshal(jsonData, &deserializedResponse)
	assert.NoError(t, err)
	assert.Equal(t, response, deserializedResponse)
	assert.Len(t, deserializedResponse.Products, 2)
	assert.Equal(t, "产品1", deserializedResponse.Products[0].Name)
	assert.Equal(t, "产品2", deserializedResponse.Products[1].Name)
}
