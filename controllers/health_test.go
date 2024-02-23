package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pavva91/tezos-delegation-service/controllers"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type HealthTestSuite struct {
	suite.Suite
	GinContextPointer  *gin.Context
	GinEnginePointer   *gin.Engine
	HTTPResponseWriter http.ResponseWriter
	Response           *httptest.ResponseRecorder
}

// Setup Stub Values
func (suite *HealthTestSuite) SetupTest() {
	// not strictly required to unit test (will run also without this line)
	gin.SetMode(gin.TestMode)
	suite.Response = httptest.NewRecorder()

	suite.GinContextPointer, _ = gin.CreateTestContext(suite.Response)
}

func (suite *HealthTestSuite) Test_Status_Return200() {
	expectedHTTPStatus := http.StatusOK
	expectedHTTPBody := "Working!"

	controllers.Health.Status(suite.GinContextPointer)

	actualHTTPStatus := suite.GinContextPointer.Writer.Status()
	actualHTTPBody := suite.Response.Body.String()

	suite.Equal(expectedHTTPStatus, actualHTTPStatus)
	suite.Equal(expectedHTTPBody, actualHTTPBody)
	// assert.Equal(suite.T(), expectedHTTPStatus,  actualHTTPStatus)
	// assert.Equal(suite.T(), expectedHTTPBody,  actualHTTPBody)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSuiteHealthController(t *testing.T) {
	suite.Run(t, new(HealthTestSuite))
}
