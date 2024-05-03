package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/juancortelezzi/gatorserver/pkg/clog"
	"github.com/juancortelezzi/gatorserver/pkg/server"
)

func main() {
	ctx := context.Background()
	logger := clog.NewLogger(os.Stdout, slog.LevelInfo)

	if err := server.Run(ctx, logger, os.LookupEnv); err != nil {
		logger.ErrorContext(ctx, "error in top level", "err", err)
		os.Exit(1)
	}
}
