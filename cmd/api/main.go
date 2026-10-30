package main

import (
	"context"
	"log"

	"github.com/dandolya/corp-app/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("connected")
}
