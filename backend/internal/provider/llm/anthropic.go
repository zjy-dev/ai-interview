package llm

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AnthropicProvider 实现 Anthropic Claude LLM
type AnthropicProvider struct {
	baseURL    string
	httpClient *http.Client
}

func NewAnthropicProvider() *AnthropicProvider {
	return &AnthropicProvider{
		baseURL:    "https://api.anthropic.com",
		httpClient: &http.Client{},
	}
}

func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
	System    string             `json:"system,omitempty"`
	Stream    bool               `json:"stream"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicStreamEvent struct {
	Type  string `json:"type"`
	Delta *struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta,omitempty"`
}

func (p *AnthropicProvider) ChatStream(ctx context.Context, req *ChatRequest) (<-chan StreamEvent, error) {
	if req.APIKey == "" {
		return nil, fmt.Errorf("anthropic llm: api key is required")
	}

	baseURL := p.baseURL
	if req.BaseURL != "" {
		baseURL = req.BaseURL
	}

	// 分离 system 消息和对话消息
	var systemPrompt string
	messages := make([]anthropicMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		if m.Role == "system" {
			systemPrompt = m.Content
			continue
		}
		messages = append(messages, anthropicMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	model := req.Model
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	apiReq := anthropicRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  messages,
		System:    systemPrompt,
		Stream:    true,
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic llm: marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v1/messages", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("anthropic llm: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", req.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic llm: request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("anthropic llm: status %d: %s", resp.StatusCode, string(respBody))
	}

	ch := make(chan StreamEvent, 32)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				ch <- StreamEvent{Done: true}
				return
			}

			var event anthropicStreamEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				continue
			}

			switch event.Type {
			case "content_block_delta":
				if event.Delta != nil && event.Delta.Text != "" {
					ch <- StreamEvent{Content: event.Delta.Text}
				}
			case "message_stop":
				ch <- StreamEvent{Done: true}
				return
			case "error":
				ch <- StreamEvent{Err: fmt.Errorf("anthropic llm: stream error")}
				return
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- StreamEvent{Err: fmt.Errorf("anthropic llm: read stream: %w", err)}
		}
	}()

	return ch, nil
}
