// package template contains code that will execute the verification and benchmark
// test suites in the [`madkins23/go-slog`](https://github.com/madkins23/go-slog) repository.
//
// # Installation
//
// Copy this directory to own repository with an appropriate name.
// Alter the code:
//
//   - match the package name to the directory name
//   - fix the Creator function to return one of your `slog.Handler` objects
//
// # Usage
//
// Run verification suite:
//
//	go test -v -args -useWarnings
//
// Run Benchmark suite:
//
//	go test -v -run=none -bench=. -args -useWarnings
package template

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/suite"

	benchtests "github.com/madkins23/go-slog/bench/tests"
	"github.com/madkins23/go-slog/infra"
	"github.com/madkins23/go-slog/infra/warning"
	verifytests "github.com/madkins23/go-slog/verify/tests"
)

// BenchmarkSlogJSON runs benchmarks for the slog/JSONHandler JSON handler.
func BenchmarkSlogJSON(b *testing.B) {
	slogSuite := benchtests.NewSlogBenchmarkSuite(Creator())
	benchtests.Run(b, slogSuite)
}

// TestVerifySlogJSON runs tests for the slog/JSONHandler JSON handler.
func TestVerifySlogJSON(t *testing.T) {
	slogSuite := verifytests.NewSlogTestSuite(Creator())
	slogSuite.WarnOnly(warning.Duplicates)
	suite.Run(t, slogSuite)
}

// Creator returns a Creator object for the [slog/JSONHandler] handler.
//
// [slog/JSONHandler]: https://pkg.go.dev/log/slog#JSONHandler
func Creator() infra.Creator {
	return infra.NewCreator("slog/JSONHandler", handlerFn, nil,
		`^slog/JSONHandler^ is the JSON handler provided with the ^slog^ library.
		It is fast and as a part of the Go distribution it is used
		along with published documentation as a model for ^slog.Handler^ behavior.`,
		map[string]string{
			"slog/JSONHandler": "https://pkg.go.dev/log/slog#JSONHandler",
		})
}

func handlerFn(w io.Writer, options *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, options)
}

func TestMain(m *testing.M) {
	warning.WithWarnings(m)
}
