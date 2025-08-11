package main

import (
	"log"
	"os"

	"github.com/ddteam/drink-master/internal/config"
	"github.com/ddteam/drink-master/internal/models"
	"github.com/ddteam/drink-master/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// 初始化数据库
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
