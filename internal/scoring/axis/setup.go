package axis

import "fmt"

func Setup() error {
	if err := setupBenchmarks(); err != nil {
		return fmt.Errorf("axis.setupBenchmarks: %w", err)
	}
	if err := setupWarnings(); err != nil {
		return fmt.Errorf("axis.setupWarnings: %w", err)
	}
	return nil
}
