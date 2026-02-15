package llm

import (
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	r.Register(NewOpenAIProvider())

	got, err := r.Get("openai")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if got.Name() != "openai" {
		t.Errorf("expected name 'openai', got %q", got.Name())
	}
}

func TestRegistry_GetNotFound(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent provider")
	}
}

func TestProviderNames(t *testing.T) {
	tests := []struct {
		provider Provider
		name     string
	}{
		{NewOpenAIProvider(), "openai"},
		{NewDeepSeekProvider(), "deepseek"},
		{NewAnthropicProvider(), "anthropic"},
		{NewGeminiProvider(), "gemini"},
	}

	for _, tc := range tests {
		if tc.provider.Name() != tc.name {
			t.Errorf("expected name %q, got %q", tc.name, tc.provider.Name())
		}
	}
}
