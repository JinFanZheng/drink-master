package repositories

import (
	"github.com/ddteam/drink-master/internal/models"
	"gorm.io/gorm"
)

// DrinkRepository 饮品存储库
type DrinkRepository struct {
	db *gorm.DB
}

// NewDrinkRepository 创建饮品存储库
func NewDrinkRepository(db *gorm.DB) *DrinkRepository {
	return &DrinkRepository{db: db}
}

// TODO: 实现饮品存储库方法

// DrinkCategoryRepository 饮品分类存储库
type DrinkCategoryRepository struct {
	db *gorm.DB
}

// NewDrinkCategoryRepository 创建饮品分类存储库
func NewDrinkCategoryRepository(db *gorm.DB) *DrinkCategoryRepository {
	return &DrinkCategoryRepository{db: db}
}

// TODO: 实现饮品分类存储库方法

// ConsumptionLogRepository 消费记录存储库
type ConsumptionLogRepository struct {
	db *gorm.DB
}

// NewConsumptionLogRepository 创建消费记录存储库
func NewConsumptionLogRepository(db *gorm.DB) *ConsumptionLogRepository {
	return &ConsumptionLogRepository{db: db}
}

// TODO: 实现消费记录存储库方法