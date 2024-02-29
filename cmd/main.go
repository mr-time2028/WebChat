package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mr-time2028/WebChat/internal/commands"
	"github.com/mr-time2028/WebChat/internal/config"
	"github.com/mr-time2028/WebChat/internal/database"
	"github.com/mr-time2028/WebChat/internal/handlers"
	"github.com/mr-time2028/WebChat/internal/helpers"
	"github.com/mr-time2028/WebChat/internal/models"
	"github.com/mr-time2028/WebChat/internal/routes"
	"log"
	"net/http"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot loading .env file")
	}

	// run command (if user want to run a command) else run http config
	command := flag.Bool("command", false, "run specific command")
	flag.Parse()

	if *command {
		commands.RunCommands()
	} else {
		app := &config.App{}
		err = Run(app)
		if err != nil {
			log.Panic("failed to start application", err)
		}
	}
}

func Run(app *config.App) error {
	app.HTTPPort = helpers.GetEnvOrDefaultString("HTTP_PORT", "8000")
	app.Domain = helpers.GetEnvOrDefaultString("DOMAIN", "localhost")
	app.Debug = helpers.GetEnvOrDefaultBool("DEBUG", true)

	// connect to the database
	log.Println("connecting to the database...")
	DB, err := database.ConnectSQL()
	if err != nil {
		log.Fatal("connecting to the database failed! ", err)
	}
	log.Println("connected to the database successfully!")
	app.DB = DB

	// initial JWT
	auth := models.NewJWTAuth()
	app.Auth = auth

	// initial clients settings
	hub := models.NewHub()
	app.Hub = hub

	// initial models
	models.NewModels(DB)
	app.Models = models.NewModelManager()

	// initial handlers
	handlers.NewHandlers(app)

	// graceful shutdown
	go app.ListenForShutdown()

	// start application
	log.Println("application running on port", app.HTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", app.HTTPPort), routes.Routes())
	if err != nil {
		return err
	}

	return nil
}
