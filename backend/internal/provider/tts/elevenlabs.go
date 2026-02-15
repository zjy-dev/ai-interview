package tts

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// ElevenLabsProvider 实现 ElevenLabs TTS API
type ElevenLabsProvider struct {
	defaultBaseURL string
	httpClient     *http.Client
}

func NewElevenLabsProvider() *ElevenLabsProvider {
	return &ElevenLabsProvider{
		defaultBaseURL: "https://api.elevenlabs.io",
		httpClient:     &http.Client{},
	}
}

func (p *ElevenLabsProvider) Name() string {
	return "elevenlabs"
}

func (p *ElevenLabsProvider) Synthesize(ctx context.Context, req *Request, w io.Writer) error {
	if req.APIKey == "" {
		return fmt.Errorf("elevenlabs tts: api key is required")
	}

	voiceID := req.Voice
	if voiceID == "" {
		voiceID = "21m00Tcm4TlvDq8ikWAM" // Rachel (默认)
	}

	baseURL := p.defaultBaseURL
	if req.BaseURL != "" {
		baseURL = req.BaseURL
	}

	url := fmt.Sprintf("%s/v1/text-to-speech/%s/stream", baseURL, voiceID)

	outputFormat := "pcm_24000"
	switch req.Format {
	case "mp3":
		outputFormat = "mp3_44100_128"
	case "opus":
		outputFormat = "opus_24000"
	}
	url += "?output_format=" + outputFormat

	body := fmt.Sprintf(`{
		"text": %q,
		"model_id": "eleven_flash_v2_5",
		"voice_settings": {
			"stability": 0.5,
			"similarity_boost": 0.75,
			"style": 0.0,
			"use_speaker_boost": true
		}
	}`, req.Text)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, io.NopCloser(
		io.Reader(stringReader(body)),
	))
	if err != nil {
		return fmt.Errorf("elevenlabs tts: creating request: %w", err)
	}

	httpReq.Header.Set("xi-api-key", req.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "audio/mpeg")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("elevenlabs tts: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("elevenlabs tts: status %d: %s", resp.StatusCode, string(respBody))
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("elevenlabs tts: streaming: %w", err)
	}

	return nil
}

type stringReaderImpl struct {
	s string
	i int
}

func stringReader(s string) io.Reader {
	return &stringReaderImpl{s: s}
}

func (r *stringReaderImpl) Read(p []byte) (n int, err error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	n = copy(p, r.s[r.i:])
	r.i += n
	return n, nil
}
