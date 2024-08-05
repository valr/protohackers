package problem1

// https://protohackers.com/problem/1

import (
	"bufio"
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"time"
)

type Server struct {
	host string
	port string
}

const (
	connTimeout = 10
)

func NewServer(host, port string) Server {
	return Server{host, port}
}

func (srv Server) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", net.JoinHostPort(srv.host, srv.port))
	if err != nil {
		slog.Error("net.Listen failed:",
			slog.Any("error", err),
			slog.String("id", "problem1/server/E001"),
			slog.Time("time", time.Now()),
			slog.String("host", srv.host),
			slog.String("port", srv.port))
		return
	}

	defer listener.Close()
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() == nil {
				slog.Error("listener.Accept failed:",
					slog.Any("error", err),
					slog.String("id", "problem1/server/E002"),
					slog.Time("time", time.Now()))
				continue
			}
			break
		}
		go srv.handleConnection(ctx, conn)
	}
}

func (srv Server) handleConnection(ctx context.Context, conn net.Conn) {
	ctx, ctxCancel := context.WithTimeout(ctx, time.Second*connTimeout)
	defer ctxCancel()

	defer conn.Close()
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	encoder := json.NewEncoder(conn)

	for scanner.Scan() {
		request, valid, err := UnmarshalRequest(scanner.Text())
		if err != nil {
			slog.Error("UnmarshalRequest failed:",
				slog.Any("error", err),
				slog.String("id", "problem1/server/E003"),
				slog.Time("time", time.Now()))
			break
		}
		if valid {
			response := ValidResponse{Method: "isPrime", Prime: IsPrime(*request.Number)}
			encoder.Encode(response)
		} else {
			response := InvalidResponse{Status: "invalid"}
			encoder.Encode(response)
			break
		}
	}
	if err := scanner.Err(); err != nil {
		slog.Error("scanner.Scan failed:",
			slog.Any("error", err),
			slog.String("id", "problem1/server/E004"),
			slog.Time("time", time.Now()))
	}
}
