package filter

import (
	"github.com/madkins23/go-slog/internal/scoring/group"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var basic score.Filter

func Basic() score.Filter {
	if basic == nil {
		basic = score.NewExcludeFilter(
			group.ReplAttr(),
			group.Dedup(),
		)
	}
	return basic
}
