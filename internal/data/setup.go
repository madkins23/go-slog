package data

import (
	"bytes"
	"fmt"
)

// Setup parse both benchmark data and benchmark and verification warnings
// into the provided Benchmarks and Warnings objects.
// These objects do the parsing with the proper order and arguments
// and return any encountered error.
func Setup(bench *Benchmarks, warns *Warnings) error {
	if err := bench.ParseBenchmarkData(nil); err != nil {
		return fmt.Errorf("parse -bench data: %w", err)
	}

	if err := warns.ParseWarningData(
		bytes.NewReader(bench.WarningText()), "Bench", bench.HandlerLookup()); err != nil {
		return fmt.Errorf("parse -bench warnings: %w", err)
	}

	if err := warns.ParseWarningData(nil, "Verify", bench.HandlerLookup()); err != nil {
		return fmt.Errorf("parse -verify warnings: %w", err)
	}

	return nil
}
