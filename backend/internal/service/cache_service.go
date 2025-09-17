package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mall/pkg/cache"
	"mall/pkg/logger"
)

// CacheService 缓存服务接口
type CacheService interface {
	// 用户相关缓存
	SetUserInfo(userID uint64, userInfo interface{}, expiration time.Duration) error
	GetUserInfo(userID uint64) (string, error)
	DeleteUserInfo(userID uint64) error
	
	// 商品相关缓存
	SetProductInfo(productID uint64, productInfo interface{}, expiration time.Duration) error
	GetProductInfo(productID uint64) (string, error)
	DeleteProductInfo(productID uint64) error
	SetProductList(key string, products interface{}, expiration time.Duration) error
	GetProductList(key string) (string, error)
	
	// 分类相关缓存
	SetCategoryTree(categories interface{}, expiration time.Duration) error
	GetCategoryTree() (string, error)
	DeleteCategoryTree() error
	
	// 购物车缓存
	SetCartItems(userID uint64, cartItems interface{}, expiration time.Duration) error
	GetCartItems(userID uint64) (string, error)
	DeleteCartItems(userID uint64) error
	
	// 订单缓存
	SetOrderInfo(orderID uint64, orderInfo interface{}, expiration time.Duration) error
	GetOrderInfo(orderID uint64) (string, error)
	DeleteOrderInfo(orderID uint64) error
	
	// 热门数据缓存
	SetHotProducts(products interface{}, expiration time.Duration) error
	GetHotProducts() (string, error)
	SetHotCategories(categories interface{}, expiration time.Duration) error
	GetHotCategories() (string, error)
	
	// 搜索缓存
	SetSearchResult(keyword string, result interface{}, expiration time.Duration) error
	GetSearchResult(keyword string) (string, error)
	
	// 统计数据缓存
	SetStatistics(key string, data interface{}, expiration time.Duration) error
	GetStatistics(key string) (string, error)
	
	// 通用缓存方法
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(keys ...string) error
	Exists(keys ...string) (int64, error)
}

// cacheService 缓存服务实现
type cacheService struct{}

// NewCacheService 创建缓存服务
func NewCacheService() CacheService {
	return &cacheService{}
}

// 缓存key常量
const (
	// 用户相关
	UserInfoKey = "user:info:%d"
	UserTokenKey = "user:token:%d"
	
	// 商品相关
	ProductInfoKey = "product:info:%d"
	ProductListKey = "product:list:%s"
	HotProductsKey = "product:hot"
	
	// 分类相关
	CategoryTreeKey = "category:tree"
	HotCategoriesKey = "category:hot"
	
	// 购物车
	CartItemsKey = "cart:items:%d"
	
	// 订单相关
	OrderInfoKey = "order:info:%d"
	
	// 搜索相关
	SearchResultKey = "search:result:%s"
	
	// 统计数据
	StatisticsKey = "statistics:%s"
)

// SetUserInfo 设置用户信息缓存
func (s *cacheService) SetUserInfo(userID uint64, userInfo interface{}, expiration time.Duration) error {
	key := fmt.Sprintf(UserInfoKey, userID)
	return s.setJSON(key, userInfo, expiration)
}

// GetUserInfo 获取用户信息缓存
func (s *cacheService) GetUserInfo(userID uint64) (string, error) {
	key := fmt.Sprintf(UserInfoKey, userID)
	return cache.Get(context.Background(), key)
}

// DeleteUserInfo 删除用户信息缓存
func (s *cacheService) DeleteUserInfo(userID uint64) error {
	key := fmt.Sprintf(UserInfoKey, userID)
	return cache.Del(context.Background(), key)
}

// SetProductInfo 设置商品信息缓存
func (s *cacheService) SetProductInfo(productID uint64, productInfo interface{}, expiration time.Duration) error {
	key := fmt.Sprintf(ProductInfoKey, productID)
	return s.setJSON(key, productInfo, expiration)
}

// GetProductInfo 获取商品信息缓存
func (s *cacheService) GetProductInfo(productID uint64) (string, error) {
	key := fmt.Sprintf(ProductInfoKey, productID)
	return cache.Get(context.Background(), key)
}

// DeleteProductInfo 删除商品信息缓存
func (s *cacheService) DeleteProductInfo(productID uint64) error {
	key := fmt.Sprintf(ProductInfoKey, productID)
	return cache.Del(context.Background(), key)
}

// SetProductList 设置商品列表缓存
func (s *cacheService) SetProductList(key string, products interface{}, expiration time.Duration) error {
	cacheKey := fmt.Sprintf(ProductListKey, key)
	return s.setJSON(cacheKey, products, expiration)
}

// GetProductList 获取商品列表缓存
func (s *cacheService) GetProductList(key string) (string, error) {
	cacheKey := fmt.Sprintf(ProductListKey, key)
	return cache.Get(context.Background(), cacheKey)
}

