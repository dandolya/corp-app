package main

import (
	"context"
	"log"

	"github.com/dandolya/corp-app/internal/approval"
	"github.com/dandolya/corp-app/internal/config"
	"github.com/dandolya/corp-app/internal/db/postgres"
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
	reqRepo := postgres.NewRequestRepository(pool)

	err = reqRepo.ResolveRequest(ctx, 7, 1, approval.StatusApproved, "") // продавец — ErrNotAllowed
	log.Println(err)
	err = reqRepo.ResolveRequest(ctx, 2, 1, approval.StatusApproved, "ок") // директор — nil
	log.Println(err)
	err = reqRepo.ResolveRequest(ctx, 1, 1, approval.StatusRejected, "") // CEO, уже решено — ErrAlreadyResolved
	log.Println(err)
}
