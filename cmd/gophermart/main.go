package main

import (
	"log"

	"github.com/iudanet/yp-diplom-1/internal/config"
	"github.com/iudanet/yp-diplom-1/internal/repo/postgres"
)

func main() {
	cfg := config.New()
	cfg.FlagParse()
	cfg.EnvParse()
	log.Println(cfg)
	// Инициализация репозитория
	_, err := postgres.NewPostgresRepo(cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}
}
