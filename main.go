package main

import (
	"errors"
	"fmt"
	"log"
	"os"
    
	"github.com/Marcus-Gustafsson/GatorRSS/internal/config"
)

// state holds a pointer to the application's configuration.
type state struct {
    cfgPtr *config.Config
}

// command represents a CLI command with its name and argument list.
type command struct{
	Name string
	Args []string
}

// commands stores available CLI command handlers/functions mapped by command name.
type commands struct{
    FunctionMap map[string]func(*state, command) error
}

// run executes the given command using the associated handler, if it exists.
func (cmdsPtr *commands) run(stPtr *state, cmd command) error{

    function, ok := cmdsPtr.FunctionMap[cmd.Name]
    // check if a handler exists before calling it, prevents a panic from a nil map entry.
    if !ok {
        return fmt.Errorf("unknown command: %s", cmd.Name)
    }
    return function(stPtr, cmd)
}

// register adds a command handler function for a given command name to used commands structs map.
func (cmdsPtr *commands) register(name string, function func(*state, command) error){
    cmdsPtr.FunctionMap[name] = function
}

// handlerLogin sets the current user in the config file using the provided username.
func handlerLogin(stPtr *state, cmd command) error {

    // Return an error if the argument slice does not contain exactly one argument = the username.
    if len(cmd.Args) != 1{
        return errors.New("login handler expects a single argument, the username")
    }

    // Set the user in the config and return any error.
    err := stPtr.cfgPtr.SetUser(cmd.Args[0])
    if err != nil{
        return err
    }

    fmt.Printf("User has been set to: %v\n", cmd.Args[0])
    return nil

}


func main() {
	// Read configuration from file.
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize State with config pointer.
	st := state{cfgPtr: &cfg}

	// Initialize Commands with a map to store handlers.
	commands := commands{FunctionMap: make(map[string]func(*state, command) error)}

	// Register the login handler.
	commands.register("login", handlerLogin)

	// Get command-line arguments.
	args := os.Args

	// Require at least the program name and a command.
	if len(args) < 2 {
		log.Fatal("Not enough arguments. Usage: go run . <command> [args...]")
	}

	// Extract command name and arguments.
	cmdName := args[1]
	cmdArgs := args[2:]

	// Create Command instance.
	cmd := command{Name: cmdName, Args: cmdArgs}

	// Run the command; print any errors and exit on failure.
    err = commands.run(&st, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
