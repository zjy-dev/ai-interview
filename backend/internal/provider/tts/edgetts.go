package tts

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"nhooyr.io/websocket"
)

// EdgeTTSProvider 实现 Microsoft Edge 免费 TTS (非官方)
type EdgeTTSProvider struct {
	wsURL string
}

func NewEdgeTTSProvider() *EdgeTTSProvider {
	return &EdgeTTSProvider{
		wsURL: "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1",
	}
}

func (p *EdgeTTSProvider) Name() string {
	return "edgetts"
}

func (p *EdgeTTSProvider) Synthesize(ctx context.Context, req *Request, w io.Writer) error {
	voice := req.Voice
	if voice == "" {
		if strings.HasPrefix(req.Language, "zh") {
			voice = "zh-CN-XiaoxiaoNeural"
		} else {
			voice = "en-US-AriaNeural"
		}
	}

	outputFormat := "audio-24khz-48kbitrate-mono-mp3"
	switch req.Format {
	case "pcm":
		outputFormat = "raw-24khz-16bit-mono-pcm"
	case "opus":
		outputFormat = "webm-24khz-16bit-mono-opus"
	}

	// 生成连接 ID
	connectID := generateConnectID()

	wsURL := fmt.Sprintf(
		"%s?TrustedClientToken=6A5AA1D4EAFF4E9FB37E23D68491D6F4&ConnectionId=%s",
		p.wsURL, connectID,
	)

	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		Subprotocols: []string{""},
	})
	if err != nil {
		return fmt.Errorf("edgetts: websocket dial: %w", err)
	}
	defer conn.CloseNow()

	// 发送配置消息
	configMsg := fmt.Sprintf(
		"Content-Type:application/json; charset=utf-8\r\nPath:speech.config\r\n\r\n"+
			`{"context":{"synthesis":{"audio":{"metadataoptions":{"sentenceBoundaryEnabled":"false","wordBoundaryEnabled":"false"},"outputFormat":"%s"}}}}`,
		outputFormat,
	)
	if err := conn.Write(ctx, websocket.MessageText, []byte(configMsg)); err != nil {
		return fmt.Errorf("edgetts: send config: %w", err)
	}

	// 构建 SSML
	ssml := fmt.Sprintf(
		`<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='%s'>`+
			`<voice name='%s'>%s</voice></speak>`,
		req.Language, voice, escapeXML(req.Text),
	)

	requestID := generateConnectID()
	ssmlMsg := fmt.Sprintf(
		"X-RequestId:%s\r\nContent-Type:application/ssml+xml\r\nPath:ssml\r\n\r\n%s",
		requestID, ssml,
	)
	if err := conn.Write(ctx, websocket.MessageText, []byte(ssmlMsg)); err != nil {
		return fmt.Errorf("edgetts: send ssml: %w", err)
	}

	// 读取音频响应
	for {
		msgType, data, err := conn.Read(ctx)
		if err != nil {
			// 连接关闭表示音频传输完成
			if websocket.CloseStatus(err) != -1 {
				break
			}
			return fmt.Errorf("edgetts: read: %w", err)
		}

		if msgType == websocket.MessageBinary {
			// 二进制消息包含头部 + 音频数据
			// 头部格式: 2 bytes (header length) + header + audio data
			if len(data) < 2 {
				continue
			}
			headerLen := int(data[0])<<8 | int(data[1])
			if len(data) <= 2+headerLen {
				continue
			}
			audioData := data[2+headerLen:]
			if _, err := w.Write(audioData); err != nil {
				return fmt.Errorf("edgetts: write audio: %w", err)
			}
		} else if msgType == websocket.MessageText {
			// 文本消息 - 检查是否是结束标记
			msg := string(data)
			if strings.Contains(msg, "Path:turn.end") {
				break
			}
		}
	}

	conn.Close(websocket.StatusNormalClosure, "")
	return nil
}

func generateConnectID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
