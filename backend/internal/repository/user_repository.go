package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础CRUD
	Create(user *model.User) error
	GetByID(id uint64) (*model.User, error)
	GetByPhone(phone string) (*model.User, error)
	GetByWechatOpenID(openID string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint64) error

	// 复杂查询
	List(page, pageSize int) ([]*model.User, int64, error)
	Search(keyword string, page, pageSize int) ([]*model.User, int64, error)
	GetUserWithProfile(id uint64) (*model.User, error)
	GetUserWithAuth(id uint64) (*model.User, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *userRepository) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByWechatOpenID 根据微信OpenID获取用户
func (r *userRepository) GetByWechatOpenID(openID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("wechat_openid = ?", openID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (r *userRepository) List(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 计算总数
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// Search 搜索用户
func (r *userRepository) Search(keyword string, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{}).Where(
		"username LIKE ? OR phone LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%",
	)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// GetUserWithProfile 获取用户及其详细信息
func (r *userRepository) GetUserWithProfile(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Profile").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserWithAuth 获取用户及其认证信息
func (r *userRepository) GetUserWithAuth(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.Preload("AuthInfos").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}