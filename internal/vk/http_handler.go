package vk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/VladisP/media-savior/internal/core/server"
)

type httpHandler struct {
	vkRouterGroup gin.IRoutes
}

func (h *httpHandler) Mount() {
	h.vkRouterGroup.GET("/hello/", h.helloHandler)
}

func (h *httpHandler) helloHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello world")
}

type HTTPHandlerParams struct {
	fx.In

	VKRouterGroup gin.IRoutes `name:"vk_router_group"`
}

type HTTPHandlerResult struct {
	fx.Out

	Handler server.HTTPHandler `name:"vk_handler"`
}

func NewHTTPHandler(p HTTPHandlerParams) HTTPHandlerResult {
	return HTTPHandlerResult{
		Handler: &httpHandler{
			vkRouterGroup: p.VKRouterGroup,
		},
	}
}
