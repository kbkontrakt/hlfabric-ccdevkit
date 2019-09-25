package utils

import "encoding/json"

type (
	// MarshalFunc .
	MarshalFunc func(interface{}) ([]byte, error)

	// UnmarshalFunc .
	UnmarshalFunc func([]byte, interface{}) error
)

var (
	// MarshalFuncJSON .
	MarshalFuncJSON = json.Marshal

	// UnmarshalFuncJSON .
	UnmarshalFuncJSON = json.Unmarshal
)
