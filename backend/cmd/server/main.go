package main

import (
	"log"
	"os"

	"prompt-backend/internal/database"
	"prompt-backend/internal/handlers"
	"prompt-backend/internal/models"
	"prompt-backend/internal/services"
	"prompt-backend/internal/services/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	// 从环境变量获取配置
	config := database.GetConfigFromEnv()

	// 初始化数据库
	if err := database.Init(config); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(
		&models.PromptTemplate{},
		&models.TemplateVariable{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 如果是开发环境，插入一些示例数据
	if os.Getenv("ENV") == "development" {
		seedSampleData()
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

	// CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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

// seedSampleData 插入示例数据
func seedSampleData() {
	db := database.GetDB()

	sampleTemplates := []models.PromptTemplate{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Name:        "代码解释器",
			Description: "解释代码的功能和逻辑",
			Content:     "请解释以下{{language}}代码的功能和逻辑：\n\n```{{language}}\n{{code}}\n```\n\n请详细说明：\n1. 代码的主要功能\n2. 关键逻辑和算法\n3. 可能的优化建议",
			Category:    "开发",
			IsPublic:    true,
			UsageCount:  0,
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Name:        "文章摘要",
			Description: "生成文章摘要",
			Content:     "请为以下{{type}}文章生成一个{{length}}的摘要：\n\n{{content}}\n\n摘要要求：\n- 保留核心观点\n- 语言简洁明了\n- 突出重点信息",
			Category:    "写作",
			IsPublic:    true,
			UsageCount:  0,
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			Name:        "邮件回复",
			Description: "生成专业的邮件回复",
			Content:     "请为以下邮件内容生成一个{{tone}}的回复：\n\n邮件主题：{{subject}}\n邮件内容：\n{{content}}\n\n请确保回复：\n- 专业得体\n- 语气{{tone}}\n- 回复内容相关且有价值",
			Category:    "办公",
			IsPublic:    true,
			UsageCount:  0,
		},
	}

	for _, tmpl := range sampleTemplates {
		// 检查是否已存在
		var existing models.PromptTemplate
		if db.Where("name = ?", tmpl.Name).First(&existing).Error == nil {
			continue
		}
		if err := db.Create(&tmpl).Error; err != nil {
			log.Printf("Failed to create sample template: %v", err)
		}
	}

	log.Println("Sample data seeded successfully")
}
