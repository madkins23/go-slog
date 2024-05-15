package scoring

import (
	"fmt"

	"github.com/madkins23/go-slog/internal/data"
	"github.com/madkins23/go-slog/internal/scoring/axis"
	"github.com/madkins23/go-slog/internal/scoring/exhibit"
	"github.com/madkins23/go-slog/internal/scoring/keeper"
)

func Setup(bench *data.Benchmarks, warns *data.Warnings) error {
	if err := exhibit.Setup(); err != nil {
		return fmt.Errorf("exhibit.Setup: %w", err)
	}
	if err := axis.Setup(); err != nil {
		return fmt.Errorf("axis.Setup: %w", err)
	}
	if err := keeper.Setup(bench, warns); err != nil {
		return fmt.Errorf("keeper.Setup: %w", err)
	}
	return nil
}
