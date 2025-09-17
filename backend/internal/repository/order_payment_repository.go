package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// OrderPaymentRepository 订单支付仓储接口
type OrderPaymentRepository interface {
	Create(payment *model.OrderPayment) error
	GetByID(id uint64) (*model.OrderPayment, error)
	GetByPaymentNo(paymentNo string) (*model.OrderPayment, error)
	GetByOrderID(orderID uint64) ([]*model.OrderPayment, error)
	Update(payment *model.OrderPayment) error
	UpdateStatus(id uint64, status int8, tradeNo string) error
}

// orderPaymentRepository 订单支付仓储实现
type orderPaymentRepository struct {
	db *gorm.DB
}

// NewOrderPaymentRepository 创建订单支付仓储
func NewOrderPaymentRepository(db *gorm.DB) OrderPaymentRepository {
	return &orderPaymentRepository{db: db}
}

// Create 创建支付记录
func (r *orderPaymentRepository) Create(payment *model.OrderPayment) error {
	return r.db.Create(payment).Error
}

// GetByID 根据ID获取支付记录
func (r *orderPaymentRepository) GetByID(id uint64) (*model.OrderPayment, error) {
	var payment model.OrderPayment
	err := r.db.Preload("Order").Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByPaymentNo 根据支付单号获取支付记录
func (r *orderPaymentRepository) GetByPaymentNo(paymentNo string) (*model.OrderPayment, error) {
	var payment model.OrderPayment
	err := r.db.Preload("Order").Where("payment_no = ?", paymentNo).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByOrderID 根据订单ID获取支付记录
func (r *orderPaymentRepository) GetByOrderID(orderID uint64) ([]*model.OrderPayment, error) {
	var payments []*model.OrderPayment
	err := r.db.Where("order_id = ?", orderID).Find(&payments).Error
	return payments, err
}

// Update 更新支付记录
func (r *orderPaymentRepository) Update(payment *model.OrderPayment) error {
	return r.db.Save(payment).Error
}

// UpdateStatus 更新支付状态
func (r *orderPaymentRepository) UpdateStatus(id uint64, status int8, tradeNo string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if tradeNo != "" {
		updates["trade_no"] = tradeNo
	}
	if status == 1 { // 支付成功
		updates["pay_time"] = "NOW()"
	}
	
	return r.db.Model(&model.OrderPayment{}).Where("id = ?", id).Updates(updates).Error
}