package server

import (
	"ai-interview/internal/biz"
	"ai-interview/internal/middleware"
	"ai-interview/internal/provider/tts"
	"ai-interview/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"nhooyr.io/websocket"
)

// wsMessage WebSocket 消息格式
type wsMessage struct {
	Type string `json:"type"` // "text", "audio", "end", "ping"
	Data string `json:"data,omitempty"`
}

// wsResponse WebSocket 响应
type wsResponse struct {
	Type string `json:"type"` // "text_start", "text_delta", "text_end", "audio", "status", "error", "evaluation"
	Data any    `json:"data,omitempty"`
}

// WebSocketHandler 面试实时交互 WebSocket 处理器
type WebSocketHandler struct {
	interviewSvc *service.InterviewService
	interviewUC  *biz.InterviewUsecase
	logger       *log.Helper
}

// NewWebSocketHandler 创建 WebSocket handler
func NewWebSocketHandler(
	interviewSvc *service.InterviewService,
	interviewUC *biz.InterviewUsecase,
	logger log.Logger,
) *WebSocketHandler {
	return &WebSocketHandler{
		interviewSvc: interviewSvc,
		interviewUC:  interviewUC,
		logger:       log.NewHelper(logger),
	}
}

// Handle 处理 WebSocket 连接
func (h *WebSocketHandler) Handle(httpCtx http.Context) error {
	userID, ok := middleware.UserIDFromContext(httpCtx)
	if !ok {
		return httpCtx.JSON(401, map[string]string{"error": "unauthorized"})
	}

	interviewID, _ := strconv.ParseInt(httpCtx.Vars().Get("id"), 10, 64)

	// 升级 HTTP 连接为 WebSocket
	conn, err := websocket.Accept(httpCtx.Response(), httpCtx.Request(), &websocket.AcceptOptions{
		InsecureSkipVerify: true, // 开发阶段允许跨域
	})
	if err != nil {
		h.logger.Errorf("websocket accept error: %v", err)
		return err
	}
	defer conn.Close(websocket.StatusNormalClosure, "done")

	ctx := httpCtx.Request().Context()

	h.logger.Infof("WebSocket connected: user=%d interview=%d", userID, interviewID)

	// 发送连接成功状态
	h.sendJSON(ctx, conn, wsResponse{Type: "status", Data: "connected"})

	// 消息循环
	for {
		_, data, readErr := conn.Read(ctx)
		if readErr != nil {
			h.logger.Infof("WebSocket closed: user=%d err=%v", userID, readErr)
			return nil
		}

		var msg wsMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			h.sendJSON(ctx, conn, wsResponse{Type: "error", Data: "invalid message"})
			continue
		}

		switch msg.Type {
		case "text":
			h.handleTextMessage(ctx, conn, interviewID, userID, msg.Data)
		case "end":
			h.handleEndInterview(ctx, conn, interviewID, userID)
			return nil
		case "ping":
			h.sendJSON(ctx, conn, wsResponse{Type: "status", Data: "pong"})
		default:
			h.sendJSON(ctx, conn, wsResponse{Type: "error", Data: "unknown message type"})
		}
	}
}

