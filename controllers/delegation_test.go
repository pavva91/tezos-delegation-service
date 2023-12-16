package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pavva91/tezos-delegation-service/models"
	"github.com/pavva91/tezos-delegation-service/services"
	"github.com/pavva91/tezos-delegation-service/stubs"
	"github.com/stretchr/testify/assert"
)

func Test_ListDelegations_QueryParameterYearString_400(t *testing.T) {
	wrongYearQueryParameterValue := "lorem ipsum"
	response := httptest.NewRecorder()
	mockContext, _ := gin.CreateTestContext(response)
	mockContext.Request = &http.Request{
		Header: make(http.Header),
	}
	mockContext.Request.URL, _ = url.Parse("?year=" + wrongYearQueryParameterValue)

	expectedHttpStatus := http.StatusBadRequest
	expectedHttpBody := "{\"error\":\"parsing time \\\"" + wrongYearQueryParameterValue + "\\\" as \\\"2006\\\": cannot parse \\\"" + wrongYearQueryParameterValue + "\\\" as \\\"2006\\\"\"}"

	Delegation.List(mockContext)

	actualHttpStatus := mockContext.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, expectedHttpStatus, actualHttpStatus)
	assert.Equal(t, expectedHttpBody, actualHttpBody)
}

func Test_ListDelegations_QueryParameterWrongDate_400(t *testing.T) {
	wrongYearQueryParameterValue := time.Now().Format(time.RFC3339)
	response := httptest.NewRecorder()
	mockContext, _ := gin.CreateTestContext(response)
	mockContext.Request = &http.Request{
		Header: make(http.Header),
	}
	mockContext.Request.URL, _ = url.Parse("?year=" + wrongYearQueryParameterValue)

	expectedHttpStatus := http.StatusBadRequest
	wrongYearInError := strings.Replace(wrongYearQueryParameterValue, "+", " ", 1)
	expectedHttpBody := "{\"error\":\"parsing time \\\"" + wrongYearInError + "\\\": extra text: \\\"" + wrongYearInError[4:] + "\\\"\"}"

	Delegation.List(mockContext)

	actualHttpStatus := mockContext.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, expectedHttpStatus, actualHttpStatus)
	assert.Equal(t, expectedHttpBody, actualHttpBody)
}

func Test_ListDelegations_QueryParameterYearTrailingChars_400(t *testing.T) {
	wrongYearQueryParameterValue := "2000asdf"
	response := httptest.NewRecorder()
	mockContext, _ := gin.CreateTestContext(response)
	mockContext.Request = &http.Request{
		Header: make(http.Header),
	}
	mockContext.Request.URL, _ = url.Parse("?year=" + wrongYearQueryParameterValue)

	expectedHttpStatus := http.StatusBadRequest
	expectedHttpBody := "{\"error\":\"parsing time \\\"" + wrongYearQueryParameterValue + "\\\": extra text: \\\"" + wrongYearQueryParameterValue[4:] + "\\\"\"}"

	Delegation.List(mockContext)

	actualHttpStatus := mockContext.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, expectedHttpStatus, actualHttpStatus)
	assert.Equal(t, expectedHttpBody, actualHttpBody)
}

func Test_ListDelegations_ServiceInternalError_500(t *testing.T) {
	correctYearQueryParameter := "2000"
	response := httptest.NewRecorder()
	mockContext, _ := gin.CreateTestContext(response)
	mockContext.Request = &http.Request{
		Header: make(http.Header),
	}
	mockContext.Request.URL, _ = url.Parse("?year=" + correctYearQueryParameter)

	expectedHttpStatus := http.StatusInternalServerError
	expectedHttpBody := "{\"error\":\"Error to list delegations\"}"
	errorMessage := "Unexpected Internal Error"

	delegationServiceStub := stubs.DelegationServiceStub{}
	delegationServiceStub.ListDelegationsFn = func(time.Time) ([]models.Delegation, error) {
		return nil, errors.New(errorMessage)
	}
	services.Delegation = delegationServiceStub

	Delegation.List(mockContext)

	actualHttpStatus := mockContext.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, expectedHttpStatus, actualHttpStatus)
	assert.Equal(t, expectedHttpBody, actualHttpBody)
}

func Test_ListDelegations_OK_200(t *testing.T) {
	correctYearQueryParameter := "2000"
	response := httptest.NewRecorder()
	mockContext, _ := gin.CreateTestContext(response)
	mockContext.Request = &http.Request{
		Header: make(http.Header),
	}
	mockContext.Request.URL, _ = url.Parse("?year=" + correctYearQueryParameter)

	delegator := "tz1huoYxZWLXVgdfEJbfqpp1KXdPiDoyGtJQ"
	amount := "1"
	block := "BMQNYHimngWWRf2d6WfM5qscYPzFSx2BgyfnTkrf6Vp8PZc7hrJ"
	timestamp := time.Now().UTC()
	var delegations []models.Delegation
	delegation1 := models.Delegation{
		Delegator: delegator,
		Amount:    amount,
		Block:     block,
		Timestamp: timestamp,
	}
	delegations = append(delegations, delegation1)

	expectedHttpStatus := http.StatusOK

	delegationServiceStub := stubs.DelegationServiceStub{}
	delegationServiceStub.ListDelegationsFn = func(time.Time) ([]models.Delegation, error) {
		return delegations, nil
	}
	services.Delegation = delegationServiceStub

	Delegation.List(mockContext)

	actualHttpStatus := mockContext.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, expectedHttpStatus, actualHttpStatus)
	assert.Contains(t, actualHttpBody, delegator)
	assert.Contains(t, actualHttpBody, amount)
	assert.Contains(t, actualHttpBody, block)
	assert.Contains(t, actualHttpBody, timestamp.Format(time.RFC3339Nano))
}
