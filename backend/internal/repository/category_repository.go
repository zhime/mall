package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	Create(category *model.Category) error
	GetByID(id uint64) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id uint64) error
	List(parentID uint64) ([]*model.Category, error)
	GetTree() ([]*model.Category, error)
	GetByLevel(level int8) ([]*model.Category, error)
	Search(keyword string) ([]*model.Category, error)
}

// categoryRepository 分类仓储实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create 创建分类
func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// GetByID 根据ID获取分类
func (r *categoryRepository) GetByID(id uint64) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update 更新分类
func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类（软删除）
func (r *categoryRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// List 获取分类列表
func (r *categoryRepository) List(parentID uint64) ([]*model.Category, error) {
	var categories []*model.Category
	query := r.db.Where("status = ?", 1).Order("sort_order ASC, id ASC")

	if parentID > 0 {
		query = query.Where("parent_id = ?", parentID)
	}

	err := query.Find(&categories).Error
	return categories, err
}

// GetTree 获取分类树
func (r *categoryRepository) GetTree() ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.Where("status = ?", 1).
		Order("level ASC, sort_order ASC, id ASC").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	// 构建树形结构
	categoryMap := make(map[uint64]*model.Category)
	rootCategories := []*model.Category{}

	// 创建map映射
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	// 构建父子关系
	for _, category := range categories {
		if category.ParentID == 0 {
			rootCategories = append(rootCategories, category)
		} else {
			if parent, exists := categoryMap[category.ParentID]; exists {
				parent.Children = append(parent.Children, *category)
			}
		}
	}

	return rootCategories, nil
}

// GetByLevel 根据级别获取分类
func (r *categoryRepository) GetByLevel(level int8) ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.Where("level = ? AND status = ?", level, 1).
		Order("sort_order ASC, id ASC").
		Find(&categories).Error
	return categories, err
}

// Search 搜索分类
func (r *categoryRepository) Search(keyword string) ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.Where("name LIKE ? AND status = ?", "%"+keyword+"%", 1).
		Order("sort_order ASC, id ASC").
		Find(&categories).Error
	return categories, err
}
