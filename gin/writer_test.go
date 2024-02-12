package gin

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const ginLine = `[GIN] 2024/01/26 - 13:21:32 | 200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"`
const ginTraffic = `200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"`

type WriterTestSuite struct {
	suite.Suite
}

func TestWriterSuite(t *testing.T) {
	// Breakout test suite startup so that GinStartupTest() can be run first.
	sweet := new(WriterTestSuite)
	sweet.SetT(t)
	sweet.GinStartupTest()

	// Run the rest of the tests
	suite.Run(t, sweet)
}

// GinStartupTest traps and tests the initial Gin startup warning for debug mode.
func (suite *WriterTestSuite) GinStartupTest() {
	suite.testLog(
		func(t *testing.T) {
			gn := gin.New()
			require.NotNil(t, gn)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "WARN", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "Running in \"debug\" mode.")
		},
		Options{})
}

//////////////////////////////////////////////////////////////////////////

func (suite *WriterTestSuite) TestDefault() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultWriter.Write([]byte("TestDefault"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "INFO", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestDefault")
		},
		Options{})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultDebug() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[DEBUG] TestDefaultDebug"))
			require.NoError(t, err)
		},
		nil,
		Options{})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultGin() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[GIN] TestDefaultGin"))
			require.NoError(t, err)
		},
		func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "ERROR", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestDefaultGin")
		},
		Options{})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultBadLevel() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[BAD] TestDefaultBadLevel"))
			require.ErrorContains(t, err, "no level BAD")
		},
		nil,
		Options{})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultWarning() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[WARNING] TestDefaultWarning"))
			require.NoError(t, err)
		},
		func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "WARN", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestDefaultWarning")
		},
		Options{})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestError() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("TestError"))
			require.NoError(t, err)
		},
		func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "ERROR", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestError")
		},
		Options{})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestErrorWarning() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[WARNING] TestErrorWarning"))
			require.NoError(t, err)
		},
		func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "WARN", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestErrorWarning")
		},
		Options{})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestTrafficIgnore() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultWriter.Write([]byte(ginLine))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "INFO", record[slog.LevelKey])
			assert.Equal(t, ginTraffic, record[slog.MessageKey])
		},
		Options{})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestTrafficEmbed() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultWriter.Write([]byte(ginLine))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "INFO", record[slog.LevelKey])
			assert.Equal(t, DefaultTrafficMessage, record[slog.MessageKey])
			checkTraffic(t, record)
		},
		Options{Traffic: Traffic{Parse: true, Embed: true}})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestTrafficGroup() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultWriter.Write([]byte(ginLine))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "INFO", record[slog.LevelKey])
			assert.Equal(t, DefaultTrafficMessage, record[slog.MessageKey])
			group, ok := record[DefaultTrafficGroup].(map[string]any)
			assert.True(t, ok)
			checkTraffic(t, group)
		},
		Options{Traffic: Traffic{Parse: true}},
	)
}

func (suite *WriterTestSuite) TestTrafficGroupName() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultWriter.Write([]byte(ginLine))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "INFO", record[slog.LevelKey])
			assert.Equal(t, DefaultTrafficMessage, record[slog.MessageKey])
			group, ok := record["group-name"].(map[string]any)
			assert.True(t, ok)
			checkTraffic(t, group)
		},
		Options{Traffic: Traffic{Parse: true, Group: "group-name"}},
	)
}

//////////////////////////////////////////////////////////////////////////

func (suite *WriterTestSuite) testLog(
	test func(t *testing.T),
	check func(t *testing.T, record map[string]interface{}),
	options Options,
) {
	gin.DefaultWriter = NewWriter(&options)
	options.Level = slog.LevelError
	gin.DefaultErrorWriter = NewWriter(&options)
	defer func() {
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}()

	// Trap output from running log function.
	sLog := slog.Default()
	defer slog.SetDefault(sLog)
	buffer := &bytes.Buffer{}
	slog.SetDefault(slog.New(slog.NewJSONHandler(buffer, nil)))

	// Execute test.
	test(suite.T())
	if check != nil {
		// Check log output which is in JSON.
		var record map[string]interface{}
		suite.Require().NoError(json.Unmarshal(buffer.Bytes(), &record))
		check(suite.T(), record)
	}
}

