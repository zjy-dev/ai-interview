// Package conf 定义配置结构。
// 注意：此文件为手写版本，make config 会从 conf.proto 生成同名文件覆盖。
// 在 proto 工具链就绪前，使用此文件保证编译通过。
package conf

import "time"

// Duration 包装 time.Duration，兼容 proto 生成的 durationpb.Duration
type Duration struct {
	d time.Duration
}

func NewDuration(d time.Duration) *Duration {
	return &Duration{d: d}
}

func (d *Duration) AsDuration() time.Duration {
	if d == nil {
		return 0
	}
	return d.d
}

// Bootstrap 根配置
type Bootstrap struct {
	Server    *Server    `yaml:"server"`
	Data      *Data      `yaml:"data"`
	Auth      *Auth      `yaml:"auth"`
	Tts       *TTS       `yaml:"tts"`
	Llm       *LLM       `yaml:"llm"`
	Interview *Interview `yaml:"interview"`
}

// Server 服务器配置
type Server struct {
	Http *Server_HTTP `yaml:"http"`
	Grpc *Server_GRPC `yaml:"grpc"`
}

type Server_HTTP struct {
	Network string    `yaml:"network"`
	Addr    string    `yaml:"addr"`
	Timeout *Duration `yaml:"timeout"`
}

type Server_GRPC struct {
	Network string    `yaml:"network"`
	Addr    string    `yaml:"addr"`
	Timeout *Duration `yaml:"timeout"`
}

// Data 数据层配置
type Data struct {
	Database *Data_Database `yaml:"database"`
	Redis    *Data_Redis    `yaml:"redis"`
}

type Data_Database struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

type Data_Redis struct {
	Addr         string    `yaml:"addr"`
	Password     string    `yaml:"password"`
	Db           int32     `yaml:"db"`
	ReadTimeout  *Duration `yaml:"read_timeout"`
	WriteTimeout *Duration `yaml:"write_timeout"`
}

// Auth 认证配置
type Auth struct {
	JwtSecret     string    `yaml:"jwt_secret"`
	TokenExpire   *Duration `yaml:"token_expire"`
	EncryptionKey string    `yaml:"encryption_key"`
}

// TTS 配置
type TTS struct {
	DefaultProvider string `yaml:"default_provider"`
	DefaultVoice    string `yaml:"default_voice"`
	AudioFormat     string `yaml:"audio_format"`
	SampleRate      int32  `yaml:"sample_rate"`
}

// LLM 配置
type LLM struct {
	DefaultProvider string  `yaml:"default_provider"`
	DefaultModel    string  `yaml:"default_model"`
	MaxTokens       int32   `yaml:"max_tokens"`
	Temperature     float32 `yaml:"temperature"`
}

// Interview 面试配置
type Interview struct {
	MaxQuestions    int32  `yaml:"max_questions"`
	DefaultLanguage string `yaml:"default_language"`
	SystemPrompt    string `yaml:"system_prompt"`
}
