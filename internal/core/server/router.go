package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/VladisP/media-savior/internal/common/log"
	"github.com/VladisP/media-savior/internal/core/config"
)

type InputParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *config.Config
	Logger    log.ZapLogger
}

type OutputRoutes struct {
	fx.Out

	RootHandler http.Handler `name:"root_handler"`
	VKGroup     gin.IRoutes  `name:"vk_router_group"`
	UsersGroup  gin.IRoutes  `name:"users_router_group"`
}

func NewRouter(p InputParams) OutputRoutes {
	if !p.Config.HTTPServer.GinDevMode {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	rootGroup := engine.Group("/v1")

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info(ctx, "Set up routes")
			return nil
		},
	})

	return OutputRoutes{
		RootHandler: engine,
		VKGroup:     rootGroup.Group("/vk"),
		UsersGroup:  rootGroup.Group("/users"),
	}
}
