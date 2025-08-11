package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/models"
)

func setupMachineOwnerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = models.AutoMigrate(db)
	require.NoError(t, err)

	return db
}

func setupMachineOwnerTestData(t *testing.T, db *gorm.DB) (string, string, string) {
	// 创建机主
	owner := models.MachineOwner{
		ID:   "owner-001",
		Name: "Test Owner",
	}
	require.NoError(t, db.Create(&owner).Error)

	// 创建普通用户
	member := models.Member{
		ID:           "member-001",
		Nickname:     "Test User",
		WeChatOpenId: "openid-001",
		Role:         "Member",
	}
	require.NoError(t, db.Create(&member).Error)

	// 创建机器
	machine1 := models.Machine{
		ID:             "machine-001",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM001",
		Name:           "Test Machine 1",
		BusinessStatus: "Open",
	}
	require.NoError(t, db.Create(&machine1).Error)

	machine2 := models.Machine{
		ID:             "machine-002",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM002",
		Name:           "Test Machine 2",
		BusinessStatus: "Open",
	}
	require.NoError(t, db.Create(&machine2).Error)

	// 创建商品
	product := models.Product{
		ID:   "product-001",
		Name: "Test Coffee",
	}
	require.NoError(t, db.Create(&product).Error)

	return owner.ID, machine1.ID, machine2.ID
}

func createMachineOwnerContext(machineOwnerID string, isMachineOwner bool) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// 设置JWT中间件设置的值
	if machineOwnerID != "" {
		c.Set("machine_owner_id", machineOwnerID)
	}
	if isMachineOwner {
		c.Set("role", "Owner")
	} else {
		c.Set("role", "Member")
	}

	return c, w
}

func TestNewMachineOwnerHandler(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.BaseHandler)
	assert.NotNil(t, handler.machineOwnerService)
}

func TestMachineOwnerHandler_GetSales_Success(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)
	ownerID, machineID1, machineID2 := setupMachineOwnerTestData(t, db)

	// 创建今天的订单数据
	today := time.Now().Truncate(24 * time.Hour)
	orders := []models.Order{
		{
			ID:            "order-001",
			MemberId:      "member-001",
			MachineId:     machineID1,
			ProductId:     "product-001",
			OrderNo:       "ON001",
			PayAmount:     15.50,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
		{
			ID:            "order-002",
			MemberId:      "member-001",
			MachineId:     machineID2,
			ProductId:     "product-001",
			OrderNo:       "ON002",
			PayAmount:     18.00,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
	}

	for _, order := range orders {
		require.NoError(t, db.Create(&order).Error)
	}

	// 创建请求
	c, w := createMachineOwnerContext(ownerID, true)

	// 调用处理器
	handler.GetSales(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response contracts.MachineOwnerSalesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)

	// 验证销售数据
	salesMap := make(map[string]decimal.Decimal)
	for _, sale := range response.Data {
		salesMap[sale.Label] = sale.Value
	}

	assert.True(t, salesMap["Test Machine 1"].Equal(decimal.NewFromFloat(15.50)))
	assert.True(t, salesMap["Test Machine 2"].Equal(decimal.NewFromFloat(18.00)))
}

func TestMachineOwnerHandler_GetSales_NotMachineOwner(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)

	// 创建普通用户上下文
	c, w := createMachineOwnerContext("", false)

	handler.GetSales(c)

	// 验证响应
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response contracts.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Contains(t, response.Error.Message, "您不是机主")
}

func TestMachineOwnerHandler_GetSales_NoMachineOwnerID(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)

	// 创建机主上下文但没有machine_owner_id
	c, w := createMachineOwnerContext("", true)

	handler.GetSales(c)

	// 验证响应
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response contracts.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Contains(t, response.Error.Message, "无效的机主信息")
}

func TestMachineOwnerHandler_GetSales_WithCustomDate(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)
	ownerID, machineID1, _ := setupMachineOwnerTestData(t, db)

	// 创建特定日期的订单数据
	customDate := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	order := models.Order{
		ID:            "order-custom",
		MemberId:      "member-001",
		MachineId:     machineID1,
		ProductId:     "product-001",
		OrderNo:       "ON-CUSTOM",
		PayAmount:     25.00,
		PaymentStatus: "Paid",
		PaymentTime:   &customDate,
	}
	require.NoError(t, db.Create(&order).Error)

	// 创建请求，指定日期
	c, w := createMachineOwnerContext(ownerID, true)
	c.Request.URL.RawQuery = "dateTime=2025-08-10T00:00:00Z"

	handler.GetSales(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response contracts.MachineOwnerSalesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)

	// 找到机器1的销售数据
	var machine1Sales decimal.Decimal
	for _, sale := range response.Data {
		if sale.Label == "Test Machine 1" {
			machine1Sales = sale.Value
			break
		}
	}

	assert.True(t, machine1Sales.Equal(decimal.NewFromFloat(25.00)))
}

func TestMachineOwnerHandler_GetSalesStats_Success(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)
	ownerID, machineID1, machineID2 := setupMachineOwnerTestData(t, db)

	// 创建订单数据
	today := time.Now().Truncate(24 * time.Hour)
	orders := []models.Order{
		{
			ID:            "order-001",
			MemberId:      "member-001",
			MachineId:     machineID1,
			ProductId:     "product-001",
			OrderNo:       "ON001",
			PayAmount:     15.50,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
		{
			ID:            "order-002",
			MemberId:      "member-001",
			MachineId:     machineID2,
			ProductId:     "product-001",
			OrderNo:       "ON002",
			PayAmount:     18.00,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
	}

	for _, order := range orders {
		require.NoError(t, db.Create(&order).Error)
	}

	// 创建请求
	c, w := createMachineOwnerContext(ownerID, true)
	c.Request.URL.RawQuery = "startDate=2025-08-11&endDate=2025-08-11"

	handler.GetSalesStats(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response contracts.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
}

func TestMachineOwnerHandler_GetSalesStats_InvalidDateFormat(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)
	ownerID, _, _ := setupMachineOwnerTestData(t, db)

	// 创建请求，使用无效日期格式
	c, w := createMachineOwnerContext(ownerID, true)
	c.Request.URL.RawQuery = "startDate=invalid-date"

	handler.GetSalesStats(c)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response contracts.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Contains(t, response.Error.Message, "日期格式错误")
}

func TestMachineOwnerHandler_GetSalesStats_EndDateBeforeStartDate(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)
	ownerID, _, _ := setupMachineOwnerTestData(t, db)

	// 创建请求，结束日期早于开始日期
	c, w := createMachineOwnerContext(ownerID, true)
	c.Request.URL.RawQuery = "startDate=2025-08-15&endDate=2025-08-10"

	handler.GetSalesStats(c)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response contracts.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Contains(t, response.Error.Message, "结束日期不能早于开始日期")
}

func TestMachineOwnerHandler_GetSalesStats_NotMachineOwner(t *testing.T) {
	db := setupMachineOwnerTestDB(t)
	handler := NewMachineOwnerHandler(db)

	// 创建普通用户上下文
	c, w := createMachineOwnerContext("", false)

	handler.GetSalesStats(c)

	// 验证响应
	assert.Equal(t, http.StatusForbidden, w.Code)

	var response contracts.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Contains(t, response.Error.Message, "您不是机主")
}
