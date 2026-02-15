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

### 1. 克隆 & 配置

```bash
git clone https://github.com/zjy-dev/ai-interview.git
cd ai-interview
cp .env.example .env   # 编辑 .env 填入密钥
```

### 2. 启动基础设施

```bash
# 仅启动 MySQL + Redis（开发模式）
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

### 3. 安装依赖 & 构建

```bash
make init          # 安装 Go 工具链
make sqlc           # 生成数据库代码
make wire           # 生成依赖注入
make build          # 编译后端

cd frontend && npm ci   # 安装前端依赖
```

### 4. 运行

```bash
# 后端
make dev

# 前端（另一个终端）
cd frontend && npm run dev
```

访问 http://localhost:5173

### 一键容器部署

```bash
cp .env.example .env   # 编辑 .env
docker compose up -d   # 启动全部服务
```

访问 http://localhost

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
