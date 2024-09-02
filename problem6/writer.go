package problem6

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type MessageWriter interface {
	Write(w io.Writer) error
}

func WriteMessage(w io.Writer, msg MessageWriter) error {
	return msg.Write(w)
}

func (m Error) Write(w io.Writer) (err error) {
	b := new(bytes.Buffer)
	data := []any{uint8(0x10), uint8(len(m.Msg)), []byte(m.Msg)}
	for _, d := range data {
		err := binary.Write(b, binary.BigEndian, d)
		if err != nil {
			return fmt.Errorf("binary write error failed: %w", err)
		}
	}
	_, err = w.Write(b.Bytes())
	if err != nil {
		return fmt.Errorf("write error failed: %w", err)
	}
	return nil
}

func (m Ticket) Write(w io.Writer) (err error) {
	b := new(bytes.Buffer)
	data := []any{
		uint8(0x21), uint8(len(m.Plate)), []byte(m.Plate),
		m.Road, m.Mile1, m.Time1, m.Mile2, m.Time2, m.Speed,
	}
	for _, d := range data {
		err := binary.Write(b, binary.BigEndian, d)
		if err != nil {
			return fmt.Errorf("binary write ticket failed: %w", err)
		}
	}
	_, err = w.Write(b.Bytes())
	if err != nil {
		return fmt.Errorf("write ticket failed: %w", err)
	}
	return nil
}

func (m Heartbeat) Write(w io.Writer) (err error) {
	b := new(bytes.Buffer)
	err = binary.Write(b, binary.BigEndian, uint8(0x41))
	if err != nil {
		return fmt.Errorf("binary write heartbeat failed: %w", err)
	}
	_, err = w.Write(b.Bytes())
	if err != nil {
		return fmt.Errorf("write heartbeat failed: %w", err)
	}
	return err
}
