package failure

import "strconv"

// Code represents an error Code.
// The code should not have internal state, so it should be
// defined as a variable.
// StringCode or IntCode are recommended if you don't need
// custom behavior on the code.
type Code interface {
	// ErrorCode returns an error Code in string representation.
	ErrorCode() string
}

// CodeGetter error code getter
type CodeGetter interface {
	GetCode() string
}

// StringCode represents an error Code in string.
type StringCode string

// ErrorCode implements the Code interface.
func (c StringCode) ErrorCode() string {
	return string(c)
}

// IntCode represents an error Code in int64.
type IntCode int64

// ErrorCode implements the Code interface.
func (c IntCode) ErrorCode() string {
	return strconv.FormatInt(int64(c), 10)
}

// Is checks error type to match with given type
func Is(err error, codes ...Code) bool {
	if err == nil {
		return false
	}

	if _err, ok := err.(CodeGetter); ok {
		errCode := _err.GetCode()

		for _, code := range codes {
			if errCode == code.ErrorCode() {
				return true
			}
		}
	}

	return false
}
