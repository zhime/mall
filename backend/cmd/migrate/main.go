package main

import (
	"flag"
	"fmt"
	"log"

	"gorm.io/gorm"

	"mall/internal/model"
	"mall/pkg/config"
	"mall/pkg/database"
	"mall/pkg/logger"
	"mall/pkg/utils"
)

func main() {
	var action string
	flag.StringVar(&action, "action", "migrate", "Action to perform: migrate, drop, seed")
	flag.Parse()

	// 加载配置
	config.LoadConfig()

	// 初始化日志
	cfg := config.GetConfig()
	logger.InitLogger(cfg.Log.Level, cfg.Log.Filename, cfg.Log.MaxSize, cfg.Log.MaxAge, cfg.Log.MaxBackups)

	// 初始化数据库
	database.InitDB()
	defer database.CloseDB()

	db := database.GetDB()

	switch action {
	case "migrate":
		fmt.Println("Running database migrations...")
		if err := runMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migration completed successfully!")
		
	case "drop":
		fmt.Println("Dropping all tables...")
		if err := dropTables(db); err != nil {
			log.Fatalf("Drop tables failed: %v", err)
		}
		fmt.Println("All tables dropped successfully!")
		
	case "seed":
		fmt.Println("Seeding database...")
		if err := seedData(db); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		fmt.Println("Database seeded successfully!")
		
	default:
		fmt.Printf("Unknown action: %s\n", action)
		fmt.Println("Available actions: migrate, drop, seed")
	}
}

// runMigrations 执行数据库迁移
func runMigrations(db *gorm.DB) error {
	// 自动迁移所有模型
	return db.AutoMigrate(
		&model.User{},
		&model.UserProfile{},
		&model.UserAddress{},
		&model.UserAuth{},
		&model.Category{},
		&model.Product{},
		&model.ProductSKU{},
		&model.ProductImage{},
		&model.Order{},
		&model.OrderItem{},
		&model.OrderPayment{},
		&model.CartItem{},
		&model.Admin{},
	)
}

