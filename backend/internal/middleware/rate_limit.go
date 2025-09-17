package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mall/pkg/cache"
	"mall/pkg/utils"
)

// RateLimiter 限流中间件
func RateLimiter(maxRequests int, windowDuration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		
		// 生成限流key
		key := fmt.Sprintf("rate_limit:%s", clientIP)
		
		ctx := context.Background()
		
		// 获取当前请求次数
		countStr, err := cache.Get(ctx, key)
		if err != nil {
			// 如果key不存在，初始化为1
			cache.Set(ctx, key, "1", windowDuration)
			c.Next()
			return
		}
		
		count, err := strconv.Atoi(countStr)
		if err != nil {
			// 解析错误，重新初始化
			cache.Set(ctx, key, "1", windowDuration)
			c.Next()
			return
		}
		
		// 检查是否超过限制
		if count >= maxRequests {
			utils.Error(c, utils.TOO_MANY_REQUESTS, "Rate limit exceeded")
			c.Abort()
			return
		}
		
		// 增加计数
		cache.Set(ctx, key, strconv.Itoa(count+1), windowDuration)
		
		c.Next()
	}
}

// APIRateLimiter API限流中间件
func APIRateLimiter() gin.HandlerFunc {
	// 默认每分钟100次请求
	return RateLimiter(100, time.Minute)
}

// LoginRateLimiter 登录限流中间件
func LoginRateLimiter() gin.HandlerFunc {
	// 登录接口每分钟5次请求
	return RateLimiter(5, time.Minute)
}

// SMSRateLimiter 短信限流中间件
func SMSRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取手机号
		phone := c.PostForm("phone")
		if phone == "" {
			phone = c.Query("phone")
		}
		
		if phone == "" {
			utils.InvalidParams(c, "Phone number is required")
			c.Abort()
			return
		}
		
		// 生成限流key
		key := fmt.Sprintf("sms_rate_limit:%s", phone)
		
		ctx := context.Background()
		
		// 检查是否在限制时间内
		exists, err := cache.Exists(ctx, key)
		if err != nil {
			utils.InternalServerError(c, "Rate limit check failed")
			c.Abort()
			return
		}
		
		if exists > 0 {
			utils.Error(c, utils.TOO_MANY_REQUESTS, "SMS request too frequent, please try again later")
			c.Abort()
			return
		}
		
		// 设置限制（1分钟内不能重复发送）
		cache.Set(ctx, key, "1", time.Minute)
		
		c.Next()
	}
}