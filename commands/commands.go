package commands

import (
	"log"
	"os"
)

// Command is a struct for our custom commands
type Command struct {
	Description string
	Function    func()
}

// Commands contains custom commands information
var Commands = map[string]Command{
	"migrate": {
		Description: "this command simply migrate all structs to the tables in database.",
		Function:    migrate,
	},
}

// RunCommands runs any command that user determine in CLI using -command parameter
func RunCommands() {
	if len(os.Args) < 3 {
		log.Fatal("it seems you want to run a command. Usage: go run main.go -command <YOUR COMMAND NAME>")
	}

	commandName := os.Args[2]
	os.Args = os.Args[2:]

	command, ok := Commands[commandName]
	if !ok {
		log.Fatal("invalid command!")
	}
	command.Function()
}
