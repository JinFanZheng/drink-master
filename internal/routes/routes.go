package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/config"
	"github.com/ddteam/drink-master/internal/handlers"
	"github.com/ddteam/drink-master/internal/middleware"
	"github.com/ddteam/drink-master/pkg/wechat"
)

// SetupRoutes 设置所有路由 (基于MobileAPI Controllers)
func SetupRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// 中间件设置
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestLogger())

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
	orderHandler := handlers.NewOrderHandler(db)
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
	product := router.Group("/api/products")
	{
		// 公开接口
		product.GET("/select", productHandler.GetSelectList)
	}

	// 基于MachineOwnerController的路由 (机主管理功能)
	machineOwnerHandler := handlers.NewMachineOwnerHandler(db)
	machineOwner := router.Group("/api/machine-owners")
	machineOwner.Use(middleware.JWTAuth()) // 所有机主接口都需要认证
	{
		machineOwner.GET("/sales", machineOwnerHandler.GetSales)
		machineOwner.GET("/sales/stats", machineOwnerHandler.GetSalesStats)
	}

	// 回调接口 (无需认证)
	callbackHandler := handlers.NewCallbackHandler(db)
	router.POST("/api/Callback/PaymentResult", callbackHandler.PaymentResult)

	return router
}
