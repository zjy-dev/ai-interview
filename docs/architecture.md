# 架构设计

## 总体架构

```
┌──────────────────────────────────────────────────┐
│                  Frontend (Vue 3)                 │
│                  nginx:80                         │
│  ┌────────┐ ┌──────────┐ ┌────────────────────┐  │
│  │ Stores │ │   API    │ │   Composables      │  │
│  │(Pinia) │ │ (Axios)  │ │ Audio/Speech/WS    │  │
│  └────────┘ └──────────┘ └────────────────────┘  │
└─────────────────────┬────────────────────────────┘
                      │ HTTP / WebSocket
                      ▼
┌──────────────────────────────────────────────────┐
│                Backend (Go + Kratos)              │
│                0.0.0.0:8080                       │
│                                                   │
│  ┌─────────────────────────────────────────────┐  │
│  │ Server Layer                                │  │
│  │  HTTP Router → Handler (REST)               │  │
│  │  WebSocket Handler (实时面试)               │  │
│  │  JWT Middleware                              │  │
│  └──────────────────────┬──────────────────────┘  │
│                         ▼                         │
│  ┌─────────────────────────────────────────────┐  │
│  │ Service Layer (编排)                         │  │
│  │  AuthService · InterviewService             │  │
│  └──────────────────────┬──────────────────────┘  │
│                         ▼                         │
│  ┌─────────────────────────────────────────────┐  │
│  │ Biz Layer (业务逻辑)                         │  │
│  │  UserUsecase · InterviewUsecase             │  │
│  │  密码哈希 · LLM 流式分句 · 评估生成         │  │
│  └──────────────┬─────────────┬────────────────┘  │
│                 ▼             ▼                    │
│  ┌──────────────────┐ ┌───────────────────────┐   │
│  │ Data Layer       │ │ Provider Layer        │   │
│  │ MySQL (手写 SQL) │ │ LLM / TTS / STT      │   │
│  │ Redis            │ │ (Registry 模式)       │   │
│  └──────────────────┘ └───────────────────────┘   │
└──────────────────────────────────────────────────┘
         │                          │
         ▼                          ▼
    ┌─────────┐            ┌─────────────────┐
    │ MySQL 8 │            │ 外部 API        │
    │ Redis 7 │            │ OpenAI/Anthropic│
    └─────────┘            │ Gemini/DeepSeek │
                           │ ElevenLabs/Edge │
                           └─────────────────┘
```

## 分层职责

### Server 层 (`internal/server/`)

- **http.go** — Kratos HTTP 服务器创建 + 路由注册
- **handler.go** — REST API handler，处理 JSON 序列化/反序列化
- **websocket.go** — WebSocket 面试处理，协调 LLM 流 → 分句 → TTS → 音频推送
- **server.go** — 构造函数，注入依赖

路由通过 `srv.Route("/")` 注册，认证路由使用 `withAuth()` 包装。

### Service 层 (`internal/service/`)

编排层，连接 handler 和 biz 层：

- **AuthService** — Register、Login、GetProfile、GetSettings、UpdateSettings
- **InterviewService** — Create、List、Get、SendMessage、EndInterview、GetEvaluation

### Biz 层 (`internal/biz/`)

核心业务逻辑，不依赖具体框架：

- **UserUsecase** — 用户注册/登录流程、密码验证
- **InterviewUsecase** — 面试创建、消息处理、StreamMessage（LLM 流式 + 分句 + TTS）、评估报告生成

关键方法 `StreamMessage()` 流程：
1. 构建 LLM 消息历史 (`buildLLMMessages`)
2. 调用 LLM provider 流式生成
3. token 累积，`isSentenceEnd()` 检测句子边界
4. 完整句子触发 TTS 合成
5. 通过 `StreamEvent` channel 返回 text/audio 事件

### Data 层 (`internal/data/`)

数据访问，使用手写 SQL（`database/sql` + `ExecContext/QueryRowContext`）：

- **data.go** — 初始化 `*sql.DB` (MySQL) 和 `*redis.Client`
- **user.go** — users / user_settings 表操作
- **interview.go** — interviews / interview_messages / evaluations 表操作

### Provider 层 (`internal/provider/`)

外部 AI 服务适配，采用 **Registry 模式**：

每种能力 (LLM/TTS/STT) 定义统一接口 + Registry。在 `cmd/server/main.go` 中集中注册所有 Provider 实例到 Registry，运行时根据用户配置的 provider 名称查找。

### Middleware (`internal/middleware/`)

- **auth.go** — JWT 生成 (`GenerateToken`) / 验证 (`ValidateToken`) / HTTP 中间件
- **crypto.go** — AES-256-GCM 加密/解密，用于保护用户 BYOK API Key（存储时加密，使用时解密）

## 依赖注入

使用 Google Wire 编译时 DI。Provider 集声明在 `cmd/server/wire.go`，生成代码在 `wire_gen.go`。

注入链：`conf → data → biz → service + provider → server → app`

## WebSocket 实时面试协议

### 连接

```
GET /api/v1/ws/interview/{id}?token=<jwt_token>
Upgrade: websocket
```

> 注：WebSocket 不支持自定义 Header，认证通过 `?token=` 查询参数传递。

### 消息格式

**客户端 → 服务端** (Text Frame):
```json
{"type": "text", "data": "用户回答文字"}
{"type": "end"}
{"type": "ping"}
```

**服务端 → 客户端** (Text Frame):
```json
{"type": "text_start"}
{"type": "text_delta", "data": "单个 token"}
{"type": "text_end", "data": "完整回复文本"}
{"type": "status", "data": "connected"}
{"type": "evaluation", "data": {"overall_score": 85, "summary": "..."}}
{"type": "error", "data": "错误信息"}
```

**服务端 → 客户端** (Binary Frame):
- MP3 音频数据，逐句合成推送

## 数据库设计

5 张表，定义在 `sql/migrations/000001_init.up.sql`：

| 表 | 说明 |
|----|------|
| users | 用户基本信息，email 唯一索引 |
| user_settings | 1:1 用户设置，存储 provider 偏好 + 加密 API key |
| interviews | 面试会话，含 provider/model 配置快照 |
| interview_messages | 面试消息记录 (system/user/assistant) |
| evaluations | 面试评估报告，含分项 JSON + 优缺点 |

所有表使用 `utf8mb4_unicode_ci`，InnoDB 引擎，外键级联删除。
