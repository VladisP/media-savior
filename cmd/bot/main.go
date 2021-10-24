package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/VladisP/media-savior/internal/common/log"
	"github.com/VladisP/media-savior/internal/common/validator"
	"github.com/VladisP/media-savior/internal/core/config"
	"github.com/VladisP/media-savior/internal/core/server"
	"github.com/VladisP/media-savior/internal/vk"
)

func appOptions() []fx.Option {
	return []fx.Option{
		fx.Provide(
			server.NewHTTPServer,
			server.NewRouter,
			server.NewHTTPHandlers,
			config.NewConfig,
			log.NewLogger,
			validator.NewValidator,
			vk.NewHTTPHandler,
		),
		fx.Invoke(
			func(server server.HTTPServer, handlers []server.HTTPHandler) {
				server.MountHandlers(handlers)
			},
		),
		fx.WithLogger(
			func(logger log.ZapLogger) fxevent.Logger {
				return &fxevent.ZapLogger{
					Logger: logger.Zap(),
				}
			},
		),
	}
}

func main() {
	app := fx.New(appOptions()...)

	app.Run()
}
