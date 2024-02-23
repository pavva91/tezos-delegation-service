package services_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pavva91/tezos-delegation-service/internal/dto"
	"github.com/pavva91/tezos-delegation-service/internal/models"
	"github.com/pavva91/tezos-delegation-service/internal/repositories"
	"github.com/pavva91/tezos-delegation-service/internal/services"
	"github.com/pavva91/tezos-delegation-service/internal/stubs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListDelegations_YearNonZero_Error(t *testing.T) {
	nonZeroValueDate := time.Now()
	errorMessage := "Unexpected Internal Error"
	unexpectedError := errors.New(errorMessage)
	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListByYearFn = func(time.Time) ([]models.Delegation, error) {
		return nil, unexpectedError
	}
	repositories.Delegation = delegationRepositoryStub

	delegations, err := services.Delegation.List(nonZeroValueDate)

	require.Error(t, err)
	assert.Nil(t, delegations)
	assert.Equal(t, errorMessage, err.Error())
}

func Test_ListDelegations_YearIsZero_Error(t *testing.T) {
	var zeroValueDate time.Time
	errorMessage := "Unexpected Internal Error"
	unexpectedError := errors.New(errorMessage)
	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListFn = func() ([]models.Delegation, error) {
		return nil, unexpectedError
	}
	repositories.Delegation = delegationRepositoryStub

	delegations, err := services.Delegation.List(zeroValueDate)

	require.Error(t, err)
	assert.Nil(t, delegations)
	assert.Equal(t, errorMessage, err.Error())
}

func Test_ListDelegations_YearNonZeroEmptyList_Empty(t *testing.T) {
	nonZeroValueDate := time.Now()
	emptyDelegationList := []models.Delegation{}

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListByYearFn = func(time.Time) ([]models.Delegation, error) {
		return emptyDelegationList, nil
	}
	repositories.Delegation = delegationRepositoryStub

	delegations, err := services.Delegation.List(nonZeroValueDate)

	require.NoError(t, err)
	assert.NotNil(t, delegations)
	assert.Empty(t, delegations)
}

func Test_ListDelegations_YearIsZeroEmptyList_Empty(t *testing.T) {
	var zeroValueDate time.Time
	emptyDelegationList := []models.Delegation{}

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListFn = func() ([]models.Delegation, error) {
		return emptyDelegationList, nil
	}
	repositories.Delegation = delegationRepositoryStub

	delegations, err := services.Delegation.List(zeroValueDate)

	require.NoError(t, err)
	assert.NotNil(t, delegations)
	assert.Empty(t, delegations)
}

func Test_PollDelegations_WrongApiEndpointScheme_Error(t *testing.T) {
	wrongAPIEndpoint := "wrong"
	pollPeriodInSeconds := uint(1)
	expectedErrorContent1 := "Get \"" + wrongAPIEndpoint + "/operations/delegations?timestamp.ge="
	expectedErrorContent2 := "unsupported protocol scheme"
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})

	go services.Delegation.Poll(pollPeriodInSeconds, wrongAPIEndpoint, stopOnError, errorCh, interruptCh)

	time.Sleep(5 * time.Second)
	interruptCh <- struct{}{}
	err := <-errorCh
	fmt.Println(err.Error())

	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
	assert.Contains(t, err.Error(), expectedErrorContent2)
}

func Test_PollDelegations_WrongApiEndpointDomain_Error(t *testing.T) {
	wrongAPIEndpoint := "http://wrong-api-endpoint"
	pollPeriodInSeconds := uint(1)
	expectedErrorContent1 := "Get \"" + wrongAPIEndpoint + "/operations/delegations?timestamp.ge="
	expectedErrorContent2 := "dial tcp: lookup " + wrongAPIEndpoint[7:]
	// TODO: stopOnError becomes a boolean
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})
	// TODO: create channel token to send signal

	go services.Delegation.Poll(pollPeriodInSeconds, wrongAPIEndpoint, stopOnError, errorCh, interruptCh)

	time.Sleep(5 * time.Second)
	// TODO: create channel token to send signal
	interruptCh <- struct{}{}
	err := <-errorCh
	fmt.Println(err.Error())

	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
	assert.Contains(t, err.Error(), expectedErrorContent2)
}

