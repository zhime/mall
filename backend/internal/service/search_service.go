package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mall/internal/model"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type SearchService struct {
	es               *elasticsearch.Client
	productRepo      *ProductRepository
	productIndexName string
}

type SearchRequest struct {
	Keyword    string  `json:"keyword"`
	CategoryID int     `json:"category_id,omitempty"`
	PriceMin   float64 `json:"price_min,omitempty"`
	PriceMax   float64 `json:"price_max,omitempty"`
	Sort       string  `json:"sort,omitempty"` // sales, price_asc, price_desc, newest
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
}

type SearchResponse struct {
	Products []model.Product `json:"products"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

type ProductDocument struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	CategoryID  int      `json:"category_id"`
	Price       float64  `json:"price"`
	Sales       int      `json:"sales"`
	Images      []string `json:"images"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	IsFeatured  bool     `json:"is_featured"`
	IsNew       bool     `json:"is_new"`
	CreatedAt   string   `json:"created_at"`
}

func NewSearchService(es *elasticsearch.Client, productRepo *ProductRepository) *SearchService {
	return &SearchService{
		es:               es,
		productRepo:      productRepo,
		productIndexName: "products",
	}
}

// 初始化搜索索引
func (s *SearchService) InitIndex() error {
	// 创建索引映射
	mapping := `{
		"mappings": {
			"properties": {
				"id": { "type": "integer" },
				"name": { 
					"type": "text", 
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_smart"
				},
				"category_id": { "type": "integer" },
				"price": { "type": "double" },
				"sales": { "type": "integer" },
				"images": { "type": "keyword" },
				"description": { 
					"type": "text", 
					"analyzer": "ik_max_word" 
				},
				"status": { "type": "integer" },
				"is_featured": { "type": "boolean" },
				"is_new": { "type": "boolean" },
				"created_at": { "type": "date" }
			}
		},
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: s.productIndexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(context.Background(), s.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && !strings.Contains(res.String(), "resource_already_exists_exception") {
		return fmt.Errorf("failed to create index: %s", res.String())
	}

	return nil
}

// 索引商品数据
func (s *SearchService) IndexProduct(product *model.Product) error {
	doc := ProductDocument{
		ID:          product.ID,
		Name:        product.Name,
		CategoryID:  product.CategoryID,
		Price:       product.Price,
		Sales:       product.Sales,
		Images:      product.Images,
		Description: product.Description,
		Status:      product.Status,
		IsFeatured:  product.IsFeatured,
		IsNew:       product.IsNew,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      s.productIndexName,
		DocumentID: fmt.Sprintf("%d", product.ID),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), s.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index product: %s", res.String())
	}

	return nil
}

// 删除商品索引
func (s *SearchService) DeleteProduct(productID int) error {
	req := esapi.DeleteRequest{
		Index:      s.productIndexName,
		DocumentID: fmt.Sprintf("%d", productID),
	}

	res, err := req.Do(context.Background(), s.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// 搜索商品
func (s *SearchService) SearchProducts(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	query := s.buildSearchQuery(req)
	
	searchBody := map[string]interface{}{
		"query": query,
		"sort":  s.buildSortQuery(req.Sort),
		"from":  (req.Page - 1) * req.PageSize,
		"size":  req.PageSize,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return nil, err
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(s.productIndexName),
		s.es.Search.WithBody(&buf),
		s.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	return s.parseSearchResponse(res, req)
}

// 构建搜索查询
func (s *SearchService) buildSearchQuery(req *SearchRequest) map[string]interface{} {
	boolQuery := map[string]interface{}{
		"must": []map[string]interface{}{
			{"term": map[string]interface{}{"status": 1}}, // 只搜索上架商品
		},
	}

	// 关键词搜索
	if req.Keyword != "" {
		boolQuery["must"] = append(boolQuery["must"].([]map[string]interface{}), map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":    req.Keyword,
				"fields":   []string{"name^2", "description"},
				"type":     "best_fields",
				"operator": "and",
			},
		})
	}

	// 分类筛选
	if req.CategoryID > 0 {
		boolQuery["must"] = append(boolQuery["must"].([]map[string]interface{}), map[string]interface{}{
			"term": map[string]interface{}{"category_id": req.CategoryID},
		})
	}

	// 价格范围筛选
	if req.PriceMin > 0 || req.PriceMax > 0 {
		rangeQuery := map[string]interface{}{}
		if req.PriceMin > 0 {
			rangeQuery["gte"] = req.PriceMin
		}
		if req.PriceMax > 0 {
			rangeQuery["lte"] = req.PriceMax
		}
		boolQuery["must"] = append(boolQuery["must"].([]map[string]interface{}), map[string]interface{}{
			"range": map[string]interface{}{"price": rangeQuery},
		})
	}

	return map[string]interface{}{
		"bool": boolQuery,
	}
}

// 构建排序查询
func (s *SearchService) buildSortQuery(sort string) []map[string]interface{} {
	switch sort {
	case "sales":
		return []map[string]interface{}{
			{"sales": map[string]interface{}{"order": "desc"}},
		}
	case "price_asc":
		return []map[string]interface{}{
			{"price": map[string]interface{}{"order": "asc"}},
		}
	case "price_desc":
		return []map[string]interface{}{
			{"price": map[string]interface{}{"order": "desc"}},
		}
	case "newest":
		return []map[string]interface{}{
			{"created_at": map[string]interface{}{"order": "desc"}},
		}
	default:
		return []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
			{"sales": map[string]interface{}{"order": "desc"}},
		}
	}
}

// 解析搜索结果
func (s *SearchService) parseSearchResponse(res *esapi.Response, req *SearchRequest) (*SearchResponse, error) {
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))
	
	var productIDs []int
	for _, hit := range hits["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		productIDs = append(productIDs, int(source["id"].(float64)))
	}

	// 从数据库获取完整的商品信息
	products, err := s.productRepo.GetProductsByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	return &SearchResponse{
		Products: products,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// 获取搜索建议
func (s *SearchService) GetSearchSuggestions(ctx context.Context, keyword string, limit int) ([]string, error) {
	if keyword == "" {
		return []string{}, nil
	}

	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			"product_suggest": map[string]interface{}{
				"prefix": keyword,
				"completion": map[string]interface{}{
					"field": "name.suggest",
					"size":  limit,
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(s.productIndexName),
		s.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var suggestions []string
	if suggest, ok := result["suggest"].(map[string]interface{}); ok {
		if productSuggest, ok := suggest["product_suggest"].([]interface{}); ok && len(productSuggest) > 0 {
			options := productSuggest[0].(map[string]interface{})["options"].([]interface{})
			for _, option := range options {
				text := option.(map[string]interface{})["text"].(string)
				suggestions = append(suggestions, text)
			}
		}
	}

	return suggestions, nil
}

// 批量索引商品
func (s *SearchService) BulkIndexProducts() error {
	products, err := s.productRepo.GetAllProducts()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, product := range products {
		doc := ProductDocument{
			ID:          product.ID,
			Name:        product.Name,
			CategoryID:  product.CategoryID,
			Price:       product.Price,
			Sales:       product.Sales,
			Images:      product.Images,
			Description: product.Description,
			Status:      product.Status,
			IsFeatured:  product.IsFeatured,
			IsNew:       product.IsNew,
			CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": s.productIndexName,
				"_id":    product.ID,
			},
		}

		metaData, _ := json.Marshal(meta)
		docData, _ := json.Marshal(doc)

		buf.Write(metaData)
		buf.WriteByte('\n')
		buf.Write(docData)
		buf.WriteByte('\n')
	}

	req := esapi.BulkRequest{
		Body:    &buf,
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), s.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index error: %s", res.String())
	}

	return nil
}