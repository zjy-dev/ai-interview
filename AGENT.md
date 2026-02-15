# AGENT.md — AI Agent 项目上下文

> 本文件供 AI 编程代理使用。包含项目架构、开发规则和实现细节。

## 项目概述

AI 面试平台：LLM 生成面试问题 → TTS 合成语音 → WebSocket 实时推送音频 → 用户语音回答 → STT 转文字 → LLM 继续对话。

用户自带 API Key (BYOK)，密钥 AES-GCM 加密存储。未来支持付费使用平台提供的 Key。

## RULES

- 后端编程语言使用 Go，框架 Kratos v2，前端框架用 Vue 3
- Go 项目要用 Makefile 做可复现的测试/编译
- **不要生成兼容性代码!!!**
- 所有工具链用最新版本，可以用 context7 获取最新文档
- 不用 ORM，用 sqlc
- 要有单元测试和集成测试
- 敏感配置从环境变量和 .env 中读，非敏感配置写入 YAML 配置文件
- 关系型数据库用 MySQL
- 容器编排使用标准 Compose spec（Docker Compose / Podman Compose 兼容）
- CI/CD 使用 GitHub Actions，容器镜像推送到 GHCR

## 技术栈详情

| 组件 | 技术 | 版本 |
|------|------|------|
| Go | golang | 1.23+ |
| 框架 | go-kratos/kratos/v2 | v2.9.2 |
| DB 代码生成 | sqlc | latest |
| DI | google/wire | latest |
| OpenAI SDK | sashabaranov/go-openai | v1.41.2 |
| WebSocket | nhooyr.io/websocket | v1.8.17 |
| 前端 | Vue 3 + TypeScript + Vite | Vue 3.5+, Vite 6+ |
| 状态管理 | Pinia | 3.x |
| 路由 | Vue Router | 4.x |
| 数据库 | MySQL | 8.0 |
| 缓存 | Redis | 7.x |

## 项目结构

```
ai-interview/
├── AGENT.md                    # 本文件 (AI agent 上下文)
├── README.md                   # 人类用户文档
├── Makefile                    # 构建 / 测试 / 代码生成
├── .env.example                # 环境变量模板
├── docker-compose.yml          # 全栈容器部署
├── docker-compose.dev.yml      # 开发模式（仅 mysql + redis）
├── .github/workflows/
│   ├── ci.yml                  # PR/push → lint + test + build
│   └── release.yml             # v* tag → 构建镜像推送 GHCR
│
├── backend/
│   ├── Dockerfile              # 多阶段构建 (golang:1.23-alpine → alpine:3.20)
│   ├── go.mod / go.sum
│   ├── sqlc.yaml               # sqlc 配置
│   ├── cmd/server/
│   │   ├── main.go             # 入口：加载 .env → 读 YAML 配置 → Wire 注入 → 启动
│   │   ├── wire.go             # Wire provider 声明
│   │   └── wire_gen.go         # Wire 生成代码 (勿手动编辑)
│   ├── configs/
│   │   └── config.yaml         # 非敏感配置 (server, tts, llm, interview defaults)
│   ├── sql/
│   │   ├── migrations/         # DDL (000001_init.up/down.sql)
│   │   └── queries/            # sqlc 查询 (users.sql, interviews.sql)
│   └── internal/
│       ├── conf/
│       │   └── conf.go         # 配置结构体 (手写 Go struct，非 protobuf)
│       ├── biz/
│       │   ├── user.go         # UserUsecase: Register, Login, Profile, Settings
│       │   ├── interview.go    # InterviewUsecase: CRUD, streamMessage, buildLLMMessages
│       │   │                   #   StreamMessage → LLM 流式输出 → 分句 → TTS → StreamEvent
│       │   └── errors.go       # hashPassword, checkPassword (bcrypt)
│       ├── data/
│       │   ├── data.go         # NewData: MySQL (*sql.DB) + Redis (go-redis)
│       │   ├── user.go         # UserRepo 实现 (raw SQL via sqlc)
│       │   └── interview.go    # InterviewRepo 实现
│       ├── service/
│       │   ├── auth.go         # AuthService: Register, Login, Profile, UpdateSettings
│       │   └── interview.go    # InterviewService: CRUD, SendMessage, EndInterview
│       ├── server/
│       │   ├── http.go         # Kratos HTTP 路由注册
│       │   ├── handler.go      # REST handler (JSON request/response)
│       │   ├── websocket.go    # WebSocket 面试处理:
│       │   │                   #   handleInterview → 接收文字 → biz.StreamMessage →
│       │   │                   #   分句 isSentenceEnd → TTS 合成 → 音频推送客户端
│       │   └── server.go       # NewServer 构造
│       ├── middleware/
│       │   ├── auth.go         # JWT 生成 / 验证 / HTTP 中间件
│       │   └── crypto.go       # AES-GCM 加密 / 解密 (用于 BYOK API Key 存储)
│       └── provider/
│           ├── llm/
│           │   ├── llm.go      # Provider 接口 + Registry
│           │   ├── openai.go   # OpenAI (GPT-4o 等)
│           │   ├── anthropic.go # Anthropic (Claude)
│           │   ├── gemini.go   # Google Gemini
│           │   └── deepseek.go # DeepSeek (通过 go-openai SDK, 自定义 BaseURL)
│           ├── tts/
│           │   ├── tts.go      # Provider 接口 + Registry
│           │   ├── openai.go   # OpenAI TTS
│           │   ├── fishaudio.go # Fish Audio
│           │   ├── elevenlabs.go # ElevenLabs
│           │   └── edgetts.go  # Edge TTS (免费，WebSocket 协议)
│           └── stt/
│               ├── stt.go      # Provider 接口 + Registry
│               └── whisper.go  # OpenAI Whisper
│
├── frontend/
│   ├── Dockerfile              # 多阶段构建 (node:22-alpine → nginx:alpine)
│   ├── nginx.conf              # SPA history mode + /api/ 反向代理 + WebSocket
│   ├── package.json
│   ├── vite.config.ts
│   └── src/
│       ├── api/
│       │   ├── client.ts       # Axios 实例 + JWT 拦截器
│       │   ├── auth.ts         # 注册 / 登录 / profile / settings API
│       │   └── interview.ts    # 面试 CRUD + report API
│       ├── stores/
│       │   ├── auth.ts         # Pinia: user state, token, login/logout
│       │   └── interview.ts    # Pinia: interview list, current interview
│       ├── composables/
│       │   ├── useAudioPlayer.ts     # PCM 音频播放 (AudioContext + AudioWorklet)
│       │   ├── useSpeechRecognition.ts # 浏览器 Web Speech API 语音识别
│       │   └── useWebSocket.ts       # WebSocket 连接管理
│       ├── views/
│       │   ├── LoginView.vue
│       │   ├── RegisterView.vue
│       │   ├── InterviewListView.vue
│       │   ├── NewInterviewView.vue
│       │   ├── InterviewView.vue     # 核心面试页面 (WS + 音频 + 语音)
│       │   ├── ReportView.vue
│       │   └── SettingsView.vue      # BYOK API Key 配置
│       ├── components/
│       │   └── NavBar.vue
│       ├── router/index.ts           # 路由 + auth guard
│       └── App.vue
│
└── docs/
    ├── architecture.md         # 架构设计
    ├── deployment.md           # 部署指南
    ├── api.md                  # HTTP API 文档
    └── providers.md            # Provider 配置说明
```

