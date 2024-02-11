package server

import (
	"fmt"
	"github.com/mr-time2028/WebChat/apps/user"
	"github.com/mr-time2028/WebChat/apps/websocket"
	"github.com/mr-time2028/WebChat/database"
	"github.com/mr-time2028/WebChat/models"
	"github.com/mr-time2028/WebChat/routes"
	"github.com/mr-time2028/WebChat/server/config"
	"log"
	"net/http"
)

var cfg *config.Config

func HTTPServer() error {
	cfg = config.NewConfig()

	// connect to the database
	log.Println("connecting to the database...")
	DB, err := database.ConnectSQL()
	if err != nil {
		log.Fatal("connecting to the database failed! ", err)
	}
	log.Println("connected to the database successfully!")
	cfg.DB = DB

	// initial clients config
	cfg.Clients = make(map[models.Client]string)

	// register handlers
	user.RegisterHandlersConfig(cfg)
	websocket.RegisterHandlersConfig(cfg)

	// start application
	log.Println("application running on port", cfg.HTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), routes.Routes())
	if err != nil {
		return err
	}

	return nil
}
