package services

import "github.com/ddteam/drink-master/internal/repositories"

// DrinkService 饮品服务
type DrinkService struct {
	drinkRepo    *repositories.DrinkRepository
	categoryRepo *repositories.DrinkCategoryRepository
}

// NewDrinkService 创建饮品服务
func NewDrinkService(drinkRepo *repositories.DrinkRepository, categoryRepo *repositories.DrinkCategoryRepository) *DrinkService {
	return &DrinkService{
		drinkRepo:    drinkRepo,
		categoryRepo: categoryRepo,
	}
}

// TODO: 实现饮品服务方法