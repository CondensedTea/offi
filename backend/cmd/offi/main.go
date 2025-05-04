package main

import (
	"context"
	"log/slog"
	info "offi/internal/build_info"
	"os"

	"github.com/go-slog/otelslog"
	"github.com/urfave/cli/v3"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(otelslog.NewHandler(handler)).With("version", info.Version)

	slog.SetDefault(logger)

	cmd := &cli.Command{
		Name:    "offi",
		Version: info.Version,
		Commands: []*cli.Command{
			serveCommand,
			crawlCommand,
			linkCommand,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