func Test_PollDelegations_Not200FromApiEndpoint_Error(t *testing.T) {
	pollPeriodInSeconds := uint(1)
	errorHTTPStatus := http.StatusBadRequest
	expectedError := "get response different than 200: %!w(<nil>) "
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})

	// Mock outbound http request https://medium.com/zus-health/mocking-outbound-http-requests-in-go-youre-probably-doing-it-wrong-60373a38d2aa
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/operations/delegations" {
			t.Errorf("Expected to request '/operations/delegations', got: %s", r.URL.Path)
		}
		w.WriteHeader(errorHTTPStatus)
		w.Write([]byte{})
	}))
	defer server.Close()

	go services.Delegation.Poll(pollPeriodInSeconds, server.URL, stopOnError, errorCh, interruptCh)

	time.Sleep(3 * time.Second)
	interruptCh <- struct{}{}
	err := <-errorCh
	fmt.Println(err.Error())

	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func Test_PollDelegations_ReturnedUnexpectedJSON_Error(t *testing.T) {
	pollPeriodInSeconds := uint(1)
	unexpectedJSON := []byte(`{"value":"fixed"}`)
	expectedError := "json: cannot unmarshal object into Go value of type []dto.DelegationResponseFromAPI"
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/operations/delegations" {
			t.Errorf("Expected to request '/operations/delegations', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(unexpectedJSON)
	}))
	defer server.Close()

	err := services.Delegation.Poll(pollPeriodInSeconds, server.URL, stopOnError, errorCh, interruptCh)
	fmt.Println(err.Error())

	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func Test_PollDelegations_WorksThenApiGoDownAfter2Seconds_Error(t *testing.T) {
	pollPeriodInSeconds := uint(1)
	emptyListDelegations := []dto.DelegationResponseFromAPI{}
	jsonResponse, err := json.Marshal(emptyListDelegations)
	if err != nil {
		return
	}
	// NOTE: error: "Get \"http://127.0.0.1:45969/operations/delegations?timestamp.ge=2023-09-07T08:49:59Z&timestamp.lt=2023-09-07T08:50:01Z\": dial tcp 127.0.0.1:45969: connect: connection refused"
	expectedErrorContent1 := "connect: connection refused"
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})

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

	go services.Delegation.Poll(pollPeriodInSeconds, server.URL, stopOnError, errorCh, interruptCh)

	time.Sleep(5 * time.Second)
	interruptCh <- struct{}{}
	err = <-errorCh
	fmt.Println(err.Error())

	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
}

func Test_SaveBulkDelegations_EmptySlice_NoError(t *testing.T) {
	var emptyList []dto.DelegationResponseFromAPI

	savedDelegations, err := services.SaveBulkDelegations(emptyList)

	require.NoError(t, err)
	assert.Empty(t, savedDelegations)
}

func Test_SaveBulkDelegations_ErrorFromRepositoryCreate_Error(t *testing.T) {
	var emptyList []dto.DelegationResponseFromAPI
	delegation1 := dto.DelegationResponseFromAPI{}
	listOneElement := append(emptyList, delegation1)

	errorMessage := "repository error: unexpected internal error"
	unexpectedError := errors.New("unexpected internal error")

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) error {
		return unexpectedError
	}
	repositories.Delegation = delegationRepositoryStub

	savedDelegations, err := services.SaveBulkDelegations(listOneElement)
	fmt.Println(err.Error())

	require.Error(t, err)
	assert.Equal(t, errorMessage, err.Error())
	assert.Empty(t, savedDelegations)
}

func Test_SaveBulkDelegations_OKList1Element_ReturnSavedDelegation(t *testing.T) {
	var emptyList []dto.DelegationResponseFromAPI
	delegation := dto.DelegationResponseFromAPI{}
	listOneElement := append(emptyList, delegation)

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) error {
		return nil
	}
	repositories.Delegation = delegationRepositoryStub

	savedDelegations, err := services.SaveBulkDelegations(listOneElement)

	require.NoError(t, err)
	assert.Len(t, savedDelegations, 1)
}

func Test_SaveBulkDelegations_OKList2Element_ReturnSavedDelegation(t *testing.T) {
	var delegations []dto.DelegationResponseFromAPI
	delegation := dto.DelegationResponseFromAPI{}
	delegations = append(delegations, delegation)
	delegations = append(delegations, delegation)

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) error {
		return nil
	}
	repositories.Delegation = delegationRepositoryStub

	savedDelegations, err := services.SaveBulkDelegations(delegations)

	require.NoError(t, err)
	assert.Len(t, savedDelegations, 2)
}
