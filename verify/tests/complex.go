package tests

import (
	"log/slog"

	"github.com/madkins23/go-utils/msg"

	"github.com/madkins23/go-slog/infra"
)

type test struct {
	name, summary string
	definition string
}

func (t *test) execute(logger *slog.Logger) (map[string]any, error) {
	return nil, &msg.ErrNotImplemented{Name: "execute"}
}

var cases = map[string]*test

func (suite *SlogTestSuite) TestComplexCases() {
	logger := suite.Logger(infra.SimpleOptions())
	for _, test := range cases {
		// TODO: Start here if I get back to this.
		x, err := test.execute(logger)
	}
}
