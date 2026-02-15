package middleware

import (
	"testing"
	"time"
)

func TestJWTHelper_GenerateAndValidate(t *testing.T) {
	helper := NewJWTHelper("test-secret-key-32chars-long!!!!!", 24*time.Hour)

	token, err := helper.GenerateToken(42)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	if token == "" {
		t.Fatal("token should not be empty")
	}

	userID, err := helper.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}
	if userID != 42 {
		t.Errorf("expected userID=42, got %d", userID)
	}
}

func TestJWTHelper_InvalidToken(t *testing.T) {
	helper := NewJWTHelper("test-secret-key-32chars-long!!!!!", 24*time.Hour)

	_, err := helper.ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestJWTHelper_ExpiredToken(t *testing.T) {
	helper := NewJWTHelper("test-secret-key-32chars-long!!!!!", -1*time.Hour) // 已过期

	token, _ := helper.GenerateToken(1)
	_, err := helper.ValidateToken(token)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		header string
		expect string
	}{
		{"Bearer abc123", "abc123"},
		{"bearer abc123", ""},
		{"abc123", ""},
		{"", ""},
		{"Bearer ", ""},
	}

	for _, tc := range tests {
		got := ExtractTokenFromHeader(tc.header)
		if got != tc.expect {
			t.Errorf("ExtractTokenFromHeader(%q) = %q, want %q", tc.header, got, tc.expect)
		}
	}
}
