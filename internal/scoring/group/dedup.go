package group

import (
	"github.com/madkins23/go-slog/creator/snqkmeld"
	"github.com/madkins23/go-slog/creator/veqryndedup"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// ----------------------------------------------------------------------------

var dedupAll *score.Group

func DedupAll() *score.Group {
	if dedupAll == nil {
		dedupAll = score.NewFilterGroup(
			veqryndedup.Name(veqryndedup.Append),
			veqryndedup.Name(veqryndedup.Ignore),
			veqryndedup.Name(veqryndedup.Increment),
			veqryndedup.Name(veqryndedup.Overwrite),
			snqkmeld.Name,
		)
	}
	return dedupAll
}
