package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Marcus-Gustafsson/GatorRSS/internal/config"
	"github.com/Marcus-Gustafsson/GatorRSS/internal/database"
	_ "github.com/lib/pq" // The underscore tells Go that you're importing it for its side effects, not because you need to use it.
)


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


// printFeedFollow displays formatted information about a feed follow relationship,
// showing the username and feed name in a consistent format.
func printFeedFollow(username string, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

// printFeed displays detailed information about a feed and its associated user,
// including ID, timestamps, name, URL, and creator.
func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name.String)
}

// printUser displays detailed information about a user including ID, timestamps, and name.
func printUser(user database.User){
    fmt.Printf("* ID:            %s\n", user.ID)
	fmt.Printf("* Created:       %v\n", user.CreatedAt)
	fmt.Printf("* Updated:       %v\n", user.UpdatedAt)
	fmt.Printf("* Name:          %s\n", user.Name.String)
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
    // Register "reset" cmd to clean users table
    cmds.register("reset", handlerReset)
    // Register "users" cmd to retrieve all user names
    cmds.register("users", handlerGetUsers)
    // Register "agg" cmd to aggerage the retrieved feed and print it.
    cmds.register("agg", handlerAgg)
    // Register "addfeed" to be able to add feed from given url to current logged in user
    cmds.register("addfeed", handlerAddFeed)
    // Register "feeds" cmd which retrieves all the feeds for the current user
    cmds.register("feeds", handlerGetFeeds)
    // Register "follow" cmd, which creates a feed follow entry with the current user
    cmds.register("follow", handlerFollow)
    // Register "following" cmd, which retrieves and displays all feeds that the current user
    cmds.register("following", handlerListFeedFollows)



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
