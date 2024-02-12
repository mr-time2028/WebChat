package settings

import (
	"github.com/mr-time2028/WebChat/database"
	"github.com/mr-time2028/WebChat/helpers"
	"github.com/mr-time2028/WebChat/models"
)

type App struct {
	Domain   string
	Debug    bool
	HTTPPort string
	DB       *database.DB
	Clients  map[models.Client]string
	Models   *models.ModelManager
}

func NewApp() *App {
	HTTPPort := helpers.GetEnvOrDefaultString("HTTP_PORT", "8000")
	Domain := helpers.GetEnvOrDefaultString("DOMAIN", "localhost")
	Debug := helpers.GetEnvOrDefaultBool("DEBUG", true)
	return &App{
		HTTPPort: HTTPPort,
		Domain:   Domain,
		Debug:    Debug,
	}
}
