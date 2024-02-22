package tests

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/madkins23/go-slog/infra"
)

// -----------------------------------------------------------------------------

func expectedBasic() map[string]any {
	return map[string]any{
		slog.LevelKey:   slog.LevelInfo.String(),
		slog.MessageKey: message,
	}
}

// -----------------------------------------------------------------------------

// Thing is a test object used in logging tests.
// The struct and its fields are public since they must be accessible via reflection.
type Thing struct {
	Name string
	Num  int
}

const (
	valBool      = true
	valInt       = 17
	valInt64     = int64(-654)
	valUInt64    = uint64(29399459393)
	valFloat64   = math.Pi
	valDuration  = time.Hour
	valString    = "The only Thing we have to fear is fear itself."
	valGroupName = "group"
	valGroupKey1 = "alpha"
	valGroupVal1 = "omega"
	valGroupKey2 = "never"
	valGroupVal2 = "mind"
)

var (
	valKVs   = []any{valGroupKey1, valGroupVal1, valGroupKey2, valGroupVal2}
	valTime  = time.Now()
	valAny   = &Thing{Name: "Skidoo", Num: 23}
	valError = errors.New("I'm sorry, Dave. I'm afraid I can't do that.")
)

var allKeyValues = []any{
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

var allAttributes = []slog.Attr{
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

func allValuesMap() map[string]any {
	return map[string]any{
		"Bool":    valBool,
		"Int":     float64(valInt),
		"Int64":   float64(valInt64),
		"Uint64":  float64(valUInt64),
		"Float64": math.Pi,
		"String":  valString,
		valGroupName: map[string]any{
			"alpha": "omega",
			"never": "mind",
		},
		"Any": map[string]any{
			"Name": "Skidoo",
			"Num":  float64(23),
		},
	}
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

func withValuesMap() map[string]any {
	return map[string]any{
		"withBool":    valBool,
		"withInt":     float64(valInt),
		"withInt64":   float64(valInt64),
		"withUint64":  float64(valUInt64),
		"withFloat64": math.Pi,
		"withString":  valString,
		valGroupName: map[string]any{
			"alpha": "omega",
			"never": "mind",
		},
		"Any": map[string]any{
			"Name": "Skidoo",
			"Num":  float64(23),
		},
	}
}

// -----------------------------------------------------------------------------

func BigGroup() slog.Attr {
	return bigGroupBuilder(0, bigGroupLimit, valGroupName)
}

const (
	bigGroupLimit = 5
	bigGroupRand  = 5
)

func bigGroupBuilder(depth, limit uint, stem string) slog.Attr {
	if depth > limit {
		return infra.EmptyAttr()
	}
	count := rand.Intn(bigGroupRand) + 1
	others := []any{"count", count, "depth", depth, "limit", limit}
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("%s-%d", stem, i)
		other := bigGroupBuilder(depth+1, limit, name)
		if !other.Equal(infra.EmptyAttr()) {
			others = append(others, other)
		}
	}
	return slog.Group(stem, others...)
}

var nonGroupFields = map[string]bool{
	"count": true,
	"depth": true,
	"limit": true,
}

func bigGroupCheck(bigGroup map[string]any, depth, limit uint, stem string) (uint, error) {
	maxDepth := depth
	for field, value := range bigGroup {
		subMap, ok := value.(map[string]any)
		if !ok {
			if nonGroupFields[field] {
				continue
			} else {
				return 0, fmt.Errorf("non-group field %s: %v", field, value)
			}
		}
		if !strings.HasPrefix(field, stem) {
			return 0, fmt.Errorf("bad field '%s' prefix '%s'", field, stem)
		}
		if d, err := bigGroupCheck(subMap, depth+1, limit, field); err != nil {
			return 0, err
		} else if d > maxDepth {
			maxDepth = d
		}
	}
	return maxDepth, nil
}
