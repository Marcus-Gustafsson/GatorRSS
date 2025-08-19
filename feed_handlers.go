package main

import (
    "context"
    "database/sql"
    "encoding/xml"
    "errors"
    "fmt"
    "html"
    "io"
    "net/http"
    "time"

    "github.com/Marcus-Gustafsson/GatorRSS/internal/database"
    "github.com/google/uuid"
)



// fetchFeed retrieves an RSS feed from the given URL and parses it into an RSSFeed struct.
// It handles HTTP requests with proper context, sets required headers, and processes XML data.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){

    // Create a new HTTP client with timeout to prevent hanging on slow servers
    client := &http.Client{
        Timeout: time.Second * 10, // Timeout for each requests
    }

    // Create new GET request with context, Feedurl and no body (nil)
    request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // Set User agent header to identify our Go program to the server
    request.Header.Add("User-Agent", "gator")

    // Actually send the custom request with the created client
    response, err := client.Do(request)
    if err != nil {
        return nil, fmt.Errorf("error sending the request: %w", err)
    }
    defer response.Body.Close()

    // Read the response body into memory so we can later parse/unmarshal it
    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }

    // Create a placeholder pointer variable for the struct that will hold the data in the response body
    rssFeedPtr := &RSSFeed{}

    // Use xml.unmarshal due to struct expecting xml formatting (not JSON in this case.)
    err = xml.Unmarshal(body, rssFeedPtr)
    if err != nil{
        return nil, fmt.Errorf("error unmarshaling response body into RSSFeed struct: %w", err)
    }

    // Why we do this: XML often contains "escaped" characters like &amp; instead of &.
    // html.UnescapeString converts these back to normal readable characters.
    // We do this for titles and descriptions so they display properly to users.

    // Unescape HTML entities in the main channel fields for proper display
    rssFeedPtr.Channel.Title = html.UnescapeString(rssFeedPtr.Channel.Title)
    rssFeedPtr.Channel.Description = html.UnescapeString(rssFeedPtr.Channel.Description)

    // Loop through each item and unescape HTML entities in their fields
    for i := range rssFeedPtr.Channel.Item {
        rssFeedPtr.Channel.Item[i].Title = html.UnescapeString(rssFeedPtr.Channel.Item[i].Title)
        rssFeedPtr.Channel.Item[i].Description = html.UnescapeString(rssFeedPtr.Channel.Item[i].Description)
    }

    return rssFeedPtr, nil
}

// handlerAgg handles the "agg" command by fetching and displaying an RSS feed.
// This is a test function to verify our RSS parsing works correctly.
func handlerAgg(stPtr *state, cmd command) error{

    // Fetch the RSS feed from the specified URL using our fetchFeed function
    rssFeedPtr, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
    if err != nil {
        return fmt.Errorf("handlerAgg: error fetching the feed: %w", err)
    }

    // Print the entire RSS feed struct to console as required by assignment
    fmt.Println("Result:", *rssFeedPtr)

    return nil
}

// handlerAddFeed creates a new RSS feed record in the database, associated with the
// currently logged-in user. It expects exactly two arguments: the feed's name and its URL.
// On success, prints the new feed's details. Returns an error if user lookup or feed
// creation fails, or if arguments are missing.
func handlerAddFeed(stPtr *state, cmd command) error {

    if len(cmd.Args) != 2 {
        return errors.New("handlerAddFeed: expects two arguments: the feed's name and URL")
    }

    currentUser, err := stPtr.dbPtr.GetUser(
        context.Background(),
        sql.NullString{String: stPtr.cfgPtr.CurrentUserName, Valid: true},
    )
    if err != nil {
        return fmt.Errorf("handlerAddFeed: error retrieving the current user: %w", err)
    }

    newFeed, err := stPtr.dbPtr.CreateFeed(
        context.Background(),
        database.CreateFeedParams{
            ID:        uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Name:      cmd.Args[0],
            Url:       cmd.Args[1],
            UserID:    uuid.NullUUID{UUID: currentUser.ID, Valid: true},
        },
    )
    if err != nil {
        return fmt.Errorf("handlerAddFeed: failed to create new feed: %w", err)
    }


    feedFollow, err := stPtr.dbPtr.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
    ID:        uuid.New(),
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
    UserID:    currentUser.ID,
    FeedID:    newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("handlerAddFeed: couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")

    fmt.Println("DBG: New feed has been created:")
    printFeed(newFeed, currentUser)

    fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName.String, feedFollow.FeedName)
	fmt.Println("=====================================")

    return nil
}

