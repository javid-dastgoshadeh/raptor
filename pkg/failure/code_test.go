package failure

import (
	"errors"
	"testing"
)

func TestStringCode(t *testing.T) {
	const MError StringCode = "MError"

	if MError.ErrorCode() != "MError" {
		t.Error("Error code doesn't match with original")
	}
}

func TestIntCode(t *testing.T) {
	const MError IntCode = 12

	if MError.ErrorCode() != "12" {
		t.Error("Error code doesn't match with original")
	}
}

func TestIs(t *testing.T) {
	const MError StringCode = "MError"
	var OError = errors.New("MError")

	ei := New(MError, "New error")

	if !Is(ei, MError) {
		t.Error("Error type doesn't match")
	}

	if Is(OError, MError) {
		t.Error("Type doesn't match but we assert that")
	}

	if Is(nil, MError) {
		t.Error("nil type must be return false")
	}
}
