package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mall/internal/service"
	"mall/pkg/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SendSMSCode 发送短信验证码
func (h *AuthHandler) SendSMSCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.authService.SendSMSCode(req.Phone); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Verification code sent successfully", nil)
}

// LoginByPhone 手机号登录
func (h *AuthHandler) LoginByPhone(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	response, err := h.authService.LoginByPhone(req.Phone, req.Code)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// LoginByWechat 微信登录
func (h *AuthHandler) LoginByWechat(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	response, err := h.authService.LoginByWechat(req.Code)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// Register 注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	if err := h.authService.RegisterByPassword(&req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Registration successful", nil)
}

// LoginByPassword 密码登录
func (h *AuthHandler) LoginByPassword(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	response, err := h.authService.LoginByPassword(req.Username, req.Password)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// RefreshToken 刷新token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		utils.Unauthorized(c, "Missing token")
		return
	}

	// 移除Bearer前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	response, err := h.authService.RefreshToken(token)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	if err := h.authService.Logout(uint64(userID)); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Logout successful", nil)
}

// GetUserInfo 获取用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	response, err := h.authService.GetUserInfo(uint64(userID))
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetProfile 获取用户详细信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	response, err := h.authService.GetUserInfo(uint64(userID))
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, response)
}

// UpdateProfile 更新用户信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Gender   int8   `json:"gender"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	// TODO: 实现更新用户信息逻辑
	utils.SuccessWithMessage(c, "Profile updated successfully", nil)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	// TODO: 实现修改密码逻辑
	utils.SuccessWithMessage(c, "Password changed successfully", nil)
}

// BindPhone 绑定手机号
func (h *AuthHandler) BindPhone(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		utils.Unauthorized(c, "Invalid user")
		return
	}

	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.InvalidParams(c, err.Error())
		return
	}

	// TODO: 实现绑定手机号逻辑
	utils.SuccessWithMessage(c, "Phone number bound successfully", nil)
}