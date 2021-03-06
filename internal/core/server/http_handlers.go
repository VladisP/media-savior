package server

import (
	"go.uber.org/fx"
)

type HTTPHandler interface {
	Mount()
}

type InputHandlers struct {
	fx.In

	VKHTTPHandler HTTPHandler `name:"vk_handler"`
}

func NewHTTPHandlers(input InputHandlers) []HTTPHandler {
	return []HTTPHandler{
		input.VKHTTPHandler,
	}
}
