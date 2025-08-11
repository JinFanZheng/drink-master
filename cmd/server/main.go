package main

import (
	"log"
	"os"

	"github.com/ddteam/drink-master/internal/handlers"
	"github.com/ddteam/drink-master/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db := initDatabase()

	// TODO: 添加自动迁移
	// if err := db.AutoMigrate(&models.YourModel{}); err != nil {
	//     log.Fatal("Failed to migrate database:", err)
	// }

	// 初始化handlers
	healthHandler := handlers.NewHealthHandler(db)

	// 初始化路由
	router := setupRouter(healthHandler)

	// 获取端口配置
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(":" + port))
}

func initDatabase() *gorm.DB {
	// 构建数据库连接字符串
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "root")
	password := getEnvOrDefault("DB_PASSWORD", "")
	dbname := getEnvOrDefault("DB_NAME", "app_dev")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

func setupRouter(healthHandler *handlers.HealthHandler) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())

	// API路由组
	api := router.Group("/api")
	{
		// 健康检查 (无需认证)
		api.GET("/health", healthHandler.Health)
		api.GET("/health/db", healthHandler.DatabaseHealth)

		// TODO: 添加其他路由
		// auth := api.Group("/auth")
		// {
		//     auth.POST("/login", authHandler.Login)
		// }

		// TODO: 添加需要认证的路由
		// protected := api.Group("/")
		// protected.Use(middleware.JWTAuth())
		// {
		//     // 添加保护路由
		// }
	}

	return router
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
