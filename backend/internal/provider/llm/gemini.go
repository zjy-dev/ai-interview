package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// GeminiProvider 实现 Google Gemini LLM (通过 OpenAI 兼容接口)
type GeminiProvider struct {
	defaultBaseURL string
}

func NewGeminiProvider() *GeminiProvider {
	return &GeminiProvider{
		defaultBaseURL: "https://generativelanguage.googleapis.com/v1beta/openai",
	}
}

func (p *GeminiProvider) Name() string {
	return "gemini"
}

func (p *GeminiProvider) ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamEvent, error) {
	if req.APIKey == "" {
		return nil, fmt.Errorf("gemini llm: api key is required")
	}

	baseURL := p.defaultBaseURL
	if req.BaseURL != "" {
		baseURL = req.BaseURL
	}

	// Gemini 提供 OpenAI 兼容端点，复用 OpenAI 逻辑
	proxy := &OpenAIProvider{
		name:           "gemini",
		defaultBaseURL: baseURL,
	}

	model := req.Model
	if model == "" {
		model = "gemini-2.0-flash"
	}

	// 创建带覆盖 model 的新 request
	proxyReq := &ChatRequest{
		Messages:    req.Messages,
		Model:       model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		APIKey:      req.APIKey,
		BaseURL:     baseURL,
	}

	// 使用 DefaultConfig 并通过 APIKey 认证
	config := openai.DefaultConfig(req.APIKey)
	config.BaseURL = baseURL
	_ = config // proxy 内部会再创建

	return proxy.ChatStream(ctx, proxyReq)
}
