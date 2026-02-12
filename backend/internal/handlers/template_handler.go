package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"prompt-backend/internal/models"
	"prompt-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TemplateHandler 模板处理器
type TemplateHandler struct {
	service *services.TemplateService
}

// NewTemplateHandler 创建模板处理器
func NewTemplateHandler(service *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

// Generate 生成提示词
func (h *TemplateHandler) Generate(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GeneratePrompt(req.TemplateID, req.Variables)
	if err != nil {
		// 如果是模板内容或解析相关错误，返回 400 并显示错误信息
		if strings.Contains(err.Error(), "invalid template content") || strings.Contains(err.Error(), "unexpected") || strings.Contains(err.Error(), "parse") {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.GenerateResponse{
		Result: result,
		Prompt: result,
	})
}

// CreateTemplate 创建模板
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var req models.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 从上下文中获取用户ID（需要认证中间件）
	userID := uuid.New() // 临时使用，实际应该从认证上下文中获取

	template, err := h.service.CreateTemplate(req, userID)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, template)
}

// GetTemplate 获取模板
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid template ID")
		return
	}

	template, err := h.service.GetTemplate(id)
	if err != nil {
		respondError(c, http.StatusNotFound, "template not found")
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetTemplates 获取模板列表
func (h *TemplateHandler) GetTemplates(c *gin.Context) {
	category := c.Query("category")
	if err := models.ValidateCategoryValue(category); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	templates, err := h.service.GetTemplates(category, page, pageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      templates,
		"page":      page,
		"page_size": pageSize,
		"total":     len(templates),
	})
}

// GetPublicTemplates 获取公开模板
func (h *TemplateHandler) GetPublicTemplates(c *gin.Context) {
	category := c.Query("category")
	if err := models.ValidateCategoryValue(category); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 20
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	templates, err := h.service.GetPublicTemplates(category, page, pageSize)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      templates,
		"page":      page,
		"page_size": pageSize,
		"total":     len(templates),
	})
}

// UpdateTemplate 更新模板
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid template ID")
		return
	}

	var req models.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	template, err := h.service.UpdateTemplate(id, req)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate 删除模板
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid template ID")
		return
	}

	if err := h.service.DeleteTemplate(id); err != nil {
		respondInternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}

// ExtractVariables 提取模板中的变量
func (h *TemplateHandler) ExtractVariables(c *gin.Context) {
	var body struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := models.ValidateTemplateContent(body.Content); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	variables := services.ExtractVariables(body.Content)
	c.JSON(http.StatusOK, gin.H{"variables": variables})
}
