package server

import (
	"ai-interview/internal/conf"
	"ai-interview/internal/middleware"
	"ai-interview/internal/service"
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer 创建 HTTP server
func NewHTTPServer(
	c *conf.Server,
	authSvc *service.AuthService,
	interviewSvc *service.InterviewService,
	wsHandler *WebSocketHandler,
	jwtHelper *middleware.JWTHelper,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
		http.Filter(corsFilter()),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	// 注册路由
	registerRoutes(srv, authSvc, interviewSvc, wsHandler, jwtHelper, logger)

	return srv
}

// registerRoutes 注册所有 HTTP 路由
func registerRoutes(
	srv *http.Server,
	authSvc *service.AuthService,
	interviewSvc *service.InterviewService,
	wsHandler *WebSocketHandler,
	jwtHelper *middleware.JWTHelper,
	logger log.Logger,
) {
	router := srv.Route("/")

	// 健康检查
	router.GET("/health", func(ctx http.Context) error {
		return ctx.JSON(200, map[string]string{"status": "ok"})
	})

	// Auth 路由 (无需认证)
	router.POST("/api/v1/auth/register", authHandler(authSvc).Register)
	router.POST("/api/v1/auth/login", authHandler(authSvc).Login)

	// 需要认证的路由
	router.GET("/api/v1/auth/profile", withAuth(jwtHelper, authHandler(authSvc).GetProfile))
	router.GET("/api/v1/auth/settings", withAuth(jwtHelper, authHandler(authSvc).GetSettings))
	router.PUT("/api/v1/auth/settings", withAuth(jwtHelper, authHandler(authSvc).UpdateSettings))

	// Interview 路由
	router.POST("/api/v1/interviews", withAuth(jwtHelper, interviewHandler(interviewSvc).Create))
	router.GET("/api/v1/interviews", withAuth(jwtHelper, interviewHandler(interviewSvc).List))
	router.GET("/api/v1/interviews/{id}", withAuth(jwtHelper, interviewHandler(interviewSvc).Get))
	router.POST("/api/v1/interviews/{id}/messages", withAuth(jwtHelper, interviewHandler(interviewSvc).SendMessage))
	router.POST("/api/v1/interviews/{id}/end", withAuth(jwtHelper, interviewHandler(interviewSvc).End))
	router.GET("/api/v1/interviews/{id}/evaluation", withAuth(jwtHelper, interviewHandler(interviewSvc).GetEvaluation))

	// WebSocket 路由 (面试实时交互)
	router.GET("/api/v1/ws/interview/{id}", withAuth(jwtHelper, wsHandler.Handle))
}

// withAuth JWT 认证包装器
func withAuth(jwtHelper *middleware.JWTHelper, handler func(http.Context) error) func(http.Context) error {
	return func(ctx http.Context) error {
		authHeader := ctx.Header().Get("Authorization")
		token := middleware.ExtractTokenFromHeader(authHeader)
		if token == "" {
			return ctx.JSON(401, map[string]string{"error": "unauthorized"})
		}

		userID, err := jwtHelper.ValidateToken(token)
		if err != nil {
			return ctx.JSON(401, map[string]string{"error": "invalid token"})
		}

		// 将 userID 存入 context
		newCtx := middleware.ContextWithUserID(ctx, userID)
		ctx.Reset(ctx.Response(), ctx.Request().WithContext(newCtx))

		return handler(ctx)
	}
}

// corsFilter 返回 CORS 过滤器
func corsFilter() http.FilterFunc {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			if r.Method == nethttp.MethodOptions {
				w.WriteHeader(nethttp.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
