# 提示词模板系统

一个功能完整的提示词模板生成和管理系统，使用 Go 后端和 Next.js 前端构建。

## 功能特性

- 📝 模板管理：创建、编辑、删除提示词模板
- 🎯 变量支持：灵活的变量定义和替换
- 🌐 公开/私有模板：支持模板共享
- 🔍 模板搜索：快速查找所需模板
- 📊 使用统计：追踪模板使用次数
- 🎨 现代化 UI：基于 Tailwind CSS 的美观界面

## 技术栈

### 后端

- Go 1.21
- Gin Web 框架
- GORM (PostgreSQL)
- UUID 生成

### 前端

- Next.js 14 (App Router)
- React 18
- TypeScript
- Tailwind CSS
- Axios

### 数据库

- PostgreSQL 15

## 快速开始

### 使用 Docker Compose（推荐）

1. 克隆项目并进入目录：

    ```bash
    cd prompt_generator
    ```

1. 配置环境变量（可选但推荐）：

    ```bash
    cp backend/.env.example backend/.env
    cp frontend/.env.example frontend/.env
    # 修改 APP_DOMAIN 等配置
    ```

1. 启用 HTTPS（Let's Encrypt）：

    确保域名已解析到服务器公网 IP，然后在根目录创建或更新 `.env`：

    ```bash
    APP_DOMAIN=your-domain.com
    CERTBOT_EMAIL=your-email@example.com
    ```

    启动服务并申请证书：

    ```bash
    docker-compose up -d nginx backend frontend
    docker-compose up -d certbot
    docker-compose restart nginx
    ```

1. 启动所有服务：

    ```bash
    docker-compose up -d
    ```

1. 访问应用：

    默认 NGINX 不映射端口，仅供容器内部访问。如果需要从宿主机访问，请在 NGINX 服务中添加端口映射（如 `80:80`），然后通过域名访问：

    - 前端: <https://your-domain>
    - 后端 API: <https://your-domain/api>

### 手动启动

#### 后端（手动启动）

1. 进入后端目录：

    ```bash
    cd backend
    ```

1. 配置环境变量：

    ```bash
    cp .env.example .env
    # 编辑 .env 文件，设置数据库连接信息
    ```

1. 安装依赖：

    ```bash
    go mod download
    ```

1. 运行服务：

    ```bash
    go run cmd/server/main.go
    ```

#### 前端（手动启动）

1. 进入前端目录：

    ```bash
    cd frontend
    ```

1. 安装依赖：

    ```bash
    npm install
    ```

1. 运行开发服务器：

    ```bash
    npm run dev
    ```

1. 访问 <http://localhost:3000>

## API 端点

### 健康检查

- `GET /api/health` - 检查服务状态

### 模板管理

- `GET /api/templates` - 获取模板列表（支持分类筛选）
- `GET /api/templates/public` - 获取公开模板
- `GET /api/templates/:id` - 获取单个模板
- `POST /api/templates` - 创建新模板
- `PUT /api/templates/:id` - 更新模板
- `DELETE /api/templates/:id` - 删除模板

### 提示词生成

- `POST /api/generate` - 生成提示词
- `POST /api/generate/extract-variables` - 从模板内容提取变量

## 数据库结构

### prompt_templates

- `id` (UUID) - 主键
- `user_id` (UUID) - 用户ID
- `name` (VARCHAR) - 模板名称
- `description` (TEXT) - 模板描述
- `content` (TEXT) - 模板内容
- `variables` (JSONB) - 变量定义
- `category` (VARCHAR) - 分类
- `is_public` (BOOLEAN) - 是否公开
- `usage_count` (INTEGER) - 使用次数
- `created_at` (TIMESTAMP) - 创建时间
- `updated_at` (TIMESTAMP) - 更新时间

### template_variables

- `id` (UUID) - 主键
- `template_id` (UUID) - 模板ID
- `name` (VARCHAR) - 变量名
- `display_name` (VARCHAR) - 显示名称
- `description` (TEXT) - 描述
- `default_value` (TEXT) - 默认值
- `required` (BOOLEAN) - 是否必填
- `sort_order` (INTEGER) - 排序

## 示例模板

系统预置了几个示例模板：

1. **代码解释器** - 解释代码的功能和逻辑
2. **文章摘要** - 生成文章摘要
3. **邮件回复** - 生成专业的邮件回复

## 开发

### 后端开发

```bash
cd backend
# 运行测试
go test ./...

# 格式化代码
go fmt ./...

# 代码检查
go vet ./...
```

### 前端开发

```bash
cd frontend
# 运行开发服务器
npm run dev

# 构建生产版本
npm run build

# 类型检查
npm run lint
```

## 部署

### 生产环境配置

1. 修改 `.env` 文件，设置生产环境变量
2. 使用 `docker-compose.prod.yml` 进行部署
3. 配置 HTTPS（建议使用 Nginx 反向代理）
4. 设置数据库备份策略

### 性能优化

- 启用数据库连接池
- 使用 CDN 加速静态资源
- 配置 Redis 缓存
- 启用 Gzip 压缩

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
