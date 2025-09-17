package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mall/internal/service"
	"mall/pkg/utils"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	response, err := h.orderService.CreateOrder(uint64(userID), &req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetOrderDetail 获取订单详情
func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid order ID")
		return
	}

	response, err := h.orderService.GetOrderDetail(uint64(userID), orderID)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetUserOrders 获取用户订单列表
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req service.OrderListRequest
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

	response, err := h.orderService.GetUserOrders(uint64(userID), &req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid order ID")
		return
	}

	if err := h.orderService.CancelOrder(uint64(userID), orderID); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Order cancelled successfully", nil)
}

// ConfirmOrder 确认收货
func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid order ID")
		return
	}

	if err := h.orderService.ConfirmOrder(uint64(userID), orderID); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Order confirmed successfully", nil)
}

// GetOrders 获取订单列表（管理员）
func (h *OrderHandler) GetOrders(c *gin.Context) {
	var req service.AdminOrderListRequest
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

	response, err := h.orderService.GetOrders(&req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// UpdateOrderStatus 更新订单状态
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		utils.InvalidParams(c, "Invalid order ID")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.orderService.UpdateOrderStatus(orderID, req.Status); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Order status updated successfully", nil)
}

// SearchOrders 搜索订单
func (h *OrderHandler) SearchOrders(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.InvalidParams(c, "Keyword is required")
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

	response, err := h.orderService.SearchOrders(keyword, page, pageSize)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}