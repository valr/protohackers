package problem4

// https://protohackers.com/problem/4

import (
	"context"
	"log/slog"
	"net"
	"strings"
)

type Server struct {
	host string
	port string
	db   map[string]string
}

const (
	ioBufferSize   = 65535
	maxRequestSize = 999
	versionKey     = "version"
	versionValue   = "Ken's k-v Store 1.0"
)

func NewServer(host, port string) *Server {
	server := Server{host, port, make(map[string]string)}
	server.db[versionKey] = versionValue
	return &server
}

func (srv *Server) Run(ctx context.Context) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(srv.host, srv.port))
	if err != nil {
		slog.Error("net.ResolveUDPAddr failed:",
			slog.Any("error", err),
			slog.String("host", srv.host),
			slog.String("port", srv.port))
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("net.ListenUDP failed:", slog.Any("error", err))
		return
	}

	defer conn.Close()
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	buffer := make([]byte, ioBufferSize)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			slog.Error("conn.ReadFromUDP failed:", slog.Any("error", err))
			continue
		}
		if n > maxRequestSize {
			slog.Error("request too big")
			continue
		}
		resp := srv.handleRequest(string(buffer[:n]))
		if len(resp) > 0 {
			_, err = conn.WriteToUDP([]byte(resp), addr)
			if err != nil {
				slog.Error("conn.WriteToUDP failed:", slog.Any("error", err))
				continue
			}
		}
	}
}

func (srv *Server) handleRequest(request string) (response string) {
	k, v, ok := strings.Cut(request, "=")
	if ok {
		if k != versionKey {
			srv.db[k] = v
		}
	} else {
		response = k + "=" + srv.db[k]
	}
	return response
}
