package handler

import (
	"mall/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService service.SearchService
}

func NewSearchHandler(searchService service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// 搜索商品
// @Summary 搜索商品
// @Description 根据关键词搜索商品
// @Tags 搜索
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param category_id query int false "分类ID"
// @Param price_min query number false "最低价格"
// @Param price_max query number false "最高价格"
// @Param sort query string false "排序方式" Enums(sales, price_asc, price_desc, newest)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=service.SearchResponse}
// @Router /api/search [get]
func (h *SearchHandler) SearchProducts(c *gin.Context) {
	var req service.SearchRequest
	
	req.Keyword = c.Query("keyword")
	req.Sort = c.Query("sort")
	
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			req.CategoryID = categoryID
		}
	}
	
	if priceMinStr := c.Query("price_min"); priceMinStr != "" {
		if priceMin, err := strconv.ParseFloat(priceMinStr, 64); err == nil {
			req.PriceMin = priceMin
		}
	}
	
	if priceMaxStr := c.Query("price_max"); priceMaxStr != "" {
		if priceMax, err := strconv.ParseFloat(priceMaxStr, 64); err == nil {
			req.PriceMax = priceMax
		}
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	req.Page = page
	req.PageSize = pageSize

	result, err := h.searchService.SearchProducts(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索成功",
		"data":    result,
	})
}

// 获取搜索建议
// @Summary 获取搜索建议
// @Description 根据输入的关键词获取搜索建议
// @Tags 搜索
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param limit query int false "建议数量" default(10)
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/search/suggest [get]
func (h *SearchHandler) GetSearchSuggestions(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "关键词不能为空",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	suggestions, err := h.searchService.GetSearchSuggestions(c.Request.Context(), keyword, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取搜索建议失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取搜索建议成功",
		"data":    suggestions,
	})
}

// 获取热门搜索关键词
// @Summary 获取热门搜索关键词
// @Description 获取热门搜索关键词列表
// @Tags 搜索
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]string}
// @Router /api/search/hot [get]
func (h *SearchHandler) GetHotKeywords(c *gin.Context) {
	// 这里可以从缓存或数据库获取热门搜索关键词
	hotKeywords := []string{
		"手机",
		"电脑",
		"耳机",
		"充电器",
		"数据线",
		"移动电源",
		"蓝牙音箱",
		"智能手表",
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取热门搜索关键词成功",
		"data":    hotKeywords,
	})
}