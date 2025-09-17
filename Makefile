# Makefile for Mall Project

.PHONY: help dev build test clean docker

# 默认目标
help:
	@echo "Mall Project Makefile"
	@echo ""
	@echo "Available commands:"
	@echo "  dev        - 启动开发环境"
	@echo "  build      - 构建所有项目"
	@echo "  test       - 运行测试"
	@echo "  clean      - 清理构建文件"
	@echo "  docker     - 使用Docker启动服务"
	@echo "  docker-down - 停止Docker服务"

# 启动开发环境
dev:
	@echo "启动开发环境..."
	@echo "1. 启动后端服务..."
	cd backend && go run cmd/api/main.go &
	@echo "2. 启动管理后台..."
	cd admin-web && npm run dev &
	@echo "3. 启动商城前端..."
	cd mall-web && npm run dev &

# 构建所有项目
build:
	@echo "构建项目..."
	@echo "1. 构建后端..."
	cd backend && go build -o bin/api cmd/api/main.go
	@echo "2. 构建管理后台..."
	cd admin-web && npm run build
	@echo "3. 构建商城前端..."
	cd mall-web && npm run build

# 运行测试
test:
	@echo "运行测试..."
	cd backend && go test ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf backend/bin
	rm -rf admin-web/dist
	rm -rf mall-web/dist

# 使用Docker启动服务
docker:
	@echo "使用Docker启动服务..."
	docker-compose up -d

# 停止Docker服务
docker-down:
	@echo "停止Docker服务..."
	docker-compose down

# 初始化项目依赖
install:
	@echo "安装项目依赖..."
	@echo "1. 初始化Go模块..."
	cd backend && go mod init mall && go mod tidy
	@echo "2. 安装管理后台依赖..."
	cd admin-web && npm install
	@echo "3. 安装商城前端依赖..."
	cd mall-web && npm install

# 数据库迁移
migrate:
	@echo "执行数据库迁移..."
	cd backend && go run cmd/migrate/main.go

# 生成API文档
docs:
	@echo "生成API文档..."
	cd backend && swag init -g cmd/api/main.go -o docs