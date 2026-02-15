package server

import (
	"testing"
)

func TestIsSentenceEnd(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Hello world.", true},
		{"你好！", true},
		{"What?", true},
		{"这是一句话。", true},
		{"继续说", false},
		{"Hello", false},
		{"", false},
		{"  ", false},
		{"多行\n", false}, // trimmed then checked, newline lost
		{"end；", true},
	}

	for _, tc := range tests {
		got := isSentenceEnd(tc.input)
		if got != tc.expected {
			t.Errorf("isSentenceEnd(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}
