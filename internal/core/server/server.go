package server

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/VladisP/media-savior/internal/common/log"
	"github.com/VladisP/media-savior/internal/core/config"
)

type HTTPServer interface {
	MountHandlers([]HTTPHandler)
}

type httpServer struct {
	server *http.Server
	logger log.ZapLogger
}

func (s *httpServer) MountHandlers(handlers []HTTPHandler) {
	for _, handler := range handlers {
		handler.Mount()
	}
}

type HTTPServerParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *config.Config
	RootHandler http.Handler `name:"root_handler"`
	Logger      log.ZapLogger
}

func NewHTTPServer(p HTTPServerParams) HTTPServer {
	s := &httpServer{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", p.Config.HTTPServer.Host, p.Config.HTTPServer.Port),
			Handler: p.RootHandler,
		},
		logger: p.Logger.NestedLogger("http_server"),
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			s.logger.Info(ctx, fmt.Sprintf("Starting HTTP server on port %s", p.Config.HTTPServer.Port))
			go func() {
				if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					s.logger.FatalWithoutContext("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.logger.Info(ctx, "Stopping HTTP server")
			return s.server.Shutdown(ctx)
		},
	})

	return s
}
