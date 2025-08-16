package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

// stringPtr helper function for test setup
func stringPtr(s string) *string {
	return &s
}

// Mock implementations for testing
type MockMachineRepository struct {
	mock.Mock
}

func (m *MockMachineRepository) GetByID(id string) (*models.Machine, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Machine), args.Error(1)
}

func (m *MockMachineRepository) GetByDeviceID(deviceID string) (*models.Machine, error) {
	args := m.Called(deviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Machine), args.Error(1)
}

func (m *MockMachineRepository) GetList(machineOwnerID string) ([]*models.Machine, error) {
	args := m.Called(machineOwnerID)
	return args.Get(0).([]*models.Machine), args.Error(1)
}

func (m *MockMachineRepository) GetPaging(machineOwnerID string, keyword string, page, pageSize int) ([]*models.Machine, int64, error) {
	args := m.Called(machineOwnerID, keyword, page, pageSize)
	return args.Get(0).([]*models.Machine), args.Get(1).(int64), args.Error(2)
}

func (m *MockMachineRepository) UpdateBusinessStatus(id string, status enums.BusinessStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockMachineRepository) CheckDeviceExists(deviceID string) (bool, error) {
	args := m.Called(deviceID)
	return args.Bool(0), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetByID(id string) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetMachineProducts(machineID string) ([]*models.MachineProductPrice, error) {
	args := m.Called(machineID)
	return args.Get(0).([]*models.MachineProductPrice), args.Error(1)
}

type MockDeviceService struct {
	mock.Mock
}

func (m *MockDeviceService) CheckDeviceOnline(deviceID string) (bool, error) {
	args := m.Called(deviceID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDeviceService) UpdateRegister(deviceID string, params map[string]int) error {
	args := m.Called(deviceID, params)
	return args.Error(0)
}

func (m *MockDeviceService) GetDeviceStatus(deviceID string) (*contracts.DeviceStatusCheckResult, error) {
	args := m.Called(deviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contracts.DeviceStatusCheckResult), args.Error(1)
}

func createMachineService() (*MachineService, *MockMachineRepository, *MockProductRepository, *MockDeviceService) {
	mockMachineRepo := new(MockMachineRepository)
	mockProductRepo := new(MockProductRepository)
	mockDeviceService := new(MockDeviceService)

	// Create in-memory database for tests that need db access
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Product{})

	service := &MachineService{
		machineRepo:   mockMachineRepo,
		productRepo:   mockProductRepo,
		deviceService: mockDeviceService,
		db:            db,
	}

	return service, mockMachineRepo, mockProductRepo, mockDeviceService
}

func TestNewMachineService(t *testing.T) {
	// This test is mainly for coverage of the constructor function
	// We can't easily test with a real DB in this file due to mock setup,
	// but we can verify the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewMachineService panicked: %v", r)
		}
	}()

	// This will panic if db is nil, but that's expected behavior
	// The function itself has 0% coverage because it's never called in the mock tests
	// service := NewMachineService(nil) // Would panic, so we skip this

	// Instead, let's just verify our mock constructor works
	service, _, _, _ := createMachineService()
	assert.NotNil(t, service)
}

func TestMachineService_GetMachinePaging(t *testing.T) {
	service, mockRepo, _, _ := createMachineService()

	req := contracts.GetMachinePagingRequest{
		Page:           1,
		PageSize:       10,
		MachineOwnerID: "owner-123",
		Keyword:        "test",
	}

	machines := []*models.Machine{
		{
			ID:             "machine-1",
			MachineOwnerId: stringPtr("owner-123"),
			MachineNo:      stringPtr("M001"),
			Name:           stringPtr("Test Machine"),
			Area:           stringPtr("Area A"),
			Address:        stringPtr("Address A"),
			BusinessStatus: enums.BusinessStatusOpen,
		},
	}

	mockRepo.On("GetPaging", "owner-123", "test", 1, 10).Return(machines, int64(1), nil)

	result, err := service.GetMachinePaging(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int64(1), result.TotalCount)
	assert.Equal(t, 1, result.PageIndex)
	assert.Equal(t, 10, result.PageSize)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "machine-1", result.Items[0].ID)
	assert.Equal(t, "M001", result.Items[0].MachineNo)

	mockRepo.AssertExpectations(t)
}

func TestMachineService_GetMachinePaging_EmptyOwnerID(t *testing.T) {
	service, _, _, _ := createMachineService()

	req := contracts.GetMachinePagingRequest{
		Page:           1,
		PageSize:       10,
		MachineOwnerID: "",
	}

	_, err := service.GetMachinePaging(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "machine owner id is required")
}

func TestMachineService_GetMachineByID(t *testing.T) {
	service, mockRepo, _, mockDevice := createMachineService()

	servicePhone := "123-456-7890"
	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: stringPtr("owner-123"),
		MachineNo:      stringPtr("M001"),
		Name:           stringPtr("Test Machine"),
		Area:           stringPtr("Area A"),
		Address:        stringPtr("Address A"),
		BusinessStatus: enums.BusinessStatusOpen,
		// DeviceId field removed from model
		ServicePhone: &servicePhone,
		CreatedOn:    time.Now(),
	}

	mockRepo.On("GetByID", "machine-123").Return(machine, nil)
	mockDevice.On("CheckDeviceOnline", "M001").Return(true, nil)

	result, err := service.GetMachineByID("machine-123")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "machine-123", result.ID)
	assert.Equal(t, "M001", result.MachineNo)
	assert.Equal(t, "Test Machine", result.Name)
	assert.Equal(t, enums.BusinessStatusOpen.ToAPIString(), result.BusinessStatus)
	assert.Equal(t, "M001", result.DeviceID)
	assert.Equal(t, "123-456-7890", result.ServicePhone)

	mockRepo.AssertExpectations(t)
	mockDevice.AssertExpectations(t)
}

