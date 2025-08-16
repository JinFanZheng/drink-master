package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutoMigrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)
}

func TestMemberModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)

	nickname := "测试用户"
	avatar := "https://example.com/avatar.jpg"
	weChatOpenId := "wx_test_openid_123"

	member := Member{
		ID:           "test-member-001",
		Nickname:     &nickname,
		Avatar:       &avatar,
		WeChatOpenId: &weChatOpenId,
		Role:         1,
		IsAdmin:      BitBool(0),
		CreatedOn:    time.Now(),
	}

	// Test creation
	err = db.Create(&member).Error
	assert.NoError(t, err)

	// Test finding
	var foundMember Member
	err = db.First(&foundMember, "id = ?", member.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, member.ID, foundMember.ID)
	assert.NotNil(t, foundMember.Nickname)
	assert.Equal(t, "测试用户", *foundMember.Nickname)
}

func TestMachineModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)

	machineNo := "VM001"
	name := "测试售货机"

	machine := Machine{
		ID:          "test-machine-001",
		MachineNo:   &machineNo,
		Name:        &name,
		IsDebugMode: BitBool(0),
		CreatedOn:   time.Now(),
	}

	// Test creation
	err = db.Create(&machine).Error
	assert.NoError(t, err)

	// Test finding
	var foundMachine Machine
	err = db.First(&foundMachine, "id = ?", machine.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, machine.ID, foundMachine.ID)
	assert.NotNil(t, foundMachine.MachineNo)
	assert.Equal(t, "VM001", *foundMachine.MachineNo)
}

func TestOrderModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)

	orderNo := "ORDER_20250815_001"

	order := Order{
		ID:            "test-order-001",
		OrderNo:       &orderNo,
		TotalAmount:   15.50,
		PayAmount:     15.50,
		PaymentStatus: 1,
		MakeStatus:    1,
		HasCup:        BitBool(1),
		CreatedOn:     time.Now(),
	}

	// Test creation
	err = db.Create(&order).Error
	assert.NoError(t, err)

	// Test finding
	var foundOrder Order
	err = db.First(&foundOrder, "id = ?", order.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, order.ID, foundOrder.ID)
	assert.Equal(t, BitBool(1), foundOrder.HasCup)
}

func TestMaterialSiloModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)

	no := "01"

	silo := MaterialSilo{
		ID:        "test-silo-001",
		No:        &no,
		Type:      1,
		IsSale:    BitBool(1),
		Total:     100,
		Stock:     50,
		CreatedOn: time.Now(),
	}

	// Test creation
	err = db.Create(&silo).Error
	assert.NoError(t, err)

	// Test finding
	var foundSilo MaterialSilo
	err = db.First(&foundSilo, "id = ?", silo.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, silo.ID, foundSilo.ID)
	assert.Equal(t, BitBool(1), foundSilo.IsSale)
	assert.False(t, foundSilo.CanSale()) // Should be false because no ProductId
}
