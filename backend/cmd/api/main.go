package main

import (
	"mall/internal/app"
)

func main() {
	// 创建应用实例
	service := app.New()
	defer service.Close()

	// 启动应用
	_ = service.Run()
}
