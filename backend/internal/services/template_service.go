package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"prompt-backend/internal/models"
	"prompt-backend/internal/services/repository"

	"github.com/google/uuid"
)

// TemplateService 模板服务
type TemplateService struct {
	repo repository.TemplateRepository
}

// NewTemplateService 创建模板服务
func NewTemplateService(repo repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

// GeneratePrompt 生成提示词
func (s *TemplateService) GeneratePrompt(templateID uuid.UUID, variables map[string]string) (string, error) {
	// 获取模板
	tmpl, err := s.repo.GetByID(templateID)
	if err != nil {
		// 明确返回未找到的友好错误，便于 handler 返回 404 而不是 500
		if strings.Contains(strings.ToLower(err.Error()), "record not found") {
			return "", fmt.Errorf("template not found")
		}
		return "", err
	}

	// 解析模板
	normalizedContent := normalizeTemplateContent(tmpl.Content)
	// 简单检查模板占位符对是否匹配，避免 text/template 解析时出现未捕获的错误
	if strings.Count(normalizedContent, "{{") != strings.Count(normalizedContent, "}}") {
		return "", fmt.Errorf("invalid template content: unbalanced braces")
	}
	t, err := template.New(tmpl.Name).Parse(normalizedContent)
	if err != nil {
		return "", err
	}

	// 执行模板替换
	var buf bytes.Buffer
	if err := t.Execute(&buf, variables); err != nil {
		return "", err
	}

	// 异步更新使用次数
	go func() {
		if err := s.repo.IncrementUsage(templateID); err != nil {
			// 记录错误但不影响主流程
			// 在实际项目中应该使用日志库
		}
	}()

	return buf.String(), nil
}

func normalizeTemplateContent(content string) string {
	// 匹配 {{...}} 内的内容（尽量宽松），然后根据内容是否为合法 ASCII 变量名分别处理：
	// - 若为 ASCII 变量名：转换为 {{.name}} 来从传入 map 中读取值
	// - 若包含非 ASCII 字符（比如中文）：将其作为字面量文本输出，使用 printf 转义以保证模板解析安全
	re := regexp.MustCompile(`\{\{\s*([^\}\s]+)\s*\}\}`)
	asciiIdent := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	return re.ReplaceAllStringFunc(content, func(m string) string {
		sub := re.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		name := sub[1]
		if asciiIdent.MatchString(name) {
			return "{{." + name + "}}"
		}
		// 对非 ASCII 内容，作为字面量返回（安全地加引号）
		return fmt.Sprintf("{{printf %s}}", strconv.Quote(name))
	})
}

// CreateTemplate 创建模板
func (s *TemplateService) CreateTemplate(req models.CreateTemplateRequest, userID uuid.UUID) (*models.PromptTemplate, error) {
	template := &models.PromptTemplate{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
		Category:    req.Category,
		IsPublic:    req.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 将变量列表转换为 JSONB
	variablesData, err := s.marshalVariables(req.Variables)
	if err != nil {
		return nil, err
	}
	template.Variables = variablesData

	if err := s.repo.Create(template); err != nil {
		return nil, err
	}

	return template, nil
}

// GetTemplate 获取模板
func (s *TemplateService) GetTemplate(id uuid.UUID) (*models.PromptTemplate, error) {
	return s.repo.GetByID(id)
}

// GetTemplates 获取模板列表
func (s *TemplateService) GetTemplates(category string, page, pageSize int) ([]models.PromptTemplate, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return s.repo.GetAll(category, limit, offset)
}

// GetPublicTemplates 获取公开模板
func (s *TemplateService) GetPublicTemplates(category string, page, pageSize int) ([]models.PromptTemplate, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return s.repo.GetPublicTemplates(category, limit, offset)
}

// UpdateTemplate 更新模板
func (s *TemplateService) UpdateTemplate(id uuid.UUID, req models.UpdateTemplateRequest) (*models.PromptTemplate, error) {
	tmpl, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != nil {
		tmpl.Name = *req.Name
	}
	if req.Description != nil {
		tmpl.Description = *req.Description
	}
	if req.Content != nil {
		tmpl.Content = *req.Content
	}
	if req.Category != nil {
		tmpl.Category = *req.Category
	}
	if req.IsPublic != nil {
		tmpl.IsPublic = *req.IsPublic
	}
	if req.Variables != nil {
		variablesData, err := s.marshalVariables(req.Variables)
		if err != nil {
			return nil, err
		}
		tmpl.Variables = variablesData
	}
	tmpl.UpdatedAt = time.Now()

	if err := s.repo.Update(tmpl); err != nil {
		return nil, err
	}

	return tmpl, nil
}

// DeleteTemplate 删除模板
func (s *TemplateService) DeleteTemplate(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// ExtractVariables 从模板内容中提取变量
func ExtractVariables(content string) []string {
	// 使用正则表达式提取 {{variable}} 格式的变量
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	variables := make([]string, 0)
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			variables = append(variables, match[1])
			seen[match[1]] = true
		}
	}

	return variables
}

// marshalVariables 将变量列表转换为 JSONB
func (s *TemplateService) marshalVariables(variables []models.TemplateVariable) (models.JSONB, error) {
	// 简化处理，直接使用 JSON 编码
	// 在实际项目中应该使用更好的方式处理
	data, err := json.Marshal(variables)
	if err != nil {
		return nil, err
	}
	return models.JSONB(data), nil
}
