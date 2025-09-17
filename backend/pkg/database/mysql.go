package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mall/pkg/config"
	appLogger "mall/pkg/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.Charset,
		cfg.Database.Timezone,
	)

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", appLogger.Error(err.Error()))
	}

	// 获取底层sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		appLogger.Fatal("Failed to get sql.DB", appLogger.Error(err.Error()))
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	appLogger.Info("Database connected successfully")
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			appLogger.Error("Failed to get sql.DB for closing", appLogger.Error(err.Error()))
			return
		}
		sqlDB.Close()
		appLogger.Info("Database connection closed")
	}
}