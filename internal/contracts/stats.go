package contracts

import "time"

// 统计分析相关的API契约定义

// StatsRequest 统计请求参数
type StatsRequest struct {
	Period    string `form:"period" example:"week"`       // 统计周期: day, week, month, year
	StartDate string `form:"start_date" example:"2023-01-01"` // 开始日期
	EndDate   string `form:"end_date" example:"2023-01-31"`   // 结束日期
	Category  string `form:"category" example:"coffee"`       // 饮品分类过滤
}

// ConsumptionStatsResponse 消费统计响应
type ConsumptionStatsResponse struct {
	Period      string              `json:"period" example:"week"`
	TotalDrinks int                 `json:"total_drinks" example:"150"`        // 总消费杯数
	TotalAmount float64             `json:"total_amount" example:"1250.50"`    // 总消费金额
	AvgPerDay   float64             `json:"avg_per_day" example:"21.43"`       // 日均消费杯数
	Categories  []CategoryStatsItem `json:"categories"`                        // 分类统计
	Timeline    []TimelineStatsItem `json:"timeline"`                          // 时间线统计
	Meta        StatsMeta           `json:"meta"`
}

// CategoryStatsItem 分类统计项
type CategoryStatsItem struct {
	Category    string  `json:"category" example:"coffee"`
	Count       int     `json:"count" example:"80"`
	Amount      float64 `json:"amount" example:"800.00"`
	Percentage  float64 `json:"percentage" example:"53.33"`
}

// TimelineStatsItem 时间线统计项
type TimelineStatsItem struct {
	Date   string  `json:"date" example:"2023-01-01"`
	Count  int     `json:"count" example:"12"`
	Amount float64 `json:"amount" example:"120.50"`
}

// PopularDrinksRequest 热门饮品请求参数
type PopularDrinksRequest struct {
	Period   string `form:"period" example:"month"`      // 统计周期
	Category string `form:"category" example:"coffee"`   // 分类过滤
	Limit    int    `form:"limit" validate:"min=1,max=50" example:"10"` // 返回数量
}

// PopularDrinksResponse 热门饮品响应
type PopularDrinksResponse struct {
	Period string              `json:"period" example:"month"`
	Drinks []PopularDrinkItem  `json:"drinks"`
	Meta   StatsMeta           `json:"meta"`
}

// PopularDrinkItem 热门饮品项
type PopularDrinkItem struct {
	Drink       DrinkResponse `json:"drink"`
	Count       int           `json:"count" example:"25"`           // 消费次数
	TotalAmount float64       `json:"total_amount" example:"637.5"` // 总金额
	Percentage  float64       `json:"percentage" example:"16.67"`   // 占比
	Rank        int           `json:"rank" example:"1"`             // 排名
}

// TrendsRequest 趋势分析请求参数
type TrendsRequest struct {
	Period    string `form:"period" example:"month"`      // 分析周期: week, month, quarter, year
	Granularity string `form:"granularity" example:"day"` // 数据粒度: day, week, month
	Category  string `form:"category" example:"coffee"`   // 分类过滤
	DrinkID   *uint  `form:"drink_id" example:"1"`        // 特定饮品过滤
}

// TrendsResponse 趋势分析响应
type TrendsResponse struct {
	Period      string           `json:"period" example:"month"`
	Granularity string           `json:"granularity" example:"day"`
	Trends      []TrendDataPoint `json:"trends"`
	Summary     TrendSummary     `json:"summary"`
	Meta        StatsMeta        `json:"meta"`
}

// TrendDataPoint 趋势数据点
type TrendDataPoint struct {
	Date        string  `json:"date" example:"2023-01-01"`
	Count       int     `json:"count" example:"12"`
	Amount      float64 `json:"amount" example:"120.50"`
	GrowthRate  float64 `json:"growth_rate" example:"5.2"`      // 增长率%
	MovingAvg   float64 `json:"moving_avg" example:"115.30"`    // 移动平均值
}

// TrendSummary 趋势摘要
type TrendSummary struct {
	TotalGrowth    float64 `json:"total_growth" example:"15.5"`     // 总增长率%
	AvgGrowthRate  float64 `json:"avg_growth_rate" example:"2.3"`   // 平均增长率%
	PeakDate       string  `json:"peak_date" example:"2023-01-15"`  // 峰值日期
	PeakValue      float64 `json:"peak_value" example:"200.0"`      // 峰值
	LowDate        string  `json:"low_date" example:"2023-01-05"`   // 最低点日期
	LowValue       float64 `json:"low_value" example:"80.0"`        // 最低值
	Volatility     float64 `json:"volatility" example:"12.5"`       // 波动率
}

// CategoryTrendsResponse 分类趋势响应
type CategoryTrendsResponse struct {
	Categories []CategoryTrend `json:"categories"`
	Meta       StatsMeta       `json:"meta"`
}

// CategoryTrend 分类趋势
type CategoryTrend struct {
	Category string           `json:"category" example:"coffee"`
	Trends   []TrendDataPoint `json:"trends"`
	Summary  TrendSummary     `json:"summary"`
}

// UserStatsResponse 用户个人统计响应
type UserStatsResponse struct {
	UserID          uint                 `json:"user_id" example:"1"`
	TotalDrinks     int                  `json:"total_drinks" example:"150"`
	TotalAmount     float64              `json:"total_amount" example:"1250.50"`
	AvgPerDay       float64              `json:"avg_per_day" example:"5.2"`
	FavoriteDrink   *DrinkResponse       `json:"favorite_drink,omitempty"`
	FavoriteCategory string              `json:"favorite_category" example:"coffee"`
	Categories      []CategoryStatsItem  `json:"categories"`
	RecentActivity  []ConsumptionLogResponse `json:"recent_activity"`
	Achievements    []AchievementItem    `json:"achievements"`
	Meta            StatsMeta            `json:"meta"`
}

// AchievementItem 成就项目
type AchievementItem struct {
	ID          string    `json:"id" example:"coffee_lover"`
	Name        string    `json:"name" example:"咖啡爱好者"`
	Description string    `json:"description" example:"连续7天喝咖啡"`
	Icon        string    `json:"icon" example:"☕"`
	UnlockedAt  time.Time `json:"unlocked_at" example:"2023-01-15T10:30:00Z"`
	Progress    float64   `json:"progress" example:"100.0"`      // 完成进度%
}

// StatsMeta 统计元数据
type StatsMeta struct {
	GeneratedAt time.Time `json:"generated_at" example:"2023-01-01T00:00:00Z"`
	Period      string    `json:"period" example:"2023-01-01 to 2023-01-31"`
	DataPoints  int       `json:"data_points" example:"31"`    // 数据点数量
	Warnings    []string  `json:"warnings,omitempty"`          // 数据质量警告
}