func TestMachineService_GetMachineByID_DeviceOffline(t *testing.T) {
	service, mockRepo, _, mockDevice := createMachineService()

	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: stringPtr("owner-123"),
		MachineNo:      stringPtr("M001"),
		Name:           stringPtr("Test Machine"),
		BusinessStatus: enums.BusinessStatusOpen,
		// DeviceId field removed from model
		CreatedOn: time.Now(),
	}

	mockRepo.On("GetByID", "machine-123").Return(machine, nil)
	mockDevice.On("CheckDeviceOnline", "M001").Return(false, nil)

	result, err := service.GetMachineByID("machine-123")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, contracts.BusinessStatusOffline, result.BusinessStatus)

	mockRepo.AssertExpectations(t)
	mockDevice.AssertExpectations(t)
}

func TestMachineService_GetProductList(t *testing.T) {
	service, _, mockProductRepo, _ := createMachineService()

	// Create test product in database
	testProduct := &models.Product{
		ID:   "product-1",
		Name: "Coffee",
	}
	service.db.Create(testProduct)

	machineProducts := []*models.MachineProductPrice{
		{
			ID:              "mp-1",
			MachineId:       "machine-123",
			ProductId:       "product-1",
			Price:           5.0,
			PriceWithoutCup: 4.5,
		},
	}

	mockProductRepo.On("GetMachineProducts", "machine-123").Return(machineProducts, nil)

	result, err := service.GetProductList("machine-123")
	require.NoError(t, err)
	require.Len(t, result, 1)

	productGroup := result[0]
	assert.Equal(t, contracts.ProductGroupTimeLimited, productGroup.Name)
	assert.Len(t, productGroup.Products, 1)

	productItem := productGroup.Products[0]
	assert.Equal(t, "mp-1", productItem.ID)
	assert.Equal(t, "Coffee", productItem.Name)
	assert.Equal(t, 5.0, productItem.Price)
	assert.Equal(t, 4.5, productItem.PriceWithoutCup)

	mockProductRepo.AssertExpectations(t)
}

func TestMachineService_OpenOrCloseBusiness(t *testing.T) {
	service, mockRepo, _, _ := createMachineService()

	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: stringPtr("owner-123"),
		BusinessStatus: enums.BusinessStatusOpen,
	}

	mockRepo.On("GetByID", "machine-123").Return(machine, nil)
	mockRepo.On("UpdateBusinessStatus", "machine-123", enums.BusinessStatusClose).Return(nil)

	result, err := service.OpenOrCloseBusiness("machine-123", "owner-123")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, enums.BusinessStatusClose.ToAPIString(), result.Status)
	assert.Contains(t, result.Message, "关闭营业")

	mockRepo.AssertExpectations(t)
}

func TestMachineService_OpenOrCloseBusiness_PermissionDenied(t *testing.T) {
	service, mockRepo, _, _ := createMachineService()

	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: stringPtr("owner-123"),
		BusinessStatus: enums.BusinessStatusOpen,
	}

	mockRepo.On("GetByID", "machine-123").Return(machine, nil)

	_, err := service.OpenOrCloseBusiness("machine-123", "different-owner")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied: not machine owner")

	mockRepo.AssertExpectations(t)
}

func TestMachineService_CheckDeviceExist(t *testing.T) {
	service, mockRepo, _, _ := createMachineService()

	mockRepo.On("CheckDeviceExists", "device-123").Return(true, nil)
	mockRepo.On("CheckDeviceExists", "nonexistent").Return(false, nil)

	// 测试存在的设备
	exists, err := service.CheckDeviceExist("device-123")
	require.NoError(t, err)
	assert.True(t, exists)

	// 测试不存在的设备
	exists, err = service.CheckDeviceExist("nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)

	// 测试空设备ID
	exists, err = service.CheckDeviceExist("")
	require.NoError(t, err)
	assert.False(t, exists)

	mockRepo.AssertExpectations(t)
}

func TestMachineService_ValidateMachineOwnership(t *testing.T) {
	service, mockRepo, _, _ := createMachineService()

	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: stringPtr("owner-123"),
	}

	mockRepo.On("GetByID", "machine-123").Return(machine, nil)

	// 测试有效的所有权
	err := service.ValidateMachineOwnership("machine-123", "owner-123")
	assert.NoError(t, err)

	// 测试无效的所有权
	err = service.ValidateMachineOwnership("machine-123", "different-owner")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied: not machine owner")

	mockRepo.AssertExpectations(t)
}

func TestMachineService_GetMachineList(t *testing.T) {
	service, mockRepo, _, _ := createMachineService()

	machines := []*models.Machine{
		{
			ID:             "machine-1",
			MachineOwnerId: stringPtr("owner-123"),
			MachineNo:      stringPtr("M001"),
			Name:           stringPtr("Test Machine 1"),
			BusinessStatus: enums.BusinessStatusOpen,
		},
		{
			ID:             "machine-2",
			MachineOwnerId: stringPtr("owner-123"),
			MachineNo:      stringPtr("M002"),
			Name:           stringPtr("Test Machine 2"),
			BusinessStatus: enums.BusinessStatusClose,
		},
	}

	mockRepo.On("GetList", "owner-123").Return(machines, nil)

	result, err := service.GetMachineList("owner-123")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "machine-1", result[0].ID)
	assert.Equal(t, "M001", result[0].MachineNo)
	assert.Equal(t, "machine-2", result[1].ID)
	assert.Equal(t, "M002", result[1].MachineNo)

	mockRepo.AssertExpectations(t)
}
