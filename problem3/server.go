package problem3

// https://protohackers.com/problem/3

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Server struct {
	host string
	port string
	user map[string]chan string
	mux  sync.RWMutex
}

const (
	connTimeout = 60
)

func NewServer(host, port string) Server {
	return Server{host, port, make(map[string]chan string), sync.RWMutex{}}
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

func (srv *Server) handleConnection(ctx context.Context, conn net.Conn) {
	ctx, ctxCancel := context.WithTimeout(ctx, time.Second*connTimeout)
	defer ctxCancel()

	defer conn.Close()
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// prompt the user and receive the name
	name, err := srv.receiveName(reader, writer)
	if err != nil {
		slog.Error("srv.receiveName failed:", slog.Any("error", err))
		return
	}

	// reject duplicate names
	_, ok := srv.user[name]
	if ok {
		return
	}

	srv.mux.Lock()
	srv.user[name] = make(chan string)
	go func() {
		for message := range srv.user[name] {
			if err = writeLine(writer, message); err != nil {
				slog.Error("writeLine failed:", slog.Any("error", err))
				break
			}
		}
	}()
	srv.sendTo(name, fmt.Sprint("* The room contains: ", strings.Join(srv.listNames(name), ", ")))
	srv.sendToAllNolock(name, fmt.Sprintf("* %v has entered the room", name))
	srv.mux.Unlock()

	for {
		message, err := readLine(reader)
		if err != nil {
			if err != io.EOF {
				slog.Error("readLine failed:", slog.Any("error", err))
			}

			srv.mux.Lock()
			close(srv.user[name])
			delete(srv.user, name)
			srv.sendToAllNolock(name, fmt.Sprintf("* %v has left the room", name))
			srv.mux.Unlock()
			break
		}
		srv.sendToAllLock(name, fmt.Sprintf("[%v] %v", name, message))
	}
}

func (srv *Server) receiveName(reader *bufio.Reader, writer *bufio.Writer) (name string, err error) {
	if err = writeLine(writer, "coucou, asv stp"); err != nil {
		return name, err
	}
	if name, err = readLine(reader); err != nil {
		return name, err
	}
	if !regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(name) {
		return name, errors.New("invalid name")
	}
	return name, nil
}

func (srv *Server) listNames(excludedName string) []string {
	names := make([]string, 0)
	for name := range srv.user {
		if name != excludedName {
			names = append(names, name)
		}
	}
	return names
}

func (srv *Server) sendTo(receiver, message string) {
	srv.user[receiver] <- message
}

func (srv *Server) sendToAllLock(excludedReceiver, message string) {
	srv.mux.RLock()
	srv.sendToAllNolock(excludedReceiver, message)
	srv.mux.RUnlock()
}

func (srv *Server) sendToAllNolock(excludedReceiver, message string) {
	for receiver, ch := range srv.user {
		if receiver != excludedReceiver {
			ch <- message
		}
	}
}
