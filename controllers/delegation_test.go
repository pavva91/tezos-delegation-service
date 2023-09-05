package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pavva91/gin-gorm-rest/models"
	"github.com/pavva91/gin-gorm-rest/services"
	"github.com/pavva91/gin-gorm-rest/stubs"
	"github.com/stretchr/testify/assert"
)

func Test_ListEvents_Error_Error(t *testing.T) {
	expectedHttpStatus := http.StatusInternalServerError
	expectedHttpBody := "{\"error\":{},\"message\":\"Error to list events\"}"

	eventServiceStub := stubs.EventServiceStub{}
	eventServiceStub.ListEventsFn = func() ([]models.Event, error) {
		return nil, errors.New("error executing ping")
	}
	services.EventService = eventServiceStub

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	EventController.ListEvents(context)

	actualHttpStatus := context.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, actualHttpStatus, expectedHttpStatus)
	assert.Equal(t, actualHttpBody, expectedHttpBody)
}

func Test_ListEvents_Empty_Empty(t *testing.T) {
	emptyEventList := []models.Event{}

	expectedHttpStatus := http.StatusOK
	expectedHttpBody := "[]"

	eventServiceStub := stubs.EventServiceStub{}
	eventServiceStub.ListEventsFn = func() ([]models.Event, error) {
		return emptyEventList, nil
	}
	services.EventService = eventServiceStub

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	EventController.ListEvents(context)

	actualHttpStatus := context.Writer.Status()
	actualHttpBody := response.Body.String()

	assert.Equal(t, actualHttpStatus, expectedHttpStatus)
	assert.Equal(t, actualHttpBody, expectedHttpBody)
}