// handleTextMessage 处理文本消息 - LLM 流式 + TTS
func (h *WebSocketHandler) handleTextMessage(
	ctx context.Context,
	conn *websocket.Conn,
	interviewID int64,
	userID int64,
	content string,
) {
	// 1. 获取用户设置
	settings, _ := h.interviewSvc.GetUserSettings(ctx, userID)

	// 2. 通知开始回复
	h.sendJSON(ctx, conn, wsResponse{Type: "text_start"})

	// 3. 调用 LLM 流式 API
	userMsg, streamCh, err := h.interviewUC.StreamMessage(ctx, interviewID, content, settings)
	if err != nil {
		h.sendJSON(ctx, conn, wsResponse{Type: "error", Data: err.Error()})
		return
	}

	// 发送用户消息确认
	h.sendJSON(ctx, conn, wsResponse{
		Type: "status",
		Data: map[string]any{"user_message_id": userMsg.ID},
	})

	// 4. 流式读取 LLM 回复，同时做句子切分 + TTS
	var fullContent strings.Builder
	var sentenceBuffer strings.Builder
	ttsEnabled := settings != nil && settings.TTSEnabled
	var ttsProvider tts.Provider
	if ttsEnabled && settings.TTSProvider != "" {
		ttsProvider, _ = h.interviewUC.GetTTSProvider(settings.TTSProvider)
	}

	for event := range streamCh {
		if event.Err != nil {
			h.sendJSON(ctx, conn, wsResponse{Type: "error", Data: event.Err.Error()})
			break
		}
		if event.Done {
			break
		}

		// 发送文本 delta
		h.sendJSON(ctx, conn, wsResponse{Type: "text_delta", Data: event.Content})
		fullContent.WriteString(event.Content)
		sentenceBuffer.WriteString(event.Content)

		// 检测句子边界，触发 TTS
		if ttsProvider != nil && isSentenceEnd(sentenceBuffer.String()) {
			sentence := strings.TrimSpace(sentenceBuffer.String())
			sentenceBuffer.Reset()
			if sentence != "" {
				go h.synthesizeAndSend(ctx, conn, ttsProvider, sentence, settings)
			}
		}
	}

	// 处理剩余的文本
	if ttsProvider != nil && sentenceBuffer.Len() > 0 {
		sentence := strings.TrimSpace(sentenceBuffer.String())
		if sentence != "" {
			h.synthesizeAndSend(ctx, conn, ttsProvider, sentence, settings)
		}
	}

	// 发送文本结束
	h.sendJSON(ctx, conn, wsResponse{Type: "text_end", Data: fullContent.String()})
}

// handleEndInterview 处理结束面试
func (h *WebSocketHandler) handleEndInterview(
	ctx context.Context,
	conn *websocket.Conn,
	interviewID int64,
	userID int64,
) {
	h.sendJSON(ctx, conn, wsResponse{Type: "status", Data: "evaluating"})

	eval, err := h.interviewSvc.EndInterview(ctx, interviewID, userID)
	if err != nil {
		h.sendJSON(ctx, conn, wsResponse{Type: "error", Data: err.Error()})
		return
	}

	h.sendJSON(ctx, conn, wsResponse{
		Type: "evaluation",
		Data: map[string]any{
			"overall_score": eval.OverallScore,
			"summary":       eval.Summary,
		},
	})
}

// synthesizeAndSend TTS 合成并发送音频
func (h *WebSocketHandler) synthesizeAndSend(
	ctx context.Context,
	conn *websocket.Conn,
	provider tts.Provider,
	text string,
	settings *biz.UserSettings,
) {
	var buf bytes.Buffer
	req := &tts.Request{
		Text:   text,
		Voice:  settings.TTSVoice,
		APIKey: settings.TTSAPIKey,
		Speed:  1.0,
		Format: "mp3",
	}

	if err := provider.Synthesize(ctx, req, &buf); err != nil {
		h.logger.Errorf("TTS synthesize error: %v", err)
		return
	}

	// 发送音频二进制数据
	if err := conn.Write(ctx, websocket.MessageBinary, buf.Bytes()); err != nil {
		h.logger.Errorf("WebSocket write audio error: %v", err)
	}
}

// sendJSON 发送 JSON 消息
func (h *WebSocketHandler) sendJSON(ctx context.Context, conn *websocket.Conn, msg wsResponse) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	_ = conn.Write(ctx, websocket.MessageText, data)
}

// isSentenceEnd 检测句子边界
func isSentenceEnd(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	last := text[len(text)-1]
	// 中文标点和英文标点
	endings := []byte{'.', '!', '?'}
	for _, e := range endings {
		if last == e {
			return true
		}
	}
	// UTF-8 中文标点
	for _, r := range []string{"。", "！", "？", "；", "\n"} {
		if strings.HasSuffix(text, r) {
			return true
		}
	}
	return false
}
