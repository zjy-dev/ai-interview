# API 文档

Base URL: `/api/v1`

所有需要认证的接口都需要在 Header 中携带 `Authorization: Bearer <token>`。

---

## 认证

### POST /auth/register

注册新用户。

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secret123",
  "nickname": "张三"
}
```

**Response 200:**
```json
{
  "id": 1,
  "token": "eyJhbGci..."
}
```

---

### POST /auth/login

用户登录。

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

**Response 200:**
```json
{
  "id": 1,
  "token": "eyJhbGci...",
  "nickname": "张三"
}
```

---

### GET /auth/profile 🔒

获取当前用户信息。

**Response 200:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "nickname": "张三",
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

### GET /auth/settings 🔒

获取用户 AI 服务配置。

**Response 200:**
```json
{
  "llm_provider": "openai",
  "llm_api_key_set": true,
  "llm_base_url": "",
  "llm_model": "gpt-4o",
  "tts_provider": "openai",
  "tts_api_key_set": true,
  "tts_voice": "alloy",
  "tts_enabled": true,
  "stt_provider": "browser",
  "stt_api_key_set": false
}
```

> 注：API Key 不会返回明文，只返回 `_set` 布尔值。

---

### PUT /auth/settings 🔒

更新用户 AI 服务配置。

**Request:**
```json
{
  "llm_provider": "openai",
  "llm_api_key": "sk-xxx",
  "llm_base_url": "",
  "llm_model": "gpt-4o",
  "tts_provider": "openai",
  "tts_api_key": "sk-xxx",
  "tts_voice": "alloy",
  "tts_enabled": true,
  "stt_provider": "browser",
  "stt_api_key": ""
}
```

**Response 200:**
```json
{"message": "ok"}
```

---

## 面试

### POST /interviews 🔒

创建新面试。

**Request:**
```json
{
  "title": "后端工程师面试",
  "position": "Senior Go Developer",
  "resume": "5年Go开发经验...",
  "language": "zh-CN"
}
```

**Response 200:**
```json
{
  "id": 1
}
```

---

### GET /interviews 🔒

获取面试列表。

**Response 200:**
```json
{
  "interviews": [
    {
      "id": 1,
      "title": "后端工程师面试",
      "position": "Senior Go Developer",
      "status": "in_progress",
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

---

### GET /interviews/{id} 🔒

获取面试详情及消息历史。

**Response 200:**
```json
{
  "id": 1,
  "title": "后端工程师面试",
  "position": "Senior Go Developer",
  "status": "in_progress",
  "messages": [
    {
      "id": 1,
      "role": "assistant",
      "content": "你好，我是面试官...",
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

---

### POST /interviews/{id}/messages 🔒

发送面试消息（非实时模式）。

**Request:**
```json
{
  "content": "我有5年Go开发经验..."
}
```

**Response 200:**
```json
{
  "reply": "AI 面试官的回复..."
}
```

---

### POST /interviews/{id}/end 🔒

结束面试并生成评估。

**Response 200:**
```json
{"message": "ok"}
```

---

### GET /interviews/{id}/evaluation 🔒

获取面试评估报告。

**Response 200:**
```json
{
  "overall_score": 85,
  "summary": "候选人表现良好...",
  "categories": [
    {"name": "技术能力", "score": 90},
    {"name": "沟通表达", "score": 80}
  ],
  "strengths": "技术基础扎实...",
  "weaknesses": "系统设计经验不足...",
  "suggestions": "建议加强..."
}
```

---

## WebSocket 面试

### GET /ws/interview/{id} 🔒

建立 WebSocket 连接进行实时面试。

**连接:**
```
ws://host/api/v1/ws/interview/{id}
Headers: Authorization: Bearer <token>
```

**客户端 → 服务端** (Text Frame):
```json
{"type": "message", "content": "用户的回答"}
```

**服务端 → 客户端** (Text Frame):
```json
{"type": "text", "content": "单个 token"}
{"type": "done"}
{"type": "error", "content": "错误描述"}
```

**服务端 → 客户端** (Binary Frame):
- PCM 音频数据 (24kHz, 16-bit, mono)
- 逐句合成推送，不等整段回复完成

---

## 健康检查

### GET /health

**Response 200:**
```json
{"status": "ok"}
```

---

## 错误格式

所有错误统一格式：
```json
{"error": "错误描述"}
```

HTTP 状态码：
- `400` — 请求参数错误
- `401` — 未认证 / Token 无效
- `404` — 资源不存在
- `500` — 服务器内部错误
