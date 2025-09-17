package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`

	Database struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		Username        string `mapstructure:"username"`
		Password        string `mapstructure:"password"`
		DBName          string `mapstructure:"dbname"`
		Charset         string `mapstructure:"charset"`
		Timezone        string `mapstructure:"timezone"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"database"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
		PoolSize int    `mapstructure:"pool_size"`
	} `mapstructure:"redis"`

	Elasticsearch struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"elasticsearch"`

	JWT struct {
		Secret      string `mapstructure:"secret"`
		ExpireHours int    `mapstructure:"expire_hours"`
	} `mapstructure:"jwt"`

	Log struct {
		Level      string `mapstructure:"level"`
		Filename   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		MaxBackups int    `mapstructure:"max_backups"`
	} `mapstructure:"log"`

	Upload struct {
		Path         string   `mapstructure:"path"`
		MaxSize      int      `mapstructure:"max_size"`
		AllowedTypes []string `mapstructure:"allowed_types"`
	} `mapstructure:"upload"`

	Wechat struct {
		AppID     string `mapstructure:"app_id"`
		AppSecret string `mapstructure:"app_secret"`
		MchID     string `mapstructure:"mch_id"`
		APIKey    string `mapstructure:"api_key"`
	} `mapstructure:"wechat"`

	Alipay struct {
		AppID           string `mapstructure:"app_id"`
		PrivateKey      string `mapstructure:"private_key"`
		AlipayPublicKey string `mapstructure:"alipay_public_key"`
	} `mapstructure:"alipay"`
}

var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig() *Config {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	// 设置环境变量读取
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	GlobalConfig = &config
	return &config
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if GlobalConfig == nil {
		return LoadConfig()
	}
	return GlobalConfig
}