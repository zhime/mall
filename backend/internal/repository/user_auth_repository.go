package repository

import (
	"gorm.io/gorm"

	"mall/internal/model"
)

// UserAuthRepository 用户认证仓储接口
type UserAuthRepository interface {
	Create(auth *model.UserAuth) error
	GetByUserIDAndType(userID uint64, authType string) (*model.UserAuth, error)
	GetByAuthKey(authKey string) (*model.UserAuth, error)
	Update(auth *model.UserAuth) error
	Delete(id uint64) error
}

// userAuthRepository 用户认证仓储实现
type userAuthRepository struct {
	db *gorm.DB
}

// NewUserAuthRepository 创建用户认证仓储
func NewUserAuthRepository(db *gorm.DB) UserAuthRepository {
	return &userAuthRepository{db: db}
}

// Create 创建认证信息
func (r *userAuthRepository) Create(auth *model.UserAuth) error {
	return r.db.Create(auth).Error
}

// GetByUserIDAndType 根据用户ID和认证类型获取认证信息
func (r *userAuthRepository) GetByUserIDAndType(userID uint64, authType string) (*model.UserAuth, error) {
	var auth model.UserAuth
	err := r.db.Where("user_id = ? AND auth_type = ?", userID, authType).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// GetByAuthKey 根据认证key获取认证信息
func (r *userAuthRepository) GetByAuthKey(authKey string) (*model.UserAuth, error) {
	var auth model.UserAuth
	err := r.db.Preload("User").Where("auth_key = ?", authKey).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// Update 更新认证信息
func (r *userAuthRepository) Update(auth *model.UserAuth) error {
	return r.db.Save(auth).Error
}

// Delete 删除认证信息
func (r *userAuthRepository) Delete(id uint64) error {
	return r.db.Delete(&model.UserAuth{}, id).Error
}