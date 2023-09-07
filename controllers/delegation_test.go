package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pavva91/tezos-delegation-service/models"
	"github.com/pavva91/tezos-delegation-service/services"
	"github.com/pavva91/tezos-delegation-service/stubs"
	"github.com/stretchr/testify/assert"
)

func Test_ListDelegations_Error_Error(t *testing.T) {
	expectedHttpStatus := http.StatusInternalServerError
	expectedHttpBody := "{\"error\":{},\"message\":\"Error to list delegations\"}"

	delegationServiceStub := stubs.DelegationServiceStub{}
	delegationServiceStub.ListDelegationsFn = func(time.Time) ([]models.Delegation, error) {
		return nil, errors.New("error executing ping")
	}
	services.DelegationService = delegationServiceStub

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	DelegationController.ListDelegations(context)

	actualHttpStatus := context.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, actualHttpStatus, expectedHttpStatus)
	assert.Equal(t, actualHttpBody, expectedHttpBody)
}

func Test_ListDelegations_Empty_Empty(t *testing.T) {
	emptyDelegationList := []models.Delegation{}

	expectedHttpStatus := http.StatusOK
	expectedHttpBody := "[]"

	delegationServiceStub := stubs.DelegationServiceStub{}
	delegationServiceStub.ListDelegationsFn = func(time.Time) ([]models.Delegation, error) {
		return emptyDelegationList, nil
	}
	services.DelegationService = delegationServiceStub

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	DelegationController.ListDelegations(context)

	actualHttpStatus := context.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, actualHttpStatus, expectedHttpStatus)
	assert.Equal(t, actualHttpBody, expectedHttpBody)
}
