package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := &Config{
		Port: 8000,
	}

	err := cfg.serve()
	if err != nil {
		log.Panic("failed to start application", err)
	}
}


func (cfg *Config) serve() error {
	routes := routes()
	log.Println("application running on port", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), routes)
	if err != nil {
		return err
	}

	return nil
}
