# Mall 电商系统

一个基于 Go + Vue3 的现代化电商系统，包含管理后台、商城前端、微信小程序等完整功能模块。

## 项目结构

```
mall/
├── backend/                    # Go 后端服务
│   ├── cmd/api/               # 应用入口
│   ├── internal/              # 内部代码
│   │   ├── handler/           # HTTP 处理器
│   │   ├── service/           # 业务逻辑
│   │   ├── repository/        # 数据访问层
│   │   ├── model/             # 数据模型
│   │   └── middleware/        # 中间件
│   ├── pkg/                   # 公共包
│   ├── config/                # 配置文件
│   └── migrations/            # 数据库迁移
├── admin-web/                 # Vue3 管理后台
│   ├── src/
│   │   ├── views/             # 页面组件
│   │   ├── components/        # 公共组件
│   │   ├── api/              # API 接口
│   │   ├── store/            # 状态管理
│   │   └── utils/            # 工具函数
│   └── public/
├── mall-web/                  # Vue3 商城前端
│   ├── src/
│   │   ├── views/             # 页面组件
│   │   ├── components/        # 公共组件
│   │   ├── api/              # API 接口
│   │   ├── store/            # 状态管理
│   │   └── utils/            # 工具函数
│   └── public/
└── wechat-miniprogram/        # 微信小程序
    ├── pages/
    ├── components/
    └── utils/
```

## 技术栈

### 后端
- **语言**: Go 1.23.10
- **框架**: Gin
- **数据库**: MySQL 8.0
- **ORM**: GORM
- **缓存**: Redis
- **认证**: JWT
- **支付**: 微信支付 + 支付宝

### 前端
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件**: Element Plus (管理后台) + Vant (移动端)
- **状态管理**: Pinia
- **样式**: SCSS + WindiCSS

### 数据库设计
- **用户模块**: 用户信息、认证、权限管理
- **商品模块**: 商品、分类、规格、属性管理
- **订单模块**: 订单、订单项、支付、物流管理
- **营销模块**: 优惠券、活动、推荐等

## 功能特性

### 管理后台
- ✅ 用户管理 (列表、搜索、状态切换、密码重置)
- ✅ 商品管理 (分类管理、商品列表、添加编辑、规格属性)
- ✅ 订单管理 (订单列表、详情查看、状态处理、发货退款)
- ✅ 权限控制 (JWT 认证、路由守卫)
- ✅ 数据统计 (订单统计、销售分析)

### 商城前端
- ✅ 首页展示 (轮播图、分类导航、商品推荐)
- ✅ 商品浏览 (分类筛选、搜索、详情查看)
- ✅ 购物流程 (加入购物车、规格选择、下单结算)
- ✅ 用户中心 (个人信息、订单管理、收货地址)
- ✅ 移动适配 (响应式设计、手势操作)

### 后端 API
- ✅ 用户认证 (手机号登录、微信登录、JWT)
- ✅ 商品系统 (商品 CRUD、分类管理、搜索推荐)
- ✅ 订单系统 (订单创建、状态流转、支付集成)
- ✅ 支付系统 (微信支付、支付宝、回调处理)
- ✅ 缓存策略 (Redis 缓存、数据一致性)

## 快速开始

### 环境要求
- Go 1.23.10+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 后端启动
```bash
cd backend
go mod tidy
go run cmd/api/main.go
```

### 管理后台启动
```bash
cd admin-web
npm install
npm run dev
```

### 商城前端启动
```bash
cd mall-web
npm install
npm run dev
```

## 部署配置

### Docker 部署
```bash
# 构建后端镜像
docker build -t mall-backend ./backend

# 构建前端镜像
docker build -t mall-admin-web ./admin-web
docker build -t mall-web ./mall-web

# 使用 docker-compose 启动
docker-compose up -d
```

### 生产环境配置
- 数据库连接池优化
- Redis 集群配置
- Nginx 反向代理
- HTTPS 证书配置
- 日志收集和监控

## API 文档

### 用户相关
- `POST /api/user/login` - 用户登录
- `POST /api/user/register` - 用户注册
- `GET /api/user/profile` - 获取用户信息
- `PUT /api/user/profile` - 更新用户信息

### 商品相关
- `GET /api/products` - 获取商品列表
- `GET /api/products/:id` - 获取商品详情
- `GET /api/categories` - 获取分类列表
- `GET /api/products/search` - 商品搜索

### 订单相关
- `POST /api/orders` - 创建订单
- `GET /api/orders` - 获取订单列表
- `GET /api/orders/:id` - 获取订单详情
- `PUT /api/orders/:id/status` - 更新订单状态

## 开发说明

### 代码规范
- Go 代码遵循 gofmt 格式化
- TypeScript 使用 ESLint + Prettier
- Git 提交使用 Conventional Commits

### 测试
- 后端单元测试覆盖率 > 80%
- 前端组件测试
- API 集成测试
- E2E 自动化测试

### 性能优化
- 数据库索引优化
- Redis 缓存策略
- 前端代码分割
- 图片懒加载
- CDN 资源优化

## 许可证

MIT License

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 联系方式

- 作者: Mall Team
- 邮箱: contact@mall.com
- 文档: https://docs.mall.com