package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"mall/pkg/logger"
)

// CacheService 缓存服务接口
type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (int64, error)
}

// CacheManager 缓存管理器
type CacheManager struct {
	cacheService CacheService
}

// CachePolicy 缓存策略
type CachePolicy struct {
	TTL           time.Duration // 过期时间
	RefreshBefore time.Duration // 提前刷新时间
	MaxSize       int64         // 最大大小限制
}

// 预定义缓存策略
var (
	// 短期缓存 - 5分钟
	ShortTermPolicy = CachePolicy{
		TTL:           5 * time.Minute,
		RefreshBefore: 1 * time.Minute,
		MaxSize:       1024 * 1024, // 1MB
	}
	
	// 中期缓存 - 30分钟
	MediumTermPolicy = CachePolicy{
		TTL:           30 * time.Minute,
		RefreshBefore: 5 * time.Minute,
		MaxSize:       5 * 1024 * 1024, // 5MB
	}
	
	// 长期缓存 - 2小时
	LongTermPolicy = CachePolicy{
		TTL:           2 * time.Hour,
		RefreshBefore: 15 * time.Minute,
		MaxSize:       10 * 1024 * 1024, // 10MB
	}
	
	// 持久缓存 - 24小时
	PersistentPolicy = CachePolicy{
		TTL:           24 * time.Hour,
		RefreshBefore: 1 * time.Hour,
		MaxSize:       50 * 1024 * 1024, // 50MB
	}
)

// 缓存key模板
const (
	// 用户相关
	UserProfileCacheKey = "user:profile:%d"
	UserPermissionsCacheKey = "user:permissions:%d"
	
	// 商品相关
	ProductDetailCacheKey = "product:detail:%d"
	ProductListCacheKey = "product:list:%s"
	CategoryListCacheKey = "category:list"
	
	// 购物车
	UserCartCacheKey = "cart:user:%d"
	
	// 统计数据
	DailyStatsCacheKey = "stats:daily:%s"
	ProductStatsCacheKey = "stats:product:%d"
	
	// 配置信息
	SystemConfigCacheKey = "config:system"
	
	// 热点数据
	HotSearchKeywordsCacheKey = "hot:search:keywords"
	RecommendProductsCacheKey = "recommend:products:%d"
)

// NewCacheManager 创建缓存管理器
func NewCacheManager() *CacheManager {
	return &CacheManager{
		cacheService: NewCacheService(),
	}
}

// CacheWithPolicy 根据策略设置缓存
func (cm *CacheManager) CacheWithPolicy(key string, value interface{}, policy CachePolicy) error {
	// 检查数据大小
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if int64(len(jsonData)) > policy.MaxSize {
		logger.Warn("Data size exceeds cache policy limit", 
			zap.String("key", key),
			zap.Int("data_size", len(jsonData)),
			zap.Int64("max_size", policy.MaxSize))
		return fmt.Errorf("data size exceeds cache policy limit")
	}
	
	// 设置缓存
	ctx := context.Background()
	return Set(ctx, key, string(jsonData), policy.TTL)
}

// GetCachedData 获取缓存数据
func (cm *CacheManager) GetCachedData(key string) (string, error) {
	ctx := context.Background()
	return Get(ctx, key)
}

// InvalidateCache 使缓存失效
func (cm *CacheManager) InvalidateCache(pattern string) error {
	ctx := context.Background()
	
	// 获取匹配的key
	keys, err := cm.getKeysByPattern(pattern)
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return Del(ctx, keys...)
	}
	
	return nil
}

// InvalidateUserCache 使用户相关缓存失效
func (cm *CacheManager) InvalidateUserCache(userID uint64) error {
	keys := []string{
		fmt.Sprintf(UserProfileCacheKey, userID),
		fmt.Sprintf(UserPermissionsCacheKey, userID),
		fmt.Sprintf(UserCartCacheKey, userID),
	}
	
	ctx := context.Background()
	return Del(ctx, keys...)
}

// InvalidateProductCache 使商品相关缓存失效
func (cm *CacheManager) InvalidateProductCache(productID uint64) error {
	keys := []string{
		fmt.Sprintf(ProductDetailCacheKey, productID),
		fmt.Sprintf(ProductStatsCacheKey, productID),
	}
	
	// 清除商品列表缓存
	if err := cm.InvalidateCache("product:list:*"); err != nil {
		logger.Error("Failed to invalidate product list cache", zap.Error(err))
	}
	
	ctx := context.Background()
	return Del(ctx, keys...)
}

// RefreshCache 刷新缓存
func (cm *CacheManager) RefreshCache(key string, refreshFunc func() (interface{}, error), policy CachePolicy) error {
	// 执行刷新函数获取新数据
	newData, err := refreshFunc()
	if err != nil {
		return fmt.Errorf("failed to refresh data: %w", err)
	}
	
	// 设置新的缓存
	return cm.CacheWithPolicy(key, newData, policy)
}

