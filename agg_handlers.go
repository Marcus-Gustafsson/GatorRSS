package main

import (
	"context"
	"time"
	"database/sql"
	"strings"
	"fmt"
	"log"
	"github.com/Marcus-Gustafsson/gator/internal/database"
	"github.com/google/uuid"
)

// handlerAgg starts the aggregation loop that fetches RSS feeds at the interval specified
// by the user (time_between_reqs argument; e.g., "1m", "10s"). It validates arguments,
// parses duration, prints the loop interval, and runs scrapeFeeds using a ticker. Returns
// an error if argument parsing or duration parsing fails.
func handlerAgg(stPtr *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("handlerAgg: usage: %v <time_between_reqs>", cmd.Name)
	}

	// Parse the requested time interval, report error with context if invalid.
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("handlerAgg: invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	// Loop forever, scraping feeds immediately and then at each interval.
	for ; ; <-ticker.C {
		scrapeFeeds(stPtr)
	}
}

// scrapeFeeds retrieves the next feed to fetch from the database and processes it.
// Logs errors if feed retrieval fails, otherwise it calls scrapeFeed.
func scrapeFeeds(stPtr *state) {
	feed, err := stPtr.dbPtr.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("scrapeFeeds: couldn't get next feed to fetch:", err)
		return
	}
	log.Printf("scrapeFeeds: found a feed to fetch: %s", feed.Name)
	scrapeFeed(stPtr.dbPtr, feed)
}

// scrapeFeed marks a feed as fetched in the database, collects its posts via fetchFeed,
// and saves each post to the database with proper time parsing. Handles duplicate post
// URLs by continuing, logs other database errors, and reports the total number of posts
// processed. Errors in marking, fetching, or post creation are logged with context.
func scrapeFeed(dbPtr *database.Queries, feed database.Feed) {

	_, err := dbPtr.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("scrapeFeed: couldn't mark feed %s as fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("scrapeFeed: couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = dbPtr.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("scrapeFeed: couldn't create post: %v", err)
			continue
		}
	}

	log.Printf("scrapeFeed: feed %s collected, %d posts found", feed.Name, len(feedData.Channel.Item))
}