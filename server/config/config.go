package config

import (
	"github.com/mr-time2028/WebChat/database"
	"github.com/mr-time2028/WebChat/helpers"
	"github.com/mr-time2028/WebChat/models"
)

type Config struct {
	Domain   string
	Debug    bool
	HTTPPort string
	DB       *database.DB
	Clients  map[models.Client]string
}

func NewConfig() *Config {
	HTTPPort := helpers.GetEnvOrDefaultString("HTTP_PORT", "8000")
	Domain := helpers.GetEnvOrDefaultString("DOMAIN", "localhost")
	Debug := helpers.GetEnvOrDefaultBool("DEBUG", true)
	return &Config{
		HTTPPort: HTTPPort,
		Domain:   Domain,
		Debug:    Debug,
	}
}
