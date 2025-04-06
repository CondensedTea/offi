package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"offi/internal/cache"
	"offi/internal/closer"
	"offi/internal/db"
	"offi/internal/etf2l"
	gen "offi/internal/gen/api"
	"offi/internal/logstf"
	"offi/internal/service"
	"offi/internal/tracing"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/urfave/cli/v3"
)

var address string

var serveCommand = &cli.Command{
	Name: "serve",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "address",
			Value:       ":8080",
			Destination: &address,
		},
	},
	UsageText: "serve <address>",
	Usage:     "starts the api server",
	Action:    serveAction,
}

func serveAction(ctx context.Context, _ *cli.Command) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	tracing.Init(ctx)

	etf2lClient := etf2l.New()

	logsClient := logstf.NewClient()

	cacheClient, err := cache.New(os.Getenv("REDIS_URL"))
	if err != nil {
		return fmt.Errorf("failed to init redis client: %w", err)
	}

	dbClient, err := db.NewClient(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("failed to init database client: %w", err)
	}

	srv := service.NewService(cacheClient, dbClient, etf2lClient, logsClient, true)

	handler, err := gen.NewServer(srv)
	if err != nil {
		return fmt.Errorf("failed to init api server: %w", err)
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://logs.tf", "https://etf2l.org", "https://steamcommunity.com"},
		AllowedMethods: []string{http.MethodGet},
	}))
	router.Use(middleware.Recoverer)
	router.Use(tracing.NewMiddleware(handler))
	router.Use(tracing.InjectTracingHeaders)

	router.Mount("/", handler)

	httpSrv := http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	closer.AddContext(httpSrv.Shutdown)

	go func() {
		if err = httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to run server", "error", err)
		}
	}()

	slog.Info("app is running")

	<-ctx.Done()

	slog.Info("shutting down server")

	stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	closer.CloseAll(stopCtx)

	cancel()

	return nil
}
