package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"mall/internal/model"
	"mall/internal/repository"
	"mall/pkg/cache"
	"mall/pkg/utils"
)

// AuthService 认证服务接口
type AuthService interface {
	// 手机号登录相关
	SendSMSCode(phone string) error
	VerifySMSCode(phone, code string) bool
	LoginByPhone(phone, code string) (*LoginResponse, error)
	
	// 微信登录相关
	LoginByWechat(code string) (*LoginResponse, error)
	
	// 密码登录相关
	RegisterByPassword(req *RegisterRequest) error
	LoginByPassword(username, password string) (*LoginResponse, error)
	
	// 通用方法
	RefreshToken(token string) (*TokenResponse, error)
	Logout(userID uint64) error
	GetUserInfo(userID uint64) (*UserInfoResponse, error)
	
	// 用户资料相关
	UpdateProfile(userID uint64, req *UpdateProfileRequest) error
	ChangePassword(userID uint64, req *ChangePasswordRequest) error
	BindPhone(userID uint64, req *BindPhoneRequest) error
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string           `json:"token"`
	ExpiresIn int64            `json:"expires_in"`
	User      *UserInfoResponse `json:"user"`
}

// TokenResponse token响应
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
	Birthday string `json:"birthday"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// BindPhoneRequest 绑定手机号请求
type BindPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Status   int8   `json:"status"`
}

// authService 认证服务实现
type authService struct {
	userRepo     repository.UserRepository
	userAuthRepo repository.UserAuthRepository
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository, userAuthRepo repository.UserAuthRepository) AuthService {
	return &authService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}

// SendSMSCode 发送短信验证码
func (s *authService) SendSMSCode(phone string) error {
	// 验证手机号格式
	if !utils.IsValidPhone(phone) {
		return errors.New("invalid phone number format")
	}

	// 生成验证码
	code := utils.GenerateVerifyCode()

	// 存储到Redis，5分钟过期
	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)
	err := cache.Set(ctx, key, code, 5*time.Minute)
	if err != nil {
		return errors.New("failed to store verification code")
	}

	// TODO: 调用第三方短信服务发送验证码
	// 这里暂时记录日志
	fmt.Printf("SMS code for %s: %s\n", phone, code)

	return nil
}

// VerifySMSCode 验证短信验证码
func (s *authService) VerifySMSCode(phone, code string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)
	
	storedCode, err := cache.Get(ctx, key)
	if err != nil {
		return false
	}

	// 验证成功后删除验证码
	if storedCode == code {
		cache.Del(ctx, key)
		return true
	}

	return false
}

// LoginByPhone 手机号登录
func (s *authService) LoginByPhone(phone, code string) (*LoginResponse, error) {
	// 验证验证码
	if !s.VerifySMSCode(phone, code) {
		return nil, errors.New("invalid verification code")
	}

	// 查找用户
	user, err := s.userRepo.GetByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，自动注册
			user = &model.User{
				Phone:  phone,
				Status: 1,
			}
			if err := s.userRepo.Create(user); err != nil {
				return nil, errors.New("failed to create user")
			}

			// 创建认证信息
			auth := &model.UserAuth{
				UserID:   user.ID,
				AuthType: "phone",
				AuthKey:  phone,
			}
			if err := s.userAuthRepo.Create(auth); err != nil {
				return nil, errors.New("failed to create auth info")
			}
		} else {
			return nil, errors.New("failed to get user")
		}
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("user account is disabled")
	}

	// 生成token
	token, err := utils.GenerateToken(int64(user.ID), user.Username, "user", "web")
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// 构造响应
	userInfo := &UserInfoResponse{
		ID:       user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		Status:   user.Status,
	}

	return &LoginResponse{
		Token:     token,
		ExpiresIn: 168 * 3600, // 7天
		User:      userInfo,
	}, nil
}

// LoginByWechat 微信登录
func (s *authService) LoginByWechat(code string) (*LoginResponse, error) {
	// TODO: 调用微信API获取openid
	// 这里暂时模拟
	openID := fmt.Sprintf("wx_openid_%d", rand.Int63())

	// 查找用户
	user, err := s.userRepo.GetByWechatOpenID(openID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，自动注册
			user = &model.User{
				WechatOpenID: openID,
				Status:       1,
			}
			if err := s.userRepo.Create(user); err != nil {
				return nil, errors.New("failed to create user")
			}

			// 创建认证信息
			auth := &model.UserAuth{
				UserID:   user.ID,
				AuthType: "wechat",
				AuthKey:  openID,
			}
			if err := s.userAuthRepo.Create(auth); err != nil {
				return nil, errors.New("failed to create auth info")
			}
		} else {
			return nil, errors.New("failed to get user")
		}
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("user account is disabled")
	}

	// 生成token
	token, err := utils.GenerateToken(int64(user.ID), user.Username, "user", "miniprogram")
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// 构造响应
	userInfo := &UserInfoResponse{
		ID:       user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		Status:   user.Status,
	}

	return &LoginResponse{
		Token:     token,
		ExpiresIn: 168 * 3600, // 7天
		User:      userInfo,
	}, nil
}

// RegisterByPassword 密码注册
func (s *authService) RegisterByPassword(req *RegisterRequest) error {
	// 验证验证码
	if !s.VerifySMSCode(req.Phone, req.Code) {
		return errors.New("invalid verification code")
	}

	// 检查用户名是否存在
	if _, err := s.userRepo.GetByPhone(req.Phone); err == nil {
		return errors.New("phone number already registered")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Phone:    req.Phone,
		Status:   1,
	}
	if err := s.userRepo.Create(user); err != nil {
		return errors.New("failed to create user")
	}

	// 加密密码
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// 创建认证信息
	auth := &model.UserAuth{
		UserID:       user.ID,
		AuthType:     "password",
		AuthKey:      req.Phone,
		PasswordHash: passwordHash,
		Salt:         utils.GenerateSalt(),
	}
	if err := s.userAuthRepo.Create(auth); err != nil {
		return errors.New("failed to create auth info")
	}

	return nil
}

// LoginByPassword 密码登录
func (s *authService) LoginByPassword(username, password string) (*LoginResponse, error) {
	// 获取认证信息
	auth, err := s.userAuthRepo.GetByAuthKey(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 验证密码
	if !utils.CheckPassword(password, auth.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	// 获取用户信息
	user := &auth.User
	if user.Status != 1 {
		return nil, errors.New("user account is disabled")
	}

	// 生成token
	token, err := utils.GenerateToken(int64(user.ID), user.Username, "user", "web")
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// 构造响应
	userInfo := &UserInfoResponse{
		ID:       user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		Status:   user.Status,
	}

	return &LoginResponse{
		Token:     token,
		ExpiresIn: 168 * 3600, // 7天
		User:      userInfo,
	}, nil
}

// RefreshToken 刷新token
func (s *authService) RefreshToken(token string) (*TokenResponse, error) {
	newToken, err := utils.RefreshToken(token)
	if err != nil {
		return nil, errors.New("failed to refresh token")
	}

	return &TokenResponse{
		Token:     newToken,
		ExpiresIn: 168 * 3600, // 7天
	}, nil
}

// Logout 登出
func (s *authService) Logout(userID uint64) error {
	// TODO: 将token加入黑名单
	ctx := context.Background()
	key := fmt.Sprintf("token_blacklist:%d", userID)
	return cache.Set(ctx, key, "1", 24*time.Hour)
}

// GetUserInfo 获取用户信息
func (s *authService) GetUserInfo(userID uint64) (*UserInfoResponse, error) {
	user, err := s.userRepo.GetUserWithProfile(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	userInfo := &UserInfoResponse{
		ID:       user.ID,
		Username: user.Username,
		Phone:    user.Phone,
		Status:   user.Status,
	}

	// 如果有详细信息，填充nickname和avatar
	if user.Profile != nil {
		userInfo.Nickname = user.Profile.Nickname
		userInfo.Avatar = user.Profile.Avatar
	}

	return userInfo, nil
}

// UpdateProfile 更新用户资料
func (s *authService) UpdateProfile(userID uint64, req *UpdateProfileRequest) error {
	user, err := s.userRepo.GetUserWithProfile(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// 更新用户基础信息
	if req.Nickname != "" {
		user.Profile.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Profile.Avatar = req.Avatar
	}
	if req.Gender != 0 {
		user.Profile.Gender = req.Gender
	}
	if req.Birthday != "" {
		// 解析生日字符串为time.Time
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			user.Profile.Birthday = birthday
		}
	}

	return s.userRepo.Update(user)
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(userID uint64, req *ChangePasswordRequest) error {
	// 验证旧密码
	userAuth, err := s.userAuthRepo.GetByUserIDAndType(userID, "password")
	if err != nil {
		return errors.New("password not set")
	}

	if !utils.CheckPassword(req.OldPassword, userAuth.PasswordHash) {
		return errors.New("old password incorrect")
	}

	// 生成新密码哈希
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	userAuth.PasswordHash = newPasswordHash

	return s.userAuthRepo.Update(userAuth)
}

// BindPhone 绑定手机号
func (s *authService) BindPhone(userID uint64, req *BindPhoneRequest) error {
	// 验证验证码
	if !s.VerifySMSCode(req.Phone, req.Code) {
		return errors.New("invalid verification code")
	}

	// 检查手机号是否已被其他用户绑定
	if _, err := s.userRepo.GetByPhone(req.Phone); err == nil {
		return errors.New("phone number already bound to another account")
	}

	// 获取用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// 更新手机号
	user.Phone = req.Phone

	// 更新认证信息
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to update user phone")
	}

	// 更新认证信息中的手机号
	userAuth, err := s.userAuthRepo.GetByUserIDAndType(userID, "phone")
	if err != nil {
		// 如果没有手机号认证记录，创建一个新的
		userAuth = &model.UserAuth{
			UserID:   userID,
			AuthType: "phone",
			AuthKey:  req.Phone,
		}
	} else {
		userAuth.AuthKey = req.Phone
	}

	if userAuth.ID == 0 {
		return s.userAuthRepo.Create(userAuth)
	} else {
		return s.userAuthRepo.Update(userAuth)
	}
}