## 关键架构

### 分层架构 (Kratos 风格)

```
HTTP Request → server (路由/handler) → service (编排) → biz (业务逻辑) → data (数据访问)
                                                          ↓
                                                     provider (LLM/TTS/STT)
```

### WebSocket 面试流程

```
客户端文字消息
  → biz.StreamMessage()
    → LLM provider 流式生成回答
      → 分句检测 (isSentenceEnd: 。！？!?.\n)
        → TTS provider 合成音频
          → WebSocket 推送 audio binary frame
  → 同时推送 text frame (逐 token)
```

### 数据模型 (MySQL)

5 张表 (定义在 sql/migrations/000001_init.up.sql):

1. **users** — id, username, email, password_hash, created_at, updated_at
2. **user_settings** — user_id (FK), provider 偏好, 加密的 API keys
3. **interviews** — id, user_id (FK), title, position, status, config_json
4. **interview_messages** — id, interview_id (FK), role, content, audio_url, created_at
5. **interview_reports** — interview_id (FK), overall_score, summary, details_json

### 环境变量 (.env)

| 变量 | 用途 |
|------|------|
| DB_DSN | MySQL 连接串 |
| REDIS_PASSWORD | Redis 密码 |
| JWT_SECRET | JWT 签名密钥 |
| ENCRYPTION_KEY | AES-GCM 密钥 (64 hex chars = 32 bytes) |

### 配置文件 (configs/config.yaml)

非敏感默认值：server 监听地址、TTS/LLM 默认 provider 和模型、面试参数等。
运行时环境变量会覆盖 YAML 中的空值。

## 开发命令

```bash
make init          # 安装 Go 工具链
make sqlc          # 生成 sqlc 代码
make wire          # 生成 Wire DI
make build         # 编译后端
make test          # 单元测试
make lint          # golangci-lint
make dev           # 启动后端开发服务器

# 前端
cd frontend && npm ci && npm run dev

# 基础设施
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# 全栈部署
docker compose up -d
```

## 测试覆盖

现有测试文件：
- `internal/biz/errors_test.go` — bcrypt 密码哈希
- `internal/middleware/auth_test.go` — JWT 生成/验证/过期/提取
- `internal/middleware/crypto_test.go` — AES-GCM 加解密
- `internal/provider/tts/tts_test.go` — TTS Registry
- `internal/provider/llm/llm_test.go` — LLM Registry
- `internal/server/websocket_test.go` — 分句检测 isSentenceEnd

运行: `make test` (14 tests)

## CI/CD

- **ci.yml**: push/PR 到 main → backend lint (golangci-lint) → backend test (race + coverage) → backend build → frontend lint+build
- **release.yml**: push v* tag → 多阶段构建 backend/frontend 镜像 → 推送到 `ghcr.io/<owner>/ai-interview/{backend,frontend}` → semver 标签

## 注意事项

- Wire 生成代码 (`wire_gen.go`) 不要手动编辑，修改 `wire.go` 后运行 `make wire`
- sqlc 生成代码在 `internal/data/` 对应文件，修改 SQL 后运行 `make sqlc`
- Provider 采用 Registry 模式，新增 provider 只需实现接口并注册
- Edge TTS 是唯一免费 TTS provider，使用 WebSocket 连接微软服务
- DeepSeek 通过 go-openai SDK 连接（自定义 BaseURL），与 OpenAI provider 共享代码模式
- 前端 nginx 配置已处理 Vue Router history mode 和 WebSocket 升级
