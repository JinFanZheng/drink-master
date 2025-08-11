package models

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutoMigrate(t *testing.T) {
	// 创建内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 测试自动迁移
	err = AutoMigrate(db)
	if err != nil {
		t.Errorf("AutoMigrate failed: %v", err)
	}

	// 验证表是否被创建
	tables := []string{"members", "machine_owners", "machines", "products", "machine_product_prices", "orders"}
	
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("Table %s was not created", table)
		}
	}
}

func TestAllModels(t *testing.T) {
	models := AllModels()
	
	if len(models) == 0 {
		t.Error("AllModels returned empty slice")
	}

	// 验证返回的模型数量
	expectedCount := 6 // Member, MachineOwner, Machine, Product, MachineProductPrice, Order
	if len(models) != expectedCount {
		t.Errorf("Expected %d models, got %d", expectedCount, len(models))
	}
}

// 测试所有模型的TableName方法
func TestTableNames(t *testing.T) {
	tests := []struct {
		model    interface{ TableName() string }
		expected string
	}{
		{Member{}, "members"},
		{MachineOwner{}, "machine_owners"},
		{Machine{}, "machines"},
		{Product{}, "products"},
		{MachineProductPrice{}, "machine_product_prices"},
		{Order{}, "orders"},
	}

	for _, tt := range tests {
		if got := tt.model.TableName(); got != tt.expected {
			t.Errorf("TableName() = %v, want %v", got, tt.expected)
		}
	}
}

// 测试Member模型的CRUD操作
func TestMemberModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// 创建测试数据
	member := Member{
		ID:           "test-member-001",
		Nickname:     "测试用户",
		Avatar:       "https://example.com/avatar.jpg",
		WeChatOpenId: "wx_test_openid_123",
		Role:         "Member",
		IsAdmin:      false,
	}

	// 测试创建
	if err := db.Create(&member).Error; err != nil {
		t.Errorf("Failed to create member: %v", err)
	}

	// 测试查询
	var foundMember Member
	if err := db.First(&foundMember, "id = ?", member.ID).Error; err != nil {
		t.Errorf("Failed to find member: %v", err)
	}

	if foundMember.Nickname != member.Nickname {
		t.Errorf("Expected nickname %s, got %s", member.Nickname, foundMember.Nickname)
	}

	// 测试更新
	foundMember.Nickname = "更新的用户名"
	if err := db.Save(&foundMember).Error; err != nil {
		t.Errorf("Failed to update member: %v", err)
	}

	// 验证更新
	var updatedMember Member
	if err := db.First(&updatedMember, "id = ?", member.ID).Error; err != nil {
		t.Errorf("Failed to find updated member: %v", err)
	}

	if updatedMember.Nickname != "更新的用户名" {
		t.Errorf("Member nickname was not updated")
	}
}

// 测试MachineOwner模型
func TestMachineOwnerModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	receivingAccount := "alipay_account_123"
	owner := MachineOwner{
		ID:               "test-owner-001",
		Name:             "测试机主",
		ReceivingAccount: &receivingAccount,
	}

	// 测试创建
	if err := db.Create(&owner).Error; err != nil {
		t.Errorf("Failed to create machine owner: %v", err)
	}

	// 测试查询
	var foundOwner MachineOwner
	if err := db.First(&foundOwner, "id = ?", owner.ID).Error; err != nil {
		t.Errorf("Failed to find machine owner: %v", err)
	}

	if foundOwner.Name != owner.Name {
		t.Errorf("Expected name %s, got %s", owner.Name, foundOwner.Name)
	}
}

// 测试Machine模型及其关联
func TestMachineModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// 先创建机主
	owner := MachineOwner{
		ID:   "test-owner-001",
		Name: "测试机主",
	}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("Failed to create machine owner: %v", err)
	}

	// 创建售货机
	servicePhone := "400-123-4567"
	deviceId := "device_001"
	machine := Machine{
		ID:             "test-machine-001",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM001",
		Name:           "测试售货机",
		Area:           "上海浦东",
		Address:        "浦东新区张江路123号",
		ServicePhone:   &servicePhone,
		DeviceId:       &deviceId,
		BusinessStatus: "Open",
	}

	// 测试创建
	if err := db.Create(&machine).Error; err != nil {
		t.Errorf("Failed to create machine: %v", err)
	}

	// 测试查询及关联
	var foundMachine Machine
	if err := db.Preload("MachineOwner").First(&foundMachine, "id = ?", machine.ID).Error; err != nil {
		t.Errorf("Failed to find machine: %v", err)
	}

	if foundMachine.Name != machine.Name {
		t.Errorf("Expected machine name %s, got %s", machine.Name, foundMachine.Name)
	}

	if foundMachine.MachineOwner == nil {
		t.Error("Machine owner should be loaded")
	} else if foundMachine.MachineOwner.Name != owner.Name {
		t.Errorf("Expected owner name %s, got %s", owner.Name, foundMachine.MachineOwner.Name)
	}
}

