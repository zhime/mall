package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mall/internal/service"
	"mall/pkg/utils"
)

// ProductHandler 商品处理器
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler 创建商品处理器
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct 创建商品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.productService.CreateProduct(&req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Product created successfully", nil)
}

// UpdateProduct 更新商品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid product ID")
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.productService.UpdateProduct(id, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Product updated successfully", nil)
}

// DeleteProduct 删除商品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid product ID")
		return
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Product deleted successfully", nil)
}

// GetProduct 获取商品基本信息
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid product ID")
		return
	}

	product, err := h.productService.GetProduct(id)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, product)
}

// GetProductDetail 获取商品详情
func (h *ProductHandler) GetProductDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid product ID")
		return
	}

	product, err := h.productService.GetProductDetail(id)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, product)
}

// GetProductList 获取商品列表
func (h *ProductHandler) GetProductList(c *gin.Context) {
	var req service.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	result, err := h.productService.GetProductList(&req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, result)
}

// SearchProducts 搜索商品
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	var req service.SearchProductRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if req.Keyword == "" {
		utils.InvalidParams(c, "Keyword is required")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	result, err := h.productService.SearchProducts(&req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, result)
}

// GetHotProducts 获取热门商品
func (h *ProductHandler) GetHotProducts(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 10 // 默认值

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	products, err := h.productService.GetHotProducts(limit)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, products)
}

// UpdateProductStock 更新商品库存
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid product ID")
		return
	}

	var req struct {
		Stock int `json:"stock" binding:"required,gte=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.productService.UpdateProductStock(id, req.Stock); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Product stock updated successfully", nil)
}

// GetProductsByCategory 根据分类获取商品
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryId")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid category ID")
		return
	}

	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := c.Query("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	req := &service.ProductListRequest{
		Page:       page,
		PageSize:   pageSize,
		CategoryID: categoryID,
		Status:     1, // 只显示上架的商品
	}

	result, err := h.productService.GetProductList(req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, result)
}