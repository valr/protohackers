package main

import (
	"context"
	"os"
	"os/signal"

	"protohackers/config"
	"protohackers/problem1"
)

// go run ./cmd/problem1 -h
// go run ./cmd/problem1 [-host <host>] [-port <port>]
// go build -o ./cmd/problem1 ./cmd/problem1

func main() {
	cfg := config.New()
	cfg.ParseFlags()

	// context responding to ctrl+c (SIGINT)
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	srv := problem1.NewServer(cfg.ServerHost, cfg.ServerPort)
	srv.Run(ctx)
}
