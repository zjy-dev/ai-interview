package tts

import (
	"context"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

// FishAudioProvider 实现 Fish Audio TTS (OpenAI 兼容 API)
type FishAudioProvider struct {
	defaultBaseURL string
}

func NewFishAudioProvider() *FishAudioProvider {
	return &FishAudioProvider{
		defaultBaseURL: "https://api.fish.audio/v1",
	}
}

func (p *FishAudioProvider) Name() string {
	return "fishaudio"
}

func (p *FishAudioProvider) Synthesize(ctx context.Context, req *Request, w io.Writer) error {
	if req.APIKey == "" {
		return fmt.Errorf("fishaudio tts: api key is required")
	}

	baseURL := p.defaultBaseURL
	if req.BaseURL != "" {
		baseURL = req.BaseURL
	}

	// Fish Audio 提供 OpenAI 兼容的 /v1/audio/speech 接口
	config := openai.DefaultConfig(req.APIKey)
	config.BaseURL = baseURL
	client := openai.NewClientWithConfig(config)

	voice := openai.SpeechVoice(req.Voice)

	format := openai.SpeechResponseFormatPcm
	switch req.Format {
	case "mp3":
		format = openai.SpeechResponseFormatMp3
	case "opus":
		format = openai.SpeechResponseFormatOpus
	case "wav":
		format = openai.SpeechResponseFormatWav
	}

	speed := 1.0
	if req.Speed > 0 {
		speed = req.Speed
	}

	speechReq := openai.CreateSpeechRequest{
		Model:          openai.TTSModel1,
		Input:          req.Text,
		Voice:          voice,
		ResponseFormat: format,
		Speed:          speed,
	}

	resp, err := client.CreateSpeech(ctx, speechReq)
	if err != nil {
		return fmt.Errorf("fishaudio tts: %w", err)
	}
	defer resp.Close()

	if _, err := io.Copy(w, resp); err != nil {
		return fmt.Errorf("fishaudio tts: streaming: %w", err)
	}

	return nil
}
