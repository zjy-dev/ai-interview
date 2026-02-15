package tts

import (
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()

	p := NewOpenAIProvider()
	r.Register(p)

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

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()
	r.Register(NewOpenAIProvider())
	r.Register(NewEdgeTTSProvider())

	names := r.List()
	if len(names) != 2 {
		t.Errorf("expected 2 providers, got %d", len(names))
	}
}

func TestProviderNames(t *testing.T) {
	tests := []struct {
		provider Provider
		name     string
	}{
		{NewOpenAIProvider(), "openai"},
		{NewFishAudioProvider(), "fishaudio"},
		{NewElevenLabsProvider(), "elevenlabs"},
		{NewEdgeTTSProvider(), "edgetts"},
	}

	for _, tc := range tests {
		if tc.provider.Name() != tc.name {
			t.Errorf("expected name %q, got %q", tc.name, tc.provider.Name())
		}
	}
}
