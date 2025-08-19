package main

import (
    "database/sql"
    "log"
    "os"

    "github.com/Marcus-Gustafsson/gator/internal/config"
    "github.com/Marcus-Gustafsson/gator/internal/database"
    _ "github.com/lib/pq"
)

func main() {
    // Read configuration from file.
    cfg, err := config.Read()
    if err != nil {
        log.Fatal(err)
    }

    dbPtr, err := sql.Open("postgres", cfg.DbURL)
    if err != nil {
        log.Fatal(err)
    }
    defer dbPtr.Close()

    // Initialize State with config and database pointer.
    st := state{cfgPtr: &cfg, dbPtr: database.New(dbPtr)}

    // Initialize Cmds with a map to store handlers.
    cmds := cmds{FunctionMap: make(map[string]func(*state, command) error)}

    // Register all handlers
    registerCommands(&cmds)

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

// registerCommands registers all available command handlers
func registerCommands(cmds *cmds) {

    // User commands
    cmds.register("login", handlerLogin)
    cmds.register("register", handlerRegister)
    cmds.register("reset", handlerReset)
    cmds.register("users", handlerGetUsers)
    
    // Feed commands
    cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
    cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))
    cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
}