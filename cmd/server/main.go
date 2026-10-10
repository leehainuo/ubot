package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"ubot/internal/api"
	"ubot/internal/config"
	"ubot/internal/platform/iam"
)

func main() {
	cfg := config.Load()
	mux := http.NewServeMux()

	iam := iam.New(iam.NewMemoryStore(), cfg.JWT.Secret, cfg.JWT.Issuer)

	api.RegisterRoutes(
		mux,
		iam,
	)

	server := &http.Server{
		Addr:        cfg.Addr,
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		_ = server.Shutdown(shutdownCtx)
	}()

	log.Println("Ubot running...")
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
