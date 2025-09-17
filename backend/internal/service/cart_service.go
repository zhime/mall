package service

import (
	"errors"

	"mall/internal/model"
	"mall/internal/repository"
)

// CartService 购物车服务接口
type CartService interface {
	AddToCart(userID uint64, req *AddToCartRequest) error
	UpdateCartItem(userID uint64, req *UpdateCartItemRequest) error
	RemoveFromCart(userID uint64, req *RemoveFromCartRequest) error
	GetUserCart(userID uint64) (*CartResponse, error)
	ClearCart(userID uint64) error
	GetCartCount(userID uint64) (int, error)
}

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	SKUID     uint64 `json:"sku_id"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// UpdateCartItemRequest 更新购物车商品请求
type UpdateCartItemRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	SKUID     uint64 `json:"sku_id"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// RemoveFromCartRequest 从购物车移除请求
type RemoveFromCartRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	SKUID     uint64 `json:"sku_id"`
}

// CartResponse 购物车响应
type CartResponse struct {
	Items       []*CartItemResponse `json:"items"`
	TotalAmount float64             `json:"total_amount"`
	TotalCount  int                 `json:"total_count"`
}

// CartItemResponse 购物车商品响应
type CartItemResponse struct {
	ID          uint64  `json:"id"`
	ProductID   uint64  `json:"product_id"`
	ProductName string  `json:"product_name"`
	ProductImage string `json:"product_image"`
	SKUID       uint64  `json:"sku_id"`
	SKUName     string  `json:"sku_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
	Stock       int     `json:"stock"`
	Status      int8    `json:"status"` // 商品状态
}

// cartService 购物车服务实现
type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	skuRepo     repository.ProductSKURepository
}

// NewCartService 创建购物车服务
func NewCartService(
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	skuRepo repository.ProductSKURepository,
) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		skuRepo:     skuRepo,
	}
}

// AddToCart 添加商品到购物车
func (s *cartService) AddToCart(userID uint64, req *AddToCartRequest) error {
	// 验证商品是否存在
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// 检查商品状态
	if product.Status != 1 {
		return errors.New("product is not available")
	}

	// 如果指定了SKU，验证SKU
	var sku *model.ProductSKU
	if req.SKUID > 0 {
		sku, err = s.skuRepo.GetByID(req.SKUID)
		if err != nil {
			return errors.New("SKU not found")
		}
		if sku.ProductID != req.ProductID {
			return errors.New("SKU does not belong to this product")
		}
		if sku.Status != 1 {
			return errors.New("SKU is not available")
		}
		// 检查SKU库存
		if sku.Stock < req.Quantity {
			return errors.New("insufficient SKU stock")
		}
	} else {
		// 检查商品库存
		if product.Stock < req.Quantity {
			return errors.New("insufficient product stock")
		}
	}

	// 添加到购物车
	cartItem := &model.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		SKUID:     req.SKUID,
		Quantity:  req.Quantity,
	}

	return s.cartRepo.AddItem(cartItem)
}

// UpdateCartItem 更新购物车商品数量
func (s *cartService) UpdateCartItem(userID uint64, req *UpdateCartItemRequest) error {
	// 获取购物车商品
	cartItem, err := s.cartRepo.GetCartItem(userID, req.ProductID, req.SKUID)
	if err != nil {
		return errors.New("cart item not found")
	}

	// 验证库存
	if req.SKUID > 0 {
		sku, err := s.skuRepo.GetByID(req.SKUID)
		if err != nil {
			return errors.New("SKU not found")
		}
		if sku.Stock < req.Quantity {
			return errors.New("insufficient SKU stock")
		}
	} else {
		product, err := s.productRepo.GetByID(req.ProductID)
		if err != nil {
			return errors.New("product not found")
		}
		if product.Stock < req.Quantity {
			return errors.New("insufficient product stock")
		}
	}

	// 更新数量
	cartItem.Quantity = req.Quantity
	return s.cartRepo.UpdateItem(cartItem)
}

// RemoveFromCart 从购物车移除商品
func (s *cartService) RemoveFromCart(userID uint64, req *RemoveFromCartRequest) error {
	return s.cartRepo.RemoveItem(userID, req.ProductID, req.SKUID)
}

// GetUserCart 获取用户购物车
func (s *cartService) GetUserCart(userID uint64) (*CartResponse, error) {
	items, err := s.cartRepo.GetUserCartItems(userID)
	if err != nil {
		return nil, err
	}

	var cartItems []*CartItemResponse
	var totalAmount float64
	var totalCount int

	for _, item := range items {
		var price float64
		var stock int
		var skuName string

		// 获取价格和库存
		if item.SKUID > 0 && item.SKU.ID > 0 {
			price = item.SKU.Price
			stock = item.SKU.Stock
			skuName = item.SKU.Name
		} else {
			price = item.Product.Price
			stock = item.Product.Stock
		}

		// 获取主图
		var mainImage string
		for _, image := range item.Product.Images {
			if image.IsMain == 1 {
				mainImage = image.ImageURL
				break
			}
		}

		cartItem := &CartItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			ProductImage: mainImage,
			SKUID:       item.SKUID,
			SKUName:     skuName,
			Price:       price,
			Quantity:    item.Quantity,
			TotalPrice:  price * float64(item.Quantity),
			Stock:       stock,
			Status:      item.Product.Status,
		}

		cartItems = append(cartItems, cartItem)
		totalAmount += cartItem.TotalPrice
		totalCount += item.Quantity
	}

	return &CartResponse{
		Items:       cartItems,
		TotalAmount: totalAmount,
		TotalCount:  totalCount,
	}, nil
}

// ClearCart 清空购物车
func (s *cartService) ClearCart(userID uint64) error {
	return s.cartRepo.ClearUserCart(userID)
}

// GetCartCount 获取购物车商品数量
func (s *cartService) GetCartCount(userID uint64) (int, error) {
	items, err := s.cartRepo.GetUserCartItems(userID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range items {
		count += item.Quantity
	}

	return count, nil
}