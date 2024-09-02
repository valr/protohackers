package problem6

// Client -> Server

type Camera struct {
	Road  uint16
	Mile  uint16
	Limit uint16
}

type Plate struct {
	Plate string
	Time  uint32
}

type Dispatcher struct {
	Roads []uint16
}

type WantHeartbeat struct {
	Interval uint32
}

// Server -> Client

type Error struct {
	Msg string
}

type Ticket struct {
	Plate string
	Road  uint16
	Mile1 uint16
	Time1 uint32
	Mile2 uint16
	Time2 uint32
	Speed uint16
}

type Heartbeat struct{}

// Server internal

type CamObservation struct {
	Road  uint16
	Plate string
	Mile  uint16
	Time  uint32
	Limit uint16
}

type (
	KeyObservation struct {
		Road  uint16
		Plate string
	}
	ValObservation     map[uint32]uint16 // Time: uint32, Mile: uint16
	Observation        map[KeyObservation]ValObservation
	ObservationTocheck map[KeyObservation]uint16 // Limit: uint16
)

type (
	KeyTicket struct {
		Plate string
		Day   uint32
	}
	TicketComputed   map[KeyTicket]bool
	TicketToDispatch map[uint16][]Ticket // Road: uint16
)
