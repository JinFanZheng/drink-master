package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ddteam/drink-master/internal/config"
	"github.com/ddteam/drink-master/internal/handlers"
	"github.com/ddteam/drink-master/internal/middleware"
	"github.com/ddteam/drink-master/internal/repositories"
	"github.com/ddteam/drink-master/internal/services"
	"github.com/ddteam/drink-master/pkg/wechat"
)

// SetupRoutes 设置所有路由 (基于MobileAPI Controllers)
func SetupRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 中间件设置
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestLogger())

	// Swagger API文档 (开发环境)
	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 健康检查 (无需认证)
	healthHandler := handlers.NewHealthHandler(db)
	router.GET("/api/health", healthHandler.Health)
	router.GET("/api/health/db", healthHandler.DatabaseHealth)

	// 初始化微信客户端
	wechatConfig := config.NewWeChatConfig()
	wechatClient := wechat.NewClient(wechatConfig.AppID, wechatConfig.AppSecret)

	// 基于AccountController的路由
	accountHandler := handlers.NewAccountHandler(db, wechatClient)
	account := router.Group("/api/Account")
	{
		// 公开接口
		account.GET("/CheckUserInfo", accountHandler.CheckUserInfo)
		account.POST("/WeChatLogin", accountHandler.WeChatLogin)

		// 需要认证的接口
		account.GET("/CheckLogin", middleware.JWTAuth(), accountHandler.CheckLogin)
		account.GET("/GetUserInfo", middleware.JWTAuth(), accountHandler.GetUserInfo)
	}

	// 基于MemberController的路由
	memberHandler := handlers.NewMemberHandler(db)
	member := router.Group("/api/Member")
	member.Use(middleware.JWTAuth()) // 所有Member接口都需要认证
	{
		member.POST("/Update", memberHandler.Update)
		member.POST("/AddFranchiseIntention", memberHandler.AddFranchiseIntention)
		member.GET("/GetUserInfo", memberHandler.GetUserInfo)
	}

	// 基于MachineController的路由
	machineHandler := handlers.NewMachineHandler(db)
	machine := router.Group("/api/Machine")
	{
		// 公开接口
		machine.GET("/Get", machineHandler.Get)
		machine.GET("/CheckDeviceExist", machineHandler.CheckDeviceExist)
		machine.GET("/GetProductList", machineHandler.GetProductList)

		// 需要认证的接口
		machine.POST("/GetPaging", middleware.JWTAuth(), machineHandler.GetPaging)
		machine.GET("/GetList", middleware.JWTAuth(), machineHandler.GetList)
		machine.GET("/OpenOrClose", middleware.JWTAuth(), machineHandler.OpenOrCloseBusiness)
	}

	// 基于OrderController的路由
	// 初始化订单服务所需的依赖
	orderRepo := repositories.NewOrderRepository(db)
	machineRepo := repositories.NewMachineRepository(db)
	memberRepo := repositories.NewMemberRepository(db)
	deviceSvc := services.NewDeviceService()
	orderService := services.NewOrderService(orderRepo, machineRepo, memberRepo, deviceSvc)
	orderHandler := handlers.NewOrderHandler(db, orderService)
	order := router.Group("/api/Order")
	order.Use(middleware.JWTAuth()) // 所有Order接口都需要认证
	{
		order.POST("/GetPaging", orderHandler.GetPaging)
		order.GET("/Get", orderHandler.Get)
		order.POST("/Create", orderHandler.Create)
		order.POST("/Refund", orderHandler.Refund)
	}

	// 基于PaymentController的路由
	paymentHandler := handlers.NewPaymentHandler(db)
	payment := router.Group("/api/Payment")
	payment.Use(middleware.JWTAuth()) // 所有Payment接口都需要认证
	{
		payment.GET("/Get", paymentHandler.Get)
		payment.GET("/Query", paymentHandler.Query)
	}

	// 基于ProductController的路由
	productHandler := handlers.NewProductHandler(db)
	product := router.Group("/api/Product")
	{
		// 公开接口
		product.GET("/GetSelectList", productHandler.GetSelectList)
	}

	// 基于MaterialSiloController的路由 (物料槽管理)
	materialSiloHandler := handlers.NewMaterialSiloHandler(db)
	materialSilo := router.Group("/api/MaterialSilo")
	materialSilo.Use(middleware.JWTAuth()) // 物料槽管理需要认证
	{
		materialSilo.POST("/GetPaging", materialSiloHandler.GetPaging)
		materialSilo.POST("/UpdateStock", materialSiloHandler.UpdateStock)
		materialSilo.POST("/UpdateProduct", materialSiloHandler.UpdateProduct)
		materialSilo.POST("/ToggleSaleStatus", materialSiloHandler.ToggleSaleStatus)
	}

	// 基于MachineOwnerController的路由 (机主管理功能)
	machineOwnerHandler := handlers.NewMachineOwnerHandler(db)
	machineOwner := router.Group("/api/MachineOwner")
	machineOwner.Use(middleware.JWTAuth()) // 所有机主接口都需要认证
	{
		machineOwner.GET("/GetSales", machineOwnerHandler.GetSales)
		machineOwner.GET("/GetSalesStats", machineOwnerHandler.GetSalesStats)
	}

	// 基于CallbackController的路由 (无需认证)
	// 初始化回调服务所需的依赖
	paymentService := services.NewPaymentService(db)
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	callbackHandler := handlers.NewCallbackHandler(orderService, paymentService, logger)
	router.POST("/api/Callback/PaymentResult", callbackHandler.PaymentResult)

	return router
}
