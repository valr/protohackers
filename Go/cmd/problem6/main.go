package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"protohackers/config"
	"protohackers/problem6"
)

// go run ./cmd/problem6 -h
// go run ./cmd/problem6 [-host <host>] [-port <port>]
// go build -o ./cmd/problem6 ./cmd/problem6

func main() {
	cfg := config.New()
	cfg.ParseFlags()

	if cfg.DebugMode {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	// context responding to ctrl+c (SIGINT)
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	srv := problem6.NewServer(cfg.ServerHost, cfg.ServerPort)
	srv.Run(ctx)
}
