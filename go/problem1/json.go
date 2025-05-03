package problem1

import (
	"encoding/json"
	"errors"
)

type Request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type ValidResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

type InvalidResponse struct {
	Status string `json:"status"`
}

func UnmarshalRequest(data string) (request Request, valid bool, err error) {
	err = json.Unmarshal([]byte(data), &request)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &syntaxError) || errors.As(err, &unmarshalTypeError) {
			err = nil
		}
		return
	}
	if request.Method != nil && request.Number != nil &&
		*request.Method == "isPrime" {
		valid = true
	}
	return
}
