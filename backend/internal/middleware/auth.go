package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"mall/pkg/utils"
)

// JWT认证中间件
func JWT() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// 检查Bearer token格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		// 提取token
		tokenString := authHeader[7:] // 移除"Bearer "前缀
		if tokenString == "" {
			utils.Unauthorized(c, "Missing token")
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// 将用户信息设置到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("platform", claims.Platform)

		c.Next()
	})
}

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 先执行JWT中间件
		JWT()(c)
		if c.IsAborted() {
			return
		}

		// 检查是否为管理员
		role := c.GetString("role")
		platform := c.GetString("platform")
		
		if platform != "admin" {
			utils.Forbidden(c, "Access denied: admin platform required")
			c.Abort()
			return
		}

		if role != "admin" && role != "super_admin" {
			utils.Forbidden(c, "Access denied: admin role required")
			c.Abort()
			return
		}

		c.Next()
	})
}

// SuperAdminAuth 超级管理员认证中间件
func SuperAdminAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 先执行JWT中间件
		JWT()(c)
		if c.IsAborted() {
			return
		}

		// 检查是否为超级管理员
		role := c.GetString("role")
		platform := c.GetString("platform")
		
		if platform != "admin" {
			utils.Forbidden(c, "Access denied: admin platform required")
			c.Abort()
			return
		}

		if role != "super_admin" {
			utils.Forbidden(c, "Access denied: super admin role required")
			c.Abort()
			return
		}

		c.Next()
	})
}

// UserAuth 用户认证中间件
func UserAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 先执行JWT中间件
		JWT()(c)
		if c.IsAborted() {
			return
		}

		// 检查平台类型
		platform := c.GetString("platform")
		allowedPlatforms := []string{"web", "miniprogram"}
		
		platformAllowed := false
		for _, p := range allowedPlatforms {
			if platform == p {
				platformAllowed = true
				break
			}
		}

		if !platformAllowed {
			utils.Forbidden(c, "Access denied: invalid platform")
			c.Abort()
			return
		}

		c.Next()
	})
}