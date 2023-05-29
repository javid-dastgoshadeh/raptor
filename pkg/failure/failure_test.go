package failure

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestE(t *testing.T) {
	const ErrUndefined StringCode = "ErrUndefined"

	a := New(ErrUndefined, "Error message")

	if a.GetCode() != "ErrUndefined" {
		t.Error("Error code doesn't match")
	}

	if a.Error() != "Error message" {
		t.Error("Error value doesn't match")
	}

	a = New(ErrUndefined, map[string]string{"name": "Ali"})

	m, _ := json.Marshal(map[string]string{"name": "Ali"})

	if a.Error() != string(m) {
		t.Error("marshaled map doesn't match with original")
	}

	a = New(ErrUndefined, errors.New("This is error"))

	if a.Error() != "This is error" {
		t.Error("error output doesn't match with original")
	}
}
