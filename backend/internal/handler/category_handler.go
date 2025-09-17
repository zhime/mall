package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mall/internal/service"
	"mall/pkg/utils"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.categoryService.CreateCategory(&req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Category created successfully", nil)
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid category ID")
		return
	}

	var req service.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.categoryService.UpdateCategory(id, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Category updated successfully", nil)
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid category ID")
		return
	}

	if err := h.categoryService.DeleteCategory(id); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Category deleted successfully", nil)
}

// GetCategory 获取分类详情
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid category ID")
		return
	}

	category, err := h.categoryService.GetCategory(id)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, category)
}

// GetCategoryTree 获取分类树
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	categories, err := h.categoryService.GetCategoryTree()
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, categories)
}

// GetCategoriesByLevel 根据级别获取分类
func (h *CategoryHandler) GetCategoriesByLevel(c *gin.Context) {
	levelStr := c.Query("level")
	if levelStr == "" {
		utils.InvalidParams(c, "Level parameter is required")
		return
	}

	level, err := strconv.ParseInt(levelStr, 10, 8)
	if err != nil {
		utils.InvalidParams(c, "Invalid level parameter")
		return
	}

	categories, err := h.categoryService.GetCategoriesByLevel(int8(level))
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, categories)
}

// SearchCategories 搜索分类
func (h *CategoryHandler) SearchCategories(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.InvalidParams(c, "Keyword parameter is required")
		return
	}

	categories, err := h.categoryService.SearchCategories(keyword)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, categories)
}