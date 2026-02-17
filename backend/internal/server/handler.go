package server

import (
	"ai-interview/internal/biz"
	"ai-interview/internal/middleware"
	"ai-interview/internal/service"
	"strconv"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// authHandlerImpl Auth HTTP handler
type authHandlerImpl struct {
	svc *service.AuthService
}

func authHandler(svc *service.AuthService) *authHandlerImpl {
	return &authHandlerImpl{svc: svc}
}

func (h *authHandlerImpl) Register(ctx http.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request"})
	}

	id, token, err := h.svc.Register(ctx, req.Email, req.Password, req.Nickname)
	if err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{
		"id":       id,
		"token":    token,
		"nickname": req.Nickname,
	})
}

func (h *authHandlerImpl) Login(ctx http.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request"})
	}

	id, token, nickname, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return ctx.JSON(401, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{
		"id":       id,
		"token":    token,
		"nickname": nickname,
	})
}

func (h *authHandlerImpl) GetProfile(ctx http.Context) error {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ctx.JSON(401, map[string]string{"error": "unauthorized"})
	}

	user, err := h.svc.GetProfile(ctx, userID)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"nickname":   user.Nickname,
		"created_at": user.CreatedAt,
	})
}

func (h *authHandlerImpl) GetSettings(ctx http.Context) error {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ctx.JSON(401, map[string]string{"error": "unauthorized"})
	}

	settings, err := h.svc.GetSettings(ctx, userID)
	if err != nil {
		return ctx.JSON(200, map[string]any{
			"llm_provider":    "",
			"llm_api_key_set": false,
			"tts_provider":    "",
			"tts_api_key_set": false,
			"tts_enabled":     true,
			"stt_provider":    "browser",
		})
	}

	return ctx.JSON(200, map[string]any{
		"llm_provider":    settings.LLMProvider,
		"llm_api_key_set": settings.LLMAPIKey != "",
		"llm_base_url":    settings.LLMBaseURL,
		"llm_model":       settings.LLMModel,
		"tts_provider":    settings.TTSProvider,
		"tts_api_key_set": settings.TTSAPIKey != "",
		"tts_voice":       settings.TTSVoice,
		"tts_enabled":     settings.TTSEnabled,
		"stt_provider":    settings.STTProvider,
		"stt_api_key_set": settings.STTAPIKey != "",
	})
}

func (h *authHandlerImpl) UpdateSettings(ctx http.Context) error {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ctx.JSON(401, map[string]string{"error": "unauthorized"})
	}

	var req struct {
		LLMProvider string `json:"llm_provider"`
		LLMAPIKey   string `json:"llm_api_key"`
		LLMBaseURL  string `json:"llm_base_url"`
		LLMModel    string `json:"llm_model"`
		TTSProvider string `json:"tts_provider"`
		TTSAPIKey   string `json:"tts_api_key"`
		TTSVoice    string `json:"tts_voice"`
		TTSEnabled  bool   `json:"tts_enabled"`
		STTProvider string `json:"stt_provider"`
		STTAPIKey   string `json:"stt_api_key"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request"})
	}

	// API Keys 加密由 Service 层处理 (AuthService.UpdateSettings)
	settings := &biz.UserSettings{
		UserID:      userID,
		LLMProvider: req.LLMProvider,
		LLMAPIKey:   req.LLMAPIKey,
		LLMBaseURL:  req.LLMBaseURL,
		LLMModel:    req.LLMModel,
		TTSProvider: req.TTSProvider,
		TTSAPIKey:   req.TTSAPIKey,
		TTSVoice:    req.TTSVoice,
		TTSEnabled:  req.TTSEnabled,
		STTProvider: req.STTProvider,
		STTAPIKey:   req.STTAPIKey,
	}

	if err := h.svc.UpdateSettings(ctx, settings); err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{"success": true})
}

// interviewHandlerImpl Interview HTTP handler
type interviewHandlerImpl struct {
	svc *service.InterviewService
}

func interviewHandler(svc *service.InterviewService) *interviewHandlerImpl {
	return &interviewHandlerImpl{svc: svc}
}

func (h *interviewHandlerImpl) Create(ctx http.Context) error {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ctx.JSON(401, map[string]string{"error": "unauthorized"})
	}

	var req struct {
		Title       string `json:"title"`
		Position    string `json:"position"`
		Language    string `json:"language"`
		LLMProvider string `json:"llm_provider"`
		LLMModel    string `json:"llm_model"`
		TTSProvider string `json:"tts_provider"`
		TTSVoice    string `json:"tts_voice"`
		Resume      string `json:"resume"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request"})
	}

	interview := &biz.Interview{
		Title:       req.Title,
		Position:    req.Position,
		Language:    req.Language,
		LLMProvider: req.LLMProvider,
		LLMModel:    req.LLMModel,
		TTSProvider: req.TTSProvider,
		TTSVoice:    req.TTSVoice,
		Resume:      req.Resume,
	}

	created, err := h.svc.CreateInterview(ctx, userID, interview)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{
		"id":            created.ID,
		"title":         created.Title,
		"status":        created.Status,
		"websocket_url": "/api/v1/ws/interview/" + strconv.FormatInt(created.ID, 10),
		"created_at":    created.CreatedAt,
	})
}

