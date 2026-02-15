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
}

// NewAuthService 创建认证服务
func NewAuthService(uc *biz.UserUsecase, jwtHelper *middleware.JWTHelper) *AuthService {
	return &AuthService{uc: uc, jwtHelper: jwtHelper}
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
	return s.uc.UpdateSettings(ctx, settings)
}

// GetSettings 获取用户设置
func (s *AuthService) GetSettings(ctx context.Context, userID int64) (*biz.UserSettings, error) {
	return s.uc.GetSettings(ctx, userID)
}
