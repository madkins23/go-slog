package data

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

//go:embed testdata/bench.txt
var benchTxt string

//go:embed testdata/verify.txt
var verifyTxt string

type ParserTestSuite struct {
	suite.Suite
	bench  *Warnings
	verify *Warnings
}

func (suite *ParserTestSuite) SetupSuite() {
	suite.bench = NewWarnings()
	suite.Require().NoError(suite.bench.ParseWarningData(strings.NewReader(benchTxt), "", nil))
	suite.verify = NewWarnings()
	suite.Require().NoError(suite.verify.ParseWarningData(strings.NewReader(verifyTxt), "", nil))
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (suite *ParserTestSuite) TestData_Parse_Bench_phsym_zerolog() {
	levels := suite.bench.ForHandler(HandlerTag("phsym/zeroslog"))
	suite.Assert().NotNil(levels)
	suite.Assert().Len(levels.Levels(), 1)
	implied, found := levels.lookup["Implied"]
	suite.Assert().True(found)
	suite.Assert().NotNil(implied)
	suite.Assert().Len(implied.Warnings(), 1)
	warning, found := implied.lookup["SourceKey"]
	suite.Assert().True(found)
	suite.Assert().NotNil(warning)
	suite.Assert().Len(warning.instances, 1)
	instance := warning.instances[0]
	suite.Assert().NotNil(instance)
	suite.Assert().Equal("Simple Source", instance.name)
	suite.Assert().Equal("no 'source' key", instance.extra)
	suite.Assert().Contains(instance.log, "{")
}

func (suite *ParserTestSuite) TestData_Parse_Verify_phsym_zerolog() {
	var expectedTestNames = []string{
		"Attribute With Empty",
		"Attributes Empty",
	}
	levels := suite.verify.ForHandler(HandlerTag("phsym/zeroslog"))
	suite.Assert().NotNil(levels)
	suite.Assert().Len(levels.Levels(), 3)
	required, found := levels.lookup["Required"]
	suite.Assert().True(found)
	suite.Assert().NotNil(required)
	suite.Assert().Len(required.Warnings(), 4)
	warning, found := required.lookup["EmptyAttributes"]
	suite.Assert().True(found)
	suite.Assert().NotNil(warning)
	suite.Assert().Len(warning.instances, 2)
	for i, instance := range warning.Instances() {
		suite.Assert().NotNil(instance)
		suite.Assert().Equal(expectedTestNames[i], instance.name)
		suite.Assert().Equal("", instance.extra)
		suite.Assert().Contains(instance.log, "{")
	}
}

func (suite *ParserTestSuite) TestData_Parse_Verify_samber_slog_zap() {
	levels := suite.verify.ForHandler(HandlerTag("samber/slog-zap"))
	suite.Assert().NotNil(levels)
	suite.Assert().Len(levels.Levels(), 3)
	required, found := levels.lookup["Required"]
	suite.Assert().True(found)
	suite.Assert().NotNil(required)
	suite.Assert().Len(required.Warnings(), 5)
	warning, found := required.lookup["ZeroPC"]
	suite.Assert().True(found)
	suite.Assert().NotNil(warning)
	suite.Assert().Len(warning.Instances(), 1)
	instance := warning.instances[0]
	suite.Assert().NotNil(instance)
	suite.Assert().Equal("Zero PC", instance.name)
	suite.Assert().Equal("non-standard key 'caller'", instance.extra)
	suite.Assert().Contains(instance.log, "{")
}
