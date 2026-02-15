package biz

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-interview/internal/provider/llm"
	"ai-interview/internal/provider/stt"
	"ai-interview/internal/provider/tts"

	"github.com/go-kratos/kratos/v2/log"
)

// Interview 面试领域模型
type Interview struct {
	ID          int64
	UserID      int64
	Title       string
	Position    string
	Status      string // pending, in_progress, completed
	Language    string
	LLMProvider string
	LLMModel    string
	TTSProvider string
	TTSVoice    string
	Resume      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// InterviewMessage 面试消息
type InterviewMessage struct {
	ID          int64
	InterviewID int64
	Role        string // system, user, assistant
	Content     string
	CreatedAt   time.Time
}

// Evaluation 面试评估
type Evaluation struct {
	ID           int64
	InterviewID  int64
	OverallScore int32
	Summary      string
	Categories   []CategoryScore
	Strengths    string
	Weaknesses   string
	Suggestions  string
	CreatedAt    time.Time
}

// CategoryScore 评估维度分数
type CategoryScore struct {
	Category string
	Score    int32
	Comment  string
}

// InterviewRepo 面试仓储接口
type InterviewRepo interface {
	Create(ctx context.Context, interview *Interview) (*Interview, error)
	GetByID(ctx context.Context, id int64) (*Interview, error)
	ListByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*Interview, int, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	CreateMessage(ctx context.Context, msg *InterviewMessage) (*InterviewMessage, error)
	ListMessages(ctx context.Context, interviewID int64) ([]*InterviewMessage, error)
	CreateEvaluation(ctx context.Context, eval *Evaluation) (*Evaluation, error)
	GetEvaluation(ctx context.Context, interviewID int64) (*Evaluation, error)
}

// InterviewUsecase 面试业务逻辑
type InterviewUsecase struct {
	repo        InterviewRepo
	userRepo    UserRepo
	llmRegistry *llm.Registry
	ttsRegistry *tts.Registry
	sttRegistry *stt.Registry
	log         *log.Helper
}

// NewInterviewUsecase 创建面试 UseCase
func NewInterviewUsecase(
	repo InterviewRepo,
	userRepo UserRepo,
	llmRegistry *llm.Registry,
	ttsRegistry *tts.Registry,
	sttRegistry *stt.Registry,
	logger log.Logger,
) *InterviewUsecase {
	return &InterviewUsecase{
		repo:        repo,
		userRepo:    userRepo,
		llmRegistry: llmRegistry,
		ttsRegistry: ttsRegistry,
		sttRegistry: sttRegistry,
		log:         log.NewHelper(logger),
	}
}

// CreateInterview 创建面试会话
func (uc *InterviewUsecase) CreateInterview(ctx context.Context, userID int64, interview *Interview) (*Interview, error) {
	interview.UserID = userID
	interview.Status = "pending"

	created, err := uc.repo.Create(ctx, interview)
	if err != nil {
		return nil, fmt.Errorf("create interview: %w", err)
	}

	return created, nil
}

// GetInterview 获取面试详情
func (uc *InterviewUsecase) GetInterview(ctx context.Context, id int64) (*Interview, []*InterviewMessage, error) {
	interview, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, ErrInterviewNotFound
	}

	messages, err := uc.repo.ListMessages(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("list messages: %w", err)
	}

	return interview, messages, nil
}

// ListInterviews 列出用户的面试记录
func (uc *InterviewUsecase) ListInterviews(ctx context.Context, userID int64, page, pageSize int) ([]*Interview, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 20
	}
	return uc.repo.ListByUserID(ctx, userID, page, pageSize)
}

