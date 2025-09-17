package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint64         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// User 用户基础信息
type User struct {
	BaseModel
	Username     string `json:"username" gorm:"size:50;uniqueIndex"`
	Phone        string `json:"phone" gorm:"size:20;uniqueIndex"`
	Email        string `json:"email" gorm:"size:100"`
	WechatOpenID string `json:"wechat_openid" gorm:"column:wechat_openid;size:100;uniqueIndex"`
	Status       int8   `json:"status" gorm:"default:1;comment:1正常 0禁用"`

	// 关联
	Profile   *UserProfile   `json:"profile,omitempty"`
	Addresses []UserAddress  `json:"addresses,omitempty"`
	AuthInfos []UserAuth     `json:"auth_infos,omitempty"`
	Orders    []Order        `json:"orders,omitempty"`
	CartItems []CartItem     `json:"cart_items,omitempty"`
}

// UserProfile 用户详细信息
type UserProfile struct {
	BaseModel
	UserID   uint64    `json:"user_id" gorm:"not null;index"`
	Nickname string    `json:"nickname" gorm:"size:50"`
	Avatar   string    `json:"avatar" gorm:"size:500"`
	Gender   int8      `json:"gender" gorm:"default:0;comment:0未知 1男 2女"`
	Birthday time.Time `json:"birthday"`

	// 关联
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserAddress 用户地址
type UserAddress struct {
	BaseModel
	UserID    uint64 `json:"user_id" gorm:"not null;index"`
	Name      string `json:"name" gorm:"size:50;not null"`
	Phone     string `json:"phone" gorm:"size:20;not null"`
	Province  string `json:"province" gorm:"size:50;not null"`
	City      string `json:"city" gorm:"size:50;not null"`
	District  string `json:"district" gorm:"size:50;not null"`
	Address   string `json:"address" gorm:"size:200;not null"`
	IsDefault int8   `json:"is_default" gorm:"default:0;comment:0否 1是"`

	// 关联
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserAuth 用户认证信息
type UserAuth struct {
	BaseModel
	UserID       uint64 `json:"user_id" gorm:"not null;uniqueIndex:uk_user_auth"`
	AuthType     string `json:"auth_type" gorm:"size:20;not null;uniqueIndex:uk_user_auth;comment:password,wechat"`
	AuthKey      string `json:"auth_key" gorm:"size:100;not null;comment:手机号或openid"`
	PasswordHash string `json:"password_hash" gorm:"size:255"`
	Salt         string `json:"salt" gorm:"size:32"`

	// 关联
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Category 商品分类
type Category struct {
	BaseModel
	ParentID    uint64 `json:"parent_id" gorm:"default:0;index"`
	Name        string `json:"name" gorm:"size:100;not null"`
	Level       int8   `json:"level" gorm:"not null;default:1"`
	SortOrder   int    `json:"sort_order" gorm:"default:0;index"`
	Icon        string `json:"icon" gorm:"size:500"`
	Description string `json:"description" gorm:"type:text"`
	Status      int8   `json:"status" gorm:"default:1;comment:1启用 0禁用"`

	// 关联
	Products []Product  `json:"products,omitempty"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
}

// Product 商品基础信息
type Product struct {
	BaseModel
	CategoryID    uint64  `json:"category_id" gorm:"not null;index"`
	Name          string  `json:"name" gorm:"size:200;not null"`
	Subtitle      string  `json:"subtitle" gorm:"size:500"`
	Description   string  `json:"description" gorm:"type:text"`
	Price         float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	OriginalPrice float64 `json:"original_price" gorm:"type:decimal(10,2)"`
	Stock         int     `json:"stock" gorm:"default:0"`
	Sales         int     `json:"sales" gorm:"default:0"`
	Status        int8    `json:"status" gorm:"default:1;comment:1上架 0下架"`
	SortOrder     int     `json:"sort_order" gorm:"default:0"`

	// 关联
	Category Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	SKUs     []ProductSKU   `json:"skus,omitempty"`
	Images   []ProductImage `json:"images,omitempty"`
}

// ProductSKU 商品SKU
type ProductSKU struct {
	BaseModel
	ProductID  uint64  `json:"product_id" gorm:"not null;index"`
	SKUCode    string  `json:"sku_code" gorm:"size:100;uniqueIndex;not null"`
	Name       string  `json:"name" gorm:"size:200"`
	Price      float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock      int     `json:"stock" gorm:"default:0"`
	AttrValues string  `json:"attr_values" gorm:"type:json;comment:属性值JSON"`
	Image      string  `json:"image" gorm:"size:500"`
	Status     int8    `json:"status" gorm:"default:1"`

	// 关联
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// ProductImage 商品图片
type ProductImage struct {
	BaseModel
	ProductID uint64 `json:"product_id" gorm:"not null;index"`
	ImageURL  string `json:"image_url" gorm:"size:500;not null"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
	IsMain    int8   `json:"is_main" gorm:"default:0;comment:0否 1是"`

	// 关联
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// Order 订单主表
type Order struct {
	BaseModel
	UserID          uint64  `json:"user_id" gorm:"not null;index"`
	OrderNo         string  `json:"order_no" gorm:"size:32;uniqueIndex;not null"`
	TotalAmount     float64 `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	PayAmount       float64 `json:"pay_amount" gorm:"type:decimal(10,2);not null"`
	FreightAmount   float64 `json:"freight_amount" gorm:"type:decimal(10,2);default:0"`
	DiscountAmount  float64 `json:"discount_amount" gorm:"type:decimal(10,2);default:0"`
	Status          int8    `json:"status" gorm:"default:1;comment:1待付款 2待发货 3已发货 4已完成 5已取消"`
	PaymentStatus   int8    `json:"payment_status" gorm:"default:0;comment:0未付款 1已付款"`
	DeliveryStatus  int8    `json:"delivery_status" gorm:"default:0;comment:0未发货 1已发货 2已收货"`
	BuyerMessage    string  `json:"buyer_message" gorm:"size:500"`
	ReceiverName    string  `json:"receiver_name" gorm:"size:50"`
	ReceiverPhone   string  `json:"receiver_phone" gorm:"size:20"`
	ReceiverAddress string  `json:"receiver_address" gorm:"size:500"`

	// 关联
	User     User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items    []OrderItem    `json:"items,omitempty"`
	Payments []OrderPayment `json:"payments,omitempty"`
}

// OrderItem 订单商品明细
type OrderItem struct {
	BaseModel
	OrderID      uint64  `json:"order_id" gorm:"not null;index"`
	ProductID    uint64  `json:"product_id" gorm:"not null"`
	SKUID        uint64  `json:"sku_id"`
	ProductName  string  `json:"product_name" gorm:"size:200;not null"`
	SKUName      string  `json:"sku_name" gorm:"size:200"`
	ProductImage string  `json:"product_image" gorm:"size:500"`
	Price        float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity     int     `json:"quantity" gorm:"not null"`
	TotalAmount  float64 `json:"total_amount" gorm:"type:decimal(10,2);not null"`

	// 关联
	Order   Order      `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Product Product    `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	SKU     ProductSKU `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// OrderPayment 订单支付记录
type OrderPayment struct {
	BaseModel
	OrderID       uint64    `json:"order_id" gorm:"not null;index"`
	PaymentNo     string    `json:"payment_no" gorm:"size:64;uniqueIndex;not null"`
	PaymentMethod string    `json:"payment_method" gorm:"size:20;not null;comment:wechat,alipay"`
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status        int8      `json:"status" gorm:"default:0;comment:0待支付 1支付成功 2支付失败"`
	TradeNo       string    `json:"trade_no" gorm:"size:64;comment:第三方交易号"`
	PayTime       time.Time `json:"pay_time"`

	// 关联
	Order Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// CartItem 购物车
type CartItem struct {
	BaseModel
	UserID    uint64 `json:"user_id" gorm:"not null;uniqueIndex:uk_user_product_sku"`
	ProductID uint64 `json:"product_id" gorm:"not null;uniqueIndex:uk_user_product_sku"`
	SKUID     uint64 `json:"sku_id" gorm:"uniqueIndex:uk_user_product_sku"`
	Quantity  int    `json:"quantity" gorm:"not null;default:1"`

	// 关联
	User    User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product Product    `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	SKU     ProductSKU `json:"sku,omitempty" gorm:"foreignKey:SKUID"`
}

// Admin 管理员
type Admin struct {
	BaseModel
	Username    string    `json:"username" gorm:"size:50;uniqueIndex;not null"`
	PasswordHash string   `json:"password_hash" gorm:"size:255;not null"`
	Salt        string    `json:"salt" gorm:"size:32;not null"`
	Nickname    string    `json:"nickname" gorm:"size:50"`
	Email       string    `json:"email" gorm:"size:100"`
	Phone       string    `json:"phone" gorm:"size:20"`
	Avatar      string    `json:"avatar" gorm:"size:500"`
	Role        string    `json:"role" gorm:"size:50;default:admin"`
	Status      int8      `json:"status" gorm:"default:1;comment:1正常 0禁用"`
	LastLoginAt time.Time `json:"last_login_at"`
}