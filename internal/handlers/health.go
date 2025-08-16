package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Health 基础健康检查
// @Summary 系统健康检查
// @Description 获取系统基础健康状态
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} contracts.HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response := contracts.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services: map[string]contracts.ServiceHealth{
			"api": {
				Status:      "ok",
				LastChecked: time.Now(),
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// DatabaseHealth 数据库健康检查
// @Summary 数据库健康检查
// @Description 检查数据库连接状态和响应时间
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} contracts.HealthResponse
// @Failure 503 {object} contracts.APIResponse
// @Router /health/db [get]
func (h *HealthHandler) DatabaseHealth(c *gin.Context) {
	start := time.Now()

	// 测试数据库连接
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, contracts.APIResponse{
			Success: false,
			Error: &contracts.APIError{
				Code:      "DATABASE_CONNECTION_ERROR",
				Message:   "数据库连接失败",
				Details:   map[string]interface{}{"error": err.Error()},
				Timestamp: time.Now(),
				Path:      c.Request.URL.Path,
				Method:    c.Request.Method,
				RequestID: "",
			},
		})
		return
	}

	// 测试数据库Ping
	err = sqlDB.Ping()
	latency := time.Since(start)

	if err != nil {
		response := contracts.HealthResponse{
			Status:    "down",
			Timestamp: time.Now(),
			Version:   "1.0.0",
			Services: map[string]contracts.ServiceHealth{
				"database": {
					Status:      "down",
					Latency:     latency.String(),
					LastChecked: time.Now(),
					Error:       err.Error(),
				},
			},
		}

		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// 获取数据库统计信息
	stats := sqlDB.Stats()

	response := contracts.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services: map[string]contracts.ServiceHealth{
			"database": {
				Status:      "ok",
				Latency:     latency.String(),
				LastChecked: time.Now(),
			},
		},
	}

	// 添加数据库连接池信息
	if _, exists := c.Get("db_meta"); !exists {
		c.Set("db_meta", map[string]interface{}{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
		})
	}

	c.JSON(http.StatusOK, response)
}
