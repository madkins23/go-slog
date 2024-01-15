package ginzero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
}

func TestWriterSuite(t *testing.T) {
	gin.DefaultWriter = NewWriter(zerolog.InfoLevel)
	gin.DefaultErrorWriter = NewWriter(zerolog.ErrorLevel)
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
			assert.Equal(t, "warn", record["level"])
			assert.Equal(t, "gin", record["sys"])
			assert.Contains(t, record["message"], "Running in \"debug\" mode.")
		})
}

//////////////////////////////////////////////////////////////////////////

func (suite *WriterTestSuite) TestDefault() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultWriter.Write([]byte("TestDefault"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "info", record["level"])
			assert.Contains(t, record["message"], "TestDefault")
		})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultDebug() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[DEBUG] TestDefaultDebug"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "debug", record["level"])
			assert.Contains(t, record["message"], "TestDefaultDebug")
		})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestDefaultGin() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[GIN] TestDefaultGin"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "error", record["level"])
			assert.Contains(t, record["message"], "TestDefaultGin")
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
			assert.Equal(t, "warn", record["level"])
			assert.Contains(t, record["message"], "TestDefaultWarning")
		})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestError() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("TestError"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "error", record["level"])
			assert.Contains(t, record["message"], "TestError")
		})
}

// ----------------------------------------------------------------------------

func (suite *WriterTestSuite) TestErrorWarning() {
	suite.testLog(
		func(t *testing.T) {
			_, err := gin.DefaultErrorWriter.Write([]byte("[WARNING] TestErrorWarning"))
			require.NoError(t, err)
		}, func(t *testing.T, record map[string]interface{}) {
			assert.Equal(t, "warn", record["level"])
			assert.Contains(t, record["message"], "TestErrorWarning")
		})
}

//////////////////////////////////////////////////////////////////////////

func (suite *WriterTestSuite) testLog(test func(t *testing.T), check func(t *testing.T, record map[string]interface{})) {
	// Trap output from running log function.
	zLog := log.Logger
	defer func() { log.Logger = zLog }()
	buffer := &bytes.Buffer{}
	log.Logger = zerolog.New(buffer)
	// Execute test.
	test(suite.T())
	if check != nil {
		// Check log output which is in JSON.
		var record map[string]interface{}
		fmt.Println("JSON ", buffer.String())
		suite.Require().NoError(json.Unmarshal(buffer.Bytes(), &record))
		check(suite.T(), record)
	}
}

//////////////////////////////////////////////////////////////////////////

func ExampleWriter() {
	// Switch zerolog to console mode.
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Local()
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:     os.Stdout,
		NoColor: true,
		// Don't show duration or time as they mess up the Output comparison.
		PartsExclude: []string{"time"},
	})

	gin.DefaultWriter = NewWriter(zerolog.InfoLevel)
	gin.DefaultErrorWriter = NewWriter(zerolog.ErrorLevel)
	defer func() {
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
	}()
	_ = gin.New()

	// Output:
	// WRN Running in "debug" mode. Switch to "release" mode in production.
	//  - using env:	export GIN_MODE=release
	//  - using code:	gin.SetMode(gin.ReleaseMode) sys=gin
}
