package llm

import (
	"context"
	"fmt"
)

// Provider 定义 LLM 供应商的统一接口
type Provider interface {
	// ChatStream 流式对话，通过 channel 返回文本片段
	ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamEvent, error)
	// Name 返回供应商标识
	Name() string
}

// ChatRequest 是 LLM 对话请求
type ChatRequest struct {
	Messages    []Message // 对话历史
	Model       string    // 模型名称
	MaxTokens   int       // 最大输出 token 数
	Temperature float64   // 温度
	APIKey      string    // BYOK API Key
	BaseURL     string    // 自定义端点 (可选)
}

// Message 是对话消息
type Message struct {
	Role    string // system, user, assistant
	Content string
}

// StreamEvent 是流式输出事件
type StreamEvent struct {
	Content string // 文本片段
	Done    bool   // 是否结束
	Err     error  // 错误信息
}

// Registry 管理所有注册的 LLM Provider
type Registry struct {
	providers map[string]Provider
}

// NewRegistry 创建 Provider 注册表
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register 注册一个 LLM Provider
func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

// Get 获取指定名称的 Provider
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("llm provider %q not found", name)
	}
	return p, nil
}

// List 列出所有可用的 Provider 名称
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}
