package service

import (
	"ai-interview/internal/biz"
	"ai-interview/internal/middleware"
	"context"
)

// InterviewService 面试服务
type InterviewService struct {
	interviewUC *biz.InterviewUsecase
	userUC      *biz.UserUsecase
	encryptor   *middleware.Encryptor
}

// NewInterviewService 创建面试服务
func NewInterviewService(interviewUC *biz.InterviewUsecase, userUC *biz.UserUsecase, encryptor *middleware.Encryptor) *InterviewService {
	return &InterviewService{interviewUC: interviewUC, userUC: userUC, encryptor: encryptor}
}

// CreateInterview 创建面试会话
func (s *InterviewService) CreateInterview(ctx context.Context, userID int64, interview *biz.Interview) (*biz.Interview, error) {
	return s.interviewUC.CreateInterview(ctx, userID, interview)
}

// GetInterview 获取面试详情
func (s *InterviewService) GetInterview(ctx context.Context, id int64) (*biz.Interview, []*biz.InterviewMessage, error) {
	return s.interviewUC.GetInterview(ctx, id)
}

// ListInterviews 列出面试记录
func (s *InterviewService) ListInterviews(ctx context.Context, userID int64, page, pageSize int) ([]*biz.Interview, int, error) {
	return s.interviewUC.ListInterviews(ctx, userID, page, pageSize)
}

// SendMessage 发送消息
func (s *InterviewService) SendMessage(ctx context.Context, interviewID int64, content string, userID int64) (*biz.InterviewMessage, string, error) {
	settings, _ := s.userUC.GetSettings(ctx, userID)
	s.decryptSettings(settings)
	return s.interviewUC.SendMessage(ctx, interviewID, content, settings)
}

// EndInterview 结束面试
func (s *InterviewService) EndInterview(ctx context.Context, id int64, userID int64) (*biz.Evaluation, error) {
	settings, _ := s.userUC.GetSettings(ctx, userID)
	s.decryptSettings(settings)
	return s.interviewUC.EndInterview(ctx, id, settings)
}

// GetEvaluation 获取评估
func (s *InterviewService) GetEvaluation(ctx context.Context, interviewID int64) (*biz.Evaluation, error) {
	return s.interviewUC.GetEvaluation(ctx, interviewID)
}

// GetUserSettings 获取用户设置 (供 WebSocket handler 使用)
func (s *InterviewService) GetUserSettings(ctx context.Context, userID int64) (*biz.UserSettings, error) {
	settings, err := s.userUC.GetSettings(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.decryptSettings(settings)
	return settings, nil
}

// decryptSettings 解密 settings 中的 API Keys
func (s *InterviewService) decryptSettings(settings *biz.UserSettings) {
	if s.encryptor == nil || settings == nil {
		return
	}
	if settings.LLMAPIKey != "" {
		if decrypted, err := s.encryptor.Decrypt(settings.LLMAPIKey); err == nil {
			settings.LLMAPIKey = decrypted
		}
	}
	if settings.TTSAPIKey != "" {
		if decrypted, err := s.encryptor.Decrypt(settings.TTSAPIKey); err == nil {
			settings.TTSAPIKey = decrypted
		}
	}
	if settings.STTAPIKey != "" {
		if decrypted, err := s.encryptor.Decrypt(settings.STTAPIKey); err == nil {
			settings.STTAPIKey = decrypted
		}
	}
}
