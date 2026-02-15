package tts

import (
	"context"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIProvider 实现 OpenAI TTS API
type OpenAIProvider struct{}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) Synthesize(ctx context.Context, req *Request, w io.Writer) error {
	if req.APIKey == "" {
		return fmt.Errorf("openai tts: api key is required")
	}

	config := openai.DefaultConfig(req.APIKey)
	if req.BaseURL != "" {
		config.BaseURL = req.BaseURL
	}
	client := openai.NewClientWithConfig(config)

	voice := openai.VoiceAlloy
	if req.Voice != "" {
		voice = openai.SpeechVoice(req.Voice)
	}

	model := openai.TTSModel1
	if req.Instructions != "" {
		// gpt-4o-mini-tts 支持 instructions 参数
		model = openai.TTSModelGPT4oMini
	}

	format := openai.SpeechResponseFormatPcm
	switch req.Format {
	case "mp3":
		format = openai.SpeechResponseFormatMp3
	case "opus":
		format = openai.SpeechResponseFormatOpus
	case "wav":
		format = openai.SpeechResponseFormatWav
	case "flac":
		format = openai.SpeechResponseFormatFlac
	}

	speed := 1.0
	if req.Speed > 0 {
		speed = req.Speed
	}

	speechReq := openai.CreateSpeechRequest{
		Model:          model,
		Input:          req.Text,
		Voice:          voice,
		ResponseFormat: format,
		Speed:          speed,
	}

	resp, err := client.CreateSpeech(ctx, speechReq)
	if err != nil {
		return fmt.Errorf("openai tts: %w", err)
	}
	defer resp.Close()

	// 流式写入音频数据
	if _, err := io.Copy(w, resp); err != nil {
		return fmt.Errorf("openai tts: streaming: %w", err)
	}

	return nil
}
