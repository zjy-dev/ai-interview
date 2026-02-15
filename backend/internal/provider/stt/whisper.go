package stt

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// WhisperProvider 实现 OpenAI Whisper STT
type WhisperProvider struct{}

func NewWhisperProvider() *WhisperProvider {
	return &WhisperProvider{}
}

func (p *WhisperProvider) Name() string {
	return "whisper"
}

func (p *WhisperProvider) Transcribe(ctx context.Context, req *Request) (*Result, error) {
	if req.APIKey == "" {
		return nil, fmt.Errorf("whisper stt: api key is required")
	}

	config := openai.DefaultConfig(req.APIKey)
	if req.BaseURL != "" {
		config.BaseURL = req.BaseURL
	}
	client := openai.NewClientWithConfig(config)

	model := req.Model
	if model == "" {
		model = "whisper-1"
	}

	language := req.Language
	if language == "" {
		language = "zh"
	}

	transcriptionReq := openai.AudioRequest{
		Model:    model,
		Reader:   req.Audio,
		FilePath: "audio." + req.Format,
		Language: language,
	}

	resp, err := client.CreateTranscription(ctx, transcriptionReq)
	if err != nil {
		return nil, fmt.Errorf("whisper stt: %w", err)
	}

	return &Result{
		Text:     resp.Text,
		Language: language,
		Duration: float64(resp.Duration),
	}, nil
}
