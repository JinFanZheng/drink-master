package services

import "github.com/ddteam/drink-master/internal/repositories"

// StatsService 统计服务
type StatsService struct {
	consumptionRepo *repositories.ConsumptionLogRepository
	drinkRepo       *repositories.DrinkRepository
}

// NewStatsService 创建统计服务
func NewStatsService(consumptionRepo *repositories.ConsumptionLogRepository, drinkRepo *repositories.DrinkRepository) *StatsService {
	return &StatsService{
		consumptionRepo: consumptionRepo,
		drinkRepo:       drinkRepo,
	}
}

// TODO: 实现统计服务方法