// WarmupCache 预热缓存
func (cm *CacheManager) WarmupCache() error {
	logger.Info("Starting cache warmup...")
	
	// 预热分类数据
	if err := cm.warmupCategoryCache(); err != nil {
		logger.Error("Failed to warmup category cache", zap.Error(err))
	}
	
	// 预热热门商品
	if err := cm.warmupHotProductsCache(); err != nil {
		logger.Error("Failed to warmup hot products cache", zap.Error(err))
	}
	
	// 预热系统配置
	if err := cm.warmupSystemConfigCache(); err != nil {
		logger.Error("Failed to warmup system config cache", zap.Error(err))
	}
	
	logger.Info("Cache warmup completed")
	return nil
}

// warmupCategoryCache 预热分类缓存
func (cm *CacheManager) warmupCategoryCache() error {
	// TODO: 从数据库加载分类数据
	// 这里模拟分类数据
	categoryData := map[string]interface{}{
		"categories": []map[string]interface{}{
			{"id": 1, "name": "服装鞋包", "level": 1},
			{"id": 2, "name": "数码家电", "level": 1},
		},
		"updated_at": time.Now(),
	}
	
	return cm.CacheWithPolicy(CategoryListCacheKey, categoryData, LongTermPolicy)
}

// warmupHotProductsCache 预热热门商品缓存
func (cm *CacheManager) warmupHotProductsCache() error {
	// TODO: 从数据库加载热门商品数据
	// 这里模拟热门商品数据
	hotProducts := map[string]interface{}{
		"products": []map[string]interface{}{
			{"id": 1, "name": "热门商品1", "sales": 1000},
			{"id": 2, "name": "热门商品2", "sales": 800},
		},
		"updated_at": time.Now(),
	}
	
	return cm.CacheWithPolicy(RecommendProductsCacheKey, hotProducts, MediumTermPolicy)
}

// warmupSystemConfigCache 预热系统配置缓存
func (cm *CacheManager) warmupSystemConfigCache() error {
	// TODO: 从数据库加载系统配置
	// 这里模拟系统配置数据
	systemConfig := map[string]interface{}{
		"site_name": "Mall商城",
		"version": "1.0.0",
		"maintenance_mode": false,
		"updated_at": time.Now(),
	}
	
	return cm.CacheWithPolicy(SystemConfigCacheKey, systemConfig, PersistentPolicy)
}

// getKeysByPattern 根据模式获取keys
func (cm *CacheManager) getKeysByPattern(pattern string) ([]string, error) {
	// Redis的KEYS命令，生产环境建议使用SCAN
	ctx := context.Background()
	rdb := GetRedis()
	
	keys, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	
	return keys, nil
}

// CacheStats 缓存统计信息
type CacheStats struct {
	TotalKeys    int64 `json:"total_keys"`
	UsedMemory   int64 `json:"used_memory"`
	HitRate      float64 `json:"hit_rate"`
	MissRate     float64 `json:"miss_rate"`
}

// GetCacheStats 获取缓存统计信息
func (cm *CacheManager) GetCacheStats() (*CacheStats, error) {
	ctx := context.Background()
	rdb := GetRedis()
	
	// 获取Redis info
	info, err := rdb.Info(ctx, "stats", "memory").Result()
	if err != nil {
		return nil, err
	}
	
	stats := &CacheStats{}
	
	// 解析info信息
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, "keyspace_hits:") {
			// 解析命中率等统计信息
			// 这里简化处理
			stats.HitRate = 0.85 // 模拟数据
			stats.MissRate = 0.15
		}
	}
	
	// 获取总key数量
	keys, err := rdb.DBSize(ctx).Result()
	if err == nil {
		stats.TotalKeys = keys
	}
	
	return stats, nil
}

// CleanupExpiredCache 清理过期缓存
func (cm *CacheManager) CleanupExpiredCache() error {
	logger.Info("Starting cleanup of expired cache entries...")
	
	// Redis会自动清理过期的key，这里主要用于统计和日志
	stats, err := cm.GetCacheStats()
	if err != nil {
		return err
	}
	
	logger.Info("Cache cleanup completed", 
		zap.Int64("total_keys", stats.TotalKeys),
		zap.Float64("hit_rate", stats.HitRate))
		
	return nil
}

// SetCacheTag 设置缓存标签（用于批量失效）
func (cm *CacheManager) SetCacheTag(key, tag string, expiration time.Duration) error {
	tagKey := fmt.Sprintf("tag:%s", tag)
	ctx := context.Background()
	
	// 将key添加到标签集合中
	return SAdd(ctx, tagKey, key)
}

// InvalidateCacheByTag 根据标签使缓存失效
func (cm *CacheManager) InvalidateCacheByTag(tag string) error {
	tagKey := fmt.Sprintf("tag:%s", tag)
	ctx := context.Background()
	
	// 获取标签下的所有key
	keys, err := SMembers(ctx, tagKey)
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		// 删除所有相关的缓存
		if err := Del(ctx, keys...); err != nil {
			return err
		}
		
		// 删除标签集合
		return Del(ctx, tagKey)
	}
	
	return nil
}