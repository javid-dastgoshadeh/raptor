package failure

import (
	"encoding/json"
	"fmt"
)

// SError general error type
type SError struct {
	code Code
	err  interface{}
}

// Error implements `error` interface
func (e *SError) Error() string {
	var m string

	switch e.err.(type) {
	case error:
		m = e.err.(error).Error()

		break
	case string:
		m = e.err.(string)

		break
	default:
		_m, err := json.Marshal(e.err)

		if err != nil {
			return fmt.Sprintf("%v", e.err)
		}

		m = string(_m)
	}

	return m
}

// MarshalJSON implement json marshal interface
func (e SError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.err)
}

// GetCode gets error code
func (e *SError) GetCode() string {
	return e.code.ErrorCode()
}

// GetErr gets error data
func (e *SError) GetErr() interface{} {
	return e.err
}

// New make new error type
func New(code Code, err interface{}) *SError {
	return &SError{code: code, err: err}
}
