package tests

import (
	"errors"
	"log/slog"
	"math"
	"time"
)

type thing = struct {
	Name string
	Num  int
	When time.Time
}

const (
	boolVal   = true
	intVal    = 17
	int64Val  = -369
	uint64Val = 134
	stringVal = "The only thing we have to fear is fear itself."
)

var (
	anything = &thing{Name: "Skidoo", Num: 23, When: time.Now()}
	anError  = errors.New("I'm sorry, Dave. I'm afraid I can't do that.")
)

func allKeyValues() []any {
	return []any{
		"Bool", boolVal,
		"Int", intVal,
		"Int64", int64Val,
		"Float64", math.Pi,
		"Uint64", uint64Val,
		"String", stringVal,
		"Time", time.Now(),
		"Duration", time.Hour,
		"Any", anything,
		"Error", anError,
	}
}

// allAttributes has an example of each slog.Attr except group.
func allAttributes() []slog.Attr {
	return []slog.Attr{
		slog.Bool("Bool", boolVal),
		slog.Int("Int", intVal),
		slog.Int64("Int64", int64Val),
		slog.Float64("Float64", math.Pi),
		slog.Uint64("Uint64", uint64Val),
		slog.String("String", stringVal),
		slog.Time("Time", time.Now()),
		slog.Duration("Duration", time.Hour),
		slog.Any("Any", anything),
		slog.Any("Error", anError),
	}
}
