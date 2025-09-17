package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// ProductRepository 商品仓储接口
type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id uint64) (*model.Product, error)
	GetWithDetails(id uint64) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id uint64) error
	List(page, pageSize int, categoryID uint64, status int8) ([]*model.Product, int64, error)
	Search(keyword string, page, pageSize int) ([]*model.Product, int64, error)
	GetByCategoryID(categoryID uint64, page, pageSize int) ([]*model.Product, int64, error)
	GetHotProducts(limit int) ([]*model.Product, error)
	UpdateStock(id uint64, stock int) error
	UpdateSales(id uint64, sales int) error
	// 为SearchService添加的方法
	GetProductsByIDs(ids []int) ([]model.Product, error)
	GetAllProducts() ([]model.Product, error)
}

// productRepository 商品仓储实现
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建商品仓储
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create 创建商品
func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

// GetByID 根据ID获取商品
func (r *productRepository) GetByID(id uint64) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetWithDetails 获取商品详情（包含分类、SKU、图片）
func (r *productRepository) GetWithDetails(id uint64) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").
		Preload("SKUs").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Update 更新商品
func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// Delete 删除商品（软删除）
func (r *productRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Product{}, id).Error
}

// List 获取商品列表
func (r *productRepository) List(page, pageSize int, categoryID uint64, status int8) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Preload("Category").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_main = ?", 1)
	})

	// 添加筛选条件
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("sort_order DESC, id DESC").
		Find(&products).Error

	return products, total, err
}

// Search 搜索商品
func (r *productRepository) Search(keyword string, page, pageSize int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	query := r.db.Model(&model.Product{}).
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_main = ?", 1)
		}).
		Where("name LIKE ? OR subtitle LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Where("status = ?", 1)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("sales DESC, id DESC").
		Find(&products).Error

	return products, total, err
}

// GetByCategoryID 根据分类ID获取商品
func (r *productRepository) GetByCategoryID(categoryID uint64, page, pageSize int) ([]*model.Product, int64, error) {
	return r.List(page, pageSize, categoryID, 1)
}

// GetHotProducts 获取热门商品
func (r *productRepository) GetHotProducts(limit int) ([]*model.Product, error) {
	var products []*model.Product
	err := r.db.Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_main = ?", 1)
		}).
		Where("status = ?", 1).
		Order("sales DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// UpdateStock 更新商品库存
func (r *productRepository) UpdateStock(id uint64, stock int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("stock", stock).Error
}

// UpdateSales 更新商品销量
func (r *productRepository) UpdateSales(id uint64, sales int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("sales", sales).Error
}

// GetProductsByIDs 根据ID列表获取商品
func (r *productRepository) GetProductsByIDs(ids []int) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_main = ?", 1)
		}).
		Where("id IN ?", ids).
		Find(&products).Error
	return products, err
}

// GetAllProducts 获取所有商品
func (r *productRepository) GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_main = ?", 1)
		}).
		Where("status = ?", 1).
		Find(&products).Error
	return products, err
}