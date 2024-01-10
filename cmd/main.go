package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/mr-time2028/WebChat/server/config"
	"github.com/mr-time2028/WebChat/server/database"
)

var cfg *config.Config

func main() {
	err := serve()
	if err != nil {
		log.Panic("failed to start application", err)
	}
}

func serve() error {
	cfg = &config.Config{
		Port: 8000,
	}

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot loading .env file")
	}

	// connect to the database
	log.Println("connecting to the database...")
	DB, err := database.ConnectSQL()
	if err != nil {
		log.Fatal("connecting to the database failed! ", err)
	}
	log.Println("connected to the database successfully!")
	cfg.DB = DB

	// start application
	log.Println("application running on port", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), routes())
	if err != nil {
		return err
	}

	return nil
}
