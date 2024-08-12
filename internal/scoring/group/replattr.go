package group

import (
	"github.com/madkins23/go-slog/creator/madkinsreplattr"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// ----------------------------------------------------------------------------

var replAttr *score.Group

func ReplAttr() *score.Group {
	if replAttr == nil {
		replAttr = score.NewFilterGroup(
			madkinsreplattr.Name,
		)
	}
	return replAttr
}
