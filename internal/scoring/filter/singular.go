package filter

import (
	"github.com/madkins23/go-slog/creator/slogjson"
	"github.com/madkins23/go-slog/creator/snqkmeld"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var singular score.Filter

func Singular() score.Filter {
	if singular == nil {
		singular = score.NewIncludeFilter(
			slogjson.Name,
			snqkmeld.Name,
		)
	}
	return singular
}
