
package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "os"
    "time"

    "github.com/Marcus-Gustafsson/gator/internal/database"
    "github.com/google/uuid"
)





// handlerLogin checks if a given username exists in the database, and sets it
// as the current user in the config if found. If the username does not exist,
// the process exits with code 1.
func handlerLogin(stPtr *state, cmd command) error {

    // Check that exactly one argument was passed (the username). If not, return an error.
    if len(cmd.Args) != 1 {
        return errors.New("handlerLogin: expects a single argument, the username")
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
        return fmt.Errorf("handlerLogin: error fetching user: %w", err)
    }

    // If we found the user, set them as the current user in the config file.
    // SetUser should persist this change.
    err = stPtr.cfgPtr.SetUser(cmd.Args[0])
    if err != nil {
        return fmt.Errorf("handlerLogin: error setting current user: %w", err)
    }

    // Helpful for debugging: print out which user has been set.
    fmt.Printf("User has been set to: %v\n", cmd.Args[0])
    return nil
}

// handlerRegister creates a new user in the database with the given username
// and sets the new user as current in the config.
// If the username already exists, the process exits with code 1.
func handlerRegister(stPtr *state, cmd command) error {

    // Make sure a username was provided and only one argument is present.
    if len(cmd.Args) != 1 {
        return errors.New("handlerRegister: expects a single argument, the username to register")
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
        return fmt.Errorf("handlerRegister: error checking if user exists: %w", err)
    }

    // No user was found with given name, proceed to register a new user.
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
        return fmt.Errorf("handlerRegister: error creating user: %w", err)
    }

    // After successful registration, set this new user as current. 
    // SetUser should update the config file as well.
    err = stPtr.cfgPtr.SetUser(newUser.Name.String)
    if err != nil {
        return fmt.Errorf("handlerRegister: error setting current user: %w", err)
    }

    // Print out the new user details for your own debugging.
    fmt.Println("User created successfully:")
    printUser((newUser))

    return nil
}

// handlerReset deletes all users from the users table in the database,
// effectively resetting the application's state for testing or development.
// Prints a confirmation message on success. Returns an error if the deletion
// fails, in which case a non-zero exit code will be produced by main().
func handlerReset(stPtr *state, cmd command) error {

    err := stPtr.dbPtr.DeleteUsers(context.Background()) // Execute the DELETE

    if err != nil {
        return fmt.Errorf("handlerReset: couldn't delete all users in 'users' table: %w", err)  
    }

    // Success message: confirms operation completed
    fmt.Println("Users table reset successfully!")  
    return nil  // Exit code 0: tells automation "success"

}


// handlerGetUsers retrieves all registered users from the database
// and prints them to the console. It marks the currently logged-in user
// with "(current)". Returns an error if the user retrieval fails.
func handlerGetUsers(stPtr *state, cmd command) error{

    userNames, err := stPtr.dbPtr.GetUsers(context.Background())
    if err != nil {
        return fmt.Errorf("handlerGetUsers: couldn't retrieve all user names in 'users' table: %w", err)  
    }

    for _, user := range userNames{
        if user.Valid{
            if stPtr.cfgPtr.CurrentUserName == user.String{
                fmt.Printf("* %v (current)\n", user.String)
            }else{
                fmt.Printf("* %v\n", user.String)
            }
        }
    }

    return nil
}