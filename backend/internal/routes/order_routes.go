package routes

import (
	"github.com/gin-gonic/gin"

	"mall/internal/handler"
	"mall/internal/middleware"
)

// OrderRoutes 订单路由组
type OrderRoutes struct {
	orderHandler   *handler.OrderHandler
	cartHandler    *handler.CartHandler
	paymentHandler *handler.PaymentHandler
}

// NewOrderRoutes 创建订单路由组
func NewOrderRoutes(orderHandler *handler.OrderHandler, cartHandler *handler.CartHandler, paymentHandler *handler.PaymentHandler) *OrderRoutes {
	return &OrderRoutes{
		orderHandler:   orderHandler,
		cartHandler:    cartHandler,
		paymentHandler: paymentHandler,
	}
}

// RegisterRoutes 注册订单和购物车相关路由
func (r *OrderRoutes) RegisterRoutes(router *gin.RouterGroup) {
	// 购物车相关路由
	cart := router.Group("/cart")
	cart.Use(middleware.UserAuth())
	{
		cart.GET("", r.cartHandler.GetUserCart)
		cart.GET("/count", r.cartHandler.GetCartCount)
		cart.POST("/add", r.cartHandler.AddToCart)
		cart.PUT("/update", r.cartHandler.UpdateCartItem)
		cart.DELETE("/remove", r.cartHandler.RemoveFromCart)
		cart.DELETE("/clear", r.cartHandler.ClearCart)
	}

	// 订单相关路由
	orders := router.Group("/orders")
	orders.Use(middleware.UserAuth())
	{
		orders.GET("", r.orderHandler.GetUserOrders)
		orders.POST("", r.orderHandler.CreateOrder)
		orders.GET("/:id", r.orderHandler.GetOrderDetail)
		orders.PUT("/:id/cancel", r.orderHandler.CancelOrder)
		orders.PUT("/:id/confirm", r.orderHandler.ConfirmOrder)
	}

	// 支付相关路由
	payment := router.Group("/payment")
	{
		payment.POST("", middleware.UserAuth(), r.paymentHandler.CreatePayment)
		payment.GET("/:paymentNo/status", r.paymentHandler.GetPaymentStatus)
		payment.PUT("/:paymentNo/cancel", middleware.UserAuth(), r.paymentHandler.CancelPayment)
		
		// 支付回调（无需认证）
		payment.POST("/wechat/callback", r.paymentHandler.WechatCallback)
		payment.POST("/alipay/callback", r.paymentHandler.AlipayCallback)
	}
}