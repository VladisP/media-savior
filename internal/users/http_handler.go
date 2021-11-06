package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/VladisP/media-savior/internal/common/server"
)

type httpHandler struct {
	usersRouterGroup gin.IRoutes
	repository       Repository
}

func (h *httpHandler) Mount() {
	h.usersRouterGroup.GET("/:UserID/", h.getUserHandler)
}

func (h *httpHandler) getUserHandler(ctx *gin.Context) {
	user, err := h.repository.GetUser(ctx.Param("UserID"))
	if err != nil {
		err = ctx.Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type HTTPHandlerParams struct {
	fx.In

	UsersRouterGroup gin.IRoutes `name:"users_router_group"`
	Repository       Repository
}

type HTTPHandlerResult struct {
	fx.Out

	Handler server.HTTPHandler `name:"users_handler"`
}

func NewHTTPHandler(p HTTPHandlerParams) HTTPHandlerResult {
	return HTTPHandlerResult{
		Handler: &httpHandler{
			usersRouterGroup: p.UsersRouterGroup,
			repository:       p.Repository,
		},
	}
}
