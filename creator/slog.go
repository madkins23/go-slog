package creator

import (
	"io"
	"log/slog"

	"github.com/madkins23/go-slog/infra"
)

// Slog returns a Creator object for the log/slog.JSONHandler.
func Slog() infra.Creator {
	return &slogCreator{CreatorData: infra.NewCreatorData("log/slog.JSONHandler")}
}

type slogCreator struct {
	infra.CreatorData
}

func (c *slogCreator) NewLogger(w io.Writer, options *slog.HandlerOptions) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, options))
}
