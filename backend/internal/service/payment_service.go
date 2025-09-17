package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"mall/internal/model"
	"mall/internal/repository"
	"mall/pkg/config"
	"mall/pkg/utils"
)

// PaymentService 支付服务接口
type PaymentService interface {
	CreatePayment(userID uint64, req *CreatePaymentRequest) (*CreatePaymentResponse, error)
	ProcessWechatCallback(req *WechatCallbackRequest) error
	ProcessAlipayCallback(req *AlipayCallbackRequest) error
	GetPaymentStatus(paymentNo string) (*PaymentStatusResponse, error)
	CancelPayment(paymentNo string) error
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID       uint64 `json:"order_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"` // wechat, alipay
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	PaymentNo     string                 `json:"payment_no"`
	PaymentMethod string                 `json:"payment_method"`
	Amount        float64                `json:"amount"`
	PaymentParams map[string]interface{} `json:"payment_params"`
}

// WechatCallbackRequest 微信支付回调请求
type WechatCallbackRequest struct {
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	TotalFee      int    `json:"total_fee"`
	ResultCode    string `json:"result_code"`
	ReturnCode    string `json:"return_code"`
}

// AlipayCallbackRequest 支付宝回调请求
type AlipayCallbackRequest struct {
	OutTradeNo  string  `json:"out_trade_no"`
	TradeNo     string  `json:"trade_no"`
	TradeStatus string  `json:"trade_status"`
	TotalAmount float64 `json:"total_amount"`
}

