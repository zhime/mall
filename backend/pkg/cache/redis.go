package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"mall/pkg/config"
	"mall/pkg/logger"
)

var RDB *redis.Client

// cacheService Redis缓存服务实现
type cacheService struct {
	rdb *redis.Client
}

// NewCacheService 创建缓存服务
func NewCacheService() CacheService {
	return &cacheService{
		rdb: RDB,
	}
}

// Set 实现CacheService接口
func (c *cacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

// Get 实现CacheService接口
func (c *cacheService) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Delete 实现CacheService接口
func (c *cacheService) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// Exists 实现CacheService接口
func (c *cacheService) Exists(ctx context.Context, key string) (int64, error) {
	return c.rdb.Exists(ctx, key).Result()
}

// InitRedis 初始化Redis连接
func InitRedis() {
	cfg := config.GetConfig()

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Info("Redis connected successfully")
}

// GetRedis 获取Redis客户端
func GetRedis() *redis.Client {
	return RDB
}

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(ctx context.Context, key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// Del 删除缓存
func Del(ctx context.Context, keys ...string) error {
	return RDB.Del(ctx, keys...).Err()
}

// Exists 检查key是否存在
func Exists(ctx context.Context, keys ...string) (int64, error) {
	return RDB.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RDB.Expire(ctx, key, expiration).Err()
}

// HSet 设置hash字段
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return RDB.HSet(ctx, key, values...).Err()
}

// HGet 获取hash字段
func HGet(ctx context.Context, key, field string) (string, error) {
	return RDB.HGet(ctx, key, field).Result()
}

// HGetAll 获取hash所有字段
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return RDB.HGetAll(ctx, key).Result()
}

// HDel 删除hash字段
func HDel(ctx context.Context, key string, fields ...string) error {
	return RDB.HDel(ctx, key, fields...).Err()
}

// LPush 从左侧推入列表
func LPush(ctx context.Context, key string, values ...interface{}) error {
	return RDB.LPush(ctx, key, values...).Err()
}

// RPop 从右侧弹出列表元素
func RPop(ctx context.Context, key string) (string, error) {
	return RDB.RPop(ctx, key).Result()
}

// SAdd 添加集合成员
func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return RDB.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(ctx context.Context, key string) ([]string, error) {
	return RDB.SMembers(ctx, key).Result()
}

// ZAdd 添加有序集合成员
func ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return RDB.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围内的成员
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return RDB.ZRange(ctx, key, start, stop).Result()
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RDB != nil {
		RDB.Close()
		logger.Info("Redis connection closed")
	}
}