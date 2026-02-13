package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxTemplateNameLen        = 200
	MaxTemplateDescriptionLen = 2000
	MaxTemplateContentLen     = 10000
	MaxCategoryLen            = 100
	MaxVariables              = 50
	MaxVariableNameLen        = 100
	MaxVariableDisplayNameLen = 100
	MaxVariableDescriptionLen = 2000
	MaxVariableValueLen       = 1000
)

var variableNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// JSONB 类型用于处理 PostgreSQL 的 JSONB 类型
type JSONB json.RawMessage

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	// 返回原始 JSON 字节，不要再次 Marshal，避免双重编码
	return []byte(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = JSONB(append([]byte(nil), v...))
		return nil
	case string:
		*j = JSONB([]byte(v))
		return nil
	default:
		// Fallback: try to marshal the value to JSON
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		*j = JSONB(b)
		return nil
	}
}

// MarshalJSON 保证 JSONB 在编码为 JSON 时以原始 JSON 字节输出（例如数组/对象），
// 而不是作为 base64 编码的字符串。
func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

// UnmarshalJSON 用于将原始 JSON 解码到 JSONB
func (j *JSONB) UnmarshalJSON(b []byte) error {
	if j == nil {
		return errors.New("nil JSONB")
	}
	// 直接复制字节
	*j = append((*j)[0:0], b...)
	return nil
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

func (r *GenerateRequest) Validate() error {
	if r.TemplateID == uuid.Nil {
		return errors.New("template_id is required")
	}
	if len(r.Variables) > MaxVariables {
		return fmt.Errorf("too many variables (max %d)", MaxVariables)
	}
	for name, value := range r.Variables {
		if err := validateVariableKey(name); err != nil {
			return err
		}
		if len(value) > MaxVariableValueLen {
			return fmt.Errorf("variable value too long (max %d)", MaxVariableValueLen)
		}
	}
	return nil
}

// GenerateResponse 生成提示词响应
type GenerateResponse struct {
	Result string `json:"result"`
	Prompt string `json:"prompt"`
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	Content     string             `json:"content" binding:"required"`
	Variables   []TemplateVariable `json:"variables"`
	Category    string             `json:"category"`
	IsPublic    bool               `json:"is_public"`
}

func (r *CreateTemplateRequest) Validate() error {
	name := strings.TrimSpace(r.Name)
	if name == "" {
		return errors.New("name is required")
	}
	if len(name) > MaxTemplateNameLen {
		return fmt.Errorf("name too long (max %d)", MaxTemplateNameLen)
	}
	if len(r.Description) > MaxTemplateDescriptionLen {
		return fmt.Errorf("description too long (max %d)", MaxTemplateDescriptionLen)
	}
	if err := ValidateTemplateContent(r.Content); err != nil {
		return err
	}
	if err := ValidateCategoryValue(r.Category); err != nil {
		return err
	}
	if len(r.Variables) > MaxVariables {
		return fmt.Errorf("too many variables (max %d)", MaxVariables)
	}
	for _, variable := range r.Variables {
		if err := validateTemplateVariable(variable); err != nil {
			return err
		}
	}
	return nil
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

func (r *UpdateTemplateRequest) Validate() error {
	if r.Name != nil {
		name := strings.TrimSpace(*r.Name)
		if name == "" {
			return errors.New("name cannot be empty")
		}
		if len(name) > MaxTemplateNameLen {
			return fmt.Errorf("name too long (max %d)", MaxTemplateNameLen)
		}
	}
	if r.Description != nil && len(*r.Description) > MaxTemplateDescriptionLen {
		return fmt.Errorf("description too long (max %d)", MaxTemplateDescriptionLen)
	}
	if r.Content != nil {
		if err := ValidateTemplateContent(*r.Content); err != nil {
			return err
		}
	}
	if r.Category != nil {
		if err := ValidateCategoryValue(*r.Category); err != nil {
			return err
		}
	}
	if r.Variables != nil {
		if len(r.Variables) > MaxVariables {
			return fmt.Errorf("too many variables (max %d)", MaxVariables)
		}
		for _, variable := range r.Variables {
			if err := validateTemplateVariable(variable); err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateTemplateContent(content string) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return errors.New("content is required")
	}
	if len(trimmed) > MaxTemplateContentLen {
		return fmt.Errorf("content too long (max %d)", MaxTemplateContentLen)
	}
	return nil
}

func ValidateCategoryValue(category string) error {
	if category == "" {
		return nil
	}
	if len(category) > MaxCategoryLen {
		return fmt.Errorf("category too long (max %d)", MaxCategoryLen)
	}
	return nil
}

func validateTemplateVariable(variable TemplateVariable) error {
	name := strings.TrimSpace(variable.Name)
	if name == "" {
		return errors.New("variable name is required")
	}
	if len(name) > MaxVariableNameLen {
		return fmt.Errorf("variable name too long (max %d)", MaxVariableNameLen)
	}
	if !variableNamePattern.MatchString(name) {
		return errors.New("invalid variable name format")
	}
	if len(variable.DisplayName) > MaxVariableDisplayNameLen {
		return fmt.Errorf("variable display_name too long (max %d)", MaxVariableDisplayNameLen)
	}
	if len(variable.Description) > MaxVariableDescriptionLen {
		return fmt.Errorf("variable description too long (max %d)", MaxVariableDescriptionLen)
	}
	if len(variable.DefaultValue) > MaxVariableValueLen {
		return fmt.Errorf("variable default_value too long (max %d)", MaxVariableValueLen)
	}
	return nil
}

func validateVariableKey(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("variable name is required")
	}
	if len(name) > MaxVariableNameLen {
		return fmt.Errorf("variable name too long (max %d)", MaxVariableNameLen)
	}
	if !variableNamePattern.MatchString(name) {
		return errors.New("invalid variable name format")
	}
	return nil
}
