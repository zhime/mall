package routes

import (
	"github.com/gin-gonic/gin"

	"mall/internal/handler"
	"mall/internal/middleware"
)

// AuthRoutes 认证路由组
type AuthRoutes struct {
	authHandler *handler.AuthHandler
}

// NewAuthRoutes 创建认证路由组
func NewAuthRoutes(authHandler *handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		authHandler: authHandler,
	}
}

// RegisterRoutes 注册认证相关路由
func (r *AuthRoutes) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/sms/send", middleware.SMSRateLimiter(), r.authHandler.SendSMSCode)
		auth.POST("/login/phone", middleware.LoginRateLimiter(), r.authHandler.LoginByPhone)
		auth.POST("/login/wechat", middleware.LoginRateLimiter(), r.authHandler.LoginByWechat)
		auth.POST("/login/password", middleware.LoginRateLimiter(), r.authHandler.LoginByPassword)
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		auth.POST("/logout", middleware.JWT(), r.authHandler.Logout)
	}
}