package main

import (
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"offi/internal/service"
	"os"

	"github.com/go-chi/cors"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	etf2lClient := etf2l.New()

	cacheClient, err := cache.New(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"))
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

	if err = http.ListenAndServe(":8080", corsMiddleware.Handler(handler)); err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}
