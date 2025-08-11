package main

import (
	"log"
	"os"

	"github.com/ddteam/drink-master/internal/handlers"
	"github.com/ddteam/drink-master/internal/middleware"
	"github.com/ddteam/drink-master/internal/models"
	"github.com/ddteam/drink-master/internal/repositories"
	"github.com/ddteam/drink-master/internal/services"
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
	
	// 自动迁移
	if err := db.AutoMigrate(&models.User{}, &models.Drink{}, &models.DrinkCategory{}, &models.ConsumptionLog{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化repositories
	userRepo := repositories.NewUserRepository(db)
	drinkRepo := repositories.NewDrinkRepository(db)
	categoryRepo := repositories.NewDrinkCategoryRepository(db)
	consumptionRepo := repositories.NewConsumptionLogRepository(db)

	// 初始化services
	userService := services.NewUserService(userRepo)
	drinkService := services.NewDrinkService(drinkRepo, categoryRepo)
	statsService := services.NewStatsService(consumptionRepo, drinkRepo)

	// 初始化handlers
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(userService)
	drinkHandler := handlers.NewDrinkHandler(drinkService)
	statsHandler := handlers.NewStatsHandler(statsService)

	// 初始化路由
	router := setupRouter(healthHandler, authHandler, drinkHandler, statsHandler)

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
	dbname := getEnvOrDefault("DB_NAME", "drink_master_dev")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

func setupRouter(
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	drinkHandler *handlers.DrinkHandler,
	statsHandler *handlers.StatsHandler,
) *gin.Engine {
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

		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			// 饮品管理
			drinks := protected.Group("/drinks")
			{
				drinks.GET("", drinkHandler.GetDrinks)
				drinks.POST("", drinkHandler.CreateDrink)
				drinks.GET("/:id", drinkHandler.GetDrink)
				drinks.PUT("/:id", drinkHandler.UpdateDrink)
				drinks.DELETE("/:id", drinkHandler.DeleteDrink)
			}

			// 统计分析
			stats := protected.Group("/stats")
			{
				stats.GET("/consumption", statsHandler.GetConsumptionStats)
				stats.GET("/popular", statsHandler.GetPopularDrinks)
				stats.GET("/trends", statsHandler.GetConsumptionTrends)
			}

			// 消费记录
			logs := protected.Group("/logs")
			{
				logs.POST("", drinkHandler.LogConsumption)
				logs.GET("", drinkHandler.GetConsumptionLogs)
			}
		}
	}

	return router
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}