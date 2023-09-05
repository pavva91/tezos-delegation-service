package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type HealthTestSuite struct {
	suite.Suite
	GinContextPointer  *gin.Context
	GinEnginePointer   *gin.Engine
	HttpResponseWriter http.ResponseWriter
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
	expectedHttpStatus := http.StatusOK
	expectedHttpBody := "Working!"

	HealthController.Status(suite.GinContextPointer)

	actualHttpStatus := suite.GinContextPointer.Writer.Status()
	actualHttpBody := suite.Response.Body.String()

	assert.Equal(suite.T(), actualHttpStatus, expectedHttpStatus)
	assert.Equal(suite.T(), actualHttpBody, expectedHttpBody)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSuiteHealthController(t *testing.T) {
	suite.Run(t, new(HealthTestSuite))
}
