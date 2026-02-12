# Prompt Generator

一个用于管理和使用提示模板（prompt templates）的开源示例项目，包含 Go 后端和 Next.js 前端。该仓库演示了模板的创建、变量输入和生成提示的工作流，适合作为构建提示/模板管理工具的起点。

## 主要特性

- 模板管理：创建、编辑、选择提示模板。
- 可配置变量输入：为模板提供动态变量以生成最终提示。
- 前后端分离：Go 后端（API、数据库迁移）、Next.js 前端（UI）。
- Docker 支持：可通过 `docker-compose` 一键运行整个应用。

## 仓库结构（概要）

- `backend/` — Go 后端源码，包含数据库、服务、处理器。
- `frontend/` — Next.js 前端源码和组件。
- `migrations/` — 数据库迁移和种子数据。
- `docker-compose.yml` — 用于本地快速启动（包含后端、前端、数据库）。

## 快速开始（使用 Docker Compose）

开发机器上只需 Docker 与 Docker Compose：

```bash
docker-compose up --build
```

服务启动后：

- 前端（Next.js）通常在 `http://localhost:3000`
- 后端 API 通常在 `http://localhost:8080`（参见 `docker-compose.yml`）

## 本地开发（不使用 Docker）

后端（Go）：

```bash
cd backend
go mod download
go run ./cmd/server
```

前端（Next.js）：

```bash
cd frontend
npm install
npm run dev
```

> 注意：后端会连接数据库（请确保环境变量正确或启动本地/Postgres 实例）。

## 数据库与迁移

迁移 SQL 存放在 `migrations/`，初始化或重建数据库时请运行这些脚本。后端内部也包含与数据库初始化/迁移逻辑（见 `backend/internal/database`）。

## 开发提示

- 后端使用模块化结构：`services`、`repository`、`handlers`，便于扩展。
- 前端组件集中在 `frontend/components`，可以快速复用或替换 UI。

## 贡献

欢迎提交 issue 或 PR。请在变更中保持风格一致、编写清晰的提交说明，并在必要时更新或添加迁移脚本。

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE)。
