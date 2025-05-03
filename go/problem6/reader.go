package problem6

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MessageReader interface {
	Read(r io.Reader) error
}

func ReadMessage(r io.Reader) (msg MessageReader, err error) {
	var typ uint8
	err = binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, fmt.Errorf("read type failed: %w", err)
	}
	switch typ {
	case 0x80:
		msg = new(Camera)
		err = msg.Read(r)
		return msg, err
	case 0x20:
		msg = new(Plate)
		err = msg.Read(r)
		return msg, err
	case 0x81:
		msg = new(Dispatcher)
		err = msg.Read(r)
		return msg, err
	case 0x40:
		msg = new(WantHeartbeat)
		err = msg.Read(r)
		return msg, err
	}
	return nil, ErrInvalidMessageType
}

func (m *Camera) Read(r io.Reader) (err error) {
	err = binary.Read(r, binary.BigEndian, m)
	if err != nil {
		return fmt.Errorf("read camera failed: %w", err)
	}
	return nil
}

func (m *Plate) Read(r io.Reader) (err error) {
	var n uint8
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return fmt.Errorf("read plate failed: %w", err)
	}
	b := make([]byte, n)
	err = binary.Read(r, binary.BigEndian, b)
	if err != nil {
		return fmt.Errorf("read plate failed: %w", err)
	}
	m.Plate = string(b)
	err = binary.Read(r, binary.BigEndian, &m.Time)
	if err != nil {
		return fmt.Errorf("read plate failed: %w", err)
	}
	return nil
}

func (m *Dispatcher) Read(r io.Reader) (err error) {
	var n uint8
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return fmt.Errorf("read dispatcher failed: %w", err)
	}
	m.Roads = make([]uint16, n)
	err = binary.Read(r, binary.BigEndian, &m.Roads)
	if err != nil {
		return fmt.Errorf("read dispatcher failed: %w", err)
	}
	return nil
}

func (m *WantHeartbeat) Read(r io.Reader) (err error) {
	err = binary.Read(r, binary.BigEndian, m)
	if err != nil {
		return fmt.Errorf("read heartbeat failed: %w", err)
	}
	return nil
}
