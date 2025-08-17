package main

import (
	"encoding/json"
	"fmt"
	"github.com/Marcus-Gustafsson/GatorRSS/internal/config"
	"log"
)

func main() {
	fmt.Println("Hello, Go!")

	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.SetUser("Marcus")
	if err != nil {
		log.Fatal(err)
	}

	// Re-read from disk to get the updated config
	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Print the updated config as pretty JSON
	prettyJSON, err := json.MarshalIndent(updatedCfg, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(prettyJSON))
}
