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
		NoTraffic)
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
		NoTraffic)
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultDebug() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[DEBUG] TestDefaultDebug"))
			require.NoError(t, err)
		},
		nil,
		NoTraffic)
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
		NoTraffic)
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultBadLevel() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[BAD] TestDefaultBadLevel"))
			require.ErrorContains(t, err, "no level BAD")
		},
		nil,
		NoTraffic)
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
		NoTraffic)
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
		NoTraffic)
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
		NoTraffic)
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
		NoTraffic)
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestTrafficSplice() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultWriter.Write([]byte(ginLine))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "INFO", record[slog.LevelKey])
			assert.Equal(t, TrafficMessage, record[slog.MessageKey])
			checkTraffic(t, record)
		},
		Traffic)
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestTrafficGroup() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultWriter.Write([]byte(ginLine))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "INFO", record[slog.LevelKey])
			assert.Equal(t, TrafficMessage, record[slog.MessageKey])
			group, ok := record[GinGroup].(map[string]any)
			assert.True(t, ok)
			checkTraffic(t, group)
		},
		GinGroup)
}

//////////////////////////////////////////////////////////////////////////

func (suite *WriterTestSuite) testLog(
	test func(t *testing.T),
	check func(t *testing.T, record map[string]interface{}),
	group string,
) {
	gin.DefaultWriter = NewWriter(slog.LevelInfo, group)
	gin.DefaultErrorWriter = NewWriter(slog.LevelError, group)
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
	// Switch slog to phsym/slog-console.
	logger := slog.New(console.NewHandler(os.Stdout, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: "<*>", // Don't want real time, too hard to match.
		NoColor:    true,
	}))
	slog.SetDefault(logger)

	gin.DefaultWriter = NewWriter(slog.LevelInfo, NoTraffic)
	gin.DefaultErrorWriter = NewWriter(slog.LevelError, NoTraffic)
	defer func() {
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}()
	_ = gin.New()
	// Output:
	// <*> WRN Running in "debug" mode. Switch to "release" mode in production.
	//  - using env:	export GIN_MODE=release
	//  - using code:	gin.SetMode(gin.ReleaseMode)
}

//////////////////////////////////////////////////////////////////////////

func checkTraffic(t *testing.T, group map[string]any) {
	assert.Equal(t, "::1", group[string(Client)])
	assert.Equal(t, float64(200), group[string(Code)])
	assert.Equal(t, "GET", group[string(Method)])
}
