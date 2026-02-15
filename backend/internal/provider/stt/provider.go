package stt

import (
	"context"
	"fmt"
	"io"
)

// Provider 定义 STT 供应商的统一接口
type Provider interface {
	// Transcribe 将音频转为文本
	Transcribe(ctx context.Context, req *Request) (*Result, error)
	// Name 返回供应商标识
	Name() string
}

// Request 是 STT 转写请求
type Request struct {
	Audio    io.Reader // 音频数据
	Format   string    // 音频格式 (webm, wav, mp3)
	Language string    // 语言代码 (zh, en)
	Model    string    // 模型名称 (whisper-1, gpt-4o-mini-transcribe)
	APIKey   string    // BYOK API Key
	BaseURL  string    // 自定义端点 (可选)
}

// Result 是 STT 转写结果
type Result struct {
	Text     string  // 转写文本
	Language string  // 检测到的语言
	Duration float64 // 音频时长 (秒)
}

// Registry 管理所有注册的 STT Provider
type Registry struct {
	providers map[string]Provider
}

// NewRegistry 创建 Provider 注册表
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register 注册一个 STT Provider
func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

// Get 获取指定名称的 Provider
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("stt provider %q not found", name)
	}
	return p, nil
}
