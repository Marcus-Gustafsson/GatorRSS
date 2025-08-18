package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Marcus-Gustafsson/GatorRSS/internal/config"
	"github.com/Marcus-Gustafsson/GatorRSS/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // The underscore tells Go that you're importing it for its side effects, not because you need to use it.
)

// state holds a pointer to the application's configuration.
type state struct {
    cfgPtr *config.Config
	dbPtr *database.Queries
}

// command represents a CLI command with its name and argument list.
type command struct{
	Name string
	Args []string
}

// cmds stores available CLI command handlers/functions mapped by command name.
type cmds struct{
    FunctionMap map[string]func(*state, command) error
}

// run executes the given command using the associated handler, if it exists.
func (cmdsPtr *cmds) run(stPtr *state, cmd command) error{

    function, ok := cmdsPtr.FunctionMap[cmd.Name]
    // check if a handler exists before calling it, prevents a panic from a nil map entry.
    if !ok {
        return fmt.Errorf("unknown command: %s", cmd.Name)
    }
    return function(stPtr, cmd)
}

// register adds a command handler function for a given command name to used cmds structs map.
func (cmdsPtr *cmds) register(name string, function func(*state, command) error){
    cmdsPtr.FunctionMap[name] = function
}

// handlerLogin checks if a given username exists in the database, and sets it
// as the current user in the config if found. If the username does not exist,
// the process exits with code 1.
func handlerLogin(stPtr *state, cmd command) error {

    // Check that exactly one argument was passed (the username). If not, return an error.
    if len(cmd.Args) != 1 {
        return errors.New("login handler expects a single argument, the username")
    }

    // Try to fetch the user by name from the database.
    _, err := stPtr.dbPtr.GetUser(
        context.Background(),
        sql.NullString{String: cmd.Args[0], Valid: true},
    )

    // If the error is sql.ErrNoRows, that means no user was found with this name.
    // In this app, "login" must fail (exit with status code 1) if the user doesn't exist.
    if errors.Is(err, sql.ErrNoRows) {
        os.Exit(1)
    }

    // If we see another error (for example, a database connection error),
    // return it so it can be handled or logged elsewhere.
    if err != nil {
        return err
    }

    // If we found the user, set them as the current user in the config file.
    // SetUser should persist this change.
    err = stPtr.cfgPtr.SetUser(cmd.Args[0])
    if err != nil {
        return err
    }

    // Helpful for debugging: print out which user has been set.
    fmt.Printf("User has been set to: %v\n", cmd.Args[0])
    return nil
}

// handlerRegister creates a new user in the database with the given username
// and sets the new user as current in the config.
// If the username already exist the process exits with code 1.
func handlerRegister(stPtr *state, cmd command) error {

    // Make sure a username was provided and only one argument is present.
    if len(cmd.Args) != 1 {
        return errors.New("register handler expects a single argument, the name of the user to register")
    }

    // Check if a user with this name already exists in the db.
    _, err := stPtr.dbPtr.GetUser(
        context.Background(),
        sql.NullString{String: cmd.Args[0], Valid: true},
    )

    // If no error, that means user already exists—so we should not register again.
    if err == nil {
        os.Exit(1)
    }
    // If the error is something other than "no rows found", it's an actual db error, return it.
    if !errors.Is(err, sql.ErrNoRows) {
        return err
    }

    // No user was found with this name, so we proceed to register a new user.
    newUser, err := stPtr.dbPtr.CreateUser(
        context.Background(),
        database.CreateUserParams{
            ID:        uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Name:      sql.NullString{String: cmd.Args[0], Valid: true},
        },
    )
    if err != nil {
        return err
    }

    // After successful registration, set this new user as current. 
    // SetUser should update the config file as well.
    stPtr.cfgPtr.SetUser(newUser.Name.String)

    // Print out the new user details for your own debugging.
    fmt.Printf("DBG: User created: %v\n", newUser)

    return nil
}


func main() {
	// Read configuration from file.
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize State with config and database pointer.
	st := state{cfgPtr: &cfg, dbPtr: database.New(db)}

	// Initialize Cmds with a map to store handlers.
	cmds := cmds{FunctionMap: make(map[string]func(*state, command) error)}

	// Register the login handler.
	cmds.register("login", handlerLogin)
	// Register the register handler.
	cmds.register("register", handlerRegister)


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
    err = cmds.run(&st, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
