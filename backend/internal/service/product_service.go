package service

import (
	"errors"
	"fmt"

	"mall/internal/model"
	"mall/internal/repository"
)

// ProductService 商品服务接口
type ProductService interface {
	CreateProduct(req *CreateProductRequest) error
	UpdateProduct(id uint64, req *UpdateProductRequest) error
	DeleteProduct(id uint64) error
	GetProduct(id uint64) (*ProductResponse, error)
	GetProductDetail(id uint64) (*ProductDetailResponse, error)
	GetProductList(req *ProductListRequest) (*ProductListResponse, error)
	SearchProducts(req *SearchProductRequest) (*ProductListResponse, error)
	GetHotProducts(limit int) ([]*ProductResponse, error)
	UpdateProductStock(id uint64, stock int) error
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	CategoryID    uint64                `json:"category_id" binding:"required"`
	Name          string                `json:"name" binding:"required"`
	Subtitle      string                `json:"subtitle"`
	Description   string                `json:"description"`
	Price         float64               `json:"price" binding:"required"`
	OriginalPrice float64               `json:"original_price"`
	Stock         int                   `json:"stock"`
	SortOrder     int                   `json:"sort_order"`
	Images        []ProductImageRequest `json:"images"`
	SKUs          []ProductSKURequest   `json:"skus"`
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	CategoryID    uint64  `json:"category_id"`
	Name          string  `json:"name"`
	Subtitle      string  `json:"subtitle"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	Stock         int     `json:"stock"`
	Status        int8    `json:"status"`
	SortOrder     int     `json:"sort_order"`
}

// ProductImageRequest 商品图片请求
type ProductImageRequest struct {
	ImageURL  string `json:"image_url" binding:"required"`
	SortOrder int    `json:"sort_order"`
	IsMain    int8   `json:"is_main"`
}

// ProductSKURequest 商品SKU请求
type ProductSKURequest struct {
	SKUCode    string  `json:"sku_code" binding:"required"`
	Name       string  `json:"name"`
	Price      float64 `json:"price" binding:"required"`
	Stock      int     `json:"stock"`
	AttrValues string  `json:"attr_values"`
	Image      string  `json:"image"`
}

// ProductListRequest 商品列表请求
type ProductListRequest struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	CategoryID uint64 `json:"category_id" form:"category_id"`
	Status     int8   `json:"status" form:"status"`
}

// SearchProductRequest 搜索商品请求
type SearchProductRequest struct {
	Keyword  string `json:"keyword" form:"keyword" binding:"required"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}

// ProductResponse 商品响应
type ProductResponse struct {
	ID            uint64  `json:"id"`
	CategoryID    uint64  `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	Name          string  `json:"name"`
	Subtitle      string  `json:"subtitle"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	Stock         int     `json:"stock"`
	Sales         int     `json:"sales"`
	Status        int8    `json:"status"`
	MainImage     string  `json:"main_image"`
}

// ProductDetailResponse 商品详情响应
type ProductDetailResponse struct {
	ID            uint64                  `json:"id"`
	CategoryID    uint64                  `json:"category_id"`
	CategoryName  string                  `json:"category_name"`
	Name          string                  `json:"name"`
	Subtitle      string                  `json:"subtitle"`
	Description   string                  `json:"description"`
	Price         float64                 `json:"price"`
	OriginalPrice float64                 `json:"original_price"`
	Stock         int                     `json:"stock"`
	Sales         int                     `json:"sales"`
	Status        int8                    `json:"status"`
	Images        []ProductImageResponse  `json:"images"`
	SKUs          []ProductSKUResponse    `json:"skus"`
}

// ProductImageResponse 商品图片响应
type ProductImageResponse struct {
	ID        uint64 `json:"id"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
	IsMain    int8   `json:"is_main"`
}

// ProductSKUResponse 商品SKU响应
type ProductSKUResponse struct {
	ID         uint64  `json:"id"`
	SKUCode    string  `json:"sku_code"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	AttrValues string  `json:"attr_values"`
	Image      string  `json:"image"`
	Status     int8    `json:"status"`
}

