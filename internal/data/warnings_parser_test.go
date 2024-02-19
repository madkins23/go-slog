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
	suite.bench = NewWarningData()
	suite.Require().NoError(suite.bench.ParseWarningData(strings.NewReader(benchTxt), "", nil))
	suite.verify = NewWarningData()
	suite.Require().NoError(suite.verify.ParseWarningData(strings.NewReader(verifyTxt), "", nil))
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (suite *ParserTestSuite) TestData_Parse_Bench_darvaza_zerolog() {
	var expectedTestNames = []string{
		"With_Attrs_Attributes",
		"With_Attrs_Key_Values",
		"With_Attrs_Simple",
		"With_Group_Attributes",
		"With_Group_Key_Values",
	}
	levels := suite.bench.ForHandler(HandlerTag("darvaza/zerolog"))
	suite.Assert().NotNil(levels)
	suite.Assert().Len(levels.Levels(), 1)
	admin, found := levels.lookup["Administrative"]
	suite.Assert().True(found)
	suite.Assert().NotNil(admin)
	suite.Assert().Len(admin.Warnings(), 1)
	warning, found := admin.lookup["NoHandlerCreation"]
	suite.Assert().True(found)
	suite.Assert().NotNil(warning)
	suite.Require().Len(warning.Instances(), 5)
	for i, instance := range warning.Instances() {
		suite.Assert().NotNil(instance)
		suite.Assert().Equal(expectedTestNames[i], instance.name)
		suite.Assert().Equal("", instance.extra)
		suite.Assert().Equal("", instance.log)
	}
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
	suite.Assert().Equal("slog_phsym_zerolog", instance.name)
	suite.Assert().Equal("Simple_Source: no 'source' key", instance.extra)
	suite.Assert().Contains(instance.log, "{")
}

func (suite *ParserTestSuite) TestData_Parse_Verify_phsym_zerolog() {
	var expectedTestNames = []string{
		"AttributeWithEmpty",
		"AttributesEmpty",
	}
	levels := suite.verify.ForHandler(HandlerTag("phsym/zeroslog"))
	suite.Assert().NotNil(levels)
	suite.Assert().Len(levels.Levels(), 3)
	required, found := levels.lookup["Required"]
	suite.Assert().True(found)
	suite.Assert().NotNil(required)
	suite.Assert().Len(required.Warnings(), 3)
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
	suite.Assert().Equal("ZeroPC", instance.name)
	suite.Assert().Equal("non-standard key 'caller'", instance.extra)
	suite.Assert().Contains(instance.log, "{")
}
