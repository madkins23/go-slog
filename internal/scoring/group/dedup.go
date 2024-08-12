package group

import (
	"github.com/madkins23/go-slog/creator/snqkmeld"
	"github.com/madkins23/go-slog/creator/veqryndedup"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// ----------------------------------------------------------------------------

var dedup *score.Group

func Dedup() *score.Group {
	if dedup == nil {
		dedup = score.NewFilterGroup(
			snqkmeld.Name,
			veqryndedup.Name(veqryndedup.Append),
			veqryndedup.Name(veqryndedup.Ignore),
			veqryndedup.Name(veqryndedup.Increment),
			veqryndedup.Name(veqryndedup.Overwrite),
		)
	}
	return dedup
}
