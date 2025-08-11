package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, contracts.APIResponse{
				Success: false,
				Error: &contracts.APIError{
					Code:      contracts.ErrorCodeUnauthorized,
					Message:   "缺少Authorization头",
					Timestamp: time.Now(),
					Path:      c.Request.URL.Path,
					Method:    c.Request.Method,
					RequestID: getRequestID(c),
				},
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, contracts.APIResponse{
				Success: false,
				Error: &contracts.APIError{
					Code:      contracts.ErrorCodeUnauthorized,
					Message:   "Authorization头格式错误，应为：Bearer <token>",
					Timestamp: time.Now(),
					Path:      c.Request.URL.Path,
					Method:    c.Request.Method,
					RequestID: getRequestID(c),
				},
			})
			c.Abort()
			return
		}

		// 验证JWT token
		claims, err := validateJWT(tokenString)
		if err != nil {
			var errorCode string
			var message string
			
			if err.Error() == "token is expired" {
				errorCode = contracts.ErrorCodeTokenExpired
				message = "Token已过期"
			} else {
				errorCode = contracts.ErrorCodeInvalidToken
				message = "无效的Token"
			}

			c.JSON(http.StatusUnauthorized, contracts.APIResponse{
				Success: false,
				Error: &contracts.APIError{
					Code:      errorCode,
					Message:   message,
					Details:   map[string]interface{}{"jwt_error": err.Error()},
					Timestamp: time.Now(),
					Path:      c.Request.URL.Path,
					Method:    c.Request.Method,
					RequestID: getRequestID(c),
				},
			})
			c.Abort()
			return
		}

		// 将用户信息添加到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// validateJWT 验证JWT token
func validateJWT(tokenString string) (*contracts.TokenClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_jwt_secret_change_this_in_production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenNotValidYet
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	// 检查token是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}

	// 从Subject中解析用户信息
	// 这里简化处理，实际项目中可能需要从数据库查询完整用户信息
	userClaims := &contracts.TokenClaims{
		UserID: parseUserID(claims.Subject),
		// Username 和 Email 可以从自定义claims中获取
		// 这里为了简化，使用Subject作为Username
		Username: claims.Subject,
	}

	return userClaims, nil
}

// parseUserID 从Subject中解析用户ID
func parseUserID(subject string) uint {
	// 这里简化处理，实际应该有更严格的解析逻辑
	// 可以考虑在JWT中直接存储用户ID
	return 1 // 简化返回，实际项目中需要正确实现
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// GetCurrentUserID 获取当前用户ID的辅助函数
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentUsername 获取当前用户名的辅助函数
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	return username.(string), true
}