// PaymentStatusResponse 支付状态响应
type PaymentStatusResponse struct {
	PaymentNo     string  `json:"payment_no"`
	OrderID       uint64  `json:"order_id"`
	Status        int8    `json:"status"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	TradeNo       string  `json:"trade_no"`
	PayTime       string  `json:"pay_time"`
}

// paymentService 支付服务实现
type paymentService struct {
	paymentRepo repository.OrderPaymentRepository
	orderRepo   repository.OrderRepository
}

// NewPaymentService 创建支付服务
func NewPaymentService(
	paymentRepo repository.OrderPaymentRepository,
	orderRepo repository.OrderRepository,
) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

// CreatePayment 创建支付
func (s *paymentService) CreatePayment(userID uint64, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// 获取订单信息
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// 验证订单归属
	if order.UserID != userID {
		return nil, errors.New("access denied")
	}

	// 验证订单状态
	if order.Status != 1 {
		return nil, errors.New("order is not pending payment")
	}

	if order.PaymentStatus == 1 {
		return nil, errors.New("order is already paid")
	}

	// 生成支付单号
	paymentNo := utils.GeneratePaymentNo()

	// 创建支付记录
	payment := &model.OrderPayment{
		OrderID:       req.OrderID,
		PaymentNo:     paymentNo,
		PaymentMethod: req.PaymentMethod,
		Amount:        order.PayAmount,
		Status:        0, // 待支付
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, errors.New("failed to create payment")
	}

	// 根据支付方式生成支付参数
	var paymentParams map[string]interface{}
	switch req.PaymentMethod {
	case "wechat":
		paymentParams = s.generateWechatPayParams(payment, order)
	case "alipay":
		paymentParams = s.generateAlipayParams(payment, order)
	default:
		return nil, errors.New("unsupported payment method")
	}

	return &CreatePaymentResponse{
		PaymentNo:     paymentNo,
		PaymentMethod: req.PaymentMethod,
		Amount:        order.PayAmount,
		PaymentParams: paymentParams,
	}, nil
}

// generateWechatPayParams 生成微信支付参数
func (s *paymentService) generateWechatPayParams(payment *model.OrderPayment, order *model.Order) map[string]interface{} {
	cfg := config.GetConfig()
	
	// 微信支付统一下单参数
	params := map[string]interface{}{
		"appid":            cfg.Wechat.AppID,
		"mch_id":           cfg.Wechat.MchID,
		"nonce_str":        utils.GenerateRandomString(32),
		"body":             fmt.Sprintf("订单支付-%s", order.OrderNo),
		"out_trade_no":     payment.PaymentNo,
		"total_fee":        int(payment.Amount * 100), // 分为单位
		"spbill_create_ip": "127.0.0.1",
		"notify_url":       "https://your-domain.com/api/v1/payment/wechat/callback",
		"trade_type":       "JSAPI",
	}

	// 生成签名
	sign := s.generateWechatSign(params, cfg.Wechat.APIKey)
	params["sign"] = sign

	// 这里应该调用微信统一下单API，获取prepay_id
	// 简化处理，直接返回模拟参数
	return map[string]interface{}{
		"prepay_id": "mock_prepay_id_" + payment.PaymentNo,
		"app_id":    cfg.Wechat.AppID,
		"time_stamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonce_str":  utils.GenerateRandomString(32),
		"package":    "prepay_id=mock_prepay_id_" + payment.PaymentNo,
		"sign_type":  "MD5",
	}
}

// generateAlipayParams 生成支付宝支付参数
func (s *paymentService) generateAlipayParams(payment *model.OrderPayment, order *model.Order) map[string]interface{} {
	cfg := config.GetConfig()
	
	// 支付宝支付参数
	params := map[string]interface{}{
		"app_id":         cfg.Alipay.AppID,
		"method":         "alipay.trade.app.pay",
		"charset":        "utf-8",
		"sign_type":      "RSA2",
		"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
		"version":        "1.0",
		"notify_url":     "https://your-domain.com/api/v1/payment/alipay/callback",
		"out_trade_no":   payment.PaymentNo,
		"total_amount":   fmt.Sprintf("%.2f", payment.Amount),
		"subject":        fmt.Sprintf("订单支付-%s", order.OrderNo),
		"product_code":   "QUICK_MSECURITY_PAY",
	}

	// 这里应该使用RSA私钥生成签名
	// 简化处理
	params["sign"] = "mock_alipay_sign"

	return params
}

// generateWechatSign 生成微信支付签名
func (s *paymentService) generateWechatSign(params map[string]interface{}, apiKey string) string {
	// 排序参数
	var keys []string
	for key := range params {
		if key != "sign" && params[key] != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	// 构建签名字符串
	var signStr strings.Builder
	for i, key := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(fmt.Sprintf("%s=%v", key, params[key]))
	}
	signStr.WriteString("&key=" + apiKey)

	// MD5签名
	hash := md5.Sum([]byte(signStr.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// ProcessWechatCallback 处理微信支付回调
func (s *paymentService) ProcessWechatCallback(req *WechatCallbackRequest) error {
	// 获取支付记录
	payment, err := s.paymentRepo.GetByPaymentNo(req.OutTradeNo)
	if err != nil {
		return errors.New("payment not found")
	}

	// 验证支付结果
	if req.ReturnCode == "SUCCESS" && req.ResultCode == "SUCCESS" {
		// 支付成功
		if err := s.paymentRepo.UpdateStatus(payment.ID, 1, req.TransactionID); err != nil {
			return err
		}

		// 更新订单状态
		return s.orderRepo.UpdateStatus(payment.OrderID, 2) // 待发货
	} else {
		// 支付失败
		return s.paymentRepo.UpdateStatus(payment.ID, 2, "")
	}
}

// ProcessAlipayCallback 处理支付宝回调
func (s *paymentService) ProcessAlipayCallback(req *AlipayCallbackRequest) error {
	// 获取支付记录
	payment, err := s.paymentRepo.GetByPaymentNo(req.OutTradeNo)
	if err != nil {
		return errors.New("payment not found")
	}

	// 验证支付结果
	if req.TradeStatus == "TRADE_SUCCESS" || req.TradeStatus == "TRADE_FINISHED" {
		// 支付成功
		if err := s.paymentRepo.UpdateStatus(payment.ID, 1, req.TradeNo); err != nil {
			return err
		}

		// 更新订单状态
		return s.orderRepo.UpdateStatus(payment.OrderID, 2) // 待发货
	} else {
		// 支付失败
		return s.paymentRepo.UpdateStatus(payment.ID, 2, "")
	}
}

// GetPaymentStatus 获取支付状态
func (s *paymentService) GetPaymentStatus(paymentNo string) (*PaymentStatusResponse, error) {
	payment, err := s.paymentRepo.GetByPaymentNo(paymentNo)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	response := &PaymentStatusResponse{
		PaymentNo:     payment.PaymentNo,
		OrderID:       payment.OrderID,
		Status:        payment.Status,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		TradeNo:       payment.TradeNo,
		PayTime:       payment.PayTime.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// CancelPayment 取消支付
func (s *paymentService) CancelPayment(paymentNo string) error {
	payment, err := s.paymentRepo.GetByPaymentNo(paymentNo)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.Status != 0 {
		return errors.New("payment cannot be cancelled")
	}

	// 更新支付状态为已取消
	return s.paymentRepo.UpdateStatus(payment.ID, 3, "")
}