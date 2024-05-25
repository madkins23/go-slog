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
	if err := setupSize(); err != nil {
		return fmt.Errorf("keeper.setupSize: %w", err)
	}
	for _, tag := range score.Keepers() {
		keeper := score.GetKeeper(tag)
		if err := keeper.Setup(bench, warns); err != nil {
			return fmt.Errorf("setup '%s': %w", keeper.Tag(), err)
		}
	}
	return nil
}
