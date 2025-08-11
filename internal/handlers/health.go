package handlers

import (
	"net/http"
	"time"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
				RequestID: getRequestID(c),
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
	if dbMeta, exists := c.Get("db_meta"); !exists {
		c.Set("db_meta", map[string]interface{}{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
		})
	}

	c.JSON(http.StatusOK, response)
}

// getRequestID 获取请求ID的辅助函数
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}