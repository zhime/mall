package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := GenerateRandomString(6)
	return fmt.Sprintf("ORD%s%s", timestamp, random)
}

// GeneratePaymentNo 生成支付单号
func GeneratePaymentNo() string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := GenerateRandomString(6)
	return fmt.Sprintf("PAY%s%s", timestamp, random)
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MD5Hash MD5加密
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SHA256Hash SHA256加密
func SHA256Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GenerateSalt 生成盐值
func GenerateSalt() string {
	return GenerateRandomString(32)
}

// GenerateVerifyCode 生成验证码
func GenerateVerifyCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

// IsValidPhone 验证手机号格式
func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(phone)
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

// FormatPrice 格式化价格
func FormatPrice(price float64) string {
	return fmt.Sprintf("%.2f", price)
}

// ParsePrice 解析价格字符串
func ParsePrice(priceStr string) (float64, error) {
	return strconv.ParseFloat(priceStr, 64)
}

// SliceContains 检查切片是否包含元素
func SliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveFromSlice 从切片中移除元素
func RemoveFromSlice(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// StringInSlice 检查字符串是否在切片中
func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// TrimAndLower 去除空格并转小写
func TrimAndLower(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return strings.ToLower(parts[len(parts)-1])
	}
	return ""
}

// IsAllowedFileType 检查文件类型是否允许
func IsAllowedFileType(contentType string, allowedTypes []string) bool {
	return StringInSlice(contentType, allowedTypes)
}

// PaginationOffset 计算分页偏移量
func PaginationOffset(page, pageSize int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * pageSize
}

// PaginationLimit 获取分页限制
func PaginationLimit(pageSize int) int {
	if pageSize <= 0 || pageSize > 100 {
		return 20 // 默认每页20条
	}
	return pageSize
}