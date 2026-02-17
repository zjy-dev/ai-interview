# Provider 配置

本平台支持多家 LLM / TTS / STT 服务商，用户在「设置」页面配置自己的 API Key (BYOK)。

## LLM Providers

### OpenAI

- **Provider 名称**: `openai`
- **支持模型**: gpt-4o, gpt-4o-mini, gpt-4-turbo, gpt-3.5-turbo 等
- **API Key**: OpenAI API Key (`sk-...`)
- **Base URL**: 默认 `https://api.openai.com/v1`，可自定义（兼容 API 代理）

### Anthropic

- **Provider 名称**: `anthropic`
- **支持模型**: claude-sonnet-4-20250514, claude-3.5-haiku 等
- **API Key**: Anthropic API Key
- **说明**: 使用官方 REST API

### Google Gemini

- **Provider 名称**: `gemini`
- **支持模型**: gemini-2.5-flash, gemini-2.5-pro 等
- **API Key**: Google AI Studio API Key
- **说明**: 通过 OpenAI 兼容接口调用

### DeepSeek

- **Provider 名称**: `deepseek`
- **支持模型**: deepseek-chat, deepseek-reasoner
- **API Key**: DeepSeek API Key
- **Base URL**: `https://api.deepseek.com/v1`
- **说明**: 通过 go-openai SDK 自定义 BaseURL 实现

## TTS Providers

### OpenAI TTS

- **Provider 名称**: `openai`
- **支持声音**: alloy, echo, fable, onyx, nova, shimmer
- **API Key**: 与 LLM 共用 OpenAI API Key
- **输出格式**: MP3

### Fish Audio

- **Provider 名称**: `fishaudio`
- **API Key**: Fish Audio API Key
- **说明**: 支持中文语音克隆，通过 REST API 调用
- **官网**: https://fish.audio

### ElevenLabs

- **Provider 名称**: `elevenlabs`
- **API Key**: ElevenLabs API Key
- **说明**: 高质量多语言语音合成
- **官网**: https://elevenlabs.io

### Edge TTS（免费）

- **Provider 名称**: `edgetts`
- **API Key**: 不需要
- **说明**: 使用微软 Edge 浏览器内置的 TTS 服务，通过 WebSocket 协议连接，完全免费
- **支持声音**: zh-CN-XiaoxiaoNeural, zh-CN-YunxiNeural, en-US-JennyNeural 等
- **限制**: 非官方 API，可能有频率限制

## STT Providers

### Whisper (OpenAI)

- **Provider 名称**: `whisper`
- **API Key**: OpenAI API Key
- **说明**: OpenAI Whisper 语音转文字 API

### Browser (浏览器内置)

- **Provider 名称**: `browser`
- **API Key**: 不需要
- **说明**: 使用浏览器 Web Speech API，无需服务端 STT。由前端 `useSpeechRecognition` composable 实现。
- **兼容性**: Chrome/Edge 支持最佳

## 配置方式

### 通过前端设置页面

用户登录后进入「设置」页面，选择 Provider 并填入 API Key。

### 默认值

系统默认值在 `backend/configs/config.yaml` 中配置：

```yaml
tts:
  default_provider: openai
  default_voice: alloy

llm:
  default_provider: openai
  default_model: gpt-4o
```

用户未配置时使用系统默认值。

## 安全性

- 用户 API Key 使用 **AES-256-GCM** 加密后存储到数据库
- 加密密钥从环境变量 `ENCRYPTION_KEY` 读取（64 hex chars = 32 bytes）
- API Key 查询接口只返回 `*_api_key_set: true/false`，不返回明文
- 每次调用外部 API 时解密使用，不缓存明文

## 新增 Provider

Provider 采用 Registry 模式，新增步骤：

1. 实现对应接口 (如 `llm.Provider`)
2. 在 `cmd/server/main.go` 的 `initLLMRegistry()`（或对应的 `initTTSRegistry()`/`initSTTRegistry()`）中注册新 Provider
3. 无需修改其他代码，运行时自动可用

示例：

```go
// internal/provider/llm/newprovider.go
package llm

type NewProvider struct{}

func (p *NewProvider) Name() string { return "newprovider" }

func (p *NewProvider) ChatStream(ctx context.Context, apiKey string, msgs []Message, opts Options) (<-chan StreamEvent, error) {
    // 实现流式聊天
}
```

```go
// cmd/server/main.go
func initLLMRegistry() *llm.Registry {
    registry := llm.NewRegistry()
    // ... 其他 provider
    registry.Register(llm.NewNewProvider())
    return registry
}
```
