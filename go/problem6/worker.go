package problem6

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"protohackers/util"
)

const (
	iAmNobody = iota
	iAmCamera
	iAmDispatcher
)

const (
	wgTimeout = 10
)

var (
	ErrInvalidMessageType   = errors.New("worker: invalid message type")
	ErrTooManyWantHeartbeat = errors.New("worker: heartbeat already on")
	ErrIAmAlreadyIdentified = errors.New("worker: i am already identified")
	ErrIAmNotYetIdentified  = errors.New("worker: i am not yet identified")
)

type Worker struct {
	srv       *Server
	iAm       int
	camera    Camera
	camObsCh  chan<- CamObservation
	heartbeat bool
}

func NewWorker(srv *Server) Worker {
	wrk := Worker{}
	wrk.srv = srv
	return wrk
}

func (wrk *Worker) Run(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	wg := &sync.WaitGroup{}
	defer util.WaitTimeout(wg, time.Duration(wgTimeout)*time.Second)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		util.WaitTimeout(wg, time.Duration(wgTimeout)*time.Second)
		conn.Close()
	}()

	r := bufio.NewReader(conn)
	w := io.Writer(conn)
MessageLoop:
	for {
		m, err := ReadMessage(r)
		if err != nil {
			switch {
			case errors.Is(err, io.EOF):
			case errors.Is(err, ErrInvalidMessageType):
				if err := WriteMessage(w, Error{err.Error()}); err != nil {
					slog.Error("write message failed:",
						slog.Any("error", err),
						slog.String("source", util.SourceInfo()))
				}
			default:
				slog.Error("read message failed:",
					slog.Any("error", err),
					slog.String("source", util.SourceInfo()))
			}
			break MessageLoop
		}
		switch m := m.(type) {
		case *Camera:
			if err := wrk.ProcessCamera(*m); err != nil {
				if err := WriteMessage(w, Error{err.Error()}); err != nil {
					slog.Error("write message failed:",
						slog.Any("error", err),
						slog.String("source", util.SourceInfo()))
				}
				break MessageLoop
			}
		case *Plate:
			if err := wrk.ProcessPlate(*m); err != nil {
				if err := WriteMessage(w, Error{err.Error()}); err != nil {
					slog.Error("write message failed:",
						slog.Any("error", err),
						slog.String("source", util.SourceInfo()))
				}
				break MessageLoop
			}
		case *Dispatcher:
			if err := wrk.ProcessDispatcher(ctx, wg, w, *m); err != nil {
				if err := WriteMessage(w, Error{err.Error()}); err != nil {
					slog.Error("write message failed:",
						slog.Any("error", err),
						slog.String("source", util.SourceInfo()))
				}
				break MessageLoop
			}
		case *WantHeartbeat:
			if err := wrk.ProcessHeartbeat(ctx, wg, w, *m); err != nil {
				if err := WriteMessage(w, Error{err.Error()}); err != nil {
					slog.Error("write message failed:",
						slog.Any("error", err),
						slog.String("source", util.SourceInfo()))
				}
				break MessageLoop
			}
		}
	}
}

func (wrk *Worker) ProcessCamera(c Camera) error {
	if wrk.iAm != iAmNobody {
		return ErrIAmAlreadyIdentified
	}
	wrk.iAm = iAmCamera
	wrk.camera = c
	wrk.camObsCh = wrk.srv.GetObsCh()
	return nil
}

func (wrk *Worker) ProcessPlate(p Plate) error {
	switch wrk.iAm {
	case iAmNobody:
		return ErrIAmNotYetIdentified
	case iAmDispatcher:
		return ErrInvalidMessageType
	}
	wrk.camObsCh <- CamObservation{wrk.camera.Road, p.Plate, wrk.camera.Mile, p.Time, wrk.camera.Limit}
	return nil
}

func (wrk *Worker) ProcessDispatcher(ctx context.Context, wg *sync.WaitGroup, w io.Writer, d Dispatcher) error {
	if wrk.iAm != iAmNobody {
		return ErrIAmAlreadyIdentified
	}
	wrk.iAm = iAmDispatcher
	for _, r := range d.Roads {
		ch := wrk.srv.GetTktCh(r)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case t := <-ch:
					if err := WriteMessage(w, t); err != nil {
						slog.Error("write ticket message failed:",
							slog.Any("error", err),
							slog.String("source", util.SourceInfo()))
						return
					}
				}
			}
		}()
	}
	return nil
}

func (wrk *Worker) ProcessHeartbeat(ctx context.Context, wg *sync.WaitGroup, w io.Writer, h WantHeartbeat) error {
	if wrk.heartbeat {
		return ErrTooManyWantHeartbeat
	}
	wrk.heartbeat = true
	if h.Interval > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(h.Interval*100) * time.Millisecond):
					if err := WriteMessage(w, Heartbeat{}); err != nil {
						slog.Error("write heartbeat message failed:",
							slog.Any("error", err),
							slog.String("source", util.SourceInfo()))
						return
					}
				}
			}
		}()
	}
	return nil
}
