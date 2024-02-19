package handlers

import "github.com/mr-time2028/WebChat/server/settings"

var HandlerRepo *HandlerRepository

type HandlerRepository struct {
	App *settings.App
}

func NewHandlerRepository(a *settings.App) *HandlerRepository {
	return &HandlerRepository{
		App: a,
	}
}

func NewHandlers(r *HandlerRepository) {
	HandlerRepo = r
}
