package main

import (
	"context"
	"os"
	"os/signal"

	"protohackers/config"
	"protohackers/problem0"
)

// go run ./cmd/problem0 -h
// go run ./cmd/problem0 [-host <host>] [-port <port>]
// go build -o ./cmd/problem0 ./cmd/problem0

func main() {
	cfg := config.New()
	cfg.ParseFlags()

	// context responding to ctrl+c (SIGINT)
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	srv := problem0.NewServer(cfg.ServerHost, cfg.ServerPort)
	srv.Run(ctx)
}
