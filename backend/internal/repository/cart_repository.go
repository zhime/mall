package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// CartRepository 购物车仓储接口
type CartRepository interface {
	AddItem(item *model.CartItem) error
	UpdateItem(item *model.CartItem) error
	RemoveItem(userID, productID, skuID uint64) error
	GetUserCartItems(userID uint64) ([]*model.CartItem, error)
	GetCartItem(userID, productID, skuID uint64) (*model.CartItem, error)
	ClearUserCart(userID uint64) error
	RemoveItems(userID uint64, itemIDs []uint64) error
}

// cartRepository 购物车仓储实现
type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository 创建购物车仓储
func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

// AddItem 添加商品到购物车
func (r *cartRepository) AddItem(item *model.CartItem) error {
	// 检查是否已存在相同商品
	var existingItem model.CartItem
	result := r.db.Where("user_id = ? AND product_id = ? AND sku_id = ?", 
		item.UserID, item.ProductID, item.SKUID).First(&existingItem)
	
	if result.Error == nil {
		// 已存在，更新数量
		existingItem.Quantity += item.Quantity
		return r.db.Save(&existingItem).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		return r.db.Create(item).Error
	}
	
	return result.Error
}

// UpdateItem 更新购物车商品
func (r *cartRepository) UpdateItem(item *model.CartItem) error {
	return r.db.Save(item).Error
}

// RemoveItem 从购物车移除商品
func (r *cartRepository) RemoveItem(userID, productID, skuID uint64) error {
	return r.db.Where("user_id = ? AND product_id = ? AND sku_id = ?", 
		userID, productID, skuID).Delete(&model.CartItem{}).Error
}

// GetUserCartItems 获取用户购物车商品列表
func (r *cartRepository) GetUserCartItems(userID uint64) ([]*model.CartItem, error) {
	var items []*model.CartItem
	err := r.db.Preload("Product").
		Preload("Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_main = ?", 1)
		}).
		Preload("SKU").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// GetCartItem 获取购物车单个商品
func (r *cartRepository) GetCartItem(userID, productID, skuID uint64) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.Where("user_id = ? AND product_id = ? AND sku_id = ?", 
		userID, productID, skuID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ClearUserCart 清空用户购物车
func (r *cartRepository) ClearUserCart(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.CartItem{}).Error
}

// RemoveItems 批量移除购物车商品
func (r *cartRepository) RemoveItems(userID uint64, itemIDs []uint64) error {
	return r.db.Where("user_id = ? AND id IN ?", userID, itemIDs).Delete(&model.CartItem{}).Error
}