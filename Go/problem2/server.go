package problem2

// https://protohackers.com/problem/2

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"net"
	"time"
)

type Server struct {
	host string
	port string
}

const (
	connTimeout  = 60
	ioBufferSize = 9
)

func NewServer(host, port string) Server {
	return Server{host, port}
}

func (srv Server) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", net.JoinHostPort(srv.host, srv.port))
	if err != nil {
		slog.Error("net.Listen failed:",
			slog.Any("error", err),
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
				slog.Error("listener.Accept failed:", slog.Any("error", err))
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

	asset := NewAsset()

	reader := bufio.NewReader(conn)
	buffer := make([]byte, ioBufferSize)

	for {
		_, err := io.ReadFull(reader, buffer)
		if err != nil {
			if err != io.EOF {
				slog.Error("io.ReadFull failed:", slog.Any("error", err))
			}
			return
		}
		query, err := NewQuery(buffer)
		if err != nil {
			slog.Error("NewQuery failed:", slog.Any("error", err))
			return
		}

		switch query.Type {
		case 'I':
			asset.AddPrice(query.Num1, query.Num2)
		case 'Q':
			response, err := NewResponse(asset.MeanPrice(query.Num1, query.Num2))
			if err != nil {
				slog.Error("NewResponse failed:", slog.Any("error", err))
				return
			}
			_, err = conn.Write(response)
			if err != nil {
				slog.Error("conn.Write failed:", slog.Any("error", err))
				return
			}
		}
	}
}
