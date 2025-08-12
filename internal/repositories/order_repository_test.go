package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo OrderRepository
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// 自动迁移
	err = db.AutoMigrate(&models.Order{}, &models.Member{}, &models.Machine{}, &models.Product{})
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = NewOrderRepository(db)

	// 创建测试数据
	suite.createTestData()
}

func (suite *OrderRepositoryTestSuite) createTestData() {
	// 创建测试会员
	member := &models.Member{
		ID:           "test-member-1",
		Nickname:     "Test Member",
		WeChatOpenId: "test-openid-1",
		Role:         "Member",
	}
	suite.db.Create(member)

	// 创建测试机器
	machine := &models.Machine{
		ID:             "test-machine-1",
		MachineOwnerId: "test-owner-1",
		MachineNo:      "TM001",
		Name:           "Test Machine",
		BusinessStatus: "Open",
	}
	suite.db.Create(machine)

	// 创建测试产品
	product := &models.Product{
		ID:   "test-product-1",
		Name: "Test Product",
	}
	suite.db.Create(product)
}

func (suite *OrderRepositoryTestSuite) TestCreate() {
	order := &models.Order{
		ID:            "test-order-1",
		MemberId:      "test-member-1",
		MachineId:     "test-machine-1",
		ProductId:     "test-product-1",
		OrderNo:       "ORD202508120001",
		HasCup:        true,
		TotalAmount:   15.80,
		PayAmount:     15.80,
		PaymentStatus: "WaitPay",
		MakeStatus:    "WaitMake",
		RefundAmount:  0,
	}

	err := suite.repo.Create(order)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), order.ID)

	// 验证订单已创建
	var createdOrder models.Order
	err = suite.db.First(&createdOrder, "id = ?", order.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), order.OrderNo, createdOrder.OrderNo)
}

func (suite *OrderRepositoryTestSuite) TestGetByID() {
	// 创建测试订单
	order := &models.Order{
		ID:            "test-order-2",
		MemberId:      "test-member-1",
		MachineId:     "test-machine-1",
		ProductId:     "test-product-1",
		OrderNo:       "ORD202508120002",
		HasCup:        true,
		TotalAmount:   15.80,
		PayAmount:     15.80,
		PaymentStatus: "Paid",
		MakeStatus:    "Made",
		RefundAmount:  0,
	}
	suite.db.Create(order)

	// 测试获取订单
	result, err := suite.repo.GetByID("test-order-2")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "test-order-2", result.ID)
	assert.Equal(suite.T(), "ORD202508120002", result.OrderNo)
}

func (suite *OrderRepositoryTestSuite) TestGetByID_NotFound() {
	result, err := suite.repo.GetByID("non-existent-id")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *OrderRepositoryTestSuite) TestGetByMemberPaging() {
	// 创建多个测试订单
	orders := []*models.Order{
		{
			ID:            "test-order-3",
			MemberId:      "test-member-1",
			MachineId:     "test-machine-1",
			ProductId:     "test-product-1",
			OrderNo:       "ORD202508120003",
			HasCup:        true,
			TotalAmount:   15.80,
			PayAmount:     15.80,
			PaymentStatus: "Paid",
			MakeStatus:    "Made",
			RefundAmount:  0,
		},
		{
			ID:            "test-order-4",
			MemberId:      "test-member-1",
			MachineId:     "test-machine-1",
			ProductId:     "test-product-1",
			OrderNo:       "ORD202508120004",
			HasCup:        false,
			TotalAmount:   12.80,
			PayAmount:     12.80,
			PaymentStatus: "WaitPay",
			MakeStatus:    "WaitMake",
			RefundAmount:  0,
		},
	}

	for _, order := range orders {
		suite.db.Create(order)
	}

	// 测试分页获取
	result, total, err := suite.repo.GetByMemberPaging("test-member-1", 1, 10)
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(result), 2)
	assert.GreaterOrEqual(suite.T(), total, int64(2))

	// 验证订单按创建时间倒序排列
	if len(result) >= 2 {
		assert.True(suite.T(), result[0].CreatedAt.After(result[1].CreatedAt) || result[0].CreatedAt.Equal(result[1].CreatedAt))
	}
}

func (suite *OrderRepositoryTestSuite) TestUpdate() {
	// 创建测试订单
	order := &models.Order{
		ID:            "test-order-5",
		MemberId:      "test-member-1",
		MachineId:     "test-machine-1",
		ProductId:     "test-product-1",
		OrderNo:       "ORD202508120005",
		HasCup:        true,
		TotalAmount:   15.80,
		PayAmount:     15.80,
		PaymentStatus: "WaitPay",
		MakeStatus:    "WaitMake",
		RefundAmount:  0,
	}
	suite.db.Create(order)

	// 更新订单
	order.PaymentStatus = "Paid"
	order.MakeStatus = "Made"
	paymentTime := time.Now()
	order.PaymentTime = &paymentTime

	err := suite.repo.Update(order)
	assert.NoError(suite.T(), err)

	// 验证更新
	var updatedOrder models.Order
	err = suite.db.First(&updatedOrder, "id = ?", order.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Paid", updatedOrder.PaymentStatus)
	assert.Equal(suite.T(), "Made", updatedOrder.MakeStatus)
	assert.NotNil(suite.T(), updatedOrder.PaymentTime)
}

func (suite *OrderRepositoryTestSuite) TestDelete() {
	// 创建测试订单
	order := &models.Order{
		ID:            "test-order-6",
		MemberId:      "test-member-1",
		MachineId:     "test-machine-1",
		ProductId:     "test-product-1",
		OrderNo:       "ORD202508120006",
		HasCup:        true,
		TotalAmount:   15.80,
		PayAmount:     15.80,
		PaymentStatus: "WaitPay",
		MakeStatus:    "WaitMake",
		RefundAmount:  0,
	}
	suite.db.Create(order)

	// 删除订单
	err := suite.repo.Delete("test-order-6")
	assert.NoError(suite.T(), err)

	// 验证软删除
	var deletedOrder models.Order
	err = suite.db.Unscoped().First(&deletedOrder, "id = ?", "test-order-6").Error
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), deletedOrder.DeletedAt)

	// 验证正常查询找不到
	_, err = suite.repo.GetByID("test-order-6")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *OrderRepositoryTestSuite) TestGetByOrderNo() {
	// 创建测试订单
	order := &models.Order{
		ID:            "test-order-7",
		MemberId:      "test-member-1",
		MachineId:     "test-machine-1",
		ProductId:     "test-product-1",
		OrderNo:       "ORD202508120007",
		HasCup:        true,
		TotalAmount:   15.80,
		PayAmount:     15.80,
		PaymentStatus: "Paid",
		MakeStatus:    "Made",
		RefundAmount:  0,
	}
	suite.db.Create(order)

	// 测试根据订单号获取订单
	result, err := suite.repo.GetByOrderNo("ORD202508120007")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "test-order-7", result.ID)
	assert.Equal(suite.T(), "ORD202508120007", result.OrderNo)
}

func (suite *OrderRepositoryTestSuite) TestGetByOrderNo_NotFound() {
	result, err := suite.repo.GetByOrderNo("NON-EXISTENT-ORDER")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

// 单独的构造函数测试
func TestNewOrderRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewOrderRepository(db)
	assert.NotNil(t, repo)
	assert.IsType(t, &orderRepository{}, repo)
}
