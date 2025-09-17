package routes

import (
	"github.com/gin-gonic/gin"

	"mall/internal/handler"
)

// ProductRoutes 商品路由组
type ProductRoutes struct {
	productHandler  *handler.ProductHandler
	categoryHandler *handler.CategoryHandler
}

// NewProductRoutes 创建商品路由组
func NewProductRoutes(productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler) *ProductRoutes {
	return &ProductRoutes{
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
	}
}

// RegisterRoutes 注册商品和分类相关路由
func (r *ProductRoutes) RegisterRoutes(router *gin.RouterGroup) {
	// 分类相关路由
	categories := router.Group("/categories")
	{
		categories.GET("", r.categoryHandler.GetCategoryTree)
		categories.GET("/level", r.categoryHandler.GetCategoriesByLevel)
		categories.GET("/search", r.categoryHandler.SearchCategories)
		categories.GET("/:id", r.categoryHandler.GetCategory)
	}

	// 商品相关路由
	products := router.Group("/products")
	{
		products.GET("", r.productHandler.GetProductList)
		products.GET("/search", r.productHandler.SearchProducts)
		products.GET("/hot", r.productHandler.GetHotProducts)
		products.GET("/category/:categoryId", r.productHandler.GetProductsByCategory)
		products.GET("/:id", r.productHandler.GetProduct)
		products.GET("/:id/detail", r.productHandler.GetProductDetail)
	}
}