package main

import (
	"context"
	"log"
	"os"

	"github.com/iudanet/yp-diplom-1/internal/config"
	"github.com/iudanet/yp-diplom-1/internal/repo/postgres"
)

func main() {
	cfg := config.New()
	cfg.FlagParse()
	cfg.EnvParse()
	log.Println(cfg)
	ctx := context.Background()
	// Инициализация репозитория
	_, err := postgres.NewPostgresRepo(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
