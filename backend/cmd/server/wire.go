//go:build wireinject
// +build wireinject

package main

import (
	"ai-interview/internal/biz"
	"ai-interview/internal/conf"
	"ai-interview/internal/data"
	"ai-interview/internal/middleware"
	"ai-interview/internal/provider/llm"
	"ai-interview/internal/provider/stt"
	"ai-interview/internal/provider/tts"
	"ai-interview/internal/server"
	"ai-interview/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(
	*conf.Server,
	*conf.Data,
	log.Logger,
	*tts.Registry,
	*llm.Registry,
	*stt.Registry,
	*middleware.JWTHelper,
	*middleware.Encryptor,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
