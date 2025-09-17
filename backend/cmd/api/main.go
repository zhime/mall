package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"mall/internal/handler"
	"mall/internal/middleware"
	"mall/internal/repository"
	"mall/internal/service"
	"mall/pkg/cache"
	"mall/pkg/config"
	"mall/pkg/database"
	"mall/pkg/logger"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志
	logger.InitLogger(
		cfg.Log.Level,
		cfg.Log.Filename,
		cfg.Log.MaxSize,
		cfg.Log.MaxAge,
		cfg.Log.MaxBackups,
	)

	// 初始化数据库
	database.InitDB()
	defer database.CloseDB()

	// 初始化Redis
	cache.InitRedis()
	defer cache.CloseRedis()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	router := setupRouter()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器
	go func() {
		logger.Info(fmt.Sprintf("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server")
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")

	// 优雅关闭服务器，等待5秒钟
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown")
	}

	logger.Info("Server exited")
}

// setupRouter 设置路由
func setupRouter() *gin.Engine {
	router := gin.New()

	// 添加中间件
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(middleware.APIRateLimiter())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 初始化仓储
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	userAuthRepo := repository.NewUserAuthRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	productSKURepo := repository.NewProductSKURepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderItemRepo := repository.NewOrderItemRepository(db)
	paymentRepo := repository.NewOrderPaymentRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, userAuthRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo, productSKURepo)
	cartService := service.NewCartService(cartRepo, productRepo, productSKURepo)
	orderService := service.NewOrderService(orderRepo, orderItemRepo, productRepo, productSKURepo, cartRepo, userRepo)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// API版本分组
	v1 := router.Group("/api/v1")
	{
		// 认证相关路由
		auth := v1.Group("/auth")
		{
			auth.POST("/sms/send", middleware.SMSRateLimiter(), authHandler.SendSMSCode)
			auth.POST("/login/phone", middleware.LoginRateLimiter(), authHandler.LoginByPhone)
			auth.POST("/login/wechat", middleware.LoginRateLimiter(), authHandler.LoginByWechat)
			auth.POST("/login/password", middleware.LoginRateLimiter(), authHandler.LoginByPassword)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", middleware.JWT(), authHandler.Logout)
		}

		// 用户相关路由
		user := v1.Group("/user")
		user.Use(middleware.UserAuth())
		{
			user.GET("/info", authHandler.GetUserInfo)
			user.GET("/profile", authHandler.GetProfile)
			user.PUT("/profile", authHandler.UpdateProfile)
			user.PUT("/password", authHandler.ChangePassword)
			user.POST("/bind-phone", authHandler.BindPhone)
		}

		// 分类相关路由
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategoryTree)
			categories.GET("/level", categoryHandler.GetCategoriesByLevel)
			categories.GET("/search", categoryHandler.SearchCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
		}

		// 商品相关路由
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetProductList)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/hot", productHandler.GetHotProducts)
			products.GET("/category/:categoryId", productHandler.GetProductsByCategory)
			products.GET("/:id", productHandler.GetProduct)
			products.GET("/:id/detail", productHandler.GetProductDetail)
		}

		// 购物车相关路由
		cart := v1.Group("/cart")
		cart.Use(middleware.UserAuth())
		{
			cart.GET("", cartHandler.GetUserCart)
			cart.GET("/count", cartHandler.GetCartCount)
			cart.POST("/add", cartHandler.AddToCart)
			cart.PUT("/update", cartHandler.UpdateCartItem)
			cart.DELETE("/remove", cartHandler.RemoveFromCart)
			cart.DELETE("/clear", cartHandler.ClearCart)
		}

		// 订单相关路由
		orders := v1.Group("/orders")
		orders.Use(middleware.UserAuth())
		{
			orders.GET("", orderHandler.GetUserOrders)
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrderDetail)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			orders.PUT("/:id/confirm", orderHandler.ConfirmOrder)
		}

		// 支付相关路由
		payment := v1.Group("/payment")
		{
			payment.POST("", middleware.UserAuth(), paymentHandler.CreatePayment)
			payment.GET("/:paymentNo/status", paymentHandler.GetPaymentStatus)
			payment.PUT("/:paymentNo/cancel", middleware.UserAuth(), paymentHandler.CancelPayment)
			
			// 支付回调（无需认证）
			payment.POST("/wechat/callback", paymentHandler.WechatCallback)
			payment.POST("/alipay/callback", paymentHandler.AlipayCallback)
		}

		// 管理员相关路由
		admin := v1.Group("/admin")
		admin.Use(middleware.AdminAuth())
		{
			// 分类管理
			adminCategories := admin.Group("/categories")
			{
				adminCategories.POST("", categoryHandler.CreateCategory)
				adminCategories.PUT("/:id", categoryHandler.UpdateCategory)
				adminCategories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// 商品管理
			adminProducts := admin.Group("/products")
			{
				adminProducts.GET("", productHandler.GetProductList)
				adminProducts.POST("", productHandler.CreateProduct)
				adminProducts.PUT("/:id", productHandler.UpdateProduct)
				adminProducts.DELETE("/:id", productHandler.DeleteProduct)
				adminProducts.PUT("/:id/stock", productHandler.UpdateProductStock)
			}

			// 订单管理
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("", orderHandler.GetOrders)
				adminOrders.GET("/search", orderHandler.SearchOrders)
				adminOrders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			}

			// 支付管理
			adminPayment := admin.Group("/payment")
			{
				adminPayment.POST("/:paymentNo/refund", paymentHandler.RefundPayment)
			}
		}
	}

	return router
}