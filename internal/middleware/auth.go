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

// JWTClaims extends jwt.RegisteredClaims with custom fields
type JWTClaims struct {
	MemberID       string `json:"member_id"`
	MachineOwnerID string `json:"machine_owner_id,omitempty"`
	Role           string `json:"role"`
	jwt.RegisteredClaims
}

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

			if strings.Contains(err.Error(), "expired") {
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
		c.Set("member_id", claims.MemberID)
		c.Set("machine_owner_id", claims.MachineOwnerID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// validateJWT 验证JWT token
func validateJWT(tokenString string) (*JWTClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_jwt_secret_change_this_in_production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenNotValidYet
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// GetCurrentMemberID 获取当前用户ID的辅助函数
func GetCurrentMemberID(c *gin.Context) (string, bool) {
	memberID, exists := c.Get("member_id")
	if !exists {
		return "", false
	}
	return memberID.(string), true
}

// GetCurrentMachineOwnerID 获取当前机主ID的辅助函数
func GetCurrentMachineOwnerID(c *gin.Context) (string, bool) {
	machineOwnerID, exists := c.Get("machine_owner_id")
	if !exists {
		return "", false
	}
	return machineOwnerID.(string), true
}

// GetCurrentRole 获取当前用户角色的辅助函数
func GetCurrentRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	return role.(string), true
}

// IsMachineOwner 检查当前用户是否为机主
func IsMachineOwner(c *gin.Context) bool {
	role, exists := GetCurrentRole(c)
	return exists && role == "Owner"
}
