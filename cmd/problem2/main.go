package main

import (
	"context"
	"os"
	"os/signal"

	"protohackers/config"
	"protohackers/problem2"
)

// go run ./cmd/problem2 -h
// go run ./cmd/problem2 [-host <host>] [-port <port>]
// go build -o ./cmd/problem2 ./cmd/problem2

func main() {
	cfg := config.New()
	cfg.ParseFlags()

	// context responding to ctrl+c (SIGINT)
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	srv := problem2.NewServer(cfg.ServerHost, cfg.ServerPort)
	srv.Run(ctx)
}
