package routes

import "github.com/gin-gonic/gin"

// RouteGroup 定义路由组接口
type RouteGroup interface {
	RegisterRoutes(router *gin.RouterGroup)
}