package biz

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 业务错误定义
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInterviewNotFound  = errors.New("interview not found")
	ErrInterviewEnded     = errors.New("interview already ended")
	ErrUnauthorized       = errors.New("unauthorized")
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
