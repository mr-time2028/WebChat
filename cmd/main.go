package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mr-time2028/WebChat/server/config"
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

	log.Println("application running on port", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), routes())
	if err != nil {
		return err
	}

	return nil
}
