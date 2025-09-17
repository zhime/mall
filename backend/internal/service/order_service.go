package service

import (
	"errors"
	"fmt"

	"mall/internal/model"
	"mall/internal/repository"
	"mall/pkg/utils"
)

// OrderService 订单服务接口
type OrderService interface {
	CreateOrder(userID uint64, req *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrderDetail(userID uint64, orderID uint64) (*OrderDetailResponse, error)
	GetUserOrders(userID uint64, req *OrderListRequest) (*OrderListResponse, error)
	CancelOrder(userID uint64, orderID uint64) error
	ConfirmOrder(userID uint64, orderID uint64) error
	
	// 管理员接口
	GetOrders(req *AdminOrderListRequest) (*OrderListResponse, error)
	UpdateOrderStatus(orderID uint64, status int8) error
	SearchOrders(keyword string, page, pageSize int) (*OrderListResponse, error)
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Items         []OrderItemRequest `json:"items" binding:"required"`
	ReceiverName  string             `json:"receiver_name" binding:"required"`
	ReceiverPhone string             `json:"receiver_phone" binding:"required"`
	ReceiverAddress string           `json:"receiver_address" binding:"required"`
	BuyerMessage  string             `json:"buyer_message"`
}

// OrderItemRequest 订单项请求
type OrderItemRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	SKUID     uint64 `json:"sku_id"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderID     uint64  `json:"order_id"`
	OrderNo     string  `json:"order_no"`
	TotalAmount float64 `json:"total_amount"`
	PayAmount   float64 `json:"pay_amount"`
}

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	Page     int  `json:"page" form:"page"`
	PageSize int  `json:"page_size" form:"page_size"`
	Status   int8 `json:"status" form:"status"`
}

// AdminOrderListRequest 管理员订单列表请求
type AdminOrderListRequest struct {
	Page     int  `json:"page" form:"page"`
	PageSize int  `json:"page_size" form:"page_size"`
	Status   int8 `json:"status" form:"status"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Items      []*OrderResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID              uint64              `json:"id"`
	OrderNo         string              `json:"order_no"`
	UserID          uint64              `json:"user_id"`
	Username        string              `json:"username,omitempty"`
	TotalAmount     float64             `json:"total_amount"`
	PayAmount       float64             `json:"pay_amount"`
	FreightAmount   float64             `json:"freight_amount"`
	DiscountAmount  float64             `json:"discount_amount"`
	Status          int8                `json:"status"`
	PaymentStatus   int8                `json:"payment_status"`
	DeliveryStatus  int8                `json:"delivery_status"`
	ReceiverName    string              `json:"receiver_name"`
	ReceiverPhone   string              `json:"receiver_phone"`
	ReceiverAddress string              `json:"receiver_address"`
	BuyerMessage    string              `json:"buyer_message"`
	Items           []*OrderItemResponse `json:"items"`
	CreatedAt       string              `json:"created_at"`
}

// OrderDetailResponse 订单详情响应
type OrderDetailResponse struct {
	*OrderResponse
	Payments []*OrderPaymentResponse `json:"payments"`
}

// OrderItemResponse 订单项响应
type OrderItemResponse struct {
	ID           uint64  `json:"id"`
	ProductID    uint64  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	SKUID        uint64  `json:"sku_id"`
	SKUName      string  `json:"sku_name"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	TotalAmount  float64 `json:"total_amount"`
}

// OrderPaymentResponse 订单支付响应
type OrderPaymentResponse struct {
	ID            uint64  `json:"id"`
	PaymentNo     string  `json:"payment_no"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Status        int8    `json:"status"`
	TradeNo       string  `json:"trade_no"`
	PayTime       string  `json:"pay_time"`
}

// orderService 订单服务实现
type orderService struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	productRepo   repository.ProductRepository
	skuRepo       repository.ProductSKURepository
	cartRepo      repository.CartRepository
	userRepo      repository.UserRepository
}

// NewOrderService 创建订单服务
func NewOrderService(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	productRepo repository.ProductRepository,
	skuRepo repository.ProductSKURepository,
	cartRepo repository.CartRepository,
	userRepo repository.UserRepository,
) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		productRepo:   productRepo,
		skuRepo:       skuRepo,
		cartRepo:      cartRepo,
		userRepo:      userRepo,
	}
}

