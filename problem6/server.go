package problem6

// https://protohackers.com/problem/6

import (
	"context"
	"log/slog"
	"maps"
	"math"
	"net"
	"slices"
	"sync"
	"time"

	"protohackers/util"
)

type Server struct {
	host          string
	port          string
	obs           Observation        // stored observations (per road/plate)
	obsToCheck    ObservationTocheck // stored observations to check for limit exceeded
	obsMu         sync.Mutex
	obsCh         chan CamObservation // channel receiving observations to store
	tktComputed   TicketComputed      // tickets computed (per plate/day)
	tktToDispatch TicketToDispatch    // stored tickets to dispatch
	tktMu         sync.Mutex
	tktCh         map[uint16]chan Ticket // channel sending tickets to dispatch
}

const (
	obsChSize  = 500
	tktChSize  = 50
	checkDelay = 5
)

func NewServer(host, port string) Server {
	return Server{
		host, port,
		make(Observation),
		make(ObservationTocheck),
		sync.Mutex{},
		make(chan CamObservation, obsChSize),
		make(TicketComputed),
		make(TicketToDispatch),
		sync.Mutex{},
		make(map[uint16]chan Ticket),
	}
}

func (srv *Server) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", net.JoinHostPort(srv.host, srv.port))
	if err != nil {
		slog.Error("net.Listen failed:",
			slog.Any("error", err),
			slog.String("host", srv.host),
			slog.String("port", srv.port),
			slog.String("source", util.SourceInfo()))
		return
	}
	defer listener.Close()
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	go srv.StoreObservation()
	go srv.CheckStoredObservation()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() == nil {
				slog.Error("listener.Accept failed:",
					slog.Any("error", err),
					slog.String("source", util.SourceInfo()))
				continue
			}
			break
		}
		wrk := NewWorker(srv)
		go wrk.Run(ctx, conn)
	}
}

func (srv *Server) GetObsCh() chan<- CamObservation {
	return srv.obsCh
}

func (srv *Server) GetTktCh(road uint16) <-chan Ticket {
	var ch chan Ticket
	srv.tktMu.Lock()
	ch, ok := srv.tktCh[road]
	if !ok {
		ch = make(chan Ticket, tktChSize)
		srv.tktCh[road] = ch
		go srv.DispatchStoredTicket(road)
	}
	srv.tktMu.Unlock()
	return ch
}

func (srv *Server) StoreObservation() {
	for o := range srv.obsCh {
		k := KeyObservation{o.Road, o.Plate}
		srv.obsMu.Lock()
		_, ok := srv.obs[k]
		if !ok {
			srv.obs[k] = make(ValObservation)
		}
		srv.obs[k][o.Time] = o.Mile
		srv.obsToCheck[k] = o.Limit
		srv.obsMu.Unlock()
	}
}

func (srv *Server) CheckStoredObservation() {
	for {
		<-time.After(time.Duration(checkDelay) * time.Second)
		srv.obsMu.Lock()
		for ko, l := range srv.obsToCheck {
			srv.ComputeNewTicket(ko, srv.obs[ko], l)
		}
		clear(srv.obsToCheck)
		srv.obsMu.Unlock()
	}
}

func (srv *Server) ComputeNewTicket(kObs KeyObservation, obs ValObservation, lim uint16) {
	obsLen := len(obs)
	if obsLen < 2 {
		return
	}
	obsKey := slices.Sorted(maps.Keys(obs))
	for i := range obsLen - 1 {
		time1, time2 := obsKey[i], obsKey[i+1]
		mile1, mile2 := obs[time1], obs[time2]
		speed := math.Abs(3600 * (float64(mile2) - float64(mile1)) / (float64(time2) - float64(time1)))
		if speed > float64(lim) {
			road, plate := kObs.Road, kObs.Plate
			day1 := uint32(math.Floor(float64(time1) / 86400))
			day2 := uint32(math.Floor(float64(time2) / 86400))

			var computed bool
			for day := day1; day <= day2; day++ {
				if srv.tktComputed[KeyTicket{plate, day}] {
					computed = true
				}
			}

			if !computed {
				for day := day1; day <= day2; day++ {
					srv.tktComputed[KeyTicket{plate, day}] = true
				}
				ticket := Ticket{
					plate, road,
					mile1, time1,
					mile2, time2,
					uint16(speed*100 + 0.5),
				}
				srv.DispatchOrStoreTicket(road, ticket)
			}
		}
	}
}

func (srv *Server) DispatchOrStoreTicket(road uint16, ticket Ticket) {
	srv.tktMu.Lock()
	ch, ok := srv.tktCh[road]
	if ok {
		ch <- ticket
	} else {
		srv.tktToDispatch[road] = append(srv.tktToDispatch[road], ticket)
	}
	srv.tktMu.Unlock()
}

func (srv *Server) DispatchStoredTicket(road uint16) {
	srv.tktMu.Lock()
	for _, ticket := range srv.tktToDispatch[road] {
		srv.tktCh[road] <- ticket
	}
	delete(srv.tktToDispatch, road)
	srv.tktMu.Unlock()
}