//////////////////////////////////////////////////////////////////////////

func ExampleWriter() {
	// Switch slog to phsym/slog-console because
	// * it's simpler to match multi-line output at Gin startup and
	// * the time format can be changed.
	logger := slog.New(console.NewHandler(os.Stdout, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: "<*>", // Don't want real time, too hard to match.
		NoColor:    true,
	}))
	slog.SetDefault(logger)

	gin.DefaultWriter = NewWriter(&Options{})
	gin.DefaultErrorWriter = NewWriter(&Options{Level: slog.LevelError})
	defer func() {
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}()
	_ = gin.New()
	_, _ = gin.DefaultWriter.Write([]byte(ginLine))
	// Output:
	// <*> WRN Running in "debug" mode. Switch to "release" mode in production.
	//  - using env:	export GIN_MODE=release
	//  - using code:	gin.SetMode(gin.ReleaseMode)
	// <*> INF 200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"
}

func ExampleWriter_embedTraffic() {
	// Switch slog to phsym/slog-console because
	// * it's simpler to match multi-line output at Gin startup and
	// * the time format can be changed.
	logger := slog.New(console.NewHandler(os.Stdout, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: "<*>", // Don't want real time, too hard to match.
		NoColor:    true,
	}))
	slog.SetDefault(logger)

	options := &Options{
		Traffic: Traffic{
			Parse: true,
			Embed: true,
		},
	}
	gin.DefaultWriter = NewWriter(options)
	gin.DefaultErrorWriter = NewWriter(&Options{Level: slog.LevelError})
	defer func() {
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}()
	_ = gin.New()
	_, _ = gin.DefaultWriter.Write([]byte(ginLine))
	// Output:
	// <*> WRN Running in "debug" mode. Switch to "release" mode in production.
	//  - using env:	export GIN_MODE=release
	//  - using code:	gin.SetMode(gin.ReleaseMode)
	// <*> INF Gin Traffic code=200 elapsed=5.529751605s client=::1 method=GET url=/chart.svg?tag=With_Attrs_Attributes&item=MemBytes
}

func ExampleWriter_groupTraffic() {
	// Switch slog to phsym/slog-console because
	// * it's simpler to match multi-line output at Gin startup and
	// * the time format can be changed.
	logger := slog.New(console.NewHandler(os.Stdout, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: "<*>", // Don't want real time, too hard to match.
		NoColor:    true,
	}))
	slog.SetDefault(logger)

	options := &Options{
		Traffic: Traffic{
			Parse: true,
		},
	}
	gin.DefaultWriter = NewWriter(options)
	gin.DefaultErrorWriter = NewWriter(&Options{Level: slog.LevelError})
	defer func() {
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}()
	_ = gin.New()
	_, _ = gin.DefaultWriter.Write([]byte(ginLine))
	// Output:
	// <*> WRN Running in "debug" mode. Switch to "release" mode in production.
	//  - using env:	export GIN_MODE=release
	//  - using code:	gin.SetMode(gin.ReleaseMode)
	// <*> INF Gin Traffic gin.code=200 gin.elapsed=5.529751605s gin.client=::1 gin.method=GET gin.url=/chart.svg?tag=With_Attrs_Attributes&item=MemBytes
}

//////////////////////////////////////////////////////////////////////////

func checkTraffic(t *testing.T, group map[string]any) {
	assert.Equal(t, "::1", group[string(Client)])
	assert.Equal(t, float64(200), group[string(Code)])
	assert.Equal(t, "GET", group[string(Method)])
}
