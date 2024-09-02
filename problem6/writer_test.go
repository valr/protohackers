package problem6

import (
	"bytes"
	"testing"
)

func TestWriteMessage(t *testing.T) {
	type args struct {
		msg MessageWriter
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			"testError1",
			args{Error{"bad"}},
			string([]byte{0x10, 0x03, 0x62, 0x61, 0x64}),
			false,
		},
		{
			"testError2",
			args{Error{"illegal msg"}},
			string([]byte{
				0x10, 0x0b, 0x69, 0x6c, 0x6c, 0x65,
				0x67, 0x61, 0x6c, 0x20, 0x6d, 0x73, 0x67,
			}),
			false,
		},
		{
			"testWriteTicket1",
			args{Ticket{"UN1X", 66, 100, 123456, 110, 123816, 10000}},
			string([]byte{
				0x21, 0x04, 0x55, 0x4e, 0x31, 0x58, 0x00, 0x42, 0x00, 0x64, 0x00,
				0x01, 0xe2, 0x40, 0x00, 0x6e, 0x00, 0x01, 0xe3, 0xa8, 0x27, 0x10,
			}),
			false,
		},
		{
			"testWriteTicket2",
			args{Ticket{"RE05BKG", 368, 1234, 1000000, 1235, 1000060, 6000}},
			string([]byte{
				0x21, 0x07, 0x52, 0x45, 0x30, 0x35, 0x42, 0x4b, 0x47, 0x01, 0x70,
				0x04, 0xd2, 0x00, 0x0f, 0x42, 0x40, 0x04, 0xd3, 0x00, 0x0f, 0x42,
				0x7c, 0x17, 0x70,
			}),
			false,
		},
		{"testHeartbeat", args{Heartbeat{}}, string([]byte{0x41}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := WriteMessage(w, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("WriteMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("WriteMessage() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
