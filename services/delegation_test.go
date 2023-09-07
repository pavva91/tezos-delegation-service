package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/pavva91/tezos-delegation-service/dto"
	"github.com/pavva91/tezos-delegation-service/models"
	"github.com/pavva91/tezos-delegation-service/repositories"
	"github.com/pavva91/tezos-delegation-service/stubs"
	"github.com/stretchr/testify/assert"
)

func Test_ListDelegations_Error_Error(t *testing.T) {
	errorMessage := "Unexpected Internal Error"
	unexpectedError := errors.New(errorMessage)
	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListFn = func() ([]models.Delegation, error) {
		return nil, unexpectedError
	}
	repositories.DelegationRepository = delegationRepositoryStub

	delegations, err := DelegationService.ListDelegations()

	assert.NotNil(t, err)
	assert.Nil(t, delegations)
	assert.Equal(t, errorMessage, err.Error())
}

func Test_ListDelegations_Empty_Empty(t *testing.T) {
	emptyDelegationList := []models.Delegation{}

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListFn = func() ([]models.Delegation, error) {
		return emptyDelegationList, nil
	}
	repositories.DelegationRepository = delegationRepositoryStub

	delegations, err := DelegationService.ListDelegations()

	assert.Nil(t, err)
	assert.NotNil(t, delegations)
	assert.Equal(t, 0, len(delegations))
}

func Test_PollDelegations_WrongApiEndpointScheme_Error(t *testing.T) {
	wrongApiEndpoint := "wrong"
	pollPeriodInSeconds := 1
	expectedErrorContent1 := "Get \"" + wrongApiEndpoint + "/operations/delegations?timestamp.ge="
	expectedErrorContent2 := "unsupported protocol scheme"

	err := DelegationService.PollDelegations(pollPeriodInSeconds, wrongApiEndpoint)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
	assert.Contains(t, err.Error(), expectedErrorContent2)
}

func Test_PollDelegations_WrongApiEndpointDomain_Error(t *testing.T) {
	wrongApiEndpoint := "http://wrong"
	pollPeriodInSeconds := 1
	expectedErrorContent1 := "Get \"" + wrongApiEndpoint + "/operations/delegations?timestamp.ge="
	expectedErrorContent2 := "dial tcp: lookup wrong: no such host"

	err := DelegationService.PollDelegations(pollPeriodInSeconds, wrongApiEndpoint)
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
	assert.Contains(t, err.Error(), expectedErrorContent2)
}

func Test_PollDelegations_Not200FromApiEndpoint_Error(t *testing.T) {
	pollPeriodInSeconds := 1
	errorHttpStatus := http.StatusBadRequest
	expectedError := "Get Response different than 200: " + strconv.Itoa(errorHttpStatus)

	// Mock outbound http request https://medium.com/zus-health/mocking-outbound-http-requests-in-go-youre-probably-doing-it-wrong-60373a38d2aa
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/operations/delegations" {
			t.Errorf("Expected to request '/operations/delegations', got: %s", r.URL.Path)
		}
		w.WriteHeader(errorHttpStatus)
		w.Write([]byte{})
	}))
	defer server.Close()

	err := DelegationService.PollDelegations(pollPeriodInSeconds, server.URL)
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func Test_PollDelegations_ReturnedUnexpectedJSON_Error(t *testing.T) {
	pollPeriodInSeconds := 1
	unexpectedJSON := []byte(`{"value":"fixed"}`)
	expectedError := "json: cannot unmarshal object into Go value of type []dto.DelegationResponseFromApi"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/operations/delegations" {
			t.Errorf("Expected to request '/operations/delegations', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(unexpectedJSON)
	}))
	defer server.Close()

	err := DelegationService.PollDelegations(pollPeriodInSeconds, server.URL)
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func Test_PollDelegations_WorksThenApiGoDownAfter2Seconds_Error(t *testing.T) {
	pollPeriodInSeconds := 1
	emptyListDelegations := []dto.DelegationResponseFromApi{}
	jsonResponse, err := json.Marshal(emptyListDelegations)
	if err != nil {
		return
	}
	// NOTE: error: "Get \"http://127.0.0.1:45969/operations/delegations?timestamp.ge=2023-09-07T08:49:59Z&timestamp.lt=2023-09-07T08:50:01Z\": dial tcp 127.0.0.1:45969: connect: connection refused"
	expectedErrorContent1 := "connect: connection refused"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/operations/delegations" {
			t.Errorf("Expected to request '/operations/delegations', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))

	go func() {
		time.Sleep(2 * time.Second)
		server.Close()
	}()

	defer server.Close()

	err = DelegationService.PollDelegations(pollPeriodInSeconds, server.URL)
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
}

func Test_SaveBulkDelegations_EmptySlice_NoError(t *testing.T) {
	var emptyList []dto.DelegationResponseFromApi

	savedDelegations, err := SaveBulkDelegations(emptyList)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(savedDelegations))
}

func Test_SaveBulkDelegations_ErrorFromRepositoryCreate_Error(t *testing.T) {
	var emptyList []dto.DelegationResponseFromApi
	delegation1 := dto.DelegationResponseFromApi{}
	listOneElement := append(emptyList, delegation1)

	errorMessage := "Unexpected Internal Error"
	unexpectedError := errors.New(errorMessage)

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) (*models.Delegation, error) {
		return nil, unexpectedError
	}
	repositories.DelegationRepository = delegationRepositoryStub

	savedDelegations, err := SaveBulkDelegations(listOneElement)
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, errorMessage, err.Error())
	assert.Equal(t, 0, len(savedDelegations))
}

func Test_SaveBulkDelegations_OKList1Element_ReturnSavedDelegation(t *testing.T) {
	var emptyList []dto.DelegationResponseFromApi
	delegation := dto.DelegationResponseFromApi{}
	delegationModel1 := models.Delegation{}
	listOneElement := append(emptyList, delegation)

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) (*models.Delegation, error) {
		return &delegationModel1, nil
	}
	repositories.DelegationRepository = delegationRepositoryStub

	savedDelegations, err := SaveBulkDelegations(listOneElement)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(savedDelegations))
}

func Test_SaveBulkDelegations_OKList2Element_ReturnSavedDelegation(t *testing.T) {
	var delegations []dto.DelegationResponseFromApi
	delegation := dto.DelegationResponseFromApi{}
	delegationModel1 := models.Delegation{}
	delegations = append(delegations, delegation)
	delegations = append(delegations, delegation)

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) (*models.Delegation, error) {
		return &delegationModel1, nil
	}
	repositories.DelegationRepository = delegationRepositoryStub

	savedDelegations, err := SaveBulkDelegations(delegations)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(savedDelegations))
}
