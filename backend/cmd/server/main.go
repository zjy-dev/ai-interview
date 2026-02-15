package main

import (
	"flag"
	"os"

	"ai-interview/internal/conf"
	"ai-interview/internal/middleware"
	"ai-interview/internal/provider/llm"
	"ai-interview/internal/provider/stt"
	"ai-interview/internal/provider/tts"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
)

var (
	Name    = "ai-interview"
	Version string

	flagconf string
	id, _    = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "configs/", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(hs),
	)
}

func main() {
	flag.Parse()

	// 加载 .env 文件 (敏感配置)
	_ = godotenv.Load()

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 从环境变量覆盖敏感配置
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		bc.Data.Database.Source = dsn
	}
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		bc.Data.Redis.Addr = redisAddr
	}
	if redisPwd := os.Getenv("REDIS_PASSWORD"); redisPwd != "" {
		bc.Data.Redis.Password = redisPwd
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		bc.Auth.JwtSecret = jwtSecret
	}
	if encKey := os.Getenv("ENCRYPTION_KEY"); encKey != "" {
		bc.Auth.EncryptionKey = encKey
	}

	// 初始化 Provider 注册表
	ttsRegistry := initTTSRegistry()
	llmRegistry := initLLMRegistry()
	sttRegistry := initSTTRegistry()

	// 初始化 JWT Helper
	jwtHelper := middleware.NewJWTHelper(bc.Auth.JwtSecret, bc.Auth.TokenExpire.AsDuration())

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger, ttsRegistry, llmRegistry, sttRegistry, jwtHelper)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initTTSRegistry() *tts.Registry {
	registry := tts.NewRegistry()
	registry.Register(tts.NewOpenAIProvider())
	registry.Register(tts.NewFishAudioProvider())
	registry.Register(tts.NewElevenLabsProvider())
	registry.Register(tts.NewEdgeTTSProvider())
	return registry
}

func initLLMRegistry() *llm.Registry {
	registry := llm.NewRegistry()
	registry.Register(llm.NewOpenAIProvider())
	registry.Register(llm.NewAnthropicProvider())
	registry.Register(llm.NewDeepSeekProvider())
	registry.Register(llm.NewGeminiProvider())
	return registry
}

func initSTTRegistry() *stt.Registry {
	registry := stt.NewRegistry()
	registry.Register(stt.NewWhisperProvider())
	return registry
}
