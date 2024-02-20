package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/mr-time2028/WebChat/commands"
	"github.com/mr-time2028/WebChat/server"
	"log"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot loading .env file")
	}

	// run command (if user want to run a command) else run http server
	command := flag.Bool("command", false, "run specific command")
	flag.Parse()

	if *command {
		commands.RunCommands()
	} else {
		err = server.HTTPServer()
		if err != nil {
			log.Panic("failed to start application", err)
		}
	}
}
