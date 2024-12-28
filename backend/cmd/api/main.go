package main

import (
	"context"
	"log/slog"
	info "offi/internal/build_info"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})).With("version", info.Version)

	slog.SetDefault(logger)

	cmd := &cli.Command{
		Name:    "offi",
		Version: info.Version,
		Commands: []*cli.Command{
			serveCommand,
		},
	}

	_ = cmd.Run(context.Background(), os.Args)
}
