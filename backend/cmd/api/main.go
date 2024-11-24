package main

import (
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"offi/internal/service"
	"os"
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

	if err = http.ListenAndServe(":8080", handler); err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}