// dropTables 删除所有表
func dropTables(db *gorm.DB) error {
	// 按照依赖关系的逆序删除表
	tables := []interface{}{
		&model.CartItem{},
		&model.OrderPayment{},
		&model.OrderItem{},
		&model.Order{},
		&model.ProductImage{},
		&model.ProductSKU{},
		&model.Product{},
		&model.Category{},
		&model.UserAuth{},
		&model.UserAddress{},
		&model.UserProfile{},
		&model.User{},
		&model.Admin{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}

// seedData 播种初始数据
func seedData(db *gorm.DB) error {
	// 创建默认管理员
	if err := seedAdmins(db); err != nil {
		return err
	}

	// 创建默认分类
	if err := seedCategories(db); err != nil {
		return err
	}

	// 创建示例商品
	if err := seedProducts(db); err != nil {
		return err
	}

	return nil
}

// seedAdmins 创建默认管理员
func seedAdmins(db *gorm.DB) error {
	// 检查是否已存在管理员
	var count int64
	db.Model(&model.Admin{}).Count(&count)
	if count > 0 {
		fmt.Println("Admins already exist, skipping...")
		return nil
	}

	// 创建超级管理员
	passwordHash, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	admin := &model.Admin{
		Username:     "admin",
		PasswordHash: passwordHash,
		Salt:         utils.GenerateSalt(),
		Nickname:     "超级管理员",
		Role:         "super_admin",
		Status:       1,
	}

	return db.Create(admin).Error
}

// seedCategories 创建默认分类
func seedCategories(db *gorm.DB) error {
	// 检查是否已存在分类
	var count int64
	db.Model(&model.Category{}).Count(&count)
	if count > 0 {
		fmt.Println("Categories already exist, skipping...")
		return nil
	}

	// 一级分类
	categories := []model.Category{
		{Name: "服装鞋包", Level: 1, SortOrder: 1, Status: 1},
		{Name: "数码家电", Level: 1, SortOrder: 2, Status: 1},
		{Name: "食品饮料", Level: 1, SortOrder: 3, Status: 1},
		{Name: "美妆护肤", Level: 1, SortOrder: 4, Status: 1},
		{Name: "家居日用", Level: 1, SortOrder: 5, Status: 1},
	}

	for i := range categories {
		if err := db.Create(&categories[i]).Error; err != nil {
			return err
		}
	}

	// 二级分类
	subCategories := []model.Category{
		// 服装鞋包
		{ParentID: categories[0].ID, Name: "男装", Level: 2, SortOrder: 1, Status: 1},
		{ParentID: categories[0].ID, Name: "女装", Level: 2, SortOrder: 2, Status: 1},
		{ParentID: categories[0].ID, Name: "鞋靴", Level: 2, SortOrder: 3, Status: 1},
		// 数码家电
		{ParentID: categories[1].ID, Name: "手机数码", Level: 2, SortOrder: 1, Status: 1},
		{ParentID: categories[1].ID, Name: "家用电器", Level: 2, SortOrder: 2, Status: 1},
		// 食品饮料
		{ParentID: categories[2].ID, Name: "休闲食品", Level: 2, SortOrder: 1, Status: 1},
		{ParentID: categories[2].ID, Name: "饮料冲调", Level: 2, SortOrder: 2, Status: 1},
		// 美妆护肤
		{ParentID: categories[3].ID, Name: "面部护肤", Level: 2, SortOrder: 1, Status: 1},
		{ParentID: categories[3].ID, Name: "彩妆", Level: 2, SortOrder: 2, Status: 1},
		// 家居日用
		{ParentID: categories[4].ID, Name: "家纺用品", Level: 2, SortOrder: 1, Status: 1},
		{ParentID: categories[4].ID, Name: "清洁用品", Level: 2, SortOrder: 2, Status: 1},
	}

	for i := range subCategories {
		if err := db.Create(&subCategories[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

// seedProducts 创建示例商品
func seedProducts(db *gorm.DB) error {
	// 检查是否已存在商品
	var count int64
	db.Model(&model.Product{}).Count(&count)
	if count > 0 {
		fmt.Println("Products already exist, skipping...")
		return nil
	}

	// 获取第一个二级分类
	var category model.Category
	if err := db.Where("level = ? AND parent_id > 0", 2).First(&category).Error; err != nil {
		return err
	}

	// 创建示例商品
	products := []model.Product{
		{
			CategoryID:    category.ID,
			Name:          "经典白色T恤",
			Subtitle:      "纯棉舒适，经典百搭",
			Description:   "采用优质纯棉面料，柔软透气，经典设计，适合各种场合穿着。",
			Price:         99.00,
			OriginalPrice: 129.00,
			Stock:         100,
			Sales:         0,
			Status:        1,
			SortOrder:     1,
		},
		{
			CategoryID:    category.ID,
			Name:          "休闲牛仔裤",
			Subtitle:      "修身显瘦，舒适自然",
			Description:   "精选优质牛仔布料，修身设计，舒适穿着，展现个人魅力。",
			Price:         199.00,
			OriginalPrice: 259.00,
			Stock:         80,
			Sales:         0,
			Status:        1,
			SortOrder:     2,
		},
	}

	for i := range products {
		if err := db.Create(&products[i]).Error; err != nil {
			return err
		}

		// 为每个商品创建示例图片
		images := []model.ProductImage{
			{
				ProductID: products[i].ID,
				ImageURL:  fmt.Sprintf("https://example.com/images/product_%d_1.jpg", products[i].ID),
				SortOrder: 1,
				IsMain:    1,
			},
			{
				ProductID: products[i].ID,
				ImageURL:  fmt.Sprintf("https://example.com/images/product_%d_2.jpg", products[i].ID),
				SortOrder: 2,
				IsMain:    0,
			},
		}

		for j := range images {
			if err := db.Create(&images[j]).Error; err != nil {
				return err
			}
		}

		// 为每个商品创建SKU
		skus := []model.ProductSKU{
			{
				ProductID:  products[i].ID,
				SKUCode:    fmt.Sprintf("SKU_%d_001", products[i].ID),
				Name:       fmt.Sprintf("%s - 小号", products[i].Name),
				Price:      products[i].Price,
				Stock:      30,
				AttrValues: `{"size": "S", "color": "白色"}`,
				Status:     1,
			},
			{
				ProductID:  products[i].ID,
				SKUCode:    fmt.Sprintf("SKU_%d_002", products[i].ID),
				Name:       fmt.Sprintf("%s - 中号", products[i].Name),
				Price:      products[i].Price,
				Stock:      40,
				AttrValues: `{"size": "M", "color": "白色"}`,
				Status:     1,
			},
			{
				ProductID:  products[i].ID,
				SKUCode:    fmt.Sprintf("SKU_%d_003", products[i].ID),
				Name:       fmt.Sprintf("%s - 大号", products[i].Name),
				Price:      products[i].Price,
				Stock:      30,
				AttrValues: `{"size": "L", "color": "白色"}`,
				Status:     1,
			},
		}

		for j := range skus {
			if err := db.Create(&skus[j]).Error; err != nil {
				return err
			}
		}
	}

	return nil
}