-- Seed data for development environment
-- Insert sample prompt templates

-- 代码解释器模板
INSERT INTO prompt_templates (id, user_id, name, description, content, category, is_public, usage_count)
VALUES (
    '6bb1d7b4-1978-4209-b223-a35e91ee56ba',
    '4d950f7b-4b9f-47a2-8c77-80438dc8cabf',
    '代码解释器',
    '解释代码的功能和逻辑',
    '请解释以下{{language}}代码的功能和逻辑：

```{{language}}
{{code}}
```

请详细说明：
1. 代码的主要功能
2. 关键逻辑和算法
3. 可能的优化建议',
    '开发',
    true,
    0
) ON CONFLICT (id) DO NOTHING;

-- 文章摘要模板
INSERT INTO prompt_templates (id, user_id, name, description, content, category, is_public, usage_count)
VALUES (
    '738ec104-5cb5-4af4-be01-b4e25abe0a10',
    'cf7d8efd-34cc-4fac-9124-86330a602911',
    '文章摘要',
    '生成文章摘要',
    '请为以下{{type}}文章生成一个{{length}}的摘要：

{{content}}

摘要要求：
- 保留核心观点
- 语言简洁明了
- 突出重点信息',
    '写作',
    true,
    0
) ON CONFLICT (id) DO NOTHING;

-- 邮件回复模板
INSERT INTO prompt_templates (id, user_id, name, description, content, category, is_public, usage_count)
VALUES (
    '8f9e2c1a-3d5b-4f7e-9a8b-1c2d3e4f5a6b',
    '5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b',
    '邮件回复',
    '生成专业的邮件回复',
    '请为以下邮件内容生成一个{{tone}}的回复：

邮件主题：{{subject}}
邮件内容：
{{content}}

请确保回复：
- 专业得体
- 语气{{tone}}
- 回复内容相关且有价值',
    '办公',
    true,
    0
) ON CONFLICT (id) DO NOTHING;