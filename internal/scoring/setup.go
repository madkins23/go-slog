package scoring

import (
	"fmt"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/score"
)
import "github.com/madkins23/go-slog/internal/scoring/keeper"

func Initialize(bench *data.Benchmarks, warns *data.Warnings) error {
	if err := keeper.Default(); err != nil {
		return fmt.Errorf("keeper.Default: %w", err)
	}
	if err := keeper.Simple(); err != nil {
		return fmt.Errorf("keeper.Simple: %w", err)
	}
	if err := score.Initialize(bench, warns); err != nil {
		return fmt.Errorf("score.Initialize: %w", err)
	}
	return nil
}
