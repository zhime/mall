package main

import (
	"mall/internal/app"
)

func main() {
	// 创建应用实例
	app := app.New()
	defer app.Close()

	// 启动应用
	app.Run()
}
