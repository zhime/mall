package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"mall/internal/handler"
	"mall/internal/middleware"
)

// Handlers 处理器容器
type Handlers struct {
	AuthHandler     *handler.AuthHandler
	CategoryHandler *handler.CategoryHandler
	ProductHandler  *handler.ProductHandler
	CartHandler     *handler.CartHandler
	OrderHandler    *handler.OrderHandler
	PaymentHandler  *handler.PaymentHandler
}

// SetupRouter 设置路由
func SetupRouter(authHandler *handler.AuthHandler, categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler, cartHandler *handler.CartHandler, orderHandler *handler.OrderHandler, paymentHandler *handler.PaymentHandler) *gin.Engine {
	router := gin.New()

	// 设置中间件
	setupMiddleware(router)

	// 注册健康检查路由
	setupHealthCheck(router)

	// 注册API路由
	handlers := &Handlers{
		AuthHandler:     authHandler,
		CategoryHandler: categoryHandler,
		ProductHandler:  productHandler,
		CartHandler:     cartHandler,
		OrderHandler:    orderHandler,
		PaymentHandler:  paymentHandler,
	}
	registerAPIRoutes(router, handlers)

	return router
}

// setupMiddleware 设置中间件
func setupMiddleware(router *gin.Engine) {
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS())
	router.Use(middleware.APIRateLimiter())
}

// setupHealthCheck 设置健康检查路由
func setupHealthCheck(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
}

// registerAPIRoutes 注册API路由
func registerAPIRoutes(router *gin.Engine, handlers *Handlers) {
	// API版本分组
	v1 := router.Group("/api/v1")

	// 创建路由组实例
	authRoutes := NewAuthRoutes(handlers.AuthHandler)
	userRoutes := NewUserRoutes(handlers.AuthHandler)
	productRoutes := NewProductRoutes(handlers.ProductHandler, handlers.CategoryHandler)
	orderRoutes := NewOrderRoutes(handlers.OrderHandler, handlers.CartHandler, handlers.PaymentHandler)
	adminRoutes := NewAdminRoutes(handlers.CategoryHandler, handlers.ProductHandler, handlers.OrderHandler, handlers.PaymentHandler)

	// 注册路由组
	authRoutes.RegisterRoutes(v1)
	userRoutes.RegisterRoutes(v1)
	productRoutes.RegisterRoutes(v1)
	orderRoutes.RegisterRoutes(v1)
	adminRoutes.RegisterRoutes(v1)
}