// ProductListResponse 商品列表响应
type ProductListResponse struct {
	Items      []*ProductResponse `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// productService 商品服务实现
type productService struct {
	productRepo    repository.ProductRepository
	categoryRepo   repository.CategoryRepository
	productSKURepo repository.ProductSKURepository
}

// NewProductService 创建商品服务
func NewProductService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	productSKURepo repository.ProductSKURepository,
) ProductService {
	return &productService{
		productRepo:    productRepo,
		categoryRepo:   categoryRepo,
		productSKURepo: productSKURepo,
	}
}

// CreateProduct 创建商品
func (s *productService) CreateProduct(req *CreateProductRequest) error {
	// 验证分类是否存在
	_, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// 创建商品
	product := &model.Product{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Subtitle:      req.Subtitle,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		Sales:         0,
		Status:        1,
		SortOrder:     req.SortOrder,
	}

	if err := s.productRepo.Create(product); err != nil {
		return err
	}

	// 创建商品图片
	for _, imgReq := range req.Images {
		image := &model.ProductImage{
			ProductID: product.ID,
			ImageURL:  imgReq.ImageURL,
			SortOrder: imgReq.SortOrder,
			IsMain:    imgReq.IsMain,
		}
		// TODO: 保存商品图片
		_ = image
	}

	// 创建商品SKU
	for _, skuReq := range req.SKUs {
		sku := &model.ProductSKU{
			ProductID:  product.ID,
			SKUCode:    skuReq.SKUCode,
			Name:       skuReq.Name,
			Price:      skuReq.Price,
			Stock:      skuReq.Stock,
			AttrValues: skuReq.AttrValues,
			Image:      skuReq.Image,
			Status:     1,
		}
		if err := s.productSKURepo.Create(sku); err != nil {
			return fmt.Errorf("failed to create SKU: %v", err)
		}
	}

	return nil
}

// UpdateProduct 更新商品
func (s *productService) UpdateProduct(id uint64, req *UpdateProductRequest) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	// 更新字段
	if req.CategoryID > 0 {
		// 验证分类是否存在
		_, err := s.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			return errors.New("category not found")
		}
		product.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Subtitle != "" {
		product.Subtitle = req.Subtitle
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.OriginalPrice > 0 {
		product.OriginalPrice = req.OriginalPrice
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.Status >= 0 {
		product.Status = req.Status
	}
	if req.SortOrder >= 0 {
		product.SortOrder = req.SortOrder
	}

	return s.productRepo.Update(product)
}

// DeleteProduct 删除商品
func (s *productService) DeleteProduct(id uint64) error {
	// TODO: 检查是否有订单使用此商品
	return s.productRepo.Delete(id)
}

// GetProduct 获取商品基本信息
func (s *productService) GetProduct(id uint64) (*ProductResponse, error) {
	product, err := s.productRepo.GetWithDetails(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return s.toProductResponse(product), nil
}

// GetProductDetail 获取商品详情
func (s *productService) GetProductDetail(id uint64) (*ProductDetailResponse, error) {
	product, err := s.productRepo.GetWithDetails(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return s.toProductDetailResponse(product), nil
}

// GetProductList 获取商品列表
func (s *productService) GetProductList(req *ProductListRequest) (*ProductListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	products, total, err := s.productRepo.List(req.Page, req.PageSize, req.CategoryID, req.Status)
	if err != nil {
		return nil, err
	}

	var items []*ProductResponse
	for _, product := range products {
		items = append(items, s.toProductResponse(product))
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &ProductListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchProducts 搜索商品
func (s *productService) SearchProducts(req *SearchProductRequest) (*ProductListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	products, total, err := s.productRepo.Search(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var items []*ProductResponse
	for _, product := range products {
		items = append(items, s.toProductResponse(product))
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &ProductListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetHotProducts 获取热门商品
func (s *productService) GetHotProducts(limit int) ([]*ProductResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	products, err := s.productRepo.GetHotProducts(limit)
	if err != nil {
		return nil, err
	}

	var result []*ProductResponse
	for _, product := range products {
		result = append(result, s.toProductResponse(product))
	}

	return result, nil
}

// UpdateProductStock 更新商品库存
func (s *productService) UpdateProductStock(id uint64, stock int) error {
	return s.productRepo.UpdateStock(id, stock)
}

// toProductResponse 转换为商品响应
func (s *productService) toProductResponse(product *model.Product) *ProductResponse {
	response := &ProductResponse{
		ID:            product.ID,
		CategoryID:    product.CategoryID,
		Name:          product.Name,
		Subtitle:      product.Subtitle,
		Price:         product.Price,
		OriginalPrice: product.OriginalPrice,
		Stock:         product.Stock,
		Sales:         product.Sales,
		Status:        product.Status,
	}

	// 设置分类名称
	if product.Category.ID > 0 {
		response.CategoryName = product.Category.Name
	}

	// 设置主图
	for _, image := range product.Images {
		if image.IsMain == 1 {
			response.MainImage = image.ImageURL
			break
		}
	}

	return response
}

// toProductDetailResponse 转换为商品详情响应
func (s *productService) toProductDetailResponse(product *model.Product) *ProductDetailResponse {
	response := &ProductDetailResponse{
		ID:            product.ID,
		CategoryID:    product.CategoryID,
		Name:          product.Name,
		Subtitle:      product.Subtitle,
		Description:   product.Description,
		Price:         product.Price,
		OriginalPrice: product.OriginalPrice,
		Stock:         product.Stock,
		Sales:         product.Sales,
		Status:        product.Status,
	}

	// 设置分类名称
	if product.Category.ID > 0 {
		response.CategoryName = product.Category.Name
	}

	// 设置图片
	for _, image := range product.Images {
		response.Images = append(response.Images, ProductImageResponse{
			ID:        image.ID,
			ImageURL:  image.ImageURL,
			SortOrder: image.SortOrder,
			IsMain:    image.IsMain,
		})
	}

	// 设置SKU
	for _, sku := range product.SKUs {
		response.SKUs = append(response.SKUs, ProductSKUResponse{
			ID:         sku.ID,
			SKUCode:    sku.SKUCode,
			Name:       sku.Name,
			Price:      sku.Price,
			Stock:      sku.Stock,
			AttrValues: sku.AttrValues,
			Image:      sku.Image,
			Status:     sku.Status,
		})
	}

	return response
}