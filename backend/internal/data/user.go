package data

import (
	"ai-interview/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo 创建用户仓储
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	result, err := r.data.db.ExecContext(ctx,
		"INSERT INTO users (email, nickname, password_hash) VALUES (?, ?, ?)",
		user.Email, user.Nickname, user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id

	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*biz.User, error) {
	user := &biz.User{}
	err := r.data.db.QueryRowContext(ctx,
		"SELECT id, email, nickname, password_hash, created_at, updated_at FROM users WHERE id = ?", id,
	).Scan(&user.ID, &user.Email, &user.Nickname, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*biz.User, error) {
	user := &biz.User{}
	err := r.data.db.QueryRowContext(ctx,
		"SELECT id, email, nickname, password_hash, created_at, updated_at FROM users WHERE email = ?", email,
	).Scan(&user.ID, &user.Email, &user.Nickname, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) UpdateSettings(ctx context.Context, settings *biz.UserSettings) error {
	_, err := r.data.db.ExecContext(ctx,
		`INSERT INTO user_settings (user_id, llm_provider, llm_api_key, llm_base_url, llm_model,
			tts_provider, tts_api_key, tts_voice, tts_enabled, stt_provider, stt_api_key)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			llm_provider = VALUES(llm_provider),
			llm_api_key = VALUES(llm_api_key),
			llm_base_url = VALUES(llm_base_url),
			llm_model = VALUES(llm_model),
			tts_provider = VALUES(tts_provider),
			tts_api_key = VALUES(tts_api_key),
			tts_voice = VALUES(tts_voice),
			tts_enabled = VALUES(tts_enabled),
			stt_provider = VALUES(stt_provider),
			stt_api_key = VALUES(stt_api_key)`,
		settings.UserID, settings.LLMProvider, settings.LLMAPIKey, settings.LLMBaseURL, settings.LLMModel,
		settings.TTSProvider, settings.TTSAPIKey, settings.TTSVoice, settings.TTSEnabled,
		settings.STTProvider, settings.STTAPIKey,
	)
	return err
}

func (r *userRepo) GetSettings(ctx context.Context, userID int64) (*biz.UserSettings, error) {
	s := &biz.UserSettings{}
	err := r.data.db.QueryRowContext(ctx,
		`SELECT user_id, llm_provider, llm_api_key, llm_base_url, llm_model,
			tts_provider, tts_api_key, tts_voice, tts_enabled, stt_provider, stt_api_key
		FROM user_settings WHERE user_id = ?`, userID,
	).Scan(&s.UserID, &s.LLMProvider, &s.LLMAPIKey, &s.LLMBaseURL, &s.LLMModel,
		&s.TTSProvider, &s.TTSAPIKey, &s.TTSVoice, &s.TTSEnabled, &s.STTProvider, &s.STTAPIKey,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}
