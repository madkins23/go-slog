package filter

import (
	"github.com/madkins23/go-slog/creator/madkinsreplattr"
	"github.com/madkins23/go-slog/creator/snqkmeld"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var basic score.Filter

func Basic() score.Filter {
	if basic == nil {
		basic = score.NewExcludeFilter(
			madkinsreplattr.Name,
			snqkmeld.Name,
		)
	}
	return basic
}