// 测试Product和MachineProductPrice模型
func TestProductAndPriceModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// 创建产品
	description := "美味可口的咖啡"
	category := "饮料"
	product := Product{
		ID:          "test-product-001",
		Name:        "拿铁咖啡",
		Description: &description,
		Category:    &category,
	}

	if err := db.Create(&product).Error; err != nil {
		t.Errorf("Failed to create product: %v", err)
	}

	// 创建机器（简化）
	owner := MachineOwner{ID: "owner-001", Name: "Owner"}
	db.Create(&owner)
	
	machine := Machine{
		ID:             "machine-001",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM001",
		Name:           "Test Machine",
	}
	db.Create(&machine)

	// 创建价格
	price := MachineProductPrice{
		ID:              "test-price-001",
		MachineId:       machine.ID,
		ProductId:       product.ID,
		Price:           15.50,
		PriceWithoutCup: 13.50,
		Stock:           100,
	}

	if err := db.Create(&price).Error; err != nil {
		t.Errorf("Failed to create machine product price: %v", err)
	}

	// 测试关联查询
	var foundPrice MachineProductPrice
	if err := db.Preload("Product").Preload("Machine").First(&foundPrice, "id = ?", price.ID).Error; err != nil {
		t.Errorf("Failed to find price: %v", err)
	}

	if foundPrice.Price != price.Price {
		t.Errorf("Expected price %f, got %f", price.Price, foundPrice.Price)
	}

	if foundPrice.Product == nil || foundPrice.Product.Name != product.Name {
		t.Error("Product association not loaded correctly")
	}
}

// 测试Order模型的完整流程
func TestOrderModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// 创建必要的关联数据
	owner := MachineOwner{ID: "owner-001", Name: "Owner"}
	db.Create(&owner)

	member := Member{
		ID:           "member-001",
		Nickname:     "Customer",
		WeChatOpenId: "wx_openid",
		Role:         "Member",
	}
	db.Create(&member)

	machine := Machine{
		ID:             "machine-001",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM001",
		Name:           "Test Machine",
	}
	db.Create(&machine)

	product := Product{
		ID:   "product-001",
		Name: "Coffee",
	}
	db.Create(&product)

	// 创建订单
	paymentTime := time.Now()
	channelOrderNo := "CHANNEL_ORDER_123"
	refundReason := "商品有问题"
	
	order := Order{
		ID:             "test-order-001",
		MemberId:       member.ID,
		MachineId:      machine.ID,
		ProductId:      product.ID,
		OrderNo:        "ORDER_20250812_001",
		HasCup:         true,
		TotalAmount:    15.50,
		PayAmount:      15.50,
		PaymentStatus:  "Paid",
		PaymentTime:    &paymentTime,
		ChannelOrderNo: &channelOrderNo,
		MakeStatus:     "Made",
		RefundAmount:   0.0,
		RefundReason:   &refundReason,
	}

	// 测试创建订单
	if err := db.Create(&order).Error; err != nil {
		t.Errorf("Failed to create order: %v", err)
	}

	// 测试查询订单及关联
	var foundOrder Order
	if err := db.Preload("Member").Preload("Machine").Preload("Product").First(&foundOrder, "id = ?", order.ID).Error; err != nil {
		t.Errorf("Failed to find order: %v", err)
	}

	if foundOrder.OrderNo != order.OrderNo {
		t.Errorf("Expected order no %s, got %s", order.OrderNo, foundOrder.OrderNo)
	}

	if foundOrder.TotalAmount != order.TotalAmount {
		t.Errorf("Expected total amount %f, got %f", order.TotalAmount, foundOrder.TotalAmount)
	}

	// 验证关联数据
	if foundOrder.Member == nil || foundOrder.Member.Nickname != member.Nickname {
		t.Error("Member association not loaded correctly")
	}

	if foundOrder.Machine == nil || foundOrder.Machine.Name != machine.Name {
		t.Error("Machine association not loaded correctly")
	}

	if foundOrder.Product == nil || foundOrder.Product.Name != product.Name {
		t.Error("Product association not loaded correctly")
	}
}

// 测试模型字段约束和验证
func TestModelConstraints(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// 测试唯一约束 - WeChatOpenId
	member1 := Member{
		ID:           "member-001",
		Nickname:     "User1",
		WeChatOpenId: "same_openid",
		Role:         "Member",
	}

	member2 := Member{
		ID:           "member-002", 
		Nickname:     "User2",
		WeChatOpenId: "same_openid", // 相同的openid应该失败
		Role:         "Member",
	}

	// 第一个应该成功
	if err := db.Create(&member1).Error; err != nil {
		t.Errorf("Failed to create first member: %v", err)
	}

	// 第二个应该失败（唯一约束）
	if err := db.Create(&member2).Error; err == nil {
		t.Error("Expected unique constraint violation for WeChatOpenId")
	}

	// 测试唯一约束 - MachineNo
	owner := MachineOwner{ID: "owner-001", Name: "Owner"}
	db.Create(&owner)

	machine1 := Machine{
		ID:             "machine-001",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM001",
		Name:           "Machine 1",
	}

	machine2 := Machine{
		ID:             "machine-002",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM001", // 相同的机器号应该失败
		Name:           "Machine 2",
	}

	// 第一个应该成功
	if err := db.Create(&machine1).Error; err != nil {
		t.Errorf("Failed to create first machine: %v", err)
	}

	// 第二个应该失败（唯一约束）
	if err := db.Create(&machine2).Error; err == nil {
		t.Error("Expected unique constraint violation for MachineNo")
	}
}