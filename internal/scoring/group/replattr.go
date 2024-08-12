package group

import (
	"github.com/madkins23/go-slog/creator/madkinsreplattr"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// ----------------------------------------------------------------------------

var replAttrOnly *score.Group

func ReplAttrOnly() *score.Group {
	if replAttrOnly == nil {
		replAttrOnly = score.NewFilterGroup(
			madkinsreplattr.Name,
		)
	}
	return replAttrOnly
}
