package routes

import (
	"github.com/gin-gonic/gin"

	"mall/internal/handler"
	"mall/internal/middleware"
)

// AdminRoutes 管理后台路由组
type AdminRoutes struct {
	categoryHandler *handler.CategoryHandler
	productHandler  *handler.ProductHandler
	orderHandler    *handler.OrderHandler
	paymentHandler  *handler.PaymentHandler
}

// NewAdminRoutes 创建管理后台路由组
func NewAdminRoutes(categoryHandler *handler.CategoryHandler, productHandler *handler.ProductHandler, 
	orderHandler *handler.OrderHandler, paymentHandler *handler.PaymentHandler) *AdminRoutes {
	return &AdminRoutes{
		categoryHandler: categoryHandler,
		productHandler:  productHandler,
		orderHandler:    orderHandler,
		paymentHandler:  paymentHandler,
	}
}

// RegisterRoutes 注册管理后台相关路由
func (r *AdminRoutes) RegisterRoutes(router *gin.RouterGroup) {
	admin := router.Group("/admin")
	admin.Use(middleware.AdminAuth())
	{
		// 分类管理
		adminCategories := admin.Group("/categories")
		{
			adminCategories.POST("", r.categoryHandler.CreateCategory)
			adminCategories.PUT("/:id", r.categoryHandler.UpdateCategory)
			adminCategories.DELETE("/:id", r.categoryHandler.DeleteCategory)
		}

		// 商品管理
		adminProducts := admin.Group("/products")
		{
			adminProducts.GET("", r.productHandler.GetProductList)
			adminProducts.POST("", r.productHandler.CreateProduct)
			adminProducts.PUT("/:id", r.productHandler.UpdateProduct)
			adminProducts.DELETE("/:id", r.productHandler.DeleteProduct)
			adminProducts.PUT("/:id/stock", r.productHandler.UpdateProductStock)
		}

		// 订单管理
		adminOrders := admin.Group("/orders")
		{
			adminOrders.GET("", r.orderHandler.GetOrders)
			adminOrders.GET("/search", r.orderHandler.SearchOrders)
			adminOrders.PUT("/:id/status", r.orderHandler.UpdateOrderStatus)
		}

		// 支付管理
		adminPayment := admin.Group("/payment")
		{
			adminPayment.POST("/:paymentNo/refund", r.paymentHandler.RefundPayment)
		}
	}
}