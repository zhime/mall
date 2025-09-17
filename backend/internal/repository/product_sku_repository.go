package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// ProductSKURepository 商品SKU仓储接口
type ProductSKURepository interface {
	Create(sku *model.ProductSKU) error
	GetByID(id uint64) (*model.ProductSKU, error)
	GetBySKUCode(skuCode string) (*model.ProductSKU, error)
	GetByProductID(productID uint64) ([]*model.ProductSKU, error)
	Update(sku *model.ProductSKU) error
	Delete(id uint64) error
	UpdateStock(id uint64, stock int) error
	BatchUpdateStock(skuIDs []uint64, stocks []int) error
}

// productSKURepository 商品SKU仓储实现
type productSKURepository struct {
	db *gorm.DB
}

// NewProductSKURepository 创建商品SKU仓储
func NewProductSKURepository(db *gorm.DB) ProductSKURepository {
	return &productSKURepository{db: db}
}

// Create 创建SKU
func (r *productSKURepository) Create(sku *model.ProductSKU) error {
	return r.db.Create(sku).Error
}

// GetByID 根据ID获取SKU
func (r *productSKURepository) GetByID(id uint64) (*model.ProductSKU, error) {
	var sku model.ProductSKU
	err := r.db.Preload("Product").Where("id = ?", id).First(&sku).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

// GetBySKUCode 根据SKU编码获取SKU
func (r *productSKURepository) GetBySKUCode(skuCode string) (*model.ProductSKU, error) {
	var sku model.ProductSKU
	err := r.db.Preload("Product").Where("sku_code = ?", skuCode).First(&sku).Error
	if err != nil {
		return nil, err
	}
	return &sku, nil
}

// GetByProductID 根据商品ID获取SKU列表
func (r *productSKURepository) GetByProductID(productID uint64) ([]*model.ProductSKU, error) {
	var skus []*model.ProductSKU
	err := r.db.Where("product_id = ? AND status = ?", productID, 1).
		Order("id ASC").
		Find(&skus).Error
	return skus, err
}

// Update 更新SKU
func (r *productSKURepository) Update(sku *model.ProductSKU) error {
	return r.db.Save(sku).Error
}

// Delete 删除SKU（软删除）
func (r *productSKURepository) Delete(id uint64) error {
	return r.db.Delete(&model.ProductSKU{}, id).Error
}

// UpdateStock 更新SKU库存
func (r *productSKURepository) UpdateStock(id uint64, stock int) error {
	return r.db.Model(&model.ProductSKU{}).Where("id = ?", id).Update("stock", stock).Error
}

// BatchUpdateStock 批量更新SKU库存
func (r *productSKURepository) BatchUpdateStock(skuIDs []uint64, stocks []int) error {
	if len(skuIDs) != len(stocks) {
		return gorm.ErrInvalidData
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, skuID := range skuIDs {
		if err := tx.Model(&model.ProductSKU{}).Where("id = ?", skuID).Update("stock", stocks[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}