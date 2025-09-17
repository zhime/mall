package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// OrderItemRepository 订单项仓储接口
type OrderItemRepository interface {
	Create(item *model.OrderItem) error
	BatchCreate(items []*model.OrderItem) error
	GetByOrderID(orderID uint64) ([]*model.OrderItem, error)
	Update(item *model.OrderItem) error
	Delete(id uint64) error
}

// orderItemRepository 订单项仓储实现
type orderItemRepository struct {
	db *gorm.DB
}

// NewOrderItemRepository 创建订单项仓储
func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepository{db: db}
}

// Create 创建订单项
func (r *orderItemRepository) Create(item *model.OrderItem) error {
	return r.db.Create(item).Error
}

// BatchCreate 批量创建订单项
func (r *orderItemRepository) BatchCreate(items []*model.OrderItem) error {
	return r.db.Create(items).Error
}

// GetByOrderID 根据订单ID获取订单项
func (r *orderItemRepository) GetByOrderID(orderID uint64) ([]*model.OrderItem, error) {
	var items []*model.OrderItem
	err := r.db.Preload("Product").
		Preload("SKU").
		Where("order_id = ?", orderID).
		Find(&items).Error
	return items, err
}

// Update 更新订单项
func (r *orderItemRepository) Update(item *model.OrderItem) error {
	return r.db.Save(item).Error
}

// Delete 删除订单项
func (r *orderItemRepository) Delete(id uint64) error {
	return r.db.Delete(&model.OrderItem{}, id).Error
}