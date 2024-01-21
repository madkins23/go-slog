package tests

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"time"
)

type Thing struct {
	Name string
	Num  int
	When time.Time
}

// -----------------------------------------------------------------------------

const (
	valBool        = true
	valInt         = 17
	valInt64       = int64(-654)
	valUInt64      = uint64(29399459393)
	valFloat64     = math.Pi
	valDuration    = time.Hour
	valString      = "The only Thing we have to fear is fear itself."
	valGroupName   = "group"
	valGroupOthers = "others"
	valGroupKey1   = "alpha"
	valGroupVal1   = "omega"
	valGroupKey2   = "never"
	valGroupVal2   = "mind"
)

var (
	valKVs   = []any{valGroupKey1, valGroupVal1, valGroupKey2, valGroupVal2}
	valGroup = slog.Group(valGroupName, valKVs...)
	valTime  = time.Now()
	valAny   = &Thing{Name: "Skidoo", Num: 23, When: time.Now()}
	valError = errors.New("I'm sorry, Dave. I'm afraid I can't do that.")
)

var allKeyValues []any = []any{
	"Bool", valBool,
	"Int", valInt,
	"Int64", valInt64,
	"Uint64", valUInt64,
	"Float64", valFloat64,
	"String", valString,
	"Time", valTime,
	"Duration", valDuration,
	slog.Group(valGroupName, valKVs...),
	"Any", valAny,
	"Error", valError,
}

var allAttributes []slog.Attr = []slog.Attr{
	slog.Bool("Bool", valBool),
	slog.Int("Int", valInt),
	slog.Int64("Int64", valInt64),
	slog.Uint64("Uint64", valUInt64),
	slog.Float64("Float64", math.Pi),
	slog.String("String", valString),
	slog.Time("Time", valTime),
	slog.Duration("Duration", valDuration),
	slog.Group(valGroupName, valKVs...),
	slog.Any("Any", valAny),
	slog.Any("Error", valError),
}

var withAttributes = []slog.Attr{
	slog.Bool("withBool", valBool),
	slog.Int("withInt", valInt),
	slog.Int64("withInt64", valInt64),
	slog.Uint64("withUint64", valUInt64),
	slog.Float64("withFloat64", math.Pi),
	slog.String("withString", valString),
	slog.Time("withTime", valTime),
	slog.Duration("withDuration", valDuration),
	slog.Group(valGroupName, valKVs...),
	slog.Any("Any", valAny),
	slog.Any("Error", valError),
}

var bigGroup slog.Attr

func BigGroup() slog.Attr {
	if bigGroup.Equal(slog.Attr{}) {
		bigGroup = bigGroupBuilder(0, 5, valGroupName)
	}

	return bigGroup
}

func bigGroupBuilder(depth, limit uint, stem string) slog.Attr {
	if depth >= limit {
		return slog.Attr{}
	}
	name := fmt.Sprintf("%s-%d", stem, depth)
	others := make([]any, 5)
	count := rand.Intn(5)
	for i := 0; i < count; i++ {
		other := bigGroupBuilder(depth+1, limit, name)
		if !other.Equal(slog.Attr{}) {
			others = append(others, other)
		}
	}
	return slog.Group(name, others...)
}
