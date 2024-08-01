package filter

import (
	"github.com/madkins23/go-slog/creator/madkinsflash"
	"github.com/madkins23/go-slog/creator/madkinsreplattr"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

var replAttr score.Filter

func ReplAttr() score.Filter {
	if replAttr == nil {
		replAttr = score.NewIncludeFilter(
			madkinsreplattr.Name,
			madkinsflash.Name,
		)
	}
	return replAttr
}
