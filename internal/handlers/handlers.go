package handlers

import (
	"github.com/mr-time2028/WebChat/internal/config"
)

var HandlerRepo *HandlerRepository

type HandlerRepository struct {
	App *config.App
}

func NewHandlers(a *config.App) {
	HandlerRepo = &HandlerRepository{
		App: a,
	}
}
