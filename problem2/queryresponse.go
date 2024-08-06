package problem2

import (
	"bytes"
	"encoding/binary"
	"unicode/utf8"
)

type Query struct {
	Type rune
	Num1 int32
	Num2 int32
}

func NewQuery(data []byte) (query Query) {
	query.Type, _ = utf8.DecodeRune(data[:1])
	_ = binary.Read(bytes.NewReader(data[1:5]), binary.BigEndian, &query.Num1)
	_ = binary.Read(bytes.NewReader(data[5:9]), binary.BigEndian, &query.Num2)
	return
}

func NewResponse(meanPrice int32) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 4))
	_ = binary.Write(buffer, binary.BigEndian, meanPrice)
	return buffer.Bytes()
}
