package service

import (
	"ai-interview/internal/biz"
	"context"
)

// InterviewService 面试服务
type InterviewService struct {
	interviewUC *biz.InterviewUsecase
	userUC      *biz.UserUsecase
}

// NewInterviewService 创建面试服务
func NewInterviewService(interviewUC *biz.InterviewUsecase, userUC *biz.UserUsecase) *InterviewService {
	return &InterviewService{interviewUC: interviewUC, userUC: userUC}
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
	return s.interviewUC.SendMessage(ctx, interviewID, content, settings)
}

// EndInterview 结束面试
func (s *InterviewService) EndInterview(ctx context.Context, id int64, userID int64) (*biz.Evaluation, error) {
	settings, _ := s.userUC.GetSettings(ctx, userID)
	return s.interviewUC.EndInterview(ctx, id, settings)
}

// GetEvaluation 获取评估
func (s *InterviewService) GetEvaluation(ctx context.Context, interviewID int64) (*biz.Evaluation, error) {
	return s.interviewUC.GetEvaluation(ctx, interviewID)
}

// GetUserSettings 获取用户设置 (供 WebSocket handler 使用)
func (s *InterviewService) GetUserSettings(ctx context.Context, userID int64) (*biz.UserSettings, error) {
	return s.userUC.GetSettings(ctx, userID)
}
