package problem2

import (
	"encoding/binary"
	"unicode/utf8"
)

type Query struct {
	Type rune
	Num1 int32
	Num2 int32
}

func NewQuery(data []byte) (query Query, err error) {
	query.Type, _ = utf8.DecodeRune(data[:1])
	_, err = binary.Decode(data[1:5], binary.BigEndian, &query.Num1)
	if err != nil {
		return query, err
	}
	_, err = binary.Decode(data[5:9], binary.BigEndian, &query.Num2)
	if err != nil {
		return query, err
	}
	return query, nil
}

func NewResponse(meanPrice int32) (buffer []byte, err error) {
	buffer = make([]byte, 4)
	_, err = binary.Encode(buffer, binary.BigEndian, meanPrice)
	if err != nil {
		return buffer, err
	}
	return buffer, nil
}
