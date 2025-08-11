package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 记录请求ID（如果有）
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 将请求ID添加到上下文
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)

		// 获取状态码
		status := c.Writer.Status()

		// 记录日志
		log.Printf("[%s] %s %s %d %v | %s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
			c.ClientIP(),
		)

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			log.Printf("[%s] Errors: %v", requestID, c.Errors)
		}
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简单的时间戳 + 随机数方式
	// 生产环境可以使用更复杂的UUID生成
	return time.Now().Format("20060102150405") + randomString(6)
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
