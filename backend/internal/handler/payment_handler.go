package handler

import (
	"github.com/gin-gonic/gin"

	"mall/internal/service"
	"mall/pkg/utils"
)

// PaymentHandler 支付处理器
type PaymentHandler struct {
	paymentService service.PaymentService
}

// NewPaymentHandler 创建支付处理器
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePayment 创建支付
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req service.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	// 验证支付方式
	if req.PaymentMethod != "wechat" && req.PaymentMethod != "alipay" {
		utils.InvalidParams(c, "Invalid payment method")
		return
	}

	response, err := h.paymentService.CreatePayment(uint64(userID), &req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// WechatCallback 微信支付回调
func (h *PaymentHandler) WechatCallback(c *gin.Context) {
	var req service.WechatCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(400, "FAIL")
		return
	}

	if err := h.paymentService.ProcessWechatCallback(&req); err != nil {
		c.String(400, "FAIL")
		return
	}

	c.String(200, "SUCCESS")
}

// AlipayCallback 支付宝回调
func (h *PaymentHandler) AlipayCallback(c *gin.Context) {
	var req service.AlipayCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(400, "failure")
		return
	}

	if err := h.paymentService.ProcessAlipayCallback(&req); err != nil {
		c.String(400, "failure")
		return
	}

	c.String(200, "success")
}

// GetPaymentStatus 获取支付状态
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentNo := c.Param("paymentNo")
	if paymentNo == "" {
		utils.InvalidParams(c, "Payment number is required")
		return
	}

	response, err := h.paymentService.GetPaymentStatus(paymentNo)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// CancelPayment 取消支付
func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	paymentNo := c.Param("paymentNo")
	if paymentNo == "" {
		utils.InvalidParams(c, "Payment number is required")
		return
	}

	if err := h.paymentService.CancelPayment(paymentNo); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Payment cancelled successfully", nil)
}

// RefundPayment 退款
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	paymentNo := c.Param("paymentNo")
	if paymentNo == "" {
		utils.InvalidParams(c, "Payment number is required")
		return
	}

	var req struct {
		RefundAmount float64 `json:"refund_amount" binding:"required,gt=0"`
		RefundReason string  `json:"refund_reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	// TODO: 实现退款逻辑
	utils.SuccessWithMessage(c, "Refund request submitted successfully", nil)
}