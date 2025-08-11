package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestLogger() (*gin.Engine, *bytes.Buffer) {
	gin.SetMode(gin.TestMode)
	
	// 创建一个buffer来捕获日志输出
	var buf bytes.Buffer
	log.SetOutput(&buf)
	
	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})
	
	return router, &buf
}

func TestRequestLogger(t *testing.T) {
	router, buf := setupTestLogger()
	defer func() {
		log.SetOutput(os.Stderr) // 恢复默认输出
	}()
	
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	// 验证日志是否包含请求信息
	logOutput := buf.String()
	if logOutput == "" {
		t.Error("Expected log output to be generated")
	}
}

func TestGenerateRequestID(t *testing.T) {
	requestID := generateRequestID()
	
	if len(requestID) == 0 {
		t.Error("Expected requestID to be generated")
	}
	
	// 验证生成的ID格式（时间戳14位 + 随机字符6位 = 20位）
	if len(requestID) != 20 {
		t.Errorf("Expected requestID length to be 20, got %d", len(requestID))
	}
}

func TestRandomString(t *testing.T) {
	str1 := randomString(6)
	str2 := randomString(6)
	
	if len(str1) != 6 {
		t.Errorf("Expected string length to be 6, got %d", len(str1))
	}
	
	if len(str2) != 6 {
		t.Errorf("Expected string length to be 6, got %d", len(str2))
	}
	
	// 验证字符串只包含字母和数字
	for _, char := range str1 {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			t.Errorf("Expected only lowercase letters and numbers, got '%c'", char)
		}
	}
}