package tts

import (
	"context"
	"fmt"
	"io"
)

// Provider 定义 TTS 供应商的统一接口
type Provider interface {
	// Synthesize 将文本转为音频流，通过 writer 流式写入
	Synthesize(ctx context.Context, req *Request, w io.Writer) error
	// Name 返回供应商标识
	Name() string
}

// Request 是 TTS 合成请求
type Request struct {
	Text         string  // 要合成的文本
	Voice        string  // 语音名称
	Language     string  // 语言代码 (zh-CN, en-US)
	Format       string  // 音频格式 (pcm, mp3, opus)
	SampleRate   int     // 采样率 (24000)
	Speed        float64 // 语速 (0.5-2.0)
	Instructions string  // 语气/情感指令 (仅 gpt-4o-mini-tts 支持)
	APIKey       string  // BYOK API Key
	BaseURL      string  // 自定义端点 (可选)
}

// Registry 管理所有注册的 TTS Provider
type Registry struct {
	providers map[string]Provider
}

// NewRegistry 创建 Provider 注册表
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register 注册一个 TTS Provider
func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

// Get 获取指定名称的 Provider
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("tts provider %q not found", name)
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
