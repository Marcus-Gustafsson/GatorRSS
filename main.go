package main

import (
	"fmt"
    "github.com/Marcus-Gustafsson/GatorRSS/internal/config"
)

func main() {
	fmt.Println("Hello, Go!")
    cfg := config.Read()
    fmt.Println(cfg)
}
