package service

import (
	"ai-interview/internal/biz"
	"ai-interview/internal/middleware"
	"context"
	"fmt"
)

// AuthService 认证服务
type AuthService struct {
	uc        *biz.UserUsecase
	jwtHelper *middleware.JWTHelper
	encryptor *middleware.Encryptor
}

// NewAuthService 创建认证服务
func NewAuthService(uc *biz.UserUsecase, jwtHelper *middleware.JWTHelper, encryptor *middleware.Encryptor) *AuthService {
	return &AuthService{uc: uc, jwtHelper: jwtHelper, encryptor: encryptor}
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, email, password, nickname string) (int64, string, error) {
	user, err := s.uc.Register(ctx, email, password, nickname)
	if err != nil {
		return 0, "", err
	}

	token, err := s.jwtHelper.GenerateToken(user.ID)
	if err != nil {
		return 0, "", fmt.Errorf("generate token: %w", err)
	}

	return user.ID, token, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, email, password string) (int64, string, string, error) {
	user, err := s.uc.Login(ctx, email, password)
	if err != nil {
		return 0, "", "", err
	}

	token, err := s.jwtHelper.GenerateToken(user.ID)
	if err != nil {
		return 0, "", "", fmt.Errorf("generate token: %w", err)
	}

	return user.ID, token, user.Nickname, nil
}

// GetProfile 获取用户信息
func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*biz.User, error) {
	return s.uc.GetProfile(ctx, userID)
}

// UpdateSettings 更新用户设置
func (s *AuthService) UpdateSettings(ctx context.Context, settings *biz.UserSettings) error {
	// 加密非空 API Keys
	if s.encryptor != nil {
		if settings.LLMAPIKey != "" {
			encrypted, err := s.encryptor.Encrypt(settings.LLMAPIKey)
			if err != nil {
				return fmt.Errorf("encrypt llm api key: %w", err)
			}
			settings.LLMAPIKey = encrypted
		}
		if settings.TTSAPIKey != "" {
			encrypted, err := s.encryptor.Encrypt(settings.TTSAPIKey)
			if err != nil {
				return fmt.Errorf("encrypt tts api key: %w", err)
			}
			settings.TTSAPIKey = encrypted
		}
		if settings.STTAPIKey != "" {
			encrypted, err := s.encryptor.Encrypt(settings.STTAPIKey)
			if err != nil {
				return fmt.Errorf("encrypt stt api key: %w", err)
			}
			settings.STTAPIKey = encrypted
		}
	}
	return s.uc.UpdateSettings(ctx, settings)
}

// GetSettings 获取用户设置
func (s *AuthService) GetSettings(ctx context.Context, userID int64) (*biz.UserSettings, error) {
	return s.uc.GetSettings(ctx, userID)
}
