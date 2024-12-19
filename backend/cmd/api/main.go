package main

import (
	"context"
	"log/slog"
	"os"

	info "offi/internal/build_info"

	"github.com/urfave/cli/v3"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cmd := &cli.Command{
		Name:    "offi",
		Version: info.Version,
		Commands: []*cli.Command{
			serveCommand,
		},
	}

	_ = cmd.Run(context.Background(), os.Args)
}
