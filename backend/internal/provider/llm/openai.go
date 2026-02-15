package llm

import (
	"context"
	"errors"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIProvider 实现 OpenAI / DeepSeek / 任意 OpenAI 兼容 LLM
type OpenAIProvider struct {
	name           string
	defaultBaseURL string
}

// NewOpenAIProvider 创建 OpenAI LLM Provider
func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{name: "openai"}
}

// NewDeepSeekProvider 创建 DeepSeek LLM Provider (OpenAI 兼容)
func NewDeepSeekProvider() *OpenAIProvider {
	return &OpenAIProvider{
		name:           "deepseek",
		defaultBaseURL: "https://api.deepseek.com/v1",
	}
}

// NewCustomProvider 创建自定义 OpenAI 兼容 LLM Provider
func NewCustomProvider(name, baseURL string) *OpenAIProvider {
	return &OpenAIProvider{
		name:           name,
		defaultBaseURL: baseURL,
	}
}

func (p *OpenAIProvider) Name() string {
	return p.name
}

func (p *OpenAIProvider) ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamEvent, error) {
	if req.APIKey == "" {
		return nil, fmt.Errorf("%s llm: api key is required", p.name)
	}

	config := openai.DefaultConfig(req.APIKey)
	if req.BaseURL != "" {
		config.BaseURL = req.BaseURL
	} else if p.defaultBaseURL != "" {
		config.BaseURL = p.defaultBaseURL
	}
	client := openai.NewClientWithConfig(config)

	messages := make([]openai.ChatCompletionMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	model := req.Model
	if model == "" {
		model = "gpt-4o"
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	temp := float32(req.Temperature)
	if temp == 0 {
		temp = 0.7
	}

	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temp,
		Stream:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s llm: create stream: %w", p.name, err)
	}

	ch := make(chan StreamEvent, 32)
	go func() {
		defer close(ch)
		defer stream.Close()

		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				ch <- StreamEvent{Done: true}
				return
			}
			if err != nil {
				ch <- StreamEvent{Err: fmt.Errorf("%s llm: stream recv: %w", p.name, err)}
				return
			}

			if len(resp.Choices) > 0 {
				delta := resp.Choices[0].Delta.Content
				if delta != "" {
					ch <- StreamEvent{Content: delta}
				}
			}
		}
	}()

	return ch, nil
}
