package main

import (
	"context"
	"os"
	"os/signal"

	"protohackers/config"
	"protohackers/problem5"
)

// go run ./cmd/problem5 -h
// go run ./cmd/problem5 [-host <host>] [-port <port>]
// go build -o ./cmd/problem5 ./cmd/problem5

func main() {
	cfg := config.New()
	cfg.ParseFlags()

	// context responding to ctrl+c (SIGINT)
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	srv := problem5.NewServer(cfg.ServerHost, cfg.ServerPort)
	srv.Run(ctx)
}
