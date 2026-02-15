package biz

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	pwd := "testpassword123"
	hash, err := hashPassword(pwd)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if hash == pwd {
		t.Fatal("hash should not equal plaintext")
	}
}

func TestCheckPassword(t *testing.T) {
	pwd := "testpassword123"
	hash, _ := hashPassword(pwd)

	if !checkPassword(pwd, hash) {
		t.Error("checkPassword should return true for correct password")
	}
	if checkPassword("wrongpassword", hash) {
		t.Error("checkPassword should return false for wrong password")
	}
}

func TestCheckPasswordEmptyHash(t *testing.T) {
	if checkPassword("anypassword", "") {
		t.Error("checkPassword should return false for empty hash")
	}
}
