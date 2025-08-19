package main

import (
    "fmt"
    "github.com/Marcus-Gustafsson/GatorRSS/internal/database"
)



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


