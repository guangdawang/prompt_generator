package repository

import (
	"prompt-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TemplateRepository 模板仓库接口
type TemplateRepository interface {
	Create(template *models.PromptTemplate) error
	GetByID(id uuid.UUID) (*models.PromptTemplate, error)
	GetAll(category string, limit, offset int) ([]models.PromptTemplate, error)
	GetByUserID(userID uuid.UUID, category string, limit, offset int) ([]models.PromptTemplate, error)
	Update(template *models.PromptTemplate) error
	Delete(id uuid.UUID) error
	IncrementUsage(id uuid.UUID) error
	GetPublicTemplates(category string, limit, offset int) ([]models.PromptTemplate, error)
}

// templateRepository 模板仓库实现
type templateRepository struct {
	db *gorm.DB
}

// NewTemplateRepository 创建模板仓库
func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{db: db}
}

// Create 创建模板
func (r *templateRepository) Create(template *models.PromptTemplate) error {
	return r.db.Create(template).Error
}

// GetByID 根据ID获取模板
func (r *templateRepository) GetByID(id uuid.UUID) (*models.PromptTemplate, error) {
	var template models.PromptTemplate
	err := r.db.Where("id = ?", id).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetAll 获取所有模板
func (r *templateRepository) GetAll(category string, limit, offset int) ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	query := r.db
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&templates).Error
	return templates, err
}

// GetByUserID 获取用户模板
func (r *templateRepository) GetByUserID(userID uuid.UUID, category string, limit, offset int) ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	query := r.db.Where("user_id = ?", userID)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&templates).Error
	return templates, err
}

// Update 更新模板
func (r *templateRepository) Update(template *models.PromptTemplate) error {
	return r.db.Save(template).Error
}

// Delete 删除模板
func (r *templateRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.PromptTemplate{}, "id = ?", id).Error
}

// IncrementUsage 增加使用次数
func (r *templateRepository) IncrementUsage(id uuid.UUID) error {
	return r.db.Model(&models.PromptTemplate{}).Where("id = ?", id).UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

// GetPublicTemplates 获取公开模板
func (r *templateRepository) GetPublicTemplates(category string, limit, offset int) ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	query := r.db.Where("is_public = ?", true)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&templates).Error
	return templates, err
}
