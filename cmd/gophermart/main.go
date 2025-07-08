package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iudanet/yp-diplom-1/internal/config"
	"github.com/iudanet/yp-diplom-1/internal/handlers"
	"github.com/iudanet/yp-diplom-1/internal/repo/postgres"
	"github.com/iudanet/yp-diplom-1/internal/service"
)

func main() {
	cfg := config.New()
	cfg.FlagParse()
	cfg.EnvParse()
	log.Println(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Инициализация репозитория
	repo, err := postgres.NewPostgresRepo(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if cfg.Migrate {
		err = repo.Migrate(ctx)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
	svc := service.NewService(repo)
	h := handlers.New(svc)
	mux := h.NewMux()
	srv := &http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: mux,
	}

	go func() {
		log.Println("Running server to", cfg.HTTPAddress)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			cancel()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigCh
	log.Println("Received signal", sig.String())
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
