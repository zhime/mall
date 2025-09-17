package routes

import (
	"github.com/gin-gonic/gin"

	"mall/internal/handler"
	"mall/internal/middleware"
)

// UserRoutes 用户路由组
type UserRoutes struct {
	authHandler *handler.AuthHandler
}

// NewUserRoutes 创建用户路由组
func NewUserRoutes(authHandler *handler.AuthHandler) *UserRoutes {
	return &UserRoutes{
		authHandler: authHandler,
	}
}

// RegisterRoutes 注册用户相关路由
func (r *UserRoutes) RegisterRoutes(router *gin.RouterGroup) {
	user := router.Group("/user")
	user.Use(middleware.UserAuth())
	{
		user.GET("/info", r.authHandler.GetUserInfo)
		user.GET("/profile", r.authHandler.GetProfile)
		user.PUT("/profile", r.authHandler.UpdateProfile)
		user.PUT("/password", r.authHandler.ChangePassword)
		user.POST("/bind-phone", r.authHandler.BindPhone)
	}
}