package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"mall/internal/handler"
	"mall/internal/repository"
	"mall/internal/routes"
	"mall/internal/service"
	"mall/pkg/cache"
	"mall/pkg/config"
	"mall/pkg/database"
	"mall/pkg/logger"
)

// App 应用结构体
type App struct {
	config   *config.Config
	server   *http.Server
	handlers *Handlers
}

// Handlers 处理器容器
type Handlers struct {
	AuthHandler     *handler.AuthHandler
	CategoryHandler *handler.CategoryHandler
	ProductHandler  *handler.ProductHandler
	CartHandler     *handler.CartHandler
	OrderHandler    *handler.OrderHandler
	PaymentHandler  *handler.PaymentHandler
}

// New 创建新的应用实例
func New() *App {
	app := &App{
		config: config.LoadConfig(),
	}
	
	app.initLogger()
	app.initDatabase()
	app.initRedis()
	app.initDependencies()
	
	return app
}

// initLogger 初始化日志
func (a *App) initLogger() {
	logger.InitLogger(
		a.config.Log.Level,
		a.config.Log.Filename,
		a.config.Log.MaxSize,
		a.config.Log.MaxAge,
		a.config.Log.MaxBackups,
	)
}

// initDatabase 初始化数据库
func (a *App) initDatabase() {
	database.InitDB()
}

// initRedis 初始化Redis
func (a *App) initRedis() {
	cache.InitRedis()
}

// initDependencies 初始化依赖
func (a *App) initDependencies() {
	db := database.GetDB()
	
	// 初始化仓储
	userRepo := repository.NewUserRepository(db)
	userAuthRepo := repository.NewUserAuthRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	productSKURepo := repository.NewProductSKURepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderItemRepo := repository.NewOrderItemRepository(db)
	paymentRepo := repository.NewOrderPaymentRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, userAuthRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo, productSKURepo)
	cartService := service.NewCartService(cartRepo, productRepo, productSKURepo)
	orderService := service.NewOrderService(orderRepo, orderItemRepo, productRepo, productSKURepo, cartRepo, userRepo)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo)

	// 初始化处理器
	a.handlers = &Handlers{
		AuthHandler:     handler.NewAuthHandler(authService),
		CategoryHandler: handler.NewCategoryHandler(categoryService),
		ProductHandler:  handler.NewProductHandler(productService),
		CartHandler:     handler.NewCartHandler(cartService),
		OrderHandler:    handler.NewOrderHandler(orderService),
		PaymentHandler:  handler.NewPaymentHandler(paymentService),
	}
}

// Run 启动应用
func (a *App) Run() error {
	// 设置Gin模式
	gin.SetMode(a.config.Server.Mode)

	// 创建路由
	router := routes.SetupRouter(a.handlers.AuthHandler, a.handlers.CategoryHandler, a.handlers.ProductHandler, a.handlers.CartHandler, a.handlers.OrderHandler, a.handlers.PaymentHandler)

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器
	go func() {
		logger.Info(fmt.Sprintf("Server starting on %s:%d", a.config.Server.Host, a.config.Server.Port))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server")
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")

	// 优雅关闭服务器，等待5秒钟
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown")
	}

	logger.Info("Server exited")
	return nil
}

// Close 关闭应用资源
func (a *App) Close() {
	database.CloseDB()
	cache.CloseRedis()
}