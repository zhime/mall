package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	Create(order *model.Order) error
	GetByID(id uint64) (*model.Order, error)
	GetByOrderNo(orderNo string) (*model.Order, error)
	GetWithDetails(id uint64) (*model.Order, error)
	Update(order *model.Order) error
	UpdateStatus(id uint64, status int8) error
	GetUserOrders(userID uint64, page, pageSize int, status int8) ([]*model.Order, int64, error)
	GetOrders(page, pageSize int, status int8) ([]*model.Order, int64, error)
	Search(keyword string, page, pageSize int) ([]*model.Order, int64, error)
}

// orderRepository 订单仓储实现
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create 创建订单
func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

// GetByID 根据ID获取订单
func (r *orderRepository) GetByID(id uint64) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号获取订单
func (r *orderRepository) GetByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetWithDetails 获取订单详情（包含订单项和支付信息）
func (r *orderRepository) GetWithDetails(id uint64) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.SKU").
		Preload("Payments").
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Update 更新订单
func (r *orderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus 更新订单状态
func (r *orderRepository) UpdateStatus(id uint64, status int8) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

// GetUserOrders 获取用户订单列表
func (r *orderRepository) GetUserOrders(userID uint64, page, pageSize int, status int8) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	query := r.db.Model(&model.Order{}).
		Preload("Items").
		Preload("Items.Product").
		Where("user_id = ?", userID)

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

// GetOrders 获取订单列表（管理员）
func (r *orderRepository) GetOrders(page, pageSize int, status int8) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	query := r.db.Model(&model.Order{}).
		Preload("User").
		Preload("Items").
		Preload("Items.Product")

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

// Search 搜索订单
func (r *orderRepository) Search(keyword string, page, pageSize int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	query := r.db.Model(&model.Order{}).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Where("order_no LIKE ? OR receiver_name LIKE ? OR receiver_phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}