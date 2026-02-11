package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONB 类型用于处理 PostgreSQL 的 JSONB 类型
type JSONB json.RawMessage

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// PromptTemplate 提示词模板
type PromptTemplate struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Variables   JSONB     `gorm:"type:jsonb;default:'[]'" json:"variables"`
	Category    string    `gorm:"size:100;index" json:"category"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	UsageCount  int       `gorm:"default:0" json:"usage_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (PromptTemplate) TableName() string {
	return "prompt_templates"
}

// TemplateVariable 模板变量
type TemplateVariable struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	TemplateID   uuid.UUID `gorm:"type:uuid;index" json:"template_id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	DisplayName  string    `gorm:"size:100;not null" json:"display_name"`
	Description  string    `gorm:"type:text" json:"description"`
	DefaultValue string    `gorm:"type:text" json:"default_value"`
	Required     bool      `gorm:"default:true" json:"required"`
	SortOrder    int       `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (TemplateVariable) TableName() string {
	return "template_variables"
}

// GenerateRequest 生成提示词请求
type GenerateRequest struct {
	TemplateID uuid.UUID         `json:"template_id" binding:"required"`
	Variables  map[string]string `json:"variables" binding:"required"`
}

// GenerateResponse 生成提示词响应
type GenerateResponse struct {
	Result string `json:"result"`
	Prompt string `json:"prompt"`
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Content     string            `json:"content" binding:"required"`
	Variables   []TemplateVariable `json:"variables"`
	Category    string            `json:"category"`
	IsPublic    bool              `json:"is_public"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name        *string            `json:"name"`
	Description *string            `json:"description"`
	Content     *string            `json:"content"`
	Variables   []TemplateVariable `json:"variables"`
	Category    *string            `json:"category"`
	IsPublic    *bool              `json:"is_public"`
}
