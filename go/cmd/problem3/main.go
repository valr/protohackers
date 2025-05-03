package main

import (
	"context"
	"os"
	"os/signal"

	"protohackers/config"
	"protohackers/problem3"
)

// go run ./cmd/problem3 -h
// go run ./cmd/problem3 [-host <host>] [-port <port>]
// go build -o ./cmd/problem3 ./cmd/problem3

func main() {
	cfg := config.New()
	cfg.ParseFlags()

	// context responding to ctrl+c (SIGINT)
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	srv := problem3.NewServer(cfg.ServerHost, cfg.ServerPort)
	srv.Run(ctx)
}
