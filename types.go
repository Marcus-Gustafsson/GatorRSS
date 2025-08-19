package main

import (
    "github.com/Marcus-Gustafsson/gator/internal/config"
    "github.com/Marcus-Gustafsson/gator/internal/database"
)

// state holds a pointer to the application's configuration.
type state struct {
    cfgPtr *config.Config
    dbPtr  *database.Queries
}

// command represents a CLI command with its name and argument list.
type command struct {
    Name string
    Args []string
}

// cmds stores available CLI command handlers/functions mapped by command name.
type cmds struct {
    FunctionMap map[string]func(*state, command) error
}

// RSSFeed represents the structure of an RSS feed with channel information and items.
type RSSFeed struct {
    Channel struct {
        Title       string    `xml:"title"`
        Link        string    `xml:"link"`
        Description string    `xml:"description"`
        Item        []RSSItem `xml:"item"`
    } `xml:"channel"`
}

// RSSItem represents a single item/article within an RSS feed.
type RSSItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
}