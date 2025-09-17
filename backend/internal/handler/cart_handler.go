package handler

import (
	"github.com/gin-gonic/gin"

	"mall/internal/service"
	"mall/pkg/utils"
)

// CartHandler 购物车处理器
type CartHandler struct {
	cartService service.CartService
}

// NewCartHandler 创建购物车处理器
func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// AddToCart 添加商品到购物车
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req service.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.cartService.AddToCart(uint64(userID), &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Product added to cart successfully", nil)
}

// UpdateCartItem 更新购物车商品数量
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req service.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.cartService.UpdateCartItem(uint64(userID), &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Cart item updated successfully", nil)
}

// RemoveFromCart 从购物车移除商品
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req service.RemoveFromCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.cartService.RemoveFromCart(uint64(userID), &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Product removed from cart successfully", nil)
}

// GetUserCart 获取用户购物车
func (h *CartHandler) GetUserCart(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	cart, err := h.cartService.GetUserCart(uint64(userID))
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, cart)
}

// ClearCart 清空购物车
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	if err := h.cartService.ClearCart(uint64(userID)); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Cart cleared successfully", nil)
}

// GetCartCount 获取购物车商品数量
func (h *CartHandler) GetCartCount(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	count, err := h.cartService.GetCartCount(uint64(userID))
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, gin.H{"count": count})
}