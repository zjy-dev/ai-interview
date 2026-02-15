package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// User 是用户领域模型
type User struct {
	ID           int64
	Email        string
	Nickname     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserSettings 是用户配置
type UserSettings struct {
	UserID      int64
	LLMProvider string
	LLMAPIKey   string // 已加密
	LLMBaseURL  string
	LLMModel    string
	TTSProvider string
	TTSAPIKey   string // 已加密
	TTSVoice    string
	TTSEnabled  bool
	STTProvider string
	STTAPIKey   string // 已加密
}

// UserRepo 用户仓储接口 (由 data 层实现)
type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdateSettings(ctx context.Context, settings *UserSettings) error
	GetSettings(ctx context.Context, userID int64) (*UserSettings, error)
}

// UserUsecase 用户业务逻辑
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase 创建用户 UseCase
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Register 用户注册
func (uc *UserUsecase) Register(ctx context.Context, email, password, nickname string) (*User, error) {
	// 检查邮箱是否已注册
	existing, _ := uc.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        email,
		Nickname:     nickname,
		PasswordHash: hash,
	}

	return uc.repo.Create(ctx, user)
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !checkPassword(password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// GetProfile 获取用户信息
func (uc *UserUsecase) GetProfile(ctx context.Context, userID int64) (*User, error) {
	return uc.repo.GetByID(ctx, userID)
}

// UpdateSettings 更新用户设置
func (uc *UserUsecase) UpdateSettings(ctx context.Context, settings *UserSettings) error {
	return uc.repo.UpdateSettings(ctx, settings)
}

// GetSettings 获取用户设置
func (uc *UserUsecase) GetSettings(ctx context.Context, userID int64) (*UserSettings, error) {
	return uc.repo.GetSettings(ctx, userID)
}
