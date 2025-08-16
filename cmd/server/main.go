// @title Drink Master API
// @version 1.0
// @description 智能售货机管理系统API文档
// @description 提供会员管理、设备管理、订单管理、支付等功能的RESTful API
//
// @contact.name API Support
// @contact.email support@drinkmaster.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/ddteam/drink-master/docs/swagger" // swagger docs
	"github.com/ddteam/drink-master/internal/config"
	"github.com/ddteam/drink-master/internal/models"
	"github.com/ddteam/drink-master/internal/routes"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	dbConfig := config.LoadDatabaseConfig()
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库模型
	if err := models.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 设置路由
	router := routes.SetupRoutes(db)

	// 获取端口配置
	port := getEnvOrDefault("PORT", "8080")

	log.Printf("Server starting on port %s", port)
	log.Printf("Database connected: %s", dbConfig.DSN())
	log.Fatal(router.Run(":" + port))
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
