package logger

import "testing"

func TestDebug(t *testing.T) {
	i := map[string]string{"key": "value"}
	s := "string"
	Debug(i)
	Debug(s)
}
func TestInfo(t *testing.T) {
	data := map[string]string{"key": "value"}
	Info(data)
}

func TestWarn(t *testing.T) {
	Warn("data")
}
