package main

import (
	"context"
	"os"
	"os/signal"

	"protohackers/config"
	"protohackers/problem4"
)

// go run ./cmd/problem4 -h
// go run ./cmd/problem4 [-host <host>] [-port <port>]
// go build -o ./cmd/problem4 ./cmd/problem4

func main() {
	cfg := config.New()
	cfg.ParseFlags()

	// context responding to ctrl+c (SIGINT)
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	srv := problem4.NewServer(cfg.ServerHost, cfg.ServerPort)
	srv.Run(ctx)
}