// SendMessage 处理用户消息并生成 AI 回复
func (uc *InterviewUsecase) SendMessage(ctx context.Context, interviewID int64, userContent string, settings *UserSettings) (*InterviewMessage, string, error) {
	interview, err := uc.repo.GetByID(ctx, interviewID)
	if err != nil {
		return nil, "", ErrInterviewNotFound
	}
	if interview.Status == "completed" {
		return nil, "", ErrInterviewEnded
	}

	// 更新状态为进行中
	if interview.Status == "pending" {
		_ = uc.repo.UpdateStatus(ctx, interviewID, "in_progress")
	}

	// 保存用户消息
	userMsg := &InterviewMessage{
		InterviewID: interviewID,
		Role:        "user",
		Content:     userContent,
	}
	userMsg, err = uc.repo.CreateMessage(ctx, userMsg)
	if err != nil {
		return nil, "", fmt.Errorf("save user message: %w", err)
	}

	// 获取对话历史
	messages, err := uc.repo.ListMessages(ctx, interviewID)
	if err != nil {
		return nil, "", fmt.Errorf("get history: %w", err)
	}

	// 构建 LLM 请求
	llmMessages := uc.buildLLMMessages(interview, messages)

	// 获取 LLM Provider
	providerName := interview.LLMProvider
	if providerName == "" && settings != nil {
		providerName = settings.LLMProvider
	}
	if providerName == "" {
		providerName = "openai"
	}

	provider, err := uc.llmRegistry.Get(providerName)
	if err != nil {
		return nil, "", fmt.Errorf("get llm provider: %w", err)
	}

	apiKey := ""
	baseURL := ""
	if settings != nil {
		apiKey = settings.LLMAPIKey
		baseURL = settings.LLMBaseURL
	}

	// 流式调用 LLM
	stream, err := provider.ChatStream(ctx, &llm.ChatRequest{
		Messages:    llmMessages,
		Model:       interview.LLMModel,
		MaxTokens:   4096,
		Temperature: 0.7,
		APIKey:      apiKey,
		BaseURL:     baseURL,
	})
	if err != nil {
		return nil, "", fmt.Errorf("llm chat: %w", err)
	}

	// 收集完整回复
	var sb strings.Builder
	for event := range stream {
		if event.Err != nil {
			return nil, "", fmt.Errorf("llm stream: %w", event.Err)
		}
		if event.Content != "" {
			sb.WriteString(event.Content)
		}
	}

	assistantContent := sb.String()

	// 保存助手消息
	assistantMsg := &InterviewMessage{
		InterviewID: interviewID,
		Role:        "assistant",
		Content:     assistantContent,
	}
	_, err = uc.repo.CreateMessage(ctx, assistantMsg)
	if err != nil {
		return nil, assistantContent, fmt.Errorf("save assistant message: %w", err)
	}

	return userMsg, assistantContent, nil
}

