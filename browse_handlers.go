package main

import (
	"context"
	"fmt"
	"strconv"
	"github.com/Marcus-Gustafsson/gator/internal/database"
)

// handlerBrowse retrieves and displays posts for the current user with an optional limit
// parameter (defaults to 2 if not provided). It validates the limit argument if given,
// queries posts from feeds the user follows, and prints formatted post details including
// publication date, feed name, title, description, and URL. Returns an error if limit
// parsing or post retrieval fails.
func handlerBrowse(stPtr *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("handlerBrowse: invalid limit argument: %w", err)
		}
	}

	posts, err := stPtr.dbPtr.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("handlerBrowse: couldn't retrieve posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name.String)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}