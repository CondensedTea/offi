package main

import (
	"context"
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"offi/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/cors"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	etf2lClient := etf2l.New()

	cacheClient, err := cache.New(os.Getenv("REDIS_URL"))
	if err != nil {
		slog.Error("failed to init redis client", "error", err)
		os.Exit(1)
	}

	srv := service.NewService(cacheClient, etf2lClient, true)

	handler, err := gen.NewServer(srv)
	if err != nil {
		slog.Error("failed to init api server", "error", err)
		os.Exit(1)
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"https://logs.tf", "https://etf2l.org", "https://steamcommunity.com"},
		AllowedMethods: []string{http.MethodGet},
	})

	httpSrv := http.Server{
		Addr:              ":8080",
		Handler:           corsMiddleware.Handler(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err = httpSrv.ListenAndServe(); err != nil {
			slog.Error("failed to run server", "error", err)
		}
	}()

	slog.Info("app is running")

	<-ctx.Done()

	slog.Info("shutting down server")

	stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	if err = httpSrv.Shutdown(stopCtx); err != nil {
		slog.Error("failed to gracefully stop http server", "error", err)
	}

	<-stopCtx.Done()
}