// CreateOrder 创建订单
func (s *orderService) CreateOrder(userID uint64, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// 验证用户是否存在
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 验证商品和计算总金额
	var totalAmount float64
	var orderItems []*model.OrderItem
	
	for _, itemReq := range req.Items {
		// 获取商品信息
		product, err := s.productRepo.GetByID(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found", itemReq.ProductID)
		}

		if product.Status != 1 {
			return nil, fmt.Errorf("product %s is not available", product.Name)
		}

		var price float64
		var skuName string
		var productImage string

		// 如果指定了SKU
		if itemReq.SKUID > 0 {
			sku, err := s.skuRepo.GetByID(itemReq.SKUID)
			if err != nil {
				return nil, fmt.Errorf("SKU %d not found", itemReq.SKUID)
			}
			if sku.ProductID != itemReq.ProductID {
				return nil, errors.New("SKU does not belong to this product")
			}
			if sku.Stock < itemReq.Quantity {
				return nil, fmt.Errorf("insufficient stock for SKU %s", sku.Name)
			}
			price = sku.Price
			skuName = sku.Name
			productImage = sku.Image
		} else {
			if product.Stock < itemReq.Quantity {
				return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
			}
			price = product.Price
		}

		// 如果没有SKU图片，使用商品主图
		if productImage == "" {
			productWithImages, err := s.productRepo.GetWithDetails(itemReq.ProductID)
			if err == nil {
				for _, image := range productWithImages.Images {
					if image.IsMain == 1 {
						productImage = image.ImageURL
						break
					}
				}
			}
		}

		itemTotalAmount := price * float64(itemReq.Quantity)
		totalAmount += itemTotalAmount

		orderItem := &model.OrderItem{
			ProductID:    itemReq.ProductID,
			SKUID:        itemReq.SKUID,
			ProductName:  product.Name,
			SKUName:      skuName,
			ProductImage: productImage,
			Price:        price,
			Quantity:     itemReq.Quantity,
			TotalAmount:  itemTotalAmount,
		}

		orderItems = append(orderItems, orderItem)
	}

	// 生成订单号
	orderNo := utils.GenerateOrderNo()

	// 创建订单
	order := &model.Order{
		UserID:          userID,
		OrderNo:         orderNo,
		TotalAmount:     totalAmount,
		PayAmount:       totalAmount,
		FreightAmount:   0,
		DiscountAmount:  0,
		Status:          1, // 待付款
		PaymentStatus:   0, // 未付款
		DeliveryStatus:  0, // 未发货
		BuyerMessage:    req.BuyerMessage,
		ReceiverName:    req.ReceiverName,
		ReceiverPhone:   req.ReceiverPhone,
		ReceiverAddress: req.ReceiverAddress,
	}

	// 使用事务创建订单
	err = s.createOrderWithTransaction(order, orderItems)
	if err != nil {
		return nil, err
	}

	return &CreateOrderResponse{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		TotalAmount: order.TotalAmount,
		PayAmount:   order.PayAmount,
	}, nil
}

// createOrderWithTransaction 使用事务创建订单
func (s *orderService) createOrderWithTransaction(order *model.Order, orderItems []*model.OrderItem) error {
	// TODO: 实现数据库事务
	// 这里暂时简化处理，实际项目中需要使用数据库事务

	// 创建订单
	if err := s.orderRepo.Create(order); err != nil {
		return err
	}

	// 创建订单项
	for _, item := range orderItems {
		item.OrderID = order.ID
	}
	if err := s.orderItemRepo.BatchCreate(orderItems); err != nil {
		return err
	}

	// 扣减库存
	for _, item := range orderItems {
		if item.SKUID > 0 {
			// 扣减SKU库存
			sku, _ := s.skuRepo.GetByID(item.SKUID)
			if sku != nil {
				s.skuRepo.UpdateStock(item.SKUID, sku.Stock-item.Quantity)
			}
		} else {
			// 扣减商品库存
			product, _ := s.productRepo.GetByID(item.ProductID)
			if product != nil {
				s.productRepo.UpdateStock(item.ProductID, product.Stock-item.Quantity)
			}
		}
	}

	return nil
}

