package websocket

import (
	"github.com/mr-time2028/WebChat/server/settings"
)

var app *settings.App

func RegisterHandlersConfig(a *settings.App) {
	app = a
}
