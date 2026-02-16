# AI Interview Platform

LLM + TTS 驱动的 AI 模拟面试平台。用户自带 API Key (BYOK)，支持多家 LLM / TTS 服务商，实时语音对话面试。

## 技术栈

| 层      | 技术                                 |
| ------- | ------------------------------------ |
| 后端    | Go · Kratos v2 · sqlc · Wire        |
| 前端    | Vue 3 · TypeScript · Vite · Pinia   |
| 数据库  | MySQL 8 · Redis 7                    |
| 部署    | Docker / Podman Compose · GitHub Actions · GHCR |

**支持的服务商**

- **LLM**: OpenAI · Anthropic · Google Gemini · DeepSeek
- **TTS**: OpenAI TTS · Fish Audio · ElevenLabs · Edge TTS (免费)
- **STT**: Whisper (OpenAI)

## 快速开始

### 前置条件

- Go 1.23+、Node 22+、Make
- Docker / Podman（运行 MySQL & Redis）

### 1. 克隆 & 配置环境变量

```bash
git clone https://github.com/zjy-dev/ai-interview.git
cd ai-interview
cp .env.example .env
```

编辑 `.env`，至少修改以下两项：

```dotenv
JWT_SECRET=<替换为一段随机字符串>
ENCRYPTION_KEY=<64 位十六进制字符串，即 32 字节>
```

> `.env` 中的 `DB_DSN` / `REDIS_ADDR` 默认值指向 `127.0.0.1`（适合本地开发）。
> 容器部署时 `docker-compose.yml` 的 `environment` 会自动覆盖这两项，指向容器内部的 `mysql` / `redis` 主机名。

---

### 方式一：本地开发

适合需要调试后端/前端代码的场景。MySQL 和 Redis 跑在容器里，应用跑在宿主机。

```bash
# 1) 启动 MySQL + Redis（仅基础设施，不构建应用镜像）
docker compose -f docker-compose.dev.yml up -d
# 或使用 podman:
podman compose -f docker-compose.dev.yml up -d

# 2) 安装工具链 & 生成代码
make init          # 安装 protoc-gen-go, wire, sqlc 等
make sqlc           # 从 SQL 生成 Go 数据访问代码
make wire           # 生成依赖注入代码
make build          # 编译后端二进制

# 3) 安装前端依赖
cd frontend && npm ci && cd ..

# 4) 运行后端（终端 1）
make dev            # 监听 :8080

# 5) 运行前端（终端 2）
cd frontend && npm run dev   # 监听 :5173，API 代理到 :8080
```

打开浏览器访问 **http://localhost:5173**

> **Podman 用户注意**：Rootless Podman 的端口映射可能因 nftables 规则不生效，导致 `make dev` 连不上 `127.0.0.1:3306`。
> 解决方法：用 `podman inspect ai-interview-mysql --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'` 获取容器 IP，
> 然后修改 `.env` 中 `DB_DSN` 和 `REDIS_ADDR` 为该 IP。

---

### 方式二：一键容器部署

适合演示或生产部署。四个服务全部运行在容器中。

```bash
# 确保 .env 已创建（见上方步骤 1）

# Docker Compose:
docker compose up -d --build

# 或 Podman Compose:
podman compose up -d --build
```

这会启动 4 个容器：

| 容器 | 端口 | 说明 |
|------|------|------|
| `ai-interview-frontend` | **:80** → 80 | Nginx 托管 Vue SPA，反代 `/api/` 到后端 |
| `ai-interview-backend` | **:8080** → 8080 | Go Kratos HTTP 服务 |
| `ai-interview-mysql` | **:3306** → 3306 | MySQL 8.0 |
| `ai-interview-redis` | **:6379** → 6379 | Redis 7 |

打开浏览器访问 **http://localhost**

停止所有服务：

```bash
docker compose down           # 保留数据卷
docker compose down -v        # 同时删除数据卷（清空数据库）
```

## 测试

```bash
make test              # 单元测试
make test-integration  # 集成测试（需要 Docker）
make lint              # 代码检查
```

## 项目结构

```
├── backend/
│   ├── cmd/server/          # 入口 + Wire DI
│   ├── configs/             # YAML 配置
│   ├── internal/
│   │   ├── biz/             # 业务逻辑
│   │   ├── data/            # 数据访问 (sqlc)
│   │   ├── service/         # gRPC/HTTP 服务
│   │   ├── server/          # Kratos 服务器 + WebSocket
│   │   ├── middleware/       # JWT 认证 · 加密
│   │   ├── provider/        # LLM / TTS / STT 适配
│   │   └── conf/            # 配置结构体
│   └── sql/                 # SQL 迁移 + 查询
├── frontend/
│   └── src/
│       ├── api/             # HTTP + WebSocket 客户端
│       ├── stores/          # Pinia 状态
│       ├── composables/     # 音频播放 · 语音识别
│       └── views/           # 页面组件
├── docker-compose.yml       # 全栈部署
├── docker-compose.dev.yml   # 开发（仅基础设施）
├── Makefile                 # 构建 / 测试 / 代码生成
└── .github/workflows/       # CI + Release
```

## 文档

详细文档见 [docs/](docs/) 目录：

- [架构设计](docs/architecture.md)
- [部署指南](docs/deployment.md)
- [API 文档](docs/api.md)
- [Provider 配置](docs/providers.md)

## 许可证

MIT
