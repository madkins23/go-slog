package keeper

import (
	"fmt"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/score"
)

func Setup(bench *data.Benchmarks, warns *data.Warnings) error {
	if err := setupDefault(); err != nil {
		return fmt.Errorf("keeper.setupDefault: %w", err)
	}
	if err := setupSimple(); err != nil {
		return fmt.Errorf("keeper.setupSimple: %w", err)
	}
	for name, tag := range score.Keepers() {
		keeper := score.GetKeeper(tag)
		if err := keeper.Setup(bench, warns); err != nil {
			return fmt.Errorf("setup '%s': %w", name, err)
		}
	}
	return nil
}