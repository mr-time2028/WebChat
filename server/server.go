package server

import (
	"fmt"
	"github.com/mr-time2028/WebChat/database"
	"github.com/mr-time2028/WebChat/handlers"
	"github.com/mr-time2028/WebChat/models"
	"github.com/mr-time2028/WebChat/routes"
	"github.com/mr-time2028/WebChat/server/settings"
	"log"
	"net/http"
)

var app *settings.App

func HTTPServer() error {
	app = settings.NewApp()

	// connect to the database
	log.Println("connecting to the database...")
	DB, err := database.ConnectSQL()
	if err != nil {
		log.Fatal("connecting to the database failed! ", err)
	}
	log.Println("connected to the database successfully!")
	app.DB = DB

	// initial clients settings
	app.Clients = make(map[models.Client]string)

	// initial models
	models.RegisterModelsConfig(DB)
	app.Models = models.NewModels()

	// initial handlers
	handlerRepo := handlers.NewHandlerRepository(app)
	handlers.NewHandlers(handlerRepo)

	// start application
	log.Println("application running on port", app.HTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", app.HTTPPort), routes.Routes())
	if err != nil {
		return err
	}

	return nil
}
