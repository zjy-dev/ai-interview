package middleware

import (
	"testing"
)

func TestEncryptor_EncryptDecrypt(t *testing.T) {
	// 32 字节密钥 (AES-256)
	// 64 hex chars = 32 bytes for AES-256
	enc, err := NewEncryptor("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("NewEncryptor error: %v", err)
	}

	plaintext := "sk-my-secret-api-key-12345"

	encrypted, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	if encrypted == plaintext {
		t.Error("encrypted should not equal plaintext")
	}

	decrypted, err := enc.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypt(%q) = %q, want %q", encrypted, decrypted, plaintext)
	}
}

func TestEncryptor_DifferentCiphertexts(t *testing.T) {
	enc, _ := NewEncryptor("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	a, _ := enc.Encrypt("hello")
	b, _ := enc.Encrypt("hello")

	// AES-GCM 使用随机 nonce，相同明文应产生不同密文
	if a == b {
		t.Error("two encryptions of same plaintext should produce different ciphertexts")
	}
}

func TestEncryptor_InvalidKey(t *testing.T) {
	_, err := NewEncryptor("0123456789abcdef") // only 8 bytes
	if err == nil {
		t.Error("expected error for short key")
	}
}

func TestEncryptor_EmptyString(t *testing.T) {
	enc, _ := NewEncryptor("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	encrypted, err := enc.Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt empty error: %v", err)
	}

	decrypted, err := enc.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt empty error: %v", err)
	}

	if decrypted != "" {
		t.Errorf("expected empty string, got %q", decrypted)
	}
}
