package main

import (
	"fmt"
	"log"

	"github.com/dandolya/corp-app/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.HTTPAddr)
}
