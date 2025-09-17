package service

import (
	"errors"

	"mall/internal/model"
	"mall/internal/repository"
)

// CategoryService 分类服务接口
type CategoryService interface {
	CreateCategory(req *CreateCategoryRequest) error
	UpdateCategory(id uint64, req *UpdateCategoryRequest) error
	DeleteCategory(id uint64) error
	GetCategory(id uint64) (*CategoryResponse, error)
	GetCategoryTree() ([]*CategoryResponse, error)
	GetCategoriesByLevel(level int8) ([]*CategoryResponse, error)
	SearchCategories(keyword string) ([]*CategoryResponse, error)
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	ParentID    uint64 `json:"parent_id"`
	Name        string `json:"name" binding:"required"`
	Level       int8   `json:"level" binding:"required"`
	SortOrder   int    `json:"sort_order"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	SortOrder   int    `json:"sort_order"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID          uint64              `json:"id"`
	ParentID    uint64              `json:"parent_id"`
	Name        string              `json:"name"`
	Level       int8                `json:"level"`
	SortOrder   int                 `json:"sort_order"`
	Icon        string              `json:"icon"`
	Description string              `json:"description"`
	Status      int8                `json:"status"`
	Children    []*CategoryResponse `json:"children,omitempty"`
}

// categoryService 分类服务实现
type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(req *CreateCategoryRequest) error {
	// 验证父分类是否存在
	if req.ParentID > 0 {
		_, err := s.categoryRepo.GetByID(req.ParentID)
		if err != nil {
			return errors.New("parent category not found")
		}
	}

	category := &model.Category{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Level:       req.Level,
		SortOrder:   req.SortOrder,
		Icon:        req.Icon,
		Description: req.Description,
		Status:      1,
	}

	return s.categoryRepo.Create(category)
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(id uint64, req *UpdateCategoryRequest) error {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.SortOrder >= 0 {
		category.SortOrder = req.SortOrder
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Status >= 0 {
		category.Status = req.Status
	}

	return s.categoryRepo.Update(category)
}

// DeleteCategory 删除分类
func (s *categoryService) DeleteCategory(id uint64) error {
	// 检查是否有子分类
	children, err := s.categoryRepo.List(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("cannot delete category with children")
	}

	// TODO: 检查是否有商品使用此分类

	return s.categoryRepo.Delete(id)
}

// GetCategory 获取分类详情
func (s *categoryService) GetCategory(id uint64) (*CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	return s.toCategoryResponse(category), nil
}

// GetCategoryTree 获取分类树
func (s *categoryService) GetCategoryTree() ([]*CategoryResponse, error) {
	categories, err := s.categoryRepo.GetTree()
	if err != nil {
		return nil, err
	}

	var result []*CategoryResponse
	for _, category := range categories {
		result = append(result, s.toCategoryResponseWithChildren(category))
	}

	return result, nil
}

// GetCategoriesByLevel 根据级别获取分类
func (s *categoryService) GetCategoriesByLevel(level int8) ([]*CategoryResponse, error) {
	categories, err := s.categoryRepo.GetByLevel(level)
	if err != nil {
		return nil, err
	}

	var result []*CategoryResponse
	for _, category := range categories {
		result = append(result, s.toCategoryResponse(category))
	}

	return result, nil
}

// SearchCategories 搜索分类
func (s *categoryService) SearchCategories(keyword string) ([]*CategoryResponse, error) {
	categories, err := s.categoryRepo.Search(keyword)
	if err != nil {
		return nil, err
	}

	var result []*CategoryResponse
	for _, category := range categories {
		result = append(result, s.toCategoryResponse(category))
	}

	return result, nil
}

// toCategoryResponse 转换为分类响应
func (s *categoryService) toCategoryResponse(category *model.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:          category.ID,
		ParentID:    category.ParentID,
		Name:        category.Name,
		Level:       category.Level,
		SortOrder:   category.SortOrder,
		Icon:        category.Icon,
		Description: category.Description,
		Status:      category.Status,
	}
}

// toCategoryResponseWithChildren 转换为分类响应（包含子分类）
func (s *categoryService) toCategoryResponseWithChildren(category *model.Category) *CategoryResponse {
	response := s.toCategoryResponse(category)
	
	if len(category.Children) > 0 {
		for i := range category.Children {
			response.Children = append(response.Children, s.toCategoryResponseWithChildren(&category.Children[i]))
		}
	}

	return response
}