// GetOrderDetail 获取订单详情
func (s *orderService) GetOrderDetail(userID uint64, orderID uint64) (*OrderDetailResponse, error) {
	order, err := s.orderRepo.GetWithDetails(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// 验证订单是否属于该用户
	if order.UserID != userID {
		return nil, errors.New("access denied")
	}

	return s.toOrderDetailResponse(order), nil
}

// GetUserOrders 获取用户订单列表
func (s *orderService) GetUserOrders(userID uint64, req *OrderListRequest) (*OrderListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	orders, total, err := s.orderRepo.GetUserOrders(userID, req.Page, req.PageSize, req.Status)
	if err != nil {
		return nil, err
	}

	var items []*OrderResponse
	for _, order := range orders {
		items = append(items, s.toOrderResponse(order))
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &OrderListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// CancelOrder 取消订单
func (s *orderService) CancelOrder(userID uint64, orderID uint64) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("access denied")
	}

	if order.Status != 1 {
		return errors.New("order cannot be cancelled")
	}

	// 更新订单状态为已取消
	return s.orderRepo.UpdateStatus(orderID, 5)
}

// ConfirmOrder 确认收货
func (s *orderService) ConfirmOrder(userID uint64, orderID uint64) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("access denied")
	}

	if order.Status != 3 {
		return errors.New("order is not shipped")
	}

	// 更新订单状态为已完成
	return s.orderRepo.UpdateStatus(orderID, 4)
}

// GetOrders 获取订单列表（管理员）
func (s *orderService) GetOrders(req *AdminOrderListRequest) (*OrderListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	orders, total, err := s.orderRepo.GetOrders(req.Page, req.PageSize, req.Status)
	if err != nil {
		return nil, err
	}

	var items []*OrderResponse
	for _, order := range orders {
		response := s.toOrderResponse(order)
		// 添加用户信息
		if order.User.ID > 0 {
			response.Username = order.User.Username
		}
		items = append(items, response)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &OrderListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateOrderStatus 更新订单状态
func (s *orderService) UpdateOrderStatus(orderID uint64, status int8) error {
	return s.orderRepo.UpdateStatus(orderID, status)
}

// SearchOrders 搜索订单
func (s *orderService) SearchOrders(keyword string, page, pageSize int) (*OrderListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := s.orderRepo.Search(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	var items []*OrderResponse
	for _, order := range orders {
		response := s.toOrderResponse(order)
		if order.User.ID > 0 {
			response.Username = order.User.Username
		}
		items = append(items, response)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &OrderListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// toOrderResponse 转换为订单响应
func (s *orderService) toOrderResponse(order *model.Order) *OrderResponse {
	response := &OrderResponse{
		ID:              order.ID,
		OrderNo:         order.OrderNo,
		UserID:          order.UserID,
		TotalAmount:     order.TotalAmount,
		PayAmount:       order.PayAmount,
		FreightAmount:   order.FreightAmount,
		DiscountAmount:  order.DiscountAmount,
		Status:          order.Status,
		PaymentStatus:   order.PaymentStatus,
		DeliveryStatus:  order.DeliveryStatus,
		ReceiverName:    order.ReceiverName,
		ReceiverPhone:   order.ReceiverPhone,
		ReceiverAddress: order.ReceiverAddress,
		BuyerMessage:    order.BuyerMessage,
		CreatedAt:       order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 添加订单项
	for _, item := range order.Items {
		response.Items = append(response.Items, &OrderItemResponse{
			ID:           item.ID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			SKUID:        item.SKUID,
			SKUName:      item.SKUName,
			Price:        item.Price,
			Quantity:     item.Quantity,
			TotalAmount:  item.TotalAmount,
		})
	}

	return response
}

// toOrderDetailResponse 转换为订单详情响应
func (s *orderService) toOrderDetailResponse(order *model.Order) *OrderDetailResponse {
	response := &OrderDetailResponse{
		OrderResponse: s.toOrderResponse(order),
	}

	// 添加支付信息
	for _, payment := range order.Payments {
		response.Payments = append(response.Payments, &OrderPaymentResponse{
			ID:            payment.ID,
			PaymentNo:     payment.PaymentNo,
			PaymentMethod: payment.PaymentMethod,
			Amount:        payment.Amount,
			Status:        payment.Status,
			TradeNo:       payment.TradeNo,
			PayTime:       payment.PayTime.Format("2006-01-02 15:04:05"),
		})
	}

	return response
}