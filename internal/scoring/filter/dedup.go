package filter

import (
	"github.com/madkins23/go-slog/creator/slogjson"
	"github.com/madkins23/go-slog/internal/scoring/group"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// ----------------------------------------------------------------------------

var dedup score.Filter

func Dedup() score.Filter {
	if dedup == nil {
		dedup = score.NewIncludeFilter(
			slogjson.Name,
			group.Dedup(),
		)
	}
	return dedup
}
