package main

import (
	"fmt"
	"log/slog"
	"os"
)

var _ fmt.Stringer = &stringer{}

type stringer struct {
	Value string
}

func (s *stringer) String() string {
	return s.Value
}

var _ slog.LogValuer = &valuer{}

type valuer struct {
	Value any
}

func (v *valuer) LogValue() slog.Value {
	return slog.AnyValue(v.Value)
}

func main() {
	s := &stringer{Value: "string"}
	v := &valuer{Value: "value"}
	fmt.Printf("fmt:  %s\n", s)
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	jsonLogger.Info("json", "stringer", s, "valuer", v)
	textLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	textLogger.Info("text", "stringer", s, "valuer", v)
}