// StreamMessage 流式处理消息，返回 LLM 文本流 channel
func (uc *InterviewUsecase) StreamMessage(ctx context.Context, interviewID int64, userContent string, settings *UserSettings) (*InterviewMessage, <-chan llm.StreamEvent, error) {
	interview, err := uc.repo.GetByID(ctx, interviewID)
	if err != nil {
		return nil, nil, ErrInterviewNotFound
	}
	if interview.Status == "completed" {
		return nil, nil, ErrInterviewEnded
	}

	if interview.Status == "pending" {
		_ = uc.repo.UpdateStatus(ctx, interviewID, "in_progress")
	}

	// 保存用户消息
	userMsg := &InterviewMessage{
		InterviewID: interviewID,
		Role:        "user",
		Content:     userContent,
	}
	userMsg, err = uc.repo.CreateMessage(ctx, userMsg)
	if err != nil {
		return nil, nil, fmt.Errorf("save user message: %w", err)
	}

	messages, err := uc.repo.ListMessages(ctx, interviewID)
	if err != nil {
		return nil, nil, fmt.Errorf("get history: %w", err)
	}

	llmMessages := uc.buildLLMMessages(interview, messages)

	providerName := interview.LLMProvider
	if providerName == "" && settings != nil {
		providerName = settings.LLMProvider
	}
	if providerName == "" {
		providerName = "openai"
	}

	provider, err := uc.llmRegistry.Get(providerName)
	if err != nil {
		return nil, nil, fmt.Errorf("get llm provider: %w", err)
	}

	apiKey := ""
	baseURL := ""
	if settings != nil {
		apiKey = settings.LLMAPIKey
		baseURL = settings.LLMBaseURL
	}

	stream, err := provider.ChatStream(ctx, &llm.ChatRequest{
		Messages:    llmMessages,
		Model:       interview.LLMModel,
		MaxTokens:   4096,
		Temperature: 0.7,
		APIKey:      apiKey,
		BaseURL:     baseURL,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("llm chat: %w", err)
	}

	return userMsg, stream, nil
}

// EndInterview 结束面试并生成评估
func (uc *InterviewUsecase) EndInterview(ctx context.Context, id int64, settings *UserSettings) (*Evaluation, error) {
	interview, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrInterviewNotFound
	}

	if interview.Status == "completed" {
		// 已结束，直接返回已有评估
		return uc.repo.GetEvaluation(ctx, id)
	}

	_ = uc.repo.UpdateStatus(ctx, id, "completed")

	messages, err := uc.repo.ListMessages(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}

	// 构建评估请求
	evalPrompt := uc.buildEvaluationPrompt(interview, messages)

	providerName := interview.LLMProvider
	if providerName == "" && settings != nil {
		providerName = settings.LLMProvider
	}
	if providerName == "" {
		providerName = "openai"
	}

	provider, err := uc.llmRegistry.Get(providerName)
	if err != nil {
		return nil, fmt.Errorf("get llm provider: %w", err)
	}

	apiKey := ""
	baseURL := ""
	if settings != nil {
		apiKey = settings.LLMAPIKey
		baseURL = settings.LLMBaseURL
	}

	stream, err := provider.ChatStream(ctx, &llm.ChatRequest{
		Messages:    evalPrompt,
		Model:       interview.LLMModel,
		MaxTokens:   4096,
		Temperature: 0.3,
		APIKey:      apiKey,
		BaseURL:     baseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("llm eval: %w", err)
	}

	var sb strings.Builder
	for event := range stream {
		if event.Err != nil {
			return nil, fmt.Errorf("llm eval stream: %w", event.Err)
		}
		if event.Content != "" {
			sb.WriteString(event.Content)
		}
	}

	// TODO: 解析 LLM 返回的 JSON 评估结果
	eval := &Evaluation{
		InterviewID:  id,
		OverallScore: 70,
		Summary:      sb.String(),
	}

	eval, err = uc.repo.CreateEvaluation(ctx, eval)
	if err != nil {
		return nil, fmt.Errorf("save evaluation: %w", err)
	}

	return eval, nil
}

// GetEvaluation 获取面试评估
func (uc *InterviewUsecase) GetEvaluation(ctx context.Context, interviewID int64) (*Evaluation, error) {
	return uc.repo.GetEvaluation(ctx, interviewID)
}

// GetTTSProvider 获取 TTS Provider
func (uc *InterviewUsecase) GetTTSProvider(providerName string) (tts.Provider, error) {
	return uc.ttsRegistry.Get(providerName)
}

// GetSTTProvider 获取 STT Provider
func (uc *InterviewUsecase) GetSTTProvider(providerName string) (stt.Provider, error) {
	return uc.sttRegistry.Get(providerName)
}

func (uc *InterviewUsecase) buildLLMMessages(interview *Interview, messages []*InterviewMessage) []llm.Message {
	result := make([]llm.Message, 0, len(messages)+1)

	// System prompt
	systemPrompt := fmt.Sprintf(
		"你是一位专业的面试官，正在面试%s岗位的候选人。\n"+
			"请根据岗位要求逐一提出面试问题，每次只问一个问题。\n"+
			"等候选人回答后，进行简短评价并提出下一个问题。\n"+
			"面试语言：%s",
		interview.Position, interview.Language,
	)
	if interview.Resume != "" {
		systemPrompt += "\n\n候选人简历：\n" + interview.Resume
	}

	result = append(result, llm.Message{Role: "system", Content: systemPrompt})

	for _, msg := range messages {
		if msg.Role == "system" {
			continue
		}
		result = append(result, llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return result
}

func (uc *InterviewUsecase) buildEvaluationPrompt(interview *Interview, messages []*InterviewMessage) []llm.Message {
	var sb strings.Builder
	sb.WriteString("请根据以下面试记录，给出综合评估。\n\n")
	sb.WriteString(fmt.Sprintf("面试岗位：%s\n\n", interview.Position))
	sb.WriteString("=== 面试记录 ===\n\n")

	for _, msg := range messages {
		if msg.Role == "system" {
			continue
		}
		role := "面试官"
		if msg.Role == "user" {
			role = "候选人"
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n\n", role, msg.Content))
	}

	sb.WriteString("=== 评估要求 ===\n")
	sb.WriteString("请从以下维度进行评分 (0-100) 并给出评语：\n")
	sb.WriteString("1. 技术能力\n2. 沟通表达\n3. 逻辑思维\n4. 问题解决\n5. 学习潜力\n\n")
	sb.WriteString("同时给出：\n- 总体评分 (0-100)\n- 总结\n- 优势\n- 不足\n- 改进建议\n")

	return []llm.Message{
		{Role: "system", Content: "你是一位资深面试评估专家。请根据面试记录给出客观、详细的评估报告。"},
		{Role: "user", Content: sb.String()},
	}
}