func (h *interviewHandlerImpl) List(ctx http.Context) error {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ctx.JSON(401, map[string]string{"error": "unauthorized"})
	}

	page, _ := strconv.Atoi(ctx.Query().Get("page"))
	pageSize, _ := strconv.Atoi(ctx.Query().Get("page_size"))

	interviews, total, err := h.svc.ListInterviews(ctx, userID, page, pageSize)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	items := make([]map[string]any, 0, len(interviews))
	for _, i := range interviews {
		items = append(items, map[string]any{
			"id":         i.ID,
			"title":      i.Title,
			"position":   i.Position,
			"status":     i.Status,
			"created_at": i.CreatedAt,
		})
	}

	return ctx.JSON(200, map[string]any{
		"interviews": items,
		"total":      total,
	})
}

func (h *interviewHandlerImpl) Get(ctx http.Context) error {
	id, _ := strconv.ParseInt(ctx.Vars().Get("id"), 10, 64)

	interview, messages, err := h.svc.GetInterview(ctx, id)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": err.Error()})
	}

	msgItems := make([]map[string]any, 0, len(messages))
	for _, m := range messages {
		msgItems = append(msgItems, map[string]any{
			"id":         m.ID,
			"role":       m.Role,
			"content":    m.Content,
			"created_at": m.CreatedAt,
		})
	}

	return ctx.JSON(200, map[string]any{
		"id":         interview.ID,
		"title":      interview.Title,
		"position":   interview.Position,
		"status":     interview.Status,
		"language":   interview.Language,
		"messages":   msgItems,
		"created_at": interview.CreatedAt,
	})
}

func (h *interviewHandlerImpl) SendMessage(ctx http.Context) error {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ctx.JSON(401, map[string]string{"error": "unauthorized"})
	}

	interviewID, _ := strconv.ParseInt(ctx.Vars().Get("id"), 10, 64)

	var req struct {
		Content string `json:"content"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request"})
	}

	userMsg, assistantContent, err := h.svc.SendMessage(ctx, interviewID, req.Content, userID)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{
		"user_message": map[string]any{
			"id":      userMsg.ID,
			"role":    userMsg.Role,
			"content": userMsg.Content,
		},
		"assistant_message": map[string]any{
			"role":    "assistant",
			"content": assistantContent,
		},
	})
}

func (h *interviewHandlerImpl) End(ctx http.Context) error {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return ctx.JSON(401, map[string]string{"error": "unauthorized"})
	}

	id, _ := strconv.ParseInt(ctx.Vars().Get("id"), 10, 64)

	eval, err := h.svc.EndInterview(ctx, id, userID)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]any{
		"status":             "completed",
		"evaluation_summary": eval.Summary,
	})
}

func (h *interviewHandlerImpl) GetEvaluation(ctx http.Context) error {
	interviewID, _ := strconv.ParseInt(ctx.Vars().Get("id"), 10, 64)

	eval, err := h.svc.GetEvaluation(ctx, interviewID)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": err.Error()})
	}

	categories := make([]map[string]any, 0, len(eval.Categories))
	for _, c := range eval.Categories {
		categories = append(categories, map[string]any{
			"category": c.Category,
			"score":    c.Score,
			"comment":  c.Comment,
		})
	}

	return ctx.JSON(200, map[string]any{
		"id":            eval.ID,
		"interview_id":  eval.InterviewID,
		"overall_score": eval.OverallScore,
		"summary":       eval.Summary,
		"categories":    categories,
		"strengths":     eval.Strengths,
		"weaknesses":    eval.Weaknesses,
		"suggestions":   eval.Suggestions,
		"created_at":    eval.CreatedAt,
	})
}

// WebSocket 处理面试实时交互
func (h *interviewHandlerImpl) WebSocket(ctx http.Context) error {
	// TODO: 实现 WebSocket 面试实时交互
	// 1. 升级 HTTP 连接为 WebSocket
	// 2. 接收用户文本/音频消息
	// 3. 调用 LLM 流式回复
	// 4. 调用 TTS 将回复转为音频流
	// 5. 通过 WebSocket 推送音频 chunk + 文本
	return ctx.JSON(501, map[string]string{"error": "websocket not yet implemented"})
}
