package problem0

// https://protohackers.com/problem/0

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
	connTimeout  = 10
	ioBufferSize = 4096
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
	// Yes, I know I could do: io.Copy(conn, conn) instead of below code

	ctx, ctxCancel := context.WithTimeout(ctx, time.Second*connTimeout)
	defer ctxCancel()

	defer conn.Close()
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	buffer := make([]byte, ioBufferSize)
	reader := bufio.NewReader(conn)

	for {
		n, errRead := reader.Read(buffer)
		if n > 0 {
			_, errWrite := conn.Write(buffer[:n])
			if errWrite != nil {
				slog.Error("conn.Write failed:", slog.Any("error", errWrite))
				break
			}
		}
		if errRead != nil {
			if errRead != io.EOF {
				slog.Error("reader.Read failed:", slog.Any("error", errRead))
			}
			break
		}
	}
}
