package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/models"
)

func setupTestDBForProduct(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 迁移表结构
	err = db.AutoMigrate(&models.Product{}, &models.Machine{}, &models.MachineProductPrice{})
	assert.NoError(t, err)

	return db
}

func TestProductHandler_GetSelectList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDBForProduct(t)

	// 创建测试数据
	product1 := &models.Product{
		ID:        "prod-1",
		Name:      "可乐",
		Status:    1,
		Price:     5.00,
		CreatedOn: time.Now(),
	}

	product2 := &models.Product{
		ID:        "prod-2",
		Name:      "橙汁",
		Status:    1,
		Price:     6.00,
		CreatedOn: time.Now(),
	}

	machine := &models.Machine{
		ID:        "machine-1",
		Name:      stringPtr("测试机器"),
		CreatedOn: time.Now(),
	}

	// 保存测试数据
	assert.NoError(t, db.Create(product1).Error)
	assert.NoError(t, db.Create(product2).Error)
	assert.NoError(t, db.Create(machine).Error)

	// 创建机器产品价格
	mpp1 := &models.MachineProductPrice{
		ID:              "mpp-1",
		MachineId:       machine.ID,
		ProductId:       product1.ID,
		Price:           3.5,
		PriceWithoutCup: 3.0,
		CreatedOn:       time.Now(),
	}

	mpp2 := &models.MachineProductPrice{
		ID:              "mpp-2",
		MachineId:       machine.ID,
		ProductId:       product2.ID,
		Price:           4.0,
		PriceWithoutCup: 3.5,
		CreatedOn:       time.Now(),
	}

	assert.NoError(t, db.Create(mpp1).Error)
	assert.NoError(t, db.Create(mpp2).Error)

	// 创建处理器
	handler := NewProductHandler(db)

	// 创建请求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/api/Product/GetSelectList", nil)
	c.Request = req

	// 执行测试
	handler.GetSelectList(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response contracts.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// 验证数据
	dataBytes, err := json.Marshal(response.Data)
	assert.NoError(t, err)

	var products []contracts.SelectViewModel
	err = json.Unmarshal(dataBytes, &products)
	assert.NoError(t, err)
	assert.Len(t, products, 2)

	// 验证产品数据
	productNames := make(map[string]float64)
	for _, p := range products {
		productNames[p.Name] = p.Price
	}

	assert.Contains(t, productNames, "可乐")
	assert.Contains(t, productNames, "橙汁")
	assert.Equal(t, 3.5, productNames["可乐"])
	assert.Equal(t, 4.0, productNames["橙汁"])
}

func TestProductHandler_GetSelectList_EmptyResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDBForProduct(t)
	handler := NewProductHandler(db)

	// 创建请求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/api/Product/GetSelectList", nil)
	c.Request = req

	// 执行测试
	handler.GetSelectList(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response contracts.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// 验证空数据
	dataBytes, err := json.Marshal(response.Data)
	assert.NoError(t, err)

	var products []contracts.SelectViewModel
	err = json.Unmarshal(dataBytes, &products)
	assert.NoError(t, err)
	assert.Len(t, products, 0)
}