// SetCategoryTree 设置分类树缓存
func (s *cacheService) SetCategoryTree(categories interface{}, expiration time.Duration) error {
	return s.setJSON(CategoryTreeKey, categories, expiration)
}

// GetCategoryTree 获取分类树缓存
func (s *cacheService) GetCategoryTree() (string, error) {
	return cache.Get(context.Background(), CategoryTreeKey)
}

// DeleteCategoryTree 删除分类树缓存
func (s *cacheService) DeleteCategoryTree() error {
	return cache.Del(context.Background(), CategoryTreeKey)
}

// SetCartItems 设置购物车缓存
func (s *cacheService) SetCartItems(userID uint64, cartItems interface{}, expiration time.Duration) error {
	key := fmt.Sprintf(CartItemsKey, userID)
	return s.setJSON(key, cartItems, expiration)
}

// GetCartItems 获取购物车缓存
func (s *cacheService) GetCartItems(userID uint64) (string, error) {
	key := fmt.Sprintf(CartItemsKey, userID)
	return cache.Get(context.Background(), key)
}

// DeleteCartItems 删除购物车缓存
func (s *cacheService) DeleteCartItems(userID uint64) error {
	key := fmt.Sprintf(CartItemsKey, userID)
	return cache.Del(context.Background(), key)
}

// SetOrderInfo 设置订单信息缓存
func (s *cacheService) SetOrderInfo(orderID uint64, orderInfo interface{}, expiration time.Duration) error {
	key := fmt.Sprintf(OrderInfoKey, orderID)
	return s.setJSON(key, orderInfo, expiration)
}

// GetOrderInfo 获取订单信息缓存
func (s *cacheService) GetOrderInfo(orderID uint64) (string, error) {
	key := fmt.Sprintf(OrderInfoKey, orderID)
	return cache.Get(context.Background(), key)
}

// DeleteOrderInfo 删除订单信息缓存
func (s *cacheService) DeleteOrderInfo(orderID uint64) error {
	key := fmt.Sprintf(OrderInfoKey, orderID)
	return cache.Del(context.Background(), key)
}

// SetHotProducts 设置热门商品缓存
func (s *cacheService) SetHotProducts(products interface{}, expiration time.Duration) error {
	return s.setJSON(HotProductsKey, products, expiration)
}

// GetHotProducts 获取热门商品缓存
func (s *cacheService) GetHotProducts() (string, error) {
	return cache.Get(context.Background(), HotProductsKey)
}

// SetHotCategories 设置热门分类缓存
func (s *cacheService) SetHotCategories(categories interface{}, expiration time.Duration) error {
	return s.setJSON(HotCategoriesKey, categories, expiration)
}

// GetHotCategories 获取热门分类缓存
func (s *cacheService) GetHotCategories() (string, error) {
	return cache.Get(context.Background(), HotCategoriesKey)
}

// SetSearchResult 设置搜索结果缓存
func (s *cacheService) SetSearchResult(keyword string, result interface{}, expiration time.Duration) error {
	key := fmt.Sprintf(SearchResultKey, keyword)
	return s.setJSON(key, result, expiration)
}

// GetSearchResult 获取搜索结果缓存
func (s *cacheService) GetSearchResult(keyword string) (string, error) {
	key := fmt.Sprintf(SearchResultKey, keyword)
	return cache.Get(context.Background(), key)
}

// SetStatistics 设置统计数据缓存
func (s *cacheService) SetStatistics(key string, data interface{}, expiration time.Duration) error {
	cacheKey := fmt.Sprintf(StatisticsKey, key)
	return s.setJSON(cacheKey, data, expiration)
}

// GetStatistics 获取统计数据缓存
func (s *cacheService) GetStatistics(key string) (string, error) {
	cacheKey := fmt.Sprintf(StatisticsKey, key)
	return cache.Get(context.Background(), cacheKey)
}

// Set 通用设置缓存
func (s *cacheService) Set(key string, value interface{}, expiration time.Duration) error {
	return cache.Set(context.Background(), key, value, expiration)
}

// Get 通用获取缓存
func (s *cacheService) Get(key string) (string, error) {
	return cache.Get(context.Background(), key)
}

// Delete 通用删除缓存
func (s *cacheService) Delete(keys ...string) error {
	return cache.Del(context.Background(), keys...)
}

// Exists 检查缓存是否存在
func (s *cacheService) Exists(keys ...string) (int64, error) {
	return cache.Exists(context.Background(), keys...)
}

// setJSON 设置JSON格式缓存
func (s *cacheService) setJSON(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		logger.Error("Failed to marshal JSON for cache", logger.String("key", key), logger.Error(err.Error()))
		return err
	}
	
	return cache.Set(context.Background(), key, string(jsonData), expiration)
}