package gin

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phsym/console-slog"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
}

func TestWriterSuite(t *testing.T) {
	gin.DefaultWriter = NewWriter(slog.LevelInfo, "")
	gin.DefaultErrorWriter = NewWriter(slog.LevelError, "")
	defer func() {
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}()

	// Breakout test suite startup so that GinStartupTest() can be run first.
	sweet := new(WriterTestSuite)
	sweet.SetT(t)
	sweet.GinStartupTest()

	// Run the rest of the tests
	suite.Run(t, sweet)
}

func TestWriterSuiteGroup(t *testing.T) {
	gin.DefaultWriter = NewWriter(slog.LevelInfo, "gin")
	gin.DefaultErrorWriter = NewWriter(slog.LevelError, "gin")
	defer func() {
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}()

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
		})
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
		})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultDebug() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[DEBUG] TestDefaultDebug"))
			require.NoError(t, err)
		},
		// Sending DEBUG to logger configured for INFO and above, no reponse.
		nil)
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultGin() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[GIN] TestDefaultGin"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "ERROR", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestDefaultGin")
		})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultBadLevel() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[BAD] TestDefaultBadLevel"))
			require.ErrorContains(t, err, "no level BAD")
		}, nil)
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultWarning() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[WARNING] TestDefaultWarning"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "WARN", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestDefaultWarning")
		})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestError() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("TestError"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "ERROR", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestError")
		})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestErrorWarning() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[WARNING] TestErrorWarning"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "WARN", record[slog.LevelKey])
			assert.Contains(t, record[slog.MessageKey], "TestErrorWarning")
		})
}

//////////////////////////////////////////////////////////////////////////

func (suite *WriterTestSuite) testLog(test func(t *testing.T), check func(t *testing.T, record map[string]interface{})) {
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
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	logger := slog.New(console.NewHandler(os.Stdout, &console.HandlerOptions{
		Level:      slog.LevelInfo,
		TimeFormat: "<*>", // Don't want real time, too hard to match.
		NoColor:    true,
	}))
	slog.SetDefault(logger)

	gin.DefaultWriter = NewWriter(slog.LevelInfo, "")
	gin.DefaultErrorWriter = NewWriter(slog.LevelError, "")
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
