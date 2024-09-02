package problem6

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestReadMessage(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantMsg MessageReader
		wantErr bool
	}{
		{
			"testIAmCamera1",
			args{bytes.NewReader([]byte{0x80, 0x00, 0x42, 0x00, 0x64, 0x00, 0x3c})},
			&Camera{66, 100, 60},
			false,
		},
		{
			"testIAmCamera2",
			args{bytes.NewReader([]byte{0x80, 0x01, 0x70, 0x04, 0xd2, 0x00, 0x28})},
			&Camera{368, 1234, 40},
			false,
		},
		{
			"testPlate1",
			args{bytes.NewReader(
				[]byte{0x20, 0x04, 0x55, 0x4e, 0x31, 0x58, 0x00, 0x00, 0x03, 0xe8})},
			&Plate{"UN1X", 1000},
			false,
		},
		{
			"testPlate2",
			args{bytes.NewReader(
				[]byte{0x20, 0x07, 0x52, 0x45, 0x30, 0x35, 0x42, 0x4b, 0x47, 0x00, 0x01, 0xe2, 0x40})},
			&Plate{"RE05BKG", 123456},
			false,
		},
		{
			"testIAmDispatcher1",
			args{bytes.NewReader([]byte{0x81, 0x01, 0x00, 0x42})},
			&Dispatcher{[]uint16{66}},
			false,
		},
		{
			"testIAmDispatcher2",
			args{bytes.NewReader([]byte{0x81, 0x03, 0x00, 0x42, 0x01, 0x70, 0x13, 0x88})},
			&Dispatcher{[]uint16{66, 368, 5000}},
			false,
		},
		{
			"testWantHeartbeat1",
			args{bytes.NewReader([]byte{0x40, 0x00, 0x00, 0x00, 0x0a})},
			&WantHeartbeat{10},
			false,
		},
		{
			"testWantHeartbeat2",
			args{bytes.NewReader([]byte{0x40, 0x00, 0x00, 0x04, 0xdb})},
			&WantHeartbeat{1243},
			false,
		},
		{
			"testInvalidType",
			args{bytes.NewReader([]byte{0xff})},
			nil,
			true,
		},
		{
			"testEmpty",
			args{bytes.NewReader([]byte{})},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, err := ReadMessage(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMsg, tt.wantMsg) {
				t.Errorf("ReadMessage() = %v, want %v", gotMsg, tt.wantMsg)
			}
		})
	}
}
