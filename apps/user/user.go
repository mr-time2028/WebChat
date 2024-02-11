package user

import (
	"github.com/mr-time2028/WebChat/server/config"
)

var cfg *config.Config

func RegisterHandlersConfig(c *config.Config) {
	cfg = c
}