// handlerGetFeeds retrieves all RSS feeds from the database and prints each feed’s
// name, URL, and its associated creator’s username to the console. It fetches the
// user for each feed using the user’s UUID field. Returns an error if retrieving
// feeds or users fails, or if no feeds are found.
func handlerGetFeeds(stPtr *state, cmd command) error{

    feeds, err := stPtr.dbPtr.GetFeeds(context.Background())
    if err != nil {
        return fmt.Errorf("handlerGetFeeds: couldn't retrieve all feeds from 'feeds' table: %w", err)  
    }

    if len(feeds) == 0{
        return errors.New("handlerGetFeeds: no feeds in the retrieved feed slice")
    }

    for _, feed := range feeds{

        feedUser, err := stPtr.dbPtr.GetUserByUUID(context.Background(), feed.UserID.UUID)
        if err != nil {
            return fmt.Errorf("handlerGetFeeds: couldn't retrieve user with uuid from 'users' table: %w", err)  
        }
        printFeed(feed, feedUser)
		fmt.Println("=====================================")
    }

    return nil
}

// handlerFollow creates a feed follow relationship between the current user and
// a feed specified by URL. It expects exactly one argument: the feed's URL.
// On success, prints confirmation with the user and feed names. Returns an error
// if the user lookup fails, the feed URL is not found, or feed follow creation fails.
func handlerFollow(stPtr *state, cmd command) error {
    // Check argument count first, like other handlers
    if len(cmd.Args) != 1 {
        return errors.New("handlerFollow: expects a single argument, the feed URL")
    }

    // Get current user with proper error wrapping
    currentUser, err := stPtr.dbPtr.GetUser(
        context.Background(),
        sql.NullString{String: stPtr.cfgPtr.CurrentUserName, Valid: true},
    )
    if err != nil {
        return fmt.Errorf("handlerFollow: error retrieving current user: %w", err)
    }

    // Get feed by URL with proper error wrapping
    feed, err := stPtr.dbPtr.GetFeedByURL(context.Background(), cmd.Args[0])
    if err != nil {
        return fmt.Errorf("handlerFollow: couldn't get feed by URL: %w", err)
    }

    // Create feed follow with proper error wrapping
    feedFollow, err := stPtr.dbPtr.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
        ID:        uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID:    currentUser.ID,
        FeedID:    feed.ID,
    })
    if err != nil {
        return fmt.Errorf("handlerFollow: couldn't create feed follow: %w", err)
    }

    // Print success message
    fmt.Println("Feed follow created:")
    printFeedFollow(feedFollow.UserName.String, feedFollow.FeedName)
    return nil
}


// handlerListFeedFollows retrieves and displays all feeds that the current user
// is following. It prints the user's name and a list of followed feed names.
// If no feeds are being followed, displays an appropriate message. Returns an
// error if user lookup or feed follow retrieval fails.
func handlerListFeedFollows(stPtr *state, cmd command) error {
    // Get current user with proper error wrapping
    currentUser, err := stPtr.dbPtr.GetUser(
        context.Background(),
        sql.NullString{String: stPtr.cfgPtr.CurrentUserName, Valid: true},
    )
    if err != nil {
        return fmt.Errorf("handlerListFeedFollows: error retrieving current user: %w", err)
    }

    // Get feed follows with proper error wrapping
    feedFollows, err := stPtr.dbPtr.GetFeedFollowsForUser(context.Background(), currentUser.ID)
    if err != nil {
        return fmt.Errorf("handlerListFeedFollows: couldn't retrieve feed follows: %w", err)
    }

    // Handle empty case
    if len(feedFollows) == 0 {
        fmt.Println("No feed follows found for this user.")
        return nil
    }

    // Print results
    fmt.Printf("Feed follows for user %s:\n", currentUser.Name.String)
    for _, ff := range feedFollows {
        fmt.Printf("* %s\n", ff.FeedName)
    }

    return nil
}