package main

import (
	"log"
	"os"

	"prompt-backend/internal/database"
	"prompt-backend/internal/handlers"
	"prompt-backend/internal/middleware"
	"prompt-backend/internal/services"
	"prompt-backend/internal/services/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	// 从环境变量获取配置
	config := database.GetConfigFromEnv()

	// 初始化数据库
	if err := database.Init(config); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 执行数据库迁移
	if err := database.RunMigrations(database.GetDB()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 如果是开发环境，插入一些示例数据（可选，因为迁移文件中已包含）
	if os.Getenv("ENV") == "development" {
		// 迁移文件已经包含了示例数据，这里可以留空或者添加额外的开发数据
		log.Println("Development environment detected - migrations include sample data")
	}

	// 创建仓库和服务
	db := database.GetDB()
	templateRepo := repository.NewTemplateRepository(db)
	templateService := services.NewTemplateService(templateRepo)

	// 创建处理器
	templateHandler := handlers.NewTemplateHandler(templateService)
	healthHandler := handlers.NewHealthHandler()

	// 创建 Gin 路由
	router := gin.Default()
	if err := router.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// 安全与稳健性中间件
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CORSFromEnv())
	router.Use(middleware.RequestSizeLimitFromEnv())
	router.Use(middleware.RateLimitFromEnv())

	// 路由组
	api := router.Group("/api")
	{
		// 健康检查
		api.GET("/health", healthHandler.Check)

		// 模板相关路由
		templates := api.Group("/templates")
		{
			templates.POST("", templateHandler.CreateTemplate)
			templates.GET("", templateHandler.GetTemplates)
			templates.GET("/public", templateHandler.GetPublicTemplates)
			templates.GET("/:id", templateHandler.GetTemplate)
			templates.PUT("/:id", templateHandler.UpdateTemplate)
			templates.DELETE("/:id", templateHandler.DeleteTemplate)
		}

		// 生成相关路由
		generate := api.Group("/generate")
		{
			generate.POST("", templateHandler.Generate)
			generate.POST("/extract-variables", templateHandler.ExtractVariables)
		}
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
