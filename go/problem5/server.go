package problem5

// https://protohackers.com/problem/5

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"net"
	"strings"
)

type Server struct {
	host string
	port string
}

const (
	upstreamHost   = "chat.protohackers.com"
	upstreamPort   = "16963"
	paymentAddress = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
)

func NewServer(host, port string) Server {
	return Server{host, port}
}

func (srv *Server) Run(ctx context.Context) {
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
		downConn, err := listener.Accept()
		if err != nil {
			if ctx.Err() == nil {
				slog.Error("listener.Accept failed:", slog.Any("error", err))
				continue
			}
			break
		}
		upConn, err := net.Dial("tcp", net.JoinHostPort(upstreamHost, upstreamPort))
		if err != nil {
			downConn.Close() // can I do better than this?
			slog.Error("net.Dial failed:", slog.Any("error", err))
			continue
		}
		go srv.handleConnection(ctx, downConn, upConn)
	}
}

func (srv *Server) handleConnection(ctx context.Context, downConn net.Conn, upConn net.Conn) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		downConn.Close()
		upConn.Close()
	}()

	go func() {
		defer cancel()
		srv.handleLines(bufio.NewReader(upConn), bufio.NewWriter(downConn))
	}()
	// defer cancel()
	srv.handleLines(bufio.NewReader(downConn), bufio.NewWriter(upConn))
}

func (srv *Server) handleLines(reader *bufio.Reader, writer *bufio.Writer) {
	for {
		line, err := readLine(reader)
		if err != nil {
			if err != io.EOF {
				slog.Error("readLine failed:", slog.Any("error", err))
			}
			break
		}
		err = writeLine(writer, rewriteAddress(line))
		if err != nil {
			slog.Error("writeLine failed:", slog.Any("error", err))
			break
		}
	}
}

func rewriteAddress(s string) string {
	t := strings.Clone(s)
	for _, str := range strings.Fields(s) {
		if strings.HasPrefix(str, "7") && len(str) >= 26 && len(str) <= 35 {
			t = strings.ReplaceAll(t, str, paymentAddress)
		}
	}
